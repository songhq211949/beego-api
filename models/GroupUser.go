package models
import "time"

type BaseTime struct{
	CreateTime time.Time `json:"createTime"`
	ModifiedTime time.Time `json:"modifiedTime"`
}

type GroupUser struct{
	Id int `json:"id"`
	GroupId int `json:"groupId"`
	Uid int `json:"uid"`
	Remark string `json:"remark"`
	LastAckMsgId int `json:"lastACkMsgId"`// 最后一次确认的消息ID
	LastMsgContent string `json:"lastMsgContent"`//最后一次的消息内容
	LastMsgTime  time.Time `json:"lastMsgTime"`//最后一次的消息时间
	UnMsgCount int `json:"unMsgCount"`
	Rank int  `json:"rank"`//等级（0：普通成员，1：管理员，2：群主）
	BaseTime
}
type UserInfoListResVO struct{
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
type User struct{
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

	//密码
	Pwd string `json:"pwd"`

	BaseTime
}
type  GroupIndexListResVO struct{
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
