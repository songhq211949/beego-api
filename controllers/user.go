package controllers

import (
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/songhq211949/beego-api/models"
	"github.com/songhq211949/beego-api/utils"
)

type UserController struct {
	beego.Controller
}

//LoginInfo  登入后的用户信息
func (c *UserController) LoginInfo() {
	// 获取登入信息
	userLoginDTO, _ := Check(c.Ctx)
	uid := userLoginDTO.Uid
	o := orm.NewOrm()
	var userProfile models.UserProfile
	var user models.User
	err := o.QueryTable("user_profile").Filter("uid", uid).One(&userProfile)
	if err != nil {
		c.Data["json"] = models.ResponseError(&models.NOT_NETWORK)
		c.ServeJSON()
		return
	}
	err = o.QueryTable("user").Filter("uid", uid).One(&user)
	if err != nil {
		c.Data["json"] = models.ResponseError(&models.NOT_NETWORK)
		c.ServeJSON()
		return
	}
	data := models.UserInfoResVO{
		Profile: userProfile,
	}
	data.Avatar = user.Avatar
	data.CreateTime = user.CreateTime
	data.ModifiedTime = user.ModifiedTime
	data.Name = user.Name
	data.Remark = user.Remark
	data.Uid = uid
	c.Data["json"] = models.ResponseOk(&data)
	c.ServeJSON()
	return

}

//FriendLists 好友列表
func (c *UserController) FriendLists() {
	// 获取登入信息
	userLoginDTO, _ := Check(c.Ctx)
	uid := userLoginDTO.Uid
	page, err := c.GetInt("page")
	if err != nil {
		page = 1
	}
	limit, err := c.GetInt("limit")
	if err != nil {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}
	page = CreateOffset(page, limit)
	var userFriends []models.UserFriend
	queryErr := UserFriendsByUid(uid, page, limit, &userFriends)
	if queryErr != nil {
		c.Data["json"] = models.ResponseError(&models.NOT_NETWORK)
		c.ServeJSON()
		return
	}
	//取出uid
	uids := []int{}
	for _, userFriend := range userFriends {
		uids = append(uids, userFriend.FriendUid)
	}
	//返回
	userMap, err := ListUserMapByUidIn(uids)
	data := []models.UserFriendListInfoResVO{}
	for _, userFriend := range userFriends {
		userFriendVo := new(models.UserFriendListInfoResVO)
		userFriendVo.User = (*userMap)[userFriend.FriendUid]
		userFriendVo.UserFriend = userFriend
		data = append(data, *userFriendVo)
	}
	c.Data["json"] = models.ResponseOk(&data)
	c.ServeJSON()
}

//查询UserFriend
func UserFriendsByUid(uid, page, limit int, userFriends *[]models.UserFriend) error {
	o := orm.NewOrm()
	_, err := o.Raw(`select id,uid,friend_uid,remark,un_msg_count,last_msg_content,modified_time
	from user_friend
	where uid = ?
	limit ?,?`, uid, page, limit).QueryRows(userFriends)
	return err
}

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
	logs.Info("从cookie获取的uid", uid, "sid", sid)
	//还为空的话，则不再登入状态
	if uid == "" || sid == "" {
		return &userLoginDTO, false
	}
	result := utils.CheckToken(uid, sid)
	number, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		logs.Error("uid 不合法", err)
	}
	userLoginDTO.Uid = int(number)
	return &userLoginDTO, result
}
