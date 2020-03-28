package models

import "time"

//websocket 通信类
type WSBaseReqVO struct {
	/**
	 * 类型
	 */
	Type int

	/**
	 * 消息实体
	 */
	Message WSMessageReqVO

	/**
	 * 发送者用户信息
	 */
	User WSUserReqVO
}

//消息体
type WSMessageReqVO struct {
	/**
	 * 接收者ID
	 */
	ReceiveId int

	/**
	 * 消息类型
	 */
	MsgType int

	/**
	 * 消息内容
	 */
	MsgContent string
}

//通信的用户信息
type WSUserReqVO struct {
	/**
	 * 用户id
	 */
	Uid int
	/**
	 * 用户昵称
	 */
	Name string
	/**
	 * 用户头像
	 */
	Avatar string
	/**
	 * 个性签名
	 */
	Remark string
}
type GroupSaveReqVO struct {
	/**
	 * 群ID
	 */
	GroupId int

	/**
	 * 群昵称
	 */
	Name string

	/**.
	 * 群头像
	 */
	Avatar string

	/**.
	 * 说明
	 */
	Remark string
}

//Group 组
type Group struct {
	/**
	* 群ID
	 */
	GroupId int `orm:"auto" json:"groupId"`
	/**
	 * 创建者用户ID
	 */
	Uid int `json:"uid"`
	/**
	 * 群昵称
	 */
	Name string `json:"name"`
	/**
	 * 群头像
	 */
	Avatar string `json:"avatar"`
	/**
	 * 成员数量
	 */
	MemberNum int `json:"memberNum"`
	/**
	 * 描述
	 */
	Remark string `json:"remark"`
	BaseTime
}

type GroupMsg struct {
	/**
	* 消息ID
	 */
	MsgId int `orm:"auto"`
	/**
	 * 群ID
	 */
	GroupId int `json:groupId`
	/**
	 * 发送消息的用户ID
	 */
	SenderUid int `json:"senderUid"`
	/**
	 * 消息类型（0：普通文字消息，1：图片消息，2：文件消息，3：语音消息，4：视频消息）
	 */
	MsgType int `json:"msgType"`
	/**
	 * 消息内容
	 */
	MsgContent string `json:"msgContent"`

	BaseTime
}

//消息列表
type GroupMsgListResVO struct {
	GroupMsg
	User UserInfoListResVO `json:"user"`
}

//GroupMsgCreateReqVO 群消息
type GroupMsgCreateReqVO struct {
	GroupId    int
	MsgType    int
	MsgContent string
}

//UserLoginPwdReqVO 账号密码
type UserLoginPwdReqVO struct {
	UserName string `form:"userName"`
	Password string `form:"password"`
}
type UserProfile struct {
	Id             int
	Uid            int
	FriendAskCount int `json:"friendAskCount"`
	FriendCount    int `json:"friendCount"`
	BaseTime
}
type UserInfoResVO struct {
	User
	Profile UserProfile `json:"profile"`
}
type UserFriend struct {
	/**
	* 自增id
	 */
	Id int `orm:"auto" json:"id"`
	/**
	 * 用户id
	 */
	Uid int `json:"uid"`
	/**
	 * 朋友的用户id
	 */
	FriendUid int `json:"friendUid"`
	/**
	 * 备注
	 */
	Remark string `json:"remark"`
	/**
	 * 未读消息数量
	 */
	UnMsgCount int `json:"unMsgCount"`
	/**
	 * 最后一次接收的消息内容
	 */
	LastMsgContent string `json:"lastMsgContent"`
	BaseTime
}
type BaseTime struct {
	CreateTime   time.Time `json:"createTime"`
	ModifiedTime time.Time `json:"modifiedTime"`
}

type GroupUser struct {
	Id             int       `json:"id" orm:"auto"`
	GroupId        int       `json:"groupId"`
	Uid            int       `json:"uid"`
	Remark         string    `json:"remark"`
	LastAckMsgId   int       `json:"lastACkMsgId"`   // 最后一次确认的消息ID
	LastMsgContent string    `json:"lastMsgContent"` //最后一次的消息内容
	LastMsgTime    time.Time `json:"lastMsgTime"`    //最后一次的消息时间
	UnMsgCount     int       `json:"unMsgCount"`
	Rank           int       `json:"rank"` //等级（0：普通成员，1：管理员，2：群主）
	BaseTime
}
type UserInfoListResVO struct {
	/**
	 * 用户id
	 */
	Uid int `json:"uid"`
	/**
	 * 用户昵称
	 */
	Name string `json:"name"`
	/**
	 * 用户头像
	 */
	Avatar string `json:"avatar"`
	/**
	 * 个性签名
	 */
	Remark string `json:"remark"`
}
type User struct {
	/**
	s	 * 用户id
	*/
	Uid int `json:"uid" orm:"pk"`
	/**
	 * 用户昵称
	 */
	Name string `json:"name"`
	/**
	 * 用户头像
	 */
	Avatar string `json:"avatar"`
	/**
	 * 个性签名
	 */
	Remark string `json:"remark"`

	//密码
	Pwd string `json:"pwd"`

	BaseTime
}
type GroupIndexListResVO struct {
	/**
	 * 群ID
	 */
	GroupId int `json:"groupId"`
	/**
	 * 说明
	 */
	Remark string `json:"remark"`
	/**
	 * 等级（0：普通成员，1：管理员，2：群主）
	 */
	Rank int `json:"rank"`

	/**
	 * 用户信息
	 */
	User UserInfoListResVO `json:"user"`
}
