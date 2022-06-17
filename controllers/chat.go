// 开启主要协程

package controllers

import (
	"chat/models"
	"container/list"
	"time"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

// 暂时没用到
type Subscription struct {
	Archive []models.Event      // All the events from the archive.
	New     <-chan models.Event // New events coming in.
}

// 新建模型的消息结构事件并返回
func newEvent(ep models.EventType, id int, username, content, headimgurl string) models.Event {
	return models.Event{ep, id, username, content, headimgurl, int(time.Now().Unix())}
}

// 用户加入时把数据填入subscribe 管道
func Join(user string, ws *websocket.Conn) {
	subscribe <- Subscriber{Name: user, Conn: ws}
}

// 用户离开时把数据填入unsubscribe 管道
func Leave(user string) {
	unsubscribe <- user
}

// 用户自己的结构
type Subscriber struct {
	Name string
	Conn *websocket.Conn // Only for WebSocket users; otherwise nil.
}

var (
	// Channel for new join users.在线用户管道
	subscribe = make(chan Subscriber, 10)
	// Channel for exit users.离线用户管道
	unsubscribe = make(chan string, 10)
	// Send events here to publish them.消息发布管道
	publish = make(chan models.Event, 10)
	// Long polling waiting list.
    // 等待消息链表
	waitingList = list.New()
    // 用户在线链表
	subscribers = list.New()
)

// This function handles all incoming chan messages.聊天协程
func chatroom() {
	for {
        // 使用select防止管道未关闭的错误
		// 	获取事件
		select {
		// 用户加入聊天事件
		case sub := <-subscribe:
			// 用户是否离线
			if !isUserExist(subscribers, sub.Name) {
				// 加入在线用户管道
				subscribers.PushBack(sub) // Add user to the end of list.
				// Publish a JOIN event.
				// 加入发布消息管道
				publish <- newEvent(models.EVENT_JOIN, 0, sub.Name, "", "")
				beego.Info("New user:", sub.Name, ";WebSocket:", sub.Conn != nil)
			} else {
				beego.Info("Old user:", sub.Name, ";WebSocket:", sub.Conn != nil)
			}
			// 发布消息事件
		case event := <-publish:
			// Notify waiting list.
			// for ch := waitingList.Back(); ch != nil; ch = ch.Prev() {
			// 	ch.Value.(chan bool) <- true
			// 	waitingList.Remove(ch)
			// }
           	// 遍历等待管道，造成ajax阻塞关键点后的处理
			for ch := waitingList.Front(); ch != nil; ch = waitingList.Front() {
                // 设置等待管道的元素为真，开始发布消息
				ch.Value.(chan bool) <- true
				// 移除元素
				waitingList.Remove(ch)
			}
			// 传入到模型的事件消息
			models.NewArchive(event)

			if event.Type == models.EVENT_MESSAGE {
				beego.Info("Message from", event.Username, ";Content:", event.Content)
			}
			// 用户退出聊天事件
		case unsub := <-unsubscribe:
			for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(Subscriber).Name == unsub {
					subscribers.Remove(sub)
					// Clone connection.
					ws := sub.Value.(Subscriber).Conn
					if ws != nil {
						ws.Close()
						beego.Error("WebSocket closed:", unsub)
					}
					publish <- newEvent(models.EVENT_LEAVE, 0, unsub, "", "") // Publish a LEAVE event.
					break
				}
			}
		}
	}
}

func init() {
	go chatroom()
}

func isUserExist(subscribers *list.List, user string) bool {
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(Subscriber).Name == user {
			return true
		}
	}
	return false
}
