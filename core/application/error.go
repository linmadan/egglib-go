package application

import (
	"os"
	"strconv"
	"strings"
)

var baseCode = 1
var baseCodeMultiple = 0

func init() {
	if os.Getenv("ERROR_BASE_CODE") != "" {
		if i, err := strconv.Atoi(os.Getenv("ERROR_BASE_CODE")); err == nil {
			baseCode = i
		}
	}
	if os.Getenv("ERROR_BASE_CODE_MULTIPLE") != "" {
		if i, err := strconv.Atoi(os.Getenv("ERROR_BASE_CODE_MULTIPLE")); err == nil {
			baseCodeMultiple = i
		}
	}
}

const (
	INTERNAL_SERVER_ERROR = iota + 1
	BAD_GATEWAY_ERROR
	TRANSACTION_ERROR
	ARG_ERROR
	RES_NO_FIND_ERROR
	BUSINESS_ERROR
)

type ServiceError struct {
	Code    int
	Message string
}

func (serviceError *ServiceError) Error() string {
	return serviceError.Message
}

func ThrowError(serviceErrorNo int, attachMessages ...string) error {
	switch serviceErrorNo {
	case INTERNAL_SERVER_ERROR:
		return &ServiceError{
			Code:    501 + baseCode*baseCodeMultiple,
			Message: strings.Join(append([]string{"内部服务出错："}, attachMessages...), ""),
		}
	case BAD_GATEWAY_ERROR:
		return &ServiceError{
			Code:    502 + baseCode*baseCodeMultiple,
			Message: strings.Join(append([]string{"服务网关错误："}, attachMessages...), ""),
		}
	case TRANSACTION_ERROR:
		return &ServiceError{
			Code:    503 + baseCode*baseCodeMultiple,
			Message: strings.Join(append([]string{"事务处理错误："}, attachMessages...), ""),
		}
	case ARG_ERROR:
		return &ServiceError{
			Code:    504 + baseCode*baseCodeMultiple,
			Message: strings.Join(append([]string{"参数认证错误："}, attachMessages...), ""),
		}
	case RES_NO_FIND_ERROR:
		return &ServiceError{
			Code:    505 + baseCode*baseCodeMultiple,
			Message: strings.Join(append([]string{"找不到相关资源："}, attachMessages...), ""),
		}
	case BUSINESS_ERROR:
		return &ServiceError{
			Code:    506 + baseCode*baseCodeMultiple,
			Message: strings.Join(append([]string{"业务逻辑错误："}, attachMessages...), ""),
		}
	default:
		return &ServiceError{
			Code:    500 + baseCode*baseCodeMultiple,
			Message: strings.Join(append([]string{"未定义错误："}, attachMessages...), ""),
		}
	}
}
