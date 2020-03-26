package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

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
	userLoginDTO, _ := Check(c.Ctx)
	uid := userLoginDTO.Uid
	var groupUser models.GroupUser
	//判断是不是在群里面
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
	if limit > 500 {
		limit = 500
	}
	//查出群里面的成员
	var groupUsers []models.GroupUser
	page = CreateOffset(page, limit)
	ListByGroupId(groupId, page, limit, &groupUsers)
	logs.Info("groupsUsers 是", groupUsers)
	//取出uid
	uids := []int{}
	for _, group := range groupUsers {
		uids = append(uids, group.Uid)
	}
	logs.Info("该组里面的用户为", uids)
	//返回
	userMap, err := ListUserMapByUidIn(uids)
	if err != nil {
		logs.Error("userMap is err", err)
	}
	data := []models.GroupIndexListResVO{}
	for _, v := range groupUsers {
		listvo := new(models.GroupIndexListResVO)
		listvo.GroupId = v.GroupId
		listvo.Rank = v.Rank
		listvo.Remark = v.Remark
		listvo.User = (*userMap)[v.Uid]
		data = append(data, *listvo)
	}
	c.Data["json"] = models.ResponseOk(data)
	c.ServeJSON()

}

//创建群
func (c *GroupController) Create() {
	var requestGroup models.GroupSaveReqVO
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &requestGroup)
	if err != nil {
		logs.Error("解析json的时候异常了", err)
		c.Data["json"] = models.ResponseError(&models.PARAM_VERIFY_FALL)
		c.ServeJSON()
		return
	}
	logs.Info("前端传过来的数据解析为", requestGroup)
	userLoginDTO, _ := Check(c.Ctx)
	uid := userLoginDTO.Uid
	var avatar string
	if requestGroup.Avatar == "" {
		avatar = `http://47.100.212.206/avatar/defaultgroup.jpg`
	} else {
		avatar = requestGroup.Avatar
	}
	var remark string
	if requestGroup.Remark == "" {
		remark = "90部落"
	} else {
		remark = requestGroup.Remark
	}
	group := models.Group{
		Uid:       uid,
		Name:      requestGroup.Name,
		Avatar:    avatar,
		MemberNum: 1,
		Remark:    remark,
	}
	group.CreateTime = time.Now()
	group.ModifiedTime = group.CreateTime
	//新增组ResponseErrorCode
	id, err := OrmInsertAotoId(&group)
	if err != nil {
		logs.Error("插入group时候异常了", err)
		c.Data["json"] = models.ResponseError(&models.NOT_NETWORK)
		c.ServeJSON()
		return
	}
	//加入群列表
	groupUser := models.GroupUser{
		GroupId: id,
		Uid:     uid,
		Rank:    2, //群主
	}
	groupUser.CreateTime = time.Now()
	groupUser.ModifiedTime = groupUser.CreateTime
	_, err = OrmInsertAotoId(&groupUser)
	if err != nil {
		logs.Error("插入groupUser时候异常了", err)
		c.Data["json"] = models.ResponseError(&models.NOT_NETWORK)
		c.ServeJSON()
		return
	}
	c.Data["json"] = models.ResponseOk(group)
	c.ServeJSON()

}

//群消息列表
func (c *GroupController) MsgLists() {
	groupId, err := c.GetInt("groupId")
	if err != nil {
		logs.Error("参数groupId没有传")
		c.Data["json"] = models.ResponseError(&models.PARAM_VERIFY_FALL)
		c.ServeJSON()
		return
	}
	// 获取登入信息
	userLoginDTO, _ := Check(c.Ctx)
	uid := userLoginDTO.Uid

	//判断是不是在群里面
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
	if limit > 25 {
		limit = 25
	}
	//查出群里面的成员
	var groupMsgs []models.GroupMsg
	page = CreateOffset(page, limit)

	data := []models.GroupMsgListResVO{}
	num, err = ListByGroupIdAndCreateTime(groupId, groupUser.CreateTime, page, limit, &groupMsgs)
	if err != nil || num == 0 {
		c.Data["json"] = models.ResponseOk(data)
		c.ServeJSON()
		return
	}

	logs.Info("groupMsgs 是", groupMsgs)
	//取出uid
	uids := []int{}
	for _, group := range groupMsgs {
		uids = append(uids, group.SenderUid)
	}
	logs.Info("该组消息列表的用户为", uids)
	//返回
	userMap, err := ListUserMapByUidIn(uids)
	if err != nil {
		logs.Error("userMap is err", err)
	}

	for _, v := range groupMsgs {
		listvo := new(models.GroupMsgListResVO)
		listvo.GroupMsg = v
		listvo.UserInfoListResVO = (*userMap)[v.SenderUid]
		data = append(data, *listvo)
	}
	c.Data["json"] = models.ResponseOk(data)
	c.ServeJSON()

}

//MsgCreate 发送群消息(创建群消息)
func (c *GroupController) MsgCreate() {
	var requestGroupMsg models.GroupMsgCreateReqVO
	data := c.Ctx.Input.RequestBody
	err := json.Unmarshal(data, &requestGroupMsg)
	if err != nil {
		logs.Error("解析json的时候异常了", err)
		c.Data["json"] = models.ResponseError(&models.PARAM_VERIFY_FALL)
		c.ServeJSON()
		return
	}
	logs.Info("前端传过来的数据解析为", requestGroupMsg)
	userLoginDTO, _ := Check(c.Ctx)
	uid := userLoginDTO.Uid
	//判断是不是在群里面
	var groupUser models.GroupUser
	num := FindByGroupIdAndUid(requestGroupMsg.GroupId, uid, &groupUser)
	logs.Info("查到的结果为，groupUser", groupUser)
	if num == 0 {
		c.Data["json"] = models.ResponseErrorCode(models.PARAM_VERIFY_FALL.Code, "请先加入群~")
		c.ServeJSON()
		return
	}
	//群发消息
	b := SendGroupMsg(uid, requestGroupMsg.GroupId, models.GROUP, requestGroupMsg.MsgType, requestGroupMsg.MsgContent)
	if !b {
		c.Data["json"] = models.ResponseError(&models.NOT_NETWORK)
		c.ServeJSON()
		return
	}
	c.Data["json"] = models.ResponseOk([]int{})
	c.ServeJSON()
}

//SendGroupMsg 群发消息
func SendGroupMsg(uid, groupId, wsType, msgType int, msgContent string) bool {
	lastMsgContent := msgContent
	groupMsg := new(models.GroupMsg)
	groupMsg.SenderUid = uid
	groupMsg.GroupId = groupId
	groupMsg.MsgType = msgType
	groupMsg.MsgContent = msgContent
	groupMsg.CreateTime = time.Now()
	groupMsg.ModifiedTime = time.Now()
	//保存群消息
	msgId, err := OrmInsertAotoId(groupMsg)
	if err != nil {
		return false
	}
	logs.Info("插入了一条群消息msgId的为", msgId)
	//保存离线消息
	groupUser := new(models.GroupUser)
	groupUser.LastMsgContent = lastMsgContent
	groupUser.GroupId = groupId
	groupUser.UnMsgCount = 1
	groupUser.LastMsgTime = time.Now()
	groupUser.ModifiedTime = groupUser.LastMsgTime
	o := orm.NewOrm()
	_, err = o.QueryTable("group_user").Filter("group_id", groupId).Update(orm.Params{
		"LastMsgContent": lastMsgContent,
		"UnMsgCount":     1,
		"ModifiedTime":   time.Now(),
		"LastMsgTime":    time.Now(),
	})
	if err != nil {
		return false
	}
	//查询用户信息
	var user models.User
	err = o.QueryTable("user").Filter("uid", uid).One(&user)
	if err != nil {
		return false
	}
	logs.Info("用户信息为", user)
	//构建消息对象
	wsBaseReqVo := models.WSBaseReqVO{}
	wsBaseReqVo.Message.MsgContent = msgContent
	wsBaseReqVo.Message.ReceiveId = groupId
	wsBaseReqVo.Message.MsgType = msgType
	wsBaseReqVo.User.Avatar = user.Avatar
	wsBaseReqVo.User.Name = user.Name
	wsBaseReqVo.User.Remark = user.Remark
	wsBaseReqVo.User.Uid = user.Uid
	wsBaseReqVo.Type = wsType
	//查出群里面的成员
	var groupUsers []models.GroupUser
	ListByGroupId(groupId, 0, 500, &groupUsers)
	logs.Info("groupsUsers 是", groupUsers)
	for _, group := range groupUsers {
		if group.Uid != uid {
			SendMsg(group.Uid, wsBaseReqVo)
		}
	}
	return true

}

func ListByGroupIdAndCreateTime(groupId int, createTime time.Time, page int, limit int, groupMsgs *[]models.GroupMsg) (int, error) {
	o := orm.NewOrm()
	num, err := o.Raw(`select msg_id,group_id,sender_uid,msg_type,msg_content,create_time
	from group_msg
	where group_id = ? and create_time >= ?
	order by create_time desc
	limit ?,?`, groupId, createTime, page, limit).QueryRows(groupMsgs)
	if err != nil {
		logs.Error("查询发生了错误", err)
		return 0, err
	}
	return int(num), nil

}

func OrmInsertAotoId(dta interface{}) (int, error) {
	o := orm.NewOrm()
	id, err := o.Insert(dta)
	if err != nil {
		logs.Error("插入数据错误")
		return 0, err
	}
	logs.Info("插入插入数据生成id为", id)
	return int(id), nil

}

func ListUserMapByUidIn(uids []int) (*map[int]models.UserInfoListResVO, error) {
	userMap := make(map[int]models.UserInfoListResVO)
	users, err := ListUserByUidIn(uids)
	if err != nil {
		return &userMap, err
	}
	for _, value := range *users {
		resvo := models.UserInfoListResVO{}
		resvo.Uid = value.Uid
		resvo.Avatar = value.Avatar
		resvo.Name = value.Name
		resvo.Remark = value.Remark
		userMap[value.Uid] = resvo
	}
	return &userMap, err
}

func ListUserByUidIn(uids []int) (*[]models.User, error) {
	var users []models.User
	if len(uids) == 0 {
		return &users, errors.New("no data")
	}
	var str string = "?"
	for i := 1; i < len(uids); i++ {
		str += ",?"
	}
	o := orm.NewOrm()
	_, err := o.Raw("SELECT uid,name,avatar,remark FROM user where uid in ("+str+")", uids).QueryRows(&users)
	if err != nil {
		return &users, err
	}
	return &users, nil
}

func CreateOffset(page, limit int) int {
	return (page - 1) * limit
}

func ListByGroupId(groupId int, page int, limit int, groupUsers *[]models.GroupUser) error {
	o := orm.NewOrm()
	_, err := o.Raw(`select id,group_id,uid,remark,rank
	from group_user
	where group_id = ?
	limit ?,?`, groupId, page, limit).QueryRows(groupUsers)
	if err != nil {
		fmt.Println("查询发生了错误", err)
		return err
	}
	return nil
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
