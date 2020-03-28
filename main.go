package main

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/songhq211949/beego-api/controllers"
	"github.com/songhq211949/beego-api/models"
	_ "github.com/songhq211949/beego-api/routers" //这里引入routers模块。会调用init方法
)

func main() {
	FilterLogin := func(ctx *context.Context) {
		logs.Info("经过了过滤器")

		_, isLogin := controllers.Check(ctx)
		if !isLogin {
			data := models.ResponseError(&models.LOGIN_VERIFY_FALL)
			result, err := json.Marshal(data)
			if err != nil {
				logs.Error("json序列化失败")
				return
			}
			ctx.ResponseWriter.Header()["Content-Type"] = []string{"application/json"}
			ctx.WriteString(string(result))
			logs.Error("登入失败")
		}
	}
	//^(?!(/user/login/*)).*$,beego不支持
	//^((?!login).)*$
	//统一校验用户需合法的登入状态，请注意，登入和在线是两回事，没有手动推出登入都可视为登入状态
	//添加过滤器会是跨域访问options请求，这里注意跨域过滤器和验证过滤器的前后顺序
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  false,
		AllowOrigins:     []string{"http://localhost:8080", ""},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	beego.InsertFilter("/api/group/*", beego.BeforeRouter, FilterLogin, true)
	beego.InsertFilter("/api/user/friendAsk/*", beego.BeforeRouter, FilterLogin, true)
	beego.InsertFilter("/api/user/friend/*", beego.BeforeRouter, FilterLogin, true)

	beego.Run()
}
