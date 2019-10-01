package ws

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"github.com/gorilla/websocket"
)


// 客户端读写消息
type WsMessage struct {
	MessageType int
	Data []byte
}

// 客户端连接
type WSConnection struct {
	wsSocket *websocket.Conn // 底层websocket
	inChan chan *WsMessage	// 读队列
	outChan chan *WsMessage // 写队列

	mutex sync.Mutex	// 避免重复关闭管道
	isClosed bool
	closeChan chan byte  // 关闭通知
}

// 链接
func InitConnection(wsSocket *websocket.Conn) (wsConn *WSConnection, err error) {
	wsConn = &WSConnection{
		wsSocket: wsSocket,
		inChan: make(chan *WsMessage, 1000),
		outChan: make(chan *WsMessage, 1000),
		closeChan: make(chan byte),
		isClosed: false,
	}

	// 处理器
	// go wsConn.procLoop()
	// 读协程
	go wsConn.wsReadLoop()
	// 写协程
	go wsConn.wsWriteLoop()
	return
}

func (wsConn *WSConnection)wsReadLoop() {
	for {
		// 读一个message
		msgType, data, err := wsConn.wsSocket.ReadMessage()
		if err != nil {
			goto error
		}
		req := &WsMessage{
			msgType,
			data,
		}
		// 放入请求队列
		select {
		case wsConn.inChan <- req:
		case <- wsConn.closeChan:
			goto closed
		}
	}
error:
	wsConn.WsClose()
closed:
}

func (wsConn *WSConnection)wsWriteLoop() {
	for {
		select {
		// 取一个应答
		case msg := <- wsConn.outChan:
			// 写给websocket
			if err := wsConn.wsSocket.WriteMessage(msg.MessageType, msg.Data); err != nil {
				goto error
			}
		case <- wsConn.closeChan:
			goto closed
		}
	}
error:
	wsConn.WsClose()
closed:
}

func (wsConn *WSConnection)procLoop() {
	// 启动一个gouroutine发送心跳
	go func() {
		for {
			time.Sleep(2 * time.Second)
			if err := wsConn.WsWrite(websocket.TextMessage, []byte("heartbeat")); err != nil {
				fmt.Println("heartbeat fail")
				wsConn.WsClose()
				break
			}
		}
	}()

	// // 这是一个同步处理模型（只是一个例子），如果希望并行处理可以每个请求一个gorutine，注意控制并发goroutine的数量!!!
	// for {
	// 	msg, err := wsConn.WsRead()
	// 	if err != nil {
	// 		fmt.Println("read fail")
	// 		break
	// 	}
	// 	fmt.Println(msg)
	// 	err = wsConn.WsWrite(msg.messageType, msg.data)
	// 	if err != nil {
	// 		fmt.Println("write fail")
	// 		break
	// 	}
	// }
}

// func (wsConn *WSConnection) SendText(message string) {
// 	err := socket.send(websocket.TextMessage, [] byte (message))
// 	if err != nil {
// 		logger.Error.Println("write:", err)
// 		return
// 	}
// }

// func (socket *Socket) SendBinary(data [] byte) {
// 	err := socket.send(websocket.BinaryMessage, data)
// 	if err != nil {
// 		logger.Error.Println("write:", err)
// 		return
// 	}
// }

func (wsConn *WSConnection) SendText(message string) error {
	// socket.sendMu.Lock()
	err := wsConn.wsSocket.WriteMessage(websocket.TextMessage, [] byte (message))
	// socket.sendMu.Unlock()
	return err
}

func (wsConn *WSConnection)WsWrite(messageType int, data []byte) error {
	select {
	case wsConn.outChan <- &WsMessage{messageType, data,}:
	case <- wsConn.closeChan:
		return errors.New("websocket closed")
	}
	return nil
}

func (wsConn *WSConnection)WsRead() (*WsMessage, error) {
	select {
	case msg := <- wsConn.inChan:
		return msg, nil
	case <- wsConn.closeChan:
	}
	return nil, errors.New("websocket closed")
}

func (wsConn *WSConnection)WsClose() {
	wsConn.wsSocket.Close()

	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		close(wsConn.closeChan)
	}
}