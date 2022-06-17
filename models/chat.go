package models

import "container/list"

type EventType int

// 定义事件状态
const (
	EVENT_JOIN = iota
	EVENT_LEAVE
	EVENT_MESSAGE
)

// 定义消息结构
type Event struct {
	Type       EventType // JOIN, LEAVE, MESSAGE
	ID         int
	Username   string
	Content    string
	Headimgurl string
	Addtime    int // Unix timestamp (secs)
}

// 消息最大存放数
const archiveSize = 20

// Event archives
// 消息链表.
var archive = list.New()

// NewArchive saves new event to archive list.
// 放入消息链表，如果链表满了就移除最前元素
func NewArchive(event Event) {
	if archive.Len() >= archiveSize {
		archive.Remove(archive.Front())
	}
	archive.PushBack(event)
}

// GetEvents returns all events after lastReceived.
// 从消息链表获取事件数组
func GetEvents(lastReceived int) []Event {
	// make 数组管道时 size 为 0 是因为可以扩容到 消息链表的长度（cap）
	events := make([]Event, 0, archive.Len())
	for event := archive.Front(); event != nil; event = event.Next() {
		e := event.Value.(Event)
		if e.Addtime > int(lastReceived) {
			events = append(events, e)
		}
	}
	return events
}
