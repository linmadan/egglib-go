package filters

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/google/uuid"
	"hash/crc32"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	timestamp  = "timestamp"
	requestId  = "requestId"
	secureSign = "secureSign"
	token      = "token"
	nonce      = "nonce"
	appKey     = "appKey"
)

var (
	SecureKeyMap = map[string]string{
		timestamp:  "x-mmm-timestamp",
		requestId:  "x-mmm-uuid",
		secureSign: "x-mmm-sign",
		token:      "x-mmm-token",
		nonce:      "x-mmm-nonce",
		appKey:     "x-mmm-key",
	}
	defaultValidDuration = 5 * 60 * time.Second
)

var (
	errorInvalidTimestamp = errors.New("timestamp is invalid")
	errorInvalidSign      = errors.New("secure sign is invalid")
	errorInvalidRequest   = errors.New("request query data is invalid")
	errorInvalidAppSecret = errors.New("app secret is invalid")
)

// OnInvalidRequestCallBack 非法的接口访问进行回调,做后续处理
type (
	OnInvalidRequestCallBack func(ctx *context.Context)
	SecureOptions            struct {
		// 回调函数,当签名验证失败
		OnInvalidRequest OnInvalidRequestCallBack
		// 验证timestamp 有效时间范围
		ValidDuration time.Duration
		// 签名计算函数
		RequestSecureSignFunc func(req *http.Request, secureOptions *SecureOptions) string
		// 启用校验请求参数 （防止参数修改）
		EnableCheckRequestQueryData bool
		// 启用检查时间 （防止请求被多次重放）
		EnableCheckTimestamp bool
		// 给外部应用配置 appKey  + appSecret 作为签名加密的一部分
		AppKeySecretStore AppKeySecretStore
	}
	secureOption func(options *SecureOptions)

	// 存取 AppKey-Secret
	AppKeySecretStore interface {
		GetSecret(appKey string) (string, error)
		AddSecret(appKey, appSecret string) error
		Remove(appKey string) error
	}
	appKeySecretStore map[string]string
)

// SecureHandler api安全拦截处理,请求头里面包含有效的签名,才允许进行后续的接口访问
func SecureHandler(options ...secureOption) func(ctx *context.Context) {
	secureOptions := NewSecureOptions(options...)
	return func(ctx *context.Context) {
		if secureOptions.EnableCheckTimestamp {
			if err := validTimestamp(ctx); err != nil {
				ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
				ctx.WriteString(err.Error())
				return
			}
		}
		if secureOptions.EnableCheckRequestQueryData {
			if err := validRequestCrc32(ctx); err != nil {
				ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
				ctx.WriteString(err.Error())
				return
			}
		}
		if err := validSecureSign(ctx, secureOptions); err != nil {
			if secureOptions.OnInvalidRequest != nil {
				secureOptions.OnInvalidRequest(ctx)
			}
			ctx.ResponseWriter.WriteHeader(http.StatusUnauthorized)
			ctx.WriteString(errorInvalidSign.Error())
			return
		}
		return
	}
}

// SecureHttpRequest api安全注入,在访问其他服务时,自动注册有效的签名
func SecureHttpRequest(req *http.Request, tokenVal string, options ...secureOption) {
	secureOptions := NewSecureOptions(options...)
	now := time.Now().Unix()
	req.Header.Set(SecureKeyMap[timestamp], strconv.Itoa(int(now)))
	req.Header.Set(SecureKeyMap[requestId], uuid.New().String())
	req.Header.Set(SecureKeyMap[token], tokenVal)
	if secureOptions.EnableCheckRequestQueryData {
		reqCrc32 := computeRequestCrc32(req)
		if reqCrc32 > 0 {
			req.Header.Set(SecureKeyMap[nonce], fmt.Sprintf("%v", reqCrc32))
		}
	}
	req.Header.Set(SecureKeyMap[secureSign], secureOptions.RequestSecureSignFunc(req, secureOptions))
}

// RequestSecureSign 根据timestamp、requestId,token 计算http请求的签名
func RequestSecureSign(req *http.Request, secureOptions *SecureOptions) string {
	timestamp := req.Header.Get(SecureKeyMap[timestamp])
	requestId := req.Header.Get(SecureKeyMap[requestId])
	token := req.Header.Get(SecureKeyMap[token])
	nonce := req.Header.Get(SecureKeyMap[nonce])
	appKey := req.Header.Get(SecureKeyMap[appKey])
	if !secureOptions.EnableCheckRequestQueryData {
		nonce = ""
	}
	requestSecureSign := fmt.Sprintf("v!(MmM%v%v%v%vMmM)i^", timestamp, requestId, token, nonce) //timestamp + requestId + token + nonce
	// 外部应用接入
	if appSecret, err := secureOptions.AppKeySecretStore.GetSecret(appKey); err == nil && len(appSecret) > 0 {
		requestSecureSign = fmt.Sprintf("%v%v%v%v%v", appSecret, timestamp, requestId, token, nonce) //appSecret  + timestamp + requestId + token + nonce
	}
	sha256 := sha256.New()
	sha256.Write([]byte(requestSecureSign))
	signHex := hex.EncodeToString(sha256.Sum(nil))
	return signHex
}

func validTimestamp(ctx *context.Context) error {
	timestamp := ctx.Request.Header.Get(SecureKeyMap[timestamp])
	if len(timestamp) == 0 {
		return errorInvalidTimestamp
	}
	timestampInt, err := strconv.Atoi(timestamp)
	if err != nil {
		return err
	}
	if (int64(timestampInt) + int64(defaultValidDuration)) < time.Now().Unix() {
		return errorInvalidTimestamp
	}
	return nil
}

func validSecureSign(ctx *context.Context, secureOptions *SecureOptions) error {
	var clientSign string
	if clientSign = ctx.Request.Header.Get(SecureKeyMap[secureSign]); len(clientSign) == 0 {
		return errorInvalidSign
	}
	serveSign := secureOptions.RequestSecureSignFunc(ctx.Request, secureOptions) //secureOptions.RequestSecureSignFunc(ctx.Request,secureOptions)
	if !strings.EqualFold(serveSign, clientSign) {
		return errorInvalidSign
	}
	return nil
}

func validRequestCrc32(ctx *context.Context) error {
	nonce := ctx.Request.Header.Get(SecureKeyMap[nonce])
	if len(nonce) == 0 {
		return nil
	}
	nonceInt, err := strconv.Atoi(nonce)
	if err != nil {
		return err
	}
	var nonceCrc32 uint32 = computeRequestCrc32(ctx.Request)
	if nonceCrc32 > 0 && nonceCrc32 != uint32(nonceInt) {
		return errorInvalidRequest
	}
	return nil
}

func computeRequestCrc32(req *http.Request) uint32 {
	var requestCrc32 uint32
	if req.Method == http.MethodGet || req.Method == http.MethodDelete {
		if len(req.URL.RawQuery) == 0 { // 没有请求参数，不需要计算
			req.Header.Set(SecureKeyMap[nonce], "")
			return 0
		}
		requestCrc32 = crc32.ChecksumIEEE([]byte(req.URL.RawQuery)) //eg: id=1&data=2
	} else if req.Method == http.MethodPost || req.Method == http.MethodPut {
		var body []byte
		body, req.Body, _ = dumpReadCloser(req.Body)
		if len(body) == 0 { // 没有请求body，不需要计算
			req.Header.Set(SecureKeyMap[nonce], "")
			return 0
		}
		requestCrc32 = crc32.ChecksumIEEE(body)
	} else {
		req.Header.Set(SecureKeyMap[nonce], "")
		return 0
	}
	return requestCrc32
}

func dumpReadCloser(reader io.ReadCloser) ([]byte, io.ReadCloser, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(reader, &buf)
	tmpRead1, tmpRead2 := ioutil.NopCloser(tee), ioutil.NopCloser(&buf)
	data, err := ioutil.ReadAll(tmpRead1)
	reader = tmpRead2
	return data, reader, err
}

func NewSecureOptions(options ...secureOption) *SecureOptions {
	secureOptions := SecureOptions{
		OnInvalidRequest:            nil,
		ValidDuration:               defaultValidDuration,
		RequestSecureSignFunc:       RequestSecureSign,
		EnableCheckRequestQueryData: false,
		EnableCheckTimestamp:        true,
		AppKeySecretStore:           appKeySecretStore(map[string]string{}),
	}
	for i := range options {
		options[i](&secureOptions)
	}
	return &secureOptions
}

//EnableCheckRequestQueryData 启用请求参数校验，防止接口请求参数被更改 ，get/delete 请求 计算url的crc32 , post,put 计算body的crc32
func WithEnableCheckRequestQueryData(flag bool) secureOption {
	return func(options *SecureOptions) {
		options.EnableCheckRequestQueryData = flag
	}
}

func WithOnInvalidRequest(callBack OnInvalidRequestCallBack) secureOption {
	return func(options *SecureOptions) {
		options.OnInvalidRequest = callBack
	}
}

func WithValidDuration(validDuration time.Duration) secureOption {
	return func(options *SecureOptions) {
		options.ValidDuration = validDuration
	}
}

func WithRequestSecureSignFunc(requestSecureSignFunc func(req *http.Request, secureOptions *SecureOptions) string) secureOption {
	return func(options *SecureOptions) {
		options.RequestSecureSignFunc = requestSecureSignFunc
	}
}

func WithEnableCheckTimestamp(flag bool) secureOption {
	return func(options *SecureOptions) {
		options.EnableCheckTimestamp = flag
	}
}

func WithAppKeySecretStore(store AppKeySecretStore) secureOption {
	return func(options *SecureOptions) {
		options.AppKeySecretStore = store
	}
}

func WithAppKeySecret(appKey, appSecret string) secureOption {
	return func(options *SecureOptions) {
		options.AppKeySecretStore.AddSecret(appKey, appSecret)
	}
}

func (store appKeySecretStore) GetSecret(appKey string) (string, error) {
	if v, ok := store[appKey]; ok {
		return v, nil
	}
	return "", errorInvalidAppSecret
}

func (store appKeySecretStore) AddSecret(appKey, appSecret string) error {
	store[appKey] = appSecret
	return nil
}

func (store appKeySecretStore) Remove(appKey string) error {
	delete(store, appKey)
	return nil
}
