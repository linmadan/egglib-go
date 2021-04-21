package utils

import (
	"github.com/beego/beego/v2/server/web/context"
	"github.com/linmadan/egglib-go/core/application"
)

type JsonResponse map[string]interface{}

func ResponseData(ctx *context.Context, data interface{}) JsonResponse {
	jsonResponse := JsonResponse{}
	jsonResponse["code"] = 0
	jsonResponse["msg"] = "ok"
	jsonResponse["data"] = data
	ctx.Input.SetData("outputData", jsonResponse)
	return jsonResponse
}

func ResponseError(ctx *context.Context, err error) JsonResponse {
	serviceError := err.(*application.ServiceError)
	jsonResponse := JsonResponse{}
	jsonResponse["code"] = serviceError.Code
	jsonResponse["msg"] = serviceError.Error()
	println(err.Error())
	ctx.Input.SetData("outputData", jsonResponse)
	return jsonResponse
}
