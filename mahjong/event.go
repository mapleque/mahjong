package mahjong

// 事件定义
type EventName int

const (
	_                EventName = iota
	EVENT_BUHUA                // 1补花
	EVENT_HU                   // 2胡
	EVENT_GANG                 // 3杠
	EVENT_CI                   // 4杠次
	EVENT_PENG                 // 5碰
	EVENT_CHI                  // 6吃
	EVENT_PUSH                 // 7出牌
	EVENT_PULL                 // 8抓牌
	EVENT_PULL1                // 9抓牌
	EVENT_PULL2                // 10抓牌
	EVENT_PULL4                // 11抓牌
	EVENT_PASS                 // 12过
	EVENT_INIT_BUHUA           // 13补花

	EVENT_EXPAND_AFTER_PULL  // 14
	EVENT_EXPAND_AFTER_PUSH  // 15
	EVENT_EXPAND_PUSH_FINISH // 16
)

type Event struct {
	PlayerId string
	Event    EventName
}

type EventQueue struct {
	Stack []*Event
}

func (es *EventQueue) put(event *Event) {
	// put on tail
	es.Stack = append(es.Stack, event)
}

// normally return (*Event, false)
// if empty return (nil, true)
func (es *EventQueue) pop() (*Event, bool) {
	// pop head
	if es.Stack != nil && len(es.Stack) > 0 {
		ret := es.Stack[0]
		es.Stack = es.Stack[1:]
		return ret, false
	}
	return nil, true
}

func (es *EventQueue) clear() {
	// TODO clear to empty
	es.Stack = []*Event{}
}
