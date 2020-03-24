package controllers

import (
	"strconv"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/songhq211949/beego-api/models"
	"github.com/songhq211949/beego-api/utils"
)

//Check 验证用户登入
func Check(ctx *context.Context) (*models.UserLoginDTO, bool) {
	userLoginDTO := models.UserLoginDTO{}
	//优先从请求参数中获取uid和sid
	uid := ctx.Request.FormValue("UID")
	sid := ctx.Request.FormValue("SID")
	logs.Info("从参数获取的uid", uid, "sid", sid)
	if uid == "" || sid == "" {
		//从cookie中获取参数
		uid = ctx.GetCookie("UID")
		sid = ctx.GetCookie("SID")
	}
	//还为空的话，则不再登入状态
	if uid == "" || sid == "" {
		return &userLoginDTO, false
	}
	result := utils.CheckToken(uid, sid)
	number,err:=strconv.ParseInt(uid, 10, 64)
	if err!=nil{
		logs.Error("uid 不合法",err)
	}
	userLoginDTO.Uid =int(number)
	return &userLoginDTO, result
}
