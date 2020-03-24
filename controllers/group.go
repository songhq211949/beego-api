package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/songhq211949/beego-api/models"
)

//GroupController 组消息管理
type GroupController struct {
	beego.Controller //继承beegoContrller
}

//Lists 群用户列表
func (c *GroupController) Lists() {
	groupId, err := c.GetInt("groupId")
	if err != nil {
		logs.Error("参数groupId没有传")
		c.Data["json"] = models.ResponseError(&models.PARAM_VERIFY_FALL)
		c.ServeJSON()
		return
	}
	// 验证登录
	userLoginDTO, isLogin := Check(c.Ctx)
	if !isLogin {
		c.Data["json"] = models.ResponseError(&models.LOGIN_VERIFY_FALL)
		c.ServeJSON()
		return
	}
	uid := userLoginDTO.Uid
	var groupUser models.GroupUser
	num := FindByGroupIdAndUid(groupId, uid, &groupUser)
	logs.Info("查到的结果为，groupUser", groupUser)
	if num == 0 {
		c.Data["json"] = models.ResponseErrorCode(models.PARAM_VERIFY_FALL.Code, "请先加入群~")
		c.ServeJSON()
		return
	}

	page, err := c.GetInt("page")
	if err != nil {
		page = 1
	}
	limit, err := c.GetInt("limit")
	if err != nil {
		limit = 20
	}
	if limit>500{
		limit =500
	}

}

//FindByGroupIdAndUid 根据groudId和uid查询group_user表
func FindByGroupIdAndUid(groupId, uid int, groupUser *models.GroupUser) int {
	o := orm.NewOrm()
	err := o.Raw(`select id,rank,create_time  
	from group_user 
	where group_id = ?  and uid = ? limit 1`, groupId, uid).QueryRow(groupUser)
	if err != nil {
		fmt.Println("查询发生了错误", err)
		return 0
	}
	if groupUser.Id == 0 {
		return 0
	}
	return 1
}
