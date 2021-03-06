package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

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

//ByQq qq第三方登入
func (c *LoginController) ByQq() {
	//code为qq第三方登入后获取的Authorization code ,redirect_uri为跳转时的uri
	//accessToken即为我们后端取到的token
	code := c.GetString("code")
	redirect_uri := c.GetString("redirect_uri")
	accessToken := GetAccessToken(code, redirect_uri)
	if accessToken == "" {
		c.Data["json"] = models.ResponseErrorCode(models.NOT_NETWORK.Code, "accessToken 获取失败~")
		c.ServeJSON()
		return
	}
	openIdVo := GetOpenID(accessToken)
	userInfo := GetUserInfo(accessToken, openIdVo.Openid)
	logs.Info("调用qq第三方接口都正常")
	o := orm.NewOrm()
	var userQQ models.UserQq
	err := o.QueryTable("user_qq").Filter("openid", openIdVo.Openid).One(&userQQ)
	if err != nil {
		//没有赋予过值
		logs.Error("查询user_qq表的发生错误，openid为", openIdVo.Openid, "发生的err错误是", err)
		//创建user和user_qq
		user := new(models.User)
		nickname := userInfo.Nickname
		user.Name = nickname
		user.Avatar = userInfo.Figureurl_qq_2
		user.Remark = "历尽千帆归来仍是少年"
		user.CreateTime = time.Now()
		user.ModifiedTime = user.CreateTime
		uid, err := OrmInsertAotoId(user)
		if err != nil {
			logs.Error("qq登入插入用户失败")
			c.Data["json"] = models.ResponseErrorCode(models.NOT_NETWORK.Code, "qq登入插入user失败~")
			c.ServeJSON()
			return
		}
		//保存user_aa
		userSaveQQ := new(models.UserQq)
		userSaveQQ.Openid = openIdVo.Openid
		userSaveQQ.Uid = uid
		userSaveQQ.CreateTime = time.Now()
		if _, err2 := OrmInsertAotoId(userSaveQQ); err2 != nil {
			logs.Error("qq登入插入user_qq失败")
			c.Data["json"] = models.ResponseErrorCode(models.NOT_NETWORK.Code, "qq登入插入user_qq失败~")
			c.ServeJSON()
			return
		}
		userResVo := new(models.UserLoginResVO)
		uidStr := strconv.Itoa(uid)
		userResVo.Sid = utils.CreateToken(uidStr)
		userResVo.Uid = uid
		c.Data["json"] = models.ResponseOk(&userResVo)
		c.ServeJSON()
		return
	} else {
		//查到了值,给出uid
		uid := userQQ.Uid
		userResVo := new(models.UserLoginResVO)
		uidStr := strconv.Itoa(uid)
		userResVo.Sid = utils.CreateToken(uidStr)
		userResVo.Uid = uid
		c.Data["json"] = models.ResponseOk(&userResVo)
		c.ServeJSON()
	}
}
func GetUserInfo(accessToken, openId string) *models.QqUserInfoResVO {
	appId := beego.AppConfig.String("qqAuthAppid")
	vo := new(models.QqUserInfoResVO)
	url := "https://graph.qq.com/user/get_user_info?" + "access_token=" + accessToken +
		"&openid=" + openId + "&oauth_consumer_key=" + appId
	resp, err := http.Get(url)
	logs.Info("GetUserInfo 的url是", url)
	if err != nil {
		logs.Error("请求", url, "失败", err)
		return vo
	}
	defer resp.Body.Close()
	logs.Info("qq返回的消息是", resp)
	//获取响应内容
	resultByte, err := ioutil.ReadAll(resp.Body)
	result := string(resultByte)
	logs.Info("GetUserInfo 的返回体是", result)
	json.Unmarshal([]byte(result), vo)
	return vo
}
func GetOpenID(accessToken string) *models.QqOpenIdResVO {
	vo := new(models.QqOpenIdResVO)
	url := "https://graph.qq.com/oauth2.0/me?" + "access_token=" + accessToken
	logs.Info("GetOpenID 的url是", url)
	resp, err := http.Get(url)
	if err != nil {
		logs.Error("请求", url, "失败", err)
		return vo
	}
	logs.Info("qq返回的消息是", resp)
	defer resp.Body.Close()
	resultByte, err := ioutil.ReadAll(resp.Body)
	result := string(resultByte)
	logs.Info("GetOpenID 的返回体是", result)
	index1 := strings.Index(result, "(")
	index2 := strings.Index(result, ")")
	jsonStr := result[index1+1 : index2]
	json.Unmarshal([]byte(jsonStr), vo)
	return vo
}
func GetAccessToken(code, redirect_uri string) string {
	appId := beego.AppConfig.String("qqAuthAppid")
	appIdKey := beego.AppConfig.String("qqAuthAppkey")
	url := "https://graph.qq.com/oauth2.0/token?" + "grant_type=authorization_code&" +
		"client_id=" + appId + "&client_secret=" + appIdKey + "&code=" + code +
		"&redirect_uri=" + redirect_uri
	logs.Info("GetAccessToken 的url是", url)
	resp, err := http.Get(url)
	logs.Info("qq返回的消息是", resp)
	if err != nil {
		logs.Error("请求", url, "失败", err)
		return ""
	}
	defer resp.Body.Close()
	resultByte, err := ioutil.ReadAll(resp.Body)
	result := string(resultByte)
	//result为响应体
	logs.Info("GetAccessToken 的返回体是", result)
	str := strings.Split(result, "&")[0]
	return strings.Split(str, "=")[1]
}
