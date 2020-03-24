package controllers

import (
	"encoding/json"
	"errors"
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
