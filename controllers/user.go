package controllers

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/dgrijalva/jwt-go"
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
	o.QueryTable("user_profile").Filter("uid", uid).One(&userProfile)
	err := o.QueryTable("user").Filter("uid", uid).One(&user)
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

//FriendAskLists 新的朋友
func (c *UserController) FriendAskLists() {
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
	var userFriendAskListResVO []models.UserFriendAskListResVO
	var userFriendASks []models.UserFriendAsk
	err1 := UserFriendsAskByUid(uid, page, limit, &userFriendASks)
	if err1 != nil {
		c.Data["json"] = models.ResponseError(&models.NOT_NETWORK)
		c.ServeJSON()
		return
	}
	if len(userFriendASks) == 0 {
		c.Data["json"] = models.ResponseOk(userFriendAskListResVO)
		c.ServeJSON()
		return
	}
	var uids []int
	for _, value := range userFriendASks {
		uids = append(uids, value.FriendUid)
	}
	userMap, err2 := ListUserMapByUidIn(uids)
	if err2 != nil {
		c.Data["json"] = models.ResponseError(&models.NOT_NETWORK)
		c.ServeJSON()
		return
	}
	for _, value := range userFriendASks {
		askVo := new(models.UserFriendAskListResVO)
		askVo.UserFriendAsk = value
		askVo.User = (*userMap)[value.FriendUid]
		userFriendAskListResVO = append(userFriendAskListResVO, *askVo)
	}
	c.Data["json"] = models.ResponseOk(userFriendAskListResVO)
	c.ServeJSON()

}

//生成二维码
func (c *UserController) GetQRCheckCode() {
	loginUser, _ := Check(c.Ctx)
	uid := loginUser.Uid
	groupIdStr := strconv.Itoa(uid)
	tokentString := utils.CreateToken(groupIdStr)
	c.Data["json"] = models.ResponseOk(tokentString)
	c.ServeJSON()
}

//FriendAskcreate 添加好友请求
func (c *UserController) FriendAskcreate() {
	checkCode := c.GetString("checkCode")
	claims, ok := utils.ParseToken(checkCode)
	if !ok {
		c.Data["json"] = models.ResponseErrorCode(models.PARAM_VERIFY_FALL.Code, "二维码已过期~")
		c.ServeJSON()
		return
	}
	uid := claims.(jwt.MapClaims)["uid"]
	uidStr := uid.(string)
	uidInt64, _ := strconv.ParseInt(uidStr, 10, 64)
	friendId := int(uidInt64)
	loginUser, _ := Check(c.Ctx)
	loginId := loginUser.Uid
	//判断要添加的好友是不是自己
	if friendId == loginId {
		c.Data["json"] = models.ResponseErrorCode(models.PARAM_VERIFY_FALL.Code, "不能自己加自己")
		c.ServeJSON()
		return
	}
	//判断好友是不是已经是好友
	userFriend, _ := FindUserFriendByUidAndFriendUid(loginId, friendId)
	//即已经是好友了
	if userFriend.Id > 0 {
		c.Data["json"] = models.ResponseErrorCode(models.DATA_REPEAT.Code, "已经是好友了")
		c.ServeJSON()
		return
	}
	//插入ask表
	userFriendAsk := models.UserFriendAsk{
		Uid:       friendId,
		FriendUid: loginId,
		Status:    0,
	}
	userFriendAsk.CreateTime = time.Now()
	userFriendAsk.ModifiedTime = userFriendAsk.CreateTime
	if _, err := OrmInsertAotoId(&userFriendAsk); err != nil {
		c.Data["json"] = models.ResponseError(&models.NOT_NETWORK)
		c.ServeJSON()
		return
	}

	//发送在线消息
	var user models.User
	o := orm.NewOrm()
	o.QueryTable("user").Filter("uid", uid).One(&user)
	baseReqVO := new(models.WSBaseReqVO) //new出来的是指针
	baseReqVO.Type = models.FRIEND_ASK   //消息为请求好友
	baseReqVO.User.Avatar = user.Avatar
	baseReqVO.User.Name = user.Name
	baseReqVO.User.Remark = user.Remark
	baseReqVO.User.Uid = user.Uid
	baseReqVO.Message.MsgContent = "请求加为好友"
	baseReqVO.Message.MsgType = 0 //文字
	baseReqVO.Message.ReceiveId = friendId
	SendMsg(friendId, *baseReqVO)
	c.Data["json"] = models.ResponseOk([]int{})
	c.ServeJSON()
}

//FriendAskAck 确认或取消好友请求
func (c *UserController) FriendAskAck() {
	userLogin, _ := Check(c.Ctx)
	uid := userLogin.Uid
	//解析请求体
	var askAckVo models.UserFriendAskAckReqVO
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &askAckVo)
	if err != nil {
		logs.Error("解析json的时候异常了", err)
		c.Data["json"] = models.ResponseError(&models.PARAM_VERIFY_FALL)
		c.ServeJSON()
		return
	}
	id := askAckVo.Id
	status := askAckVo.Status
	o := orm.NewOrm()
	var userFriendAsk models.UserFriendAsk
	if err := o.QueryTable("user_friend_ask").Filter("id", id).One(&userFriendAsk); err != nil {
		c.Data["json"] = models.ResponseError(&models.NOT_NETWORK)
		c.ServeJSON()
		return
	}
	//已经处理过了
	if userFriendAsk.Status != 0 {
		c.Data["json"] = models.ResponseErrorCode(models.PARAM_VERIFY_FALL.Code, "已经取人或拒绝过了，请勿重复操作")
		c.ServeJSON()
		return
	}
	//更新状态
	userFriendAsk.Status = status
	if _, err := o.Update(&userFriendAsk, "status"); err != nil {
		c.Data["json"] = models.ResponseError(&models.NOT_NETWORK)
		c.ServeJSON()
		return
	}
	//如果是拒绝 ，不必发送消息
	if status == 2 {
		c.Data["json"] = models.ResponseOk([]int{})
		c.ServeJSON()
		return
	}
	//如果是接受，如果已经是好友了，就不必做什么
	friendUid := userFriendAsk.FriendUid
	//查不到会报错
	if _, err := FindUserFriendByUidAndFriendUid(uid, friendUid); err == nil {
		logs.Info("已经是好友了", uid, userFriendAsk.FriendUid)
		c.Data["json"] = models.ResponseOk([]int{})
		c.ServeJSON()
		return
	}
	//建立好友关系,插入user_friend表
	msgContent := "成为好友，现在开始聊吧~"
	userFriend1 := new(models.UserFriend)
	userFriend1.Uid = uid
	userFriend1.FriendUid = friendUid
	userFriend1.Remark = ""
	userFriend1.LastMsgContent = msgContent
	userFriend1.CreateTime = time.Now()
	userFriend1.ModifiedTime = time.Now()
	OrmInsertAotoId(userFriend1)
	userFriend2 := new(models.UserFriend)
	userFriend2.Uid = friendUid
	userFriend2.FriendUid = uid
	userFriend2.Remark = ""
	userFriend2.LastMsgContent = msgContent
	userFriend2.CreateTime = time.Now()
	userFriend2.ModifiedTime = time.Now()
	OrmInsertAotoId(userFriend2)
	//追加消息
	senderUid := uid
	toUid := friendUid
	if uid > friendUid {
		toUid = uid
		uid = friendUid
	}
	msg := new(models.UserFriendMsg)
	msg.Uid = uid
	msg.ToUid = toUid
	msg.SenderUid = senderUid
	msg.MsgContent = msgContent
	msg.MsgType = 0 //0代表文字
	msg.CreateTime = time.Now()
	OrmInsertAotoId(msg)
	//发送消息
	var user models.User
	o.QueryTable("user").Filter("uid", uid).One(&user)
	baseReqVO := new(models.WSBaseReqVO) //new出来的是指针
	baseReqVO.Type = models.FRIEND_ACK   //消息为请求好友
	baseReqVO.User.Avatar = user.Avatar
	baseReqVO.User.Name = user.Name
	baseReqVO.User.Remark = user.Remark
	baseReqVO.User.Uid = user.Uid
	baseReqVO.Message.MsgContent = msgContent
	baseReqVO.Message.MsgType = 0 //文字
	baseReqVO.Message.ReceiveId = friendUid
	SendMsg(friendUid, *baseReqVO)
	c.Data["json"] = models.ResponseOk([]int{})
	c.ServeJSON()
}

func (c *UserController) ClearFriendAskCount() {
	c.Data["json"] = models.ResponseOk([]int{})
	c.ServeJSON()
}

func UserFriendsAskByUid(uid, page, limit int, userFriendsAsk *[]models.UserFriendAsk) error {
	o := orm.NewOrm()
	_, err := o.Raw(`select id,uid,friend_uid,remark,status,create_time
	from user_friend_ask
	where uid = ?
	order by id desc
	limit ?,?`, uid, page, limit).QueryRows(userFriendsAsk)
	return err
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
