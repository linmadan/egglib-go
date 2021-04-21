package filters

import (
	"github.com/beego/beego/v2/server/web/context"
	"github.com/linmadan/egglib-go/log"
)

func CreateResponseLogFilter(logger log.Logger) func(ctx *context.Context) {
	return func(ctx *context.Context) {
		var append = make(map[string]interface{})
		append["framework"] = "beego"
		append["method"] = ctx.Input.Method()
		append["url"] = ctx.Input.URL()
		append["outputData"] = ctx.Input.GetData("outputData")
		logger.Debug("http响应", append)
	}
}
