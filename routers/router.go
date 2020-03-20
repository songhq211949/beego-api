package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/songhq211949/beego-api/controllers"
)

func init() {
	//wuhan
	beego.Router("/wuhan/list", &controllers.WuhanController{}, "get:Lists")
}

var filterFunc = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302, "/login")
		return
	}
}
