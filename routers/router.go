package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/songhq211949/beego-api/controllers"
)

//引入包就会初始化执行
func init() {
	//武汉疫情数据
	beego.Router("/api/wuhan/list", &controllers.WuhanController{}, "get:Lists")
	//建立websocket连接
	beego.Router("/ws", &controllers.WebSocketController{}, "get:Connect")
	//向指定的用户发送消息,仅用来测试长连接发送为加密的消息，uid需在线 msg为明文发送的消息
	beego.Router("/api/ws/send", &controllers.WebSocketController{}, "get:SendMessage")
	//群用户列表
	beego.Router("/api/group/lists", &controllers.GroupController{}, "get:Lists")
	//群创建
	beego.Router("/api/group/create", &controllers.GroupController{}, "post:Create")
	//群消息列表
	beego.Router("/api/group/msg/lists", &controllers.GroupController{}, "get:MsgLists")
	//群发送消息
	beego.Router("/api/group/msg/create", &controllers.GroupController{}, "post:MsgCreate")
	//群列表（用户所在哪些群）
	beego.Router("/api/group/user/lists", &controllers.GroupController{}, "get:Userlists")
	//读取未读群消息
	beego.Router("/api/group/user/clearUnMsgCount", &controllers.GroupController{}, "post:ClearUnMsgCount")
	//群二维码
	beego.Router("/api/group/user/getCheckCode", &controllers.GroupController{}, "get:GetCheckCode")
	//通过群二维码加入该群
	beego.Router("/api/group/user/create", &controllers.GroupController{}, "post:GroupUserCreate")

	//登入by用户名和密码
	beego.Router("/api/user/login/byPwd", &controllers.LoginController{}, "post:ByPwd")
	//登入后用户的信息
	beego.Router("/api/user/loginInfo", &controllers.UserController{}, "get:LoginInfo")
	//用户列表
	beego.Router("/api/user/friend/lists", &controllers.UserController{}, "get:FriendLists")
	//用户消息列表
	beego.Router("/api/user/friendMsg/lists", &controllers.UserController{}, "get:FriendMsgLists")
	//朋友消息发送
	beego.Router("/api/user/friendMsg/create", &controllers.UserController{}, "post:FriendMsgCreate")
	//读取未读朋友消息 /api/user/friendMsg/clearUnMsgCount
	beego.Router("/api/user/friendMsg/clearUnMsgCount", &controllers.UserController{}, "post:ClearUnMsgCount")
}

var filterFunc = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302, "/login")
		return
	}
}
