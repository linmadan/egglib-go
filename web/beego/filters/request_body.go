package filters

import (
	"github.com/beego/beego/v2/server/web/context"
	"io/ioutil"
)

func CreateRequestBodyFilter() func(ctx *context.Context) {
	return func(ctx *context.Context) {
		if len(ctx.Input.RequestBody) == 0 {
			if ctx.Request.Method == "POST" || ctx.Request.Method == "PUT" {
				body, _ := ioutil.ReadAll(ctx.Request.Body)
				ctx.Input.SetData("requestBody", body)
				ctx.Request.Body.Close()
			}
		} else {
			ctx.Input.SetData("requestBody", ctx.Input.RequestBody)
		}
	}
}
