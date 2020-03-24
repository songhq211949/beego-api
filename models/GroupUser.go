package models
import "time"
type GroupUser struct{
	Id int
	GroupId int
	Uid int
	Remark string
	LastAckMsgId int // 最后一次确认的消息ID
	LastMsgContent string //最后一次的消息内容
	lastMsgTime  time.Time //最后一次的消息时间
	UnMsgCount int
	Rank int  //等级（0：普通成员，1：管理员，2：群主）
	CreateTime time.Time
	ModifiedTime time.Time
}
