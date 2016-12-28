package mahjong

type StandardRule struct{}

/*
func (rule *StandardRule) eventChain() map[int]*Event {
	chain := make(map[int]*Event)
	&Event{name: EVENT_BUHUA, rule: EVENT_RULE_SELF}
	&Event{name: EVENT_ZIMO, rule: EVENT_RULE_SELF}
	&Event{name: EVENT_HU, rule: EVENT_RULE_OTHER}
	&Event{name: EVENT_GANG, rule: EVENT_RULE_ALL}
	&Event{name: EVENT_PENG, rule: EVENT_RULE_OTHER}
	&Event{name: EVENT_CHI, rule: EVENT_RULE_NEXT}
	&Event{name: EVENT_PUSH, rule: EVENT_RULE_SELF}
	&Event{name: EVENT_PULL, rule: EVENT_RULE_SELF}
	return chain
}
*/

func (rule *StandardRule) initPool(set *Set) {
	pool := createAllCards()
	disOrderCards(&pool)
	set.cardPool = &pool
}

func (rule *StandardRule) doEvent(set *Set) {
	set.curPlayer().curEvent = set.curEvent
	switch set.curEvent.name {
	// 系统自动代理操作
	case EVENT_PULL4:
		getCards(4, &set.cardPool, set.curPlayer())
		set.curPlayer().doEvent = set.curEvent
		return
	case EVENT_PULL2:
		getCards(2, &set.cardPool, set.curPlayer())
		set.curPlayer().doEvent = set.curEvent
		return
	case EVENT_BUHUA:
		buhua(&set.cardPool, set.curPlayer())
		set.curPlayer().doEvent = set.curEvent
		return
	case EVENT_ZIMO, EVENT_GANG:
		if !checkEvent(set) {
			set.curPlayer().doEvent = &Event{EVENT_PASS}
		}
		return
	case EVENT_PUSH:
		return
	case EVENT_PULL:
		getCards(1, &set.cardPool, set.curPlayer)
		set.curPlayer().doEvent = set.curEvent
		return
	default:
		// wait for user op
	}
}

func (rule *StandarRule) nextEvent(set *Set) {
	// 首先看当前事件是什么，用户干了什么，按照用户的意愿处理掉
	// 然后更新牌局，下一步需要哪个用户干什么
	switch set.curPlayer().curEvent.name {
	case EVENT_PULL4:
		// 换下家
		set.curPlayerIndex = (set.curPlayerIndex + 1) % set.playerTotal
		// 如果下家是庄，并且牌已经到了12张，就该pull2了，否则依然pull4
		if len(set.curPlayer().handCards) == 12 {
			set.curEvent = &Event{EVENT_PULL2}
		} else {
			set.curEvent = &Event{EVENT_PULL4}
		}
		return
	case EVENT_PULL2:
		set.curEvent = &Event{EVENT_BUHUA}
		return
	case BUHUA:
		orderCards(&set.curPlayer().handCards)
		set.curEvent = &Event{EVENT_ZIMO}
		return
	case EVENT_ZIMO:
		set.curEvent = &Event{EVENT_GANG}
		return
	case EVENT_GANG:
		set.curEvent = &Event{EVENT_PUSH}
		return
	case EVENT_PUSH:
		// TODO 玩家出牌之后的事件整理
		// 换下家
		set.curPlayerIndex = (set.curPlayerIndex + 1) % set.playerTotal
		set.curEvent = &Event{EVENT_PULL}
		return
	case EVENT_PULL:
		set.curEvent = &Event{EVENT_BUHUA}
		return
	}
}

func getCards(num int, pool *[]*Card, player *Player) {
	poolCards := *pool
	poolSize := len(poolCards)
	cards := poolCards[poolSize-1 : poolSize]
	poolCards = poolCards[:poolSize-1]
	if num == 1 {
		player.waitCard = cards[0]
	} else {
		for _, card := range cards {
			handCards := append(player.handCards, card)
			player.handCards = handCards
		}
	}
}

func buhua(pool *[]*Card, player *Player) {
	// TODO 给当前玩家补花
}

func checkEvent(set *Set) bool {
	// TODO 检查各种事件的成立前提条件是否满足
	switch set.curPlayer().curEvent.name {
	case EVENT_ZIMO:
		return false
	}
	return false
}
