package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/songhq211949/beego-api/controllers"
)
//引入包就会初始化执行
func init() {
	//武汉疫情数据
	beego.Router("/wuhan/list", &controllers.WuhanController{}, "get:Lists")
	//建立websocket连接
	beego.Router("/ws", &controllers.WebSocketController{}, "get:Connect")
	//向指定的用户发送消息
	beego.Router("/ws/send", &controllers.WebSocketController{}, "get:SendMessage")
	//群用户列表
	beego.Router("/group/lists", &controllers.GroupController{}, "get:Lists")
}

var filterFunc = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302, "/login")
		return
	}
}
