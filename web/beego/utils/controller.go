package utils

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/linmadan/egglib-go/utils/json"
)

type (
	Controller struct {
		beego.Controller
	}
)

func (c Controller) Response(data interface{}, err error) {
	var response JsonResponse
	if err != nil {
		response = ResponseError(c.Ctx, err)
	} else {
		response = ResponseData(c.Ctx, data)
	}
	c.Data["json"] = response
	c.ServeJSON()
}

func (c Controller) Unmarshal(v interface{}) error {
	body := c.Ctx.Input.RequestBody
	if len(body) == 0 {
		body = []byte("{}")
	}
	return json.Unmarshal(body, v)
}

func NewController(ctx *context.Context) Controller {
	return Controller{
		beego.Controller{
			Ctx:  ctx,
			Data: make(map[interface{}]interface{}),
		},
	}
}
