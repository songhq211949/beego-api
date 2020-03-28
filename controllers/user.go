package controllers

import (
	"encoding/json"
	"strconv"
	"time"

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

//FriendMsgLists
func (c *UserController) FriendMsgLists() {
	// 获取登入信息
	userLoginDTO, _ := Check(c.Ctx)
	uid := userLoginDTO.Uid
	senderUid, err := c.GetInt("senderUid")
	if err != nil {
		c.Data["json"] = models.ResponseError(&models.PARAM_VERIFY_FALL)
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
	if limit > 50 {
		limit = 50
	}
	page = CreateOffset(page, limit)
	// 把最小的那个 用户ID作为查询条件
	toUid := senderUid
	if uid > senderUid {
		toUid = uid
		uid = senderUid
	}
	var userFriendMsgs []models.UserFriendMsg
	err = ListUserFriendMsgByUidAndToUid(uid, toUid, page, limit, &userFriendMsgs)
	if err != nil {
		c.Data["json"] = models.ResponseError(&models.NOT_NETWORK)
		c.ServeJSON()
		return
	}
	c.Data["json"] = models.ResponseOk(&userFriendMsgs)
	c.ServeJSON()

}

//FriendMsgCreate 发送朋友消息
func (c *UserController) FriendMsgCreate() {
	var userFriendMsgSaveReqVO models.UserFriendMsgSaveReqVO
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &userFriendMsgSaveReqVO)
	if err != nil {
		logs.Error("解析json的时候异常了", err)
		c.Data["json"] = models.ResponseError(&models.PARAM_VERIFY_FALL)
		c.ServeJSON()
		return
	}
	// 获取登入信息
	userLoginDTO, _ := Check(c.Ctx)
	uid := userLoginDTO.Uid

	lastMsgContent := userFriendMsgSaveReqVO.MsgContent
	receiverUid := userFriendMsgSaveReqVO.ReceiverUid
	//判断是不是朋友 userFriend1是以发送消息为主体的 userFriend2为接收消息者为主体的
	userFriend1, err1 := FindUserFriendByUidAndFriendUid(uid, receiverUid)
	userFriend2, err2 := FindUserFriendByUidAndFriendUid(receiverUid, uid)
	if err1 != nil || err2 != nil {
		c.Data["json"] = models.ResponseErrorCode(models.PARAM_VERIFY_FALL.Code, "该用户还不是你的好友~")
		c.ServeJSON()
		return
	}
	senderUid := uid
	toUid := receiverUid
	if uid > receiverUid {
		toUid = uid
		uid = receiverUid
	}
	userFriendMsg := models.UserFriendMsg{
		Uid:        uid,
		ToUid:      toUid,
		SenderUid:  senderUid,
		MsgContent: userFriendMsgSaveReqVO.MsgContent,
		MsgType:    userFriendMsgSaveReqVO.MsgType,
		CreateTime: time.Now(),
	}
	_, err = OrmInsertAotoId(&userFriendMsg)
	if err != nil {
		c.Data["json"] = models.ResponseError(&models.NOT_NETWORK)
		c.ServeJSON()
		return
	}
	//更新当前用户user_friend表
	userFriend1.LastMsgContent = lastMsgContent
	userFriend1.UnMsgCount = 0
	userFriend1.CreateTime = time.Now()
	userFriend1.ModifiedTime = userFriend1.CreateTime
	err3 := UpdateUserFriend(userFriend1)

	//更新接收用户user_friend表
	userFriend2.LastMsgContent = lastMsgContent
	userFriend2.UnMsgCount = userFriend1.UnMsgCount + 1 //未读消息数加1
	userFriend2.CreateTime = time.Now()
	userFriend2.ModifiedTime = userFriend1.CreateTime
	err4 := UpdateUserFriend(userFriend2)
	if err3 != nil || err4 != nil {
		logs.Error("更新user_friend表失败", err3, err4)
		c.Data["json"] = models.ResponseError(&models.NOT_NETWORK)
		c.ServeJSON()
		return
	}
	//发送在线消息
	var user models.User
	o := orm.NewOrm()
	err = o.QueryTable("user").Filter("uid", uid).One(&user)
	baseReqVO := new(models.WSBaseReqVO) //new出来的是指针
	baseReqVO.Type = models.FRIEND       //消息为朋友消息
	baseReqVO.User.Avatar = user.Avatar
	baseReqVO.User.Name = user.Name
	baseReqVO.User.Remark = user.Remark
	baseReqVO.User.Uid = user.Uid
	baseReqVO.Message.MsgContent = lastMsgContent
	baseReqVO.Message.MsgType = userFriendMsgSaveReqVO.MsgType //文字
	baseReqVO.Message.ReceiveId = receiverUid
	SendMsg(receiverUid, *baseReqVO)
	c.Data["json"] = models.ResponseOk([]int{})
	c.ServeJSON()
}

//ClearUnMsgCount 清除未读消息
func (c *UserController) ClearUnMsgCount() {
	var userFriendMsgClearMsgCountReqVO models.UserFriendMsgClearMsgCountReqVO
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &userFriendMsgClearMsgCountReqVO)
	if err != nil {
		logs.Error("解析json的时候异常了", err)
		c.Data["json"] = models.ResponseError(&models.PARAM_VERIFY_FALL)
		c.ServeJSON()
		return
	}
	receiverId := userFriendMsgClearMsgCountReqVO.ReceiverUid
	re, _ := Check(c.Ctx)
	uid := re.Uid
	userFriend, err := FindUserFriendByUidAndFriendUid(uid, receiverId)
	if err != nil {
		c.Data["json"] = models.ResponseError(&models.NOT_NETWORK)
		c.ServeJSON()
		return
	}
	userFriend.UnMsgCount = 0
	if err := UpdateUserFriend(userFriend); err != nil {
		c.Data["json"] = models.ResponseError(&models.NOT_NETWORK)
		c.ServeJSON()
		return
	}
	c.Data["json"] = models.ResponseOk([]int{})
	c.ServeJSON()
}

func UpdateUserFriend(userFriend *models.UserFriend) error {
	o := orm.NewOrm()
	_, err := o.Update(userFriend)
	return err

}
func FindUserFriendByUidAndFriendUid(uid, receiverId int) (*models.UserFriend, error) {
	o := orm.NewOrm()
	var userFriend models.UserFriend
	err := o.QueryTable("user_friend").Filter("uid", uid).Filter("friend_uid", receiverId).One(&userFriend)
	return &userFriend, err
}

//ListUserFriendMsgByUidAndToUid 查询朋友消息列表
func ListUserFriendMsgByUidAndToUid(uid, toUid, page, limit int, userFriendMsgs *[]models.UserFriendMsg) error {
	o := orm.NewOrm()
	_, err := o.Raw(`select msg_id,sender_uid,msg_type,msg_content,create_time
	from user_friend_msg
	where uid = ?
	and to_uid =?
	order by msg_id desc
	limit ?,?`, uid, toUid, page, limit).QueryRows(userFriendMsgs)
	return err
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
	//logs.Info("从参数获取的uid", uid, "sid", sid)
	if uid == "" || sid == "" {
		//从cookie中获取参数
		uid = ctx.GetCookie("UID")
		sid = ctx.GetCookie("SID")
	}
	//logs.Info("从cookie获取的uid", uid, "sid", sid)
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
