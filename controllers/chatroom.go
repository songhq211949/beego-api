package controllers

import (
	"github.com/gorilla/websocket"
)

func PutMap(uid int, ws *websocket.Conn) {
	//subscribe <- Subscriber{Uid: uid, Conn: ws}
	conMap[uid] = ws
}

func LeaveMap(uid int) {
	//从map中删除，表示下线
	delete(conMap,uid)
}
var (
	conMap = make(map[int]*websocket.Conn)

)
