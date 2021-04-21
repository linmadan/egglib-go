package beego

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/linmadan/egglib-go/utils/json"
	"github.com/linmadan/egglib-go/web/beego/utils"
)

type (
	BaseController struct {
		web.Controller
	}
)

func (c BaseController) Response(data interface{}, err error) {
	var response utils.JsonResponse
	if err != nil {
		response = utils.ResponseError(c.Ctx, err)
	} else {
		response = utils.ResponseData(c.Ctx, data)
	}
	c.Data["json"] = response
	c.ServeJSON()
}

func (c BaseController) Unmarshal(v interface{}) error {
	body := c.Ctx.Input.RequestBody
	if len(body) == 0 {
		body = []byte("{}")
	}
	return json.Unmarshal(body, v)
}
