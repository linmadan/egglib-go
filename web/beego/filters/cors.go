package filters

import (
	"github.com/beego/beego/v2/server/web/context"
	"net/http"
)

func AllowCors() func(ctx *context.Context) {
	return func(ctx *context.Context) {
		ctx.Output.Header("Access-Control-Allow-Methods", "OPTIONS,DELETE,POST,GET,PUT,PATCH")
		//ctx.Output.Header("Access-Control-Max-Age", "3600")
		ctx.Output.Header("Access-Control-Allow-Headers", "*")
		ctx.Output.Header("Access-Control-Allow-Credentials", "true")
		ctx.Output.Header("Access-Control-Allow-Origin", "*") //origin
		if ctx.Input.Method() == http.MethodOptions {
			// options请求，返回200
			ctx.Output.SetStatus(http.StatusOK)
			_ = ctx.Output.Body([]byte("options support"))
		}
	}
}
