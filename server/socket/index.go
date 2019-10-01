package socket

import (
	"net/http"
	"fmt"
	"strconv"
	"github.com/gorilla/websocket"
	ws "go-socket-demo/modules"
)

// http升级websocket协议的配置
var wsUpgrader = websocket.Upgrader{
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 保存连接的客户端
var Clients = make(map[string] *ws.WSConnection) 

// 链接
func InitConnection(token string, resp http.ResponseWriter, req *http.Request) (wsConn *ws.WSConnection)  {
	var(
		wsSocket *websocket.Conn
		err error
	)

	// 应答客户端告知升级连接为websocket
	wsSocket, err = wsUpgrader.Upgrade(resp, req, nil)
	if err != nil {
		return
	}

	defer wsSocket.Close()

	wsConn, err = ws.InitConnection(wsSocket)
	if err != nil {
		goto ERR
	}

	Clients[token] = wsConn

	// 接收客服端的消息
	for {
		msg, err := wsConn.WsRead()
		if err != nil {
			fmt.Println("read fail")
			break
		}
		// 向客户端发送消息
		fmt.Println(strconv.Itoa(msg.MessageType))
		fmt.Println(string(msg.Data))
		err = wsConn.WsWrite(msg.MessageType, msg.Data)
		if err != nil {
			fmt.Println("write fail")
			break
		}
	}

	ERR:
	  wsConn.WsClose()
		
  return wsConn
}

// 获取服务
func GetServer (token string) *ws.WSConnection {
	return Clients[token]
}