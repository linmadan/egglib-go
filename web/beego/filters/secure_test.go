package filters

import (
	"bytes"
	"fmt"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestSecureHandler(t *testing.T) {
	server := web.NewHttpSever()
	server.Get("/", func(context *context.Context) {
		context.WriteString("ok")
	})
	server.InsertFilter("/*", web.BeforeExec, SecureHandler())
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	SecureHttpRequest(r, "")
	w := httptest.NewRecorder()
	server.Handlers.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSecureHandler_Fail_Sign(t *testing.T) {
	server := web.NewHttpSever()
	server.Get("/", func(context *context.Context) {
		context.WriteString("ok")
	})
	server.InsertFilter("/*", web.BeforeExec, SecureHandler())
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	SecureHttpRequest(r, "")
	r.Header.Set("x-mmm-sign", "123456")
	w := httptest.NewRecorder()
	server.Handlers.ServeHTTP(w, r)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSecureHandler_Fail_Timestamp(t *testing.T) {
	server := web.NewHttpSever()
	server.Get("/", func(context *context.Context) {
		context.WriteString("ok")
	})
	server.InsertFilter("/*", web.BeforeExec, SecureHandler())
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	SecureHttpRequest(r, "")
	r.Header.Set("x-mmm-timestamp", fmt.Sprintf("%v", time.Now().Unix()-int64(defaultValidDuration)-1))
	w := httptest.NewRecorder()
	server.Handlers.ServeHTTP(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSecureHandler_OnFail(t *testing.T) {
	server := web.NewHttpSever()
	var count int32 = 0
	OnInvalidRequest := func(ctx *context.Context) {
		atomic.AddInt32(&count, 1)
	}
	server.Get("/", func(context *context.Context) {
		context.WriteString("ok")
	})
	server.InsertFilter("/*", web.BeforeExec, SecureHandler(WithOnInvalidRequest(OnInvalidRequest)))
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	SecureHttpRequest(r, "")
	r.Header.Set("x-mmm-sign", "123456")
	w := httptest.NewRecorder()
	server.Handlers.ServeHTTP(w, r)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, int32(1), count)
}

func TestSecureHandler_CheckRequestQuery_Get(t *testing.T) {
	server := web.NewHttpSever()
	server.Get("/*", func(context *context.Context) {
		//t.Log(context.Request.Header)
		context.WriteString("ok")
	})
	server.InsertFilter("/*", web.BeforeExec, SecureHandler(WithEnableCheckRequestQueryData(true), WithValidDuration(time.Second*100)))
	r := httptest.NewRequest(http.MethodGet, "/url?id=1&data=2", nil)
	SecureHttpRequest(r, "123456", WithEnableCheckRequestQueryData(true))
	//SecureHttpRequest(r, "123456")
	w := httptest.NewRecorder()
	server.Handlers.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSecureHandler_CheckRequestQuery_Post(t *testing.T) {
	server := web.NewHttpSever()
	server.Post("/*", func(context *context.Context) {
		context.WriteString("ok")
	})
	server.InsertFilter("/*", web.BeforeExec, SecureHandler(WithEnableCheckRequestQueryData(true)))
	r := httptest.NewRequest(http.MethodPost, "/url?id=1&data=2", nil)
	r.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(`{"data":"hello world!"}`))) // crc32 body
	SecureHttpRequest(r, "123456", WithEnableCheckRequestQueryData(true))
	w := httptest.NewRecorder()
	server.Handlers.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSecureHandler_CheckAppSecret_Post(t *testing.T) {
	var (
		key    = "app1"
		secret = "g47edfges745"
	)
	server := web.NewHttpSever()
	server.Post("/*", func(context *context.Context) {
		context.WriteString("ok")
	})
	server.InsertFilter("/*", web.BeforeExec, SecureHandler(WithEnableCheckRequestQueryData(true), WithAppKeySecret(key, secret)))
	r := httptest.NewRequest(http.MethodPost, "/url?id=1&data=2", nil)
	r.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("hello world!"))) // crc32 body
	SecureHttpRequest(r, "123456", WithEnableCheckRequestQueryData(true), WithAppKeySecret(key, secret))
	w := httptest.NewRecorder()
	server.Handlers.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}
