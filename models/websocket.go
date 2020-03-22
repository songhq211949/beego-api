package models

import "github.com/gorilla/websocket"

type Client struct {
	conn *websocket.Conn    // 用户websocket连接
	name string             // 用户名称
}

// 1.设置为公开属性(即首字母大写)，是因为属性值私有时，外包的函数无法使用或访问该属性值(如：json.Marshal())
// 2.`json:"name"` 是为了在对该结构类型进行json编码时，自定义该属性的名称
type Message struct {
	EventType byte  `json:"type"`       // 0表示用户发布消息；1表示用户进入；2表示用户退出
	Name string     `json:"name"`       // 用户名称
	Message string  `json:"message"`    // 消息内容
}
