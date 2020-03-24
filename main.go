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
	//统一校验用户需合法的登入状态，请注意，登入和在线是两回事，没有手动推出登入都可视为登入状态
	beego.InsertFilter("*", beego.BeforeRouter, FilterLogin, true)
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	beego.Run()
}
