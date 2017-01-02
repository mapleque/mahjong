package mahjong

type StandardRule struct{}

func (rule *StandardRule) initPool(set *Set) {
	pool := createAllCards()
	disOrderCards(&pool)
	set.cardPool = pool
}

func (rule *StandardRule) doEventManual(set *Set) bool {
	switch set.curEvent.name {
	case EVENT_PULL4:
		getCards(4, &set.cardPool, set.curPlayer())
		return true
	case EVENT_PULL2:
		getCards(2, &set.cardPool, set.curPlayer())
		return true
	case EVENT_BUHUA:
		buhua(&set.cardPool, set.curPlayer())
		return true
	case EVENT_PULL:
		getCards(1, &set.cardPool, set.curPlayer())
		return true
	case EVENT_PUSH:
		ret := pushCard(set.curPlayer())
		if ret == nil {
			return false
		}
		set.lastEventCard = ret
		return true
	case EVENT_ZIMO:
		return zimo(set)
	case EVENT_GANG:
		return gang(set)
	case EVENT_HU:
		return hu(set)
	case EVENT_PENG:
		return peng(set)
	case EVENT_CHI:
		return chi(set)
	}
	return false
}

func (rule *StandardRule) doEventAuto(set *Set) bool {
	if set.curEvent == nil {
		set.curEvent = &Event{EVENT_PULL4}
		set.curPlayer().curEvent = set.curEvent
	}
	switch set.curEvent.name {
	case EVENT_PULL4, EVENT_PULL2, EVENT_BUHUA, EVENT_PULL:
		set.curPlayer().doEvent = set.curEvent
		return true
	case EVENT_ZIMO, EVENT_GANG, EVENT_HU, EVENT_PENG, EVENT_CHI:
		if !checkEvent(set) {
			set.curPlayer().doEvent = &Event{EVENT_PASS}
			return true
		}
	}
	return false
}

func (rule *StandardRule) nextEvent(set *Set) {
	set.lastEvent = set.curEvent
	set.lastEventPlayerIndex = set.curEventPlayerIndex
	// 首先看当前事件是什么，用户干了什么，按照用户的意愿处理掉
	// 然后更新牌局，下一步需要哪个用户干什么
	switch set.curPlayer().curEvent.name {
	case EVENT_PULL4:
		nextPlayerIndex := (set.curPlayerIndex + 1) % set.playerTotal
		set.curPlayerIndex = nextPlayerIndex
		set.curEventPlayerIndex = nextPlayerIndex
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
	case EVENT_BUHUA:
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
		nextPlayerIndex := (set.curPlayerIndex + 1) % set.playerTotal
		set.curPlayerIndex = nextPlayerIndex
		set.curEventPlayerIndex = nextPlayerIndex
		// TODO 可能牌池空了结束
		set.curEvent = &Event{EVENT_PULL}
		return
	case EVENT_PULL:
		// TODO 可能牌池空了结束
		set.curEvent = &Event{EVENT_BUHUA}
		return
	}
}

func getCards(num int, pool *[]*Card, player *Player) {
	poolCards := *pool
	poolSize := len(poolCards)
	cards := poolCards[poolSize-num : poolSize]
	*pool = poolCards[:poolSize-num]
	if num == 1 {
		player.waitCard = cards[0]
	} else {
		for _, card := range cards {
			player.handCards = append(player.handCards, card)
		}
	}
	orderCards(&player.handCards)
	if num == 2 && len(player.handCards) == 14 {
		player.waitCard = player.handCards[13]
		player.handCards = player.handCards[:13]
	}
}

func buhua(pool *[]*Card, player *Player) {
	// TODO 给当前玩家补花
}

func pushCard(player *Player) *Card {
	index := player.doIndex
	var pushCard *Card
	// 先看是不是刚抓上来的牌
	if player.waitCard.Index == (*index)[0] {
		// 如果是，直接换指针
		pushCard = player.waitCard
	} else {
		// 如果不是，在手牌里找
		for i, card := range player.handCards {
			if card.Index == (*index)[0] {
				swapCards(player.handCards[i], player.waitCard)
				orderCards(&player.handCards)
				pushCard = player.waitCard
			}
		}
	}
	if pushCard != nil {
		player.outCards = append(player.outCards, pushCard)
		player.waitCard = nil
	}
	// 如果找到了要出的牌，就返回这张牌，否则返回nil
	return pushCard

}
func zimo(set *Set) bool {
	return false
}
func hu(set *Set) bool {
	return false
}
func gang(set *Set) bool {
	return false
}
func peng(set *Set) bool {
	return false
}
func chi(set *Set) bool {
	return false
}

func checkEvent(set *Set) bool {
	// TODO 检查各种事件的成立前提条件是否满足
	switch set.curPlayer().curEvent.name {
	case EVENT_ZIMO:
		return isHu(set)
	case EVENT_HU:
		return isHu(set)
	case EVENT_CHI:
		return isChi(set.curPlayer().handCards, set.lastEventCard)
	case EVENT_PENG:
		return isPeng(set.curPlayer().handCards, set.lastEventCard)
	case EVENT_GANG:
		return isGang(set.curPlayer().handCards,
			set.curPlayer().upCardSet,
			set.lastEventCard)
		return false
	}
	return false
}

func isHu(set *Set) bool {
	return false
}
func isGang(handCards []*Card, upCardSet []*CardSet, waitCard *Card) bool {
	return false
}
func isPeng(handCards []*Card, waitCard *Card) bool {
	return false
}
func isChi(handCards []*Card, waitCard *Card) bool {
	return false
}

func calcFan(
	poolCards, handCards, flourCards []*Card,
	upCardSet, downCardSet []*CardSet,
	waitCard *Card,
	lastEvent *Event) int {
	return 1
}
