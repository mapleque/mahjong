package mahjong

type Rule int

const (
	RULE_STANDARD Rule = iota
)

type Event struct {
	name EventName // 事件名称
}

// 事件定义
type EventName int

const (
	_           EventName = iota
	EVENT_BUHUA           // 补花
	EVENT_ZIMO            // 自摸
	EVENT_HU              // 胡
	EVENT_GANG            // 杠
	EVENT_PENG            // 碰
	EVENT_CHI             // 吃
	EVENT_PUSH            // 出牌
	EVENT_PULL            // 抓牌
	EVENT_PULL2           // 抓牌
	EVENT_PULL4           // 抓牌
	EVENT_PASS            // 过
)

type CardSet *[]*Card

type RuleImpl interface {
	initPool(*Set)     // 获取初始牌，乱序的
	doEvent(*Set) bool // 如果是系统自动处理的事件，这里会把它做掉，并返回true
	nextEvent(*Set)    // 更新下一次操作
}

func getRule(rule Rule) RuleImpl {
	switch rule {
	case RULE_STANDARD:
		return &StandardRule{}
	default:
		return nil
	}
}
