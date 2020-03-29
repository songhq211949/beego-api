package models

import "time"

const LOGIN = 1 //登入
//const  PING  = 0 //心跳

//消息的大类型
const LOGIN_OUT = -2
const WS_OUT = -1
const PING = 0
const FRIEND = 1
const GROUP = 2
const FRIEND_ASK = 3
const FRIEND_ACK = 4
const JOIN_GROUP = 5

//UserLoginDTO  用户登入
type UserLoginDTO struct {
	Uid int `json:"uid"`
}

//userLoginResVO 用户登入vo
type UserLoginResVO struct {
	Uid int    `json:"uid"`
	Sid string `json:"sid"`
}
type UserFriendListInfoResVO struct {
	UserFriend
	User UserInfoListResVO `json:"user"`
}
type GroupUserListResVO struct {
	GroupUser
	Group Group `json:"group"`
}

type UserFriendAskListResVO struct {
	UserFriendAsk
	User UserInfoListResVO `json:"user"`
}
type QqOpenIdResVO struct {
	Client_id string `json:"client_id"`
	Openid    string `json:"openid"`
}
type QqUserInfoResVO struct {
	Ret                int    `json:"ret"`
	Msg                string `json:"msg"`
	Nickname           string `json:"nickname"`
	Figureurl          string `json:"figureurl"`
	Figureurl_1        string `json:"figureurl_1"`
	Figureurl_2        string `json:"figureurl_2"`
	Figureurl_qq_1     string `json:"figureurl_qq_1"`
	Figureurl_qq_2     string `json:"figureurl_qq_2"`
	Gender             string `json:"gender"`
	Is_yellow_vip      string `json:"is_yellow_vip"`
	Vip                string `json:"vip"`
	Yellow_vip_level   string `json:"yellow_vip_level"`
	Level              string `json:"level"`
	Is_yellow_year_vip string `json:"is_yellow_year_vip"`
}
type UserQq struct {
	Id int `orm:"auto"`
	/**
	* 用户uid（关联 user 表）
	 */
	Uid int
	/**
	* 用户openID
	 */
	Openid string
	/**
	* 创建时间
	 */
	CreateTime time.Time
}
