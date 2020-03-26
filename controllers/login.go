package controllers

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/songhq211949/beego-api/models"
	"github.com/songhq211949/beego-api/utils"
)

type LoginController struct {
	beego.Controller
}

//账号密码登入
func (c *LoginController) ByPwd() {
	var loginVo models.UserLoginPwdReqVO
	if err := c.ParseForm(&loginVo); err != nil {
		logs.Error("解析json的时候异常了", err)
		c.Data["json"] = models.ResponseError(&models.PARAM_VERIFY_FALL)
		c.ServeJSON()
		return
	}
	o := orm.NewOrm()
	var user models.User
	err := o.QueryTable("user").Filter("name", loginVo.UserName).One(&user)
	if err != nil || !(user.Pwd == loginVo.Password) {
		c.Data["json"] = models.ResponseError(&models.PARAM_VERIFY_FALL)
		c.ServeJSON()
		return
	}
	userResVo := new(models.UserLoginResVO)
	uidStr := strconv.Itoa(user.Uid)
	userResVo.Sid = utils.CreateToken(uidStr)
	userResVo.Uid = user.Uid
	c.Data["json"] = models.ResponseOk(&userResVo)
	c.ServeJSON()
	return
}

