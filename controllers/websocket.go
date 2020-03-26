package controllers

import (
	"net/http"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket" //使用的是这个github项目的包
	"github.com/songhq211949/beego-api/models"
	my "github.com/songhq211949/beego-api/proto"
)

//WebSocketController 处理websocket连接的请求
type WebSocketController struct {
	beego.Controller
}

//Connect  Join method handles WebSocket requests for WebSocketController.加入websocket协议
func (this *WebSocketController) Connect() {
	// Upgrade from http request to WebSocket.升级为websocket协议，即进来一个就会创建一个连接
	//websocket.Upgrader{ReadBufferSize:1024, WriteBufferSize:1024}.Upgrade 报错 websocket.Upgrader{ReadBufferSize:1024, WriteBufferSize:1024}.Upgrade
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := (&upgrader).Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil)

	//连接判断处理
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		logs.Error("Cannot setup WebSocket connection:", err)
		return
	}
	//方法完成后最后执行的部分，通常是一些释放资源等操作
	defer ws.Close()
	// Message receive loop.
	for {
		//这里可以做心跳机制，客户端每10秒会发送数据
		//ReadMessage如果发生了异常就跳出for循环，随后即可关闭连接
		_, p, err := ws.ReadMessage()
		if err != nil {
			logs.Info("客户端掉线", err)
			return
		}
		//porto解码
		wsBaseReq := &my.WSBaseReqProto{}
		err = proto.Unmarshal(p, wsBaseReq)
		if err != nil {
			logs.Error("proto 解码失败")
			return
		}
		//proto解码成功
		uid := wsBaseReq.Uid //用户id
		sid := wsBaseReq.Sid
		msgType := wsBaseReq.Type

		if msgType == models.LOGIN { //登入操作
			this.UserLogin(int(uid), sid, ws)
		} else if msgType == models.PING { //心跳操作
			logs.Info("客户端心跳,用户id为：", uid)
		} else {
			logs.Info("未知类型的消息,用户id为", uid)
		}
		defer LeaveMap(int(wsBaseReq.Uid))
	}
}

//建立用户登入条件是用户真的是登入，并非是连接上就认为登入了
//因为此处的长连接是对外开发的，即连接上是不需要额外的校验，
//但该连接进入登入就需要校验
func (this *WebSocketController) UserLogin(uid int, sid string, ws *websocket.Conn) {
	//保存连接到map中
	PutMap(uid, ws)
}

//使用http的方式指定uid，向某个用户推送数据，真个仅用来测试
func (this *WebSocketController) SendMessage() {
	uid, err := this.GetInt("uid")
	if err != nil {
		this.Redirect("/", 302)
		return
	}
	msg := this.GetString("msg")
	if msg == "" {
		return
	}
	//根据uid获取连接
	ws, ok := conMap[uid]
	if !ok {
		logs.Info("用户不在线,用户id为", uid)
		return
	}
	if ws != nil {
		if ws.WriteMessage(websocket.TextMessage, []byte(msg)) != nil {
			logs.Info("用户不在线,用户id为", uid)
			LeaveMap(uid)
		}
	}
}

//封装一个发送消息的接口，这里不再是广播，而是对指定的对象发送数据
func SendMsg(receiverUid int, wsBaseReqVO models.WSBaseReqVO) {
	logs.Info("发送消息的接受者为用户id为", receiverUid, "发送的消为:", wsBaseReqVO)
	msgVo := wsBaseReqVO.Message
	userVO := wsBaseReqVO.User
	msgResVO := my.WSMessageResProto{
		ReceiveId:  uint64(msgVo.ReceiveId),
		MsgType:    int32(msgVo.MsgType),
		MsgContent: msgVo.MsgContent,
	}
	userResVo := my.WSUserResProto{
		Uid:    uint64(userVO.Uid),
		Name:   userVO.Name,
		Avatar: userVO.Avatar,
		Remark: userVO.Remark,
	}
	BaseResVo := my.WSBaseResProto{
		Type:       int32(wsBaseReqVO.Type),
		Message:    &msgResVO,
		User:       &userResVo,
		CreateTime: time.Now().String(),
	}
	data, err := proto.Marshal(&BaseResVo)
	if err != nil {
		logs.Error("protco 编码 错误", err)
	}
	//根据uid获取连接
	ws, ok := conMap[receiverUid]
	if !ok {
		logs.Info("用户不在线,用户id为", receiverUid)
		return
	}
	if ws != nil {
		//writeMessage,向该连接中写入数据
		if ws.WriteMessage(websocket.TextMessage, data) != nil {
			logs.Info("用户不在线,用户id为", receiverUid)
			LeaveMap(receiverUid)
		}
	}
}
