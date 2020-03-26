package models

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
	GroupId int `orm:"auto"`
	/**
	 * 创建者用户ID
	 */
	Uid int
	/**
	 * 群昵称
	 */
	Name string
	/**
	 * 群头像
	 */
	Avatar string
	/**
	 * 成员数量
	 */
	MemberNum int
	/**
	 * 描述
	 */
	Remark string
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
	GroupId int
	/**
	 * 发送消息的用户ID
	 */
	SenderUid int
	/**
	 * 消息类型（0：普通文字消息，1：图片消息，2：文件消息，3：语音消息，4：视频消息）
	 */
	MsgType int
	/**
	 * 消息内容
	 */
	MsgContent string

	BaseTime
}

//消息列表
type GroupMsgListResVO struct {
	GroupMsg
	UserInfoListResVO
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
