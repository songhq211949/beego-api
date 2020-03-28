package models

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
