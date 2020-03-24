package models

//websocket 通信类
type WSBaseReqVO struct {
	/**
	 * 类型
	 */
	Type int;

	/**
	 * 消息实体
	 */
	Message WSMessageReqVO;

	/**
	 * 发送者用户信息
	 */
	User WSUserReqVO;

}
//消息体
type WSMessageReqVO struct {
	/**
	 * 接收者ID
	 */
	ReceiveId int;

	/**
	 * 消息类型
	 */
	MsgType int;

	/**
	 * 消息内容
	 */
	MsgContent string;

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
type GroupSaveReqVO struct{
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