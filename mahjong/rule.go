package mahjong

type Rule int

const (
	RULE_STANDARD Rule = iota
)

type RuleImpl interface {
	allCards() *[]*Card
	getCards(int, *[]*Card, *Player)
	getNextEvent(*Player) *Event
}

func getRule(rule Rule) RuleImpl {
	switch rule {
	case RULE_STANDARD:
		return &StandardRule{}
	default:
		return nil
	}
}

type StandardRule struct{}

/*
func (rule *StandardRule) eventChain() map[int]*Event {
	chain := make(map[int]*Event)
	chain[0] = &Event{name: EVENT_BUHUA, rule: EVENT_RULE_SELF}
	chain[1] = &Event{name: EVENT_ZIMO, rule: EVENT_RULE_SELF}
	chain[2] = &Event{name: EVENT_HU, rule: EVENT_RULE_OTHER}
	chain[3] = &Event{name: EVENT_GANG, rule: EVENT_RULE_ALL}
	chain[4] = &Event{name: EVENT_PENG, rule: EVENT_RULE_OTHER}
	chain[5] = &Event{name: EVENT_CHI, rule: EVENT_RULE_NEXT}
	chain[6] = &Event{name: EVENT_PUSH, rule: EVENT_RULE_SELF}
	chain[7] = &Event{name: EVENT_PULL, rule: EVENT_RULE_SELF}
	return chain
}
*/

func (rule *StandardRule) allCards() *[]*Card {
	var pool []*Card
	for i := 0; i < 160; i++ {
		card := createCard(i)
		if card != nil {
			pool = append(pool, card)
		}
	}
	return &pool
}

func (rule *StandardRule) getCards(num int, pool *[]*Card, player *Player) {
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
		// TODO 排序
	}
}
func (rule *StandardRule) getNextEvent(player *Player) *Event {
	if player.curEvent == nil {
		event := &Event{name: EVENT_PULL4, rule: EVENT_RULE_SELF}
		return event
	}
	switch player.curEvent.name {
	case EVENT_PULL4:
		if len(player.handCards) < 12 {
			event := &Event{name: EVENT_PULL4, rule: EVENT_RULE_SELF}
			player.curEvent = event
			return event
		} else if player.zhuang {
			event := &Event{name: EVENT_PULL2, rule: EVENT_RULE_SELF}
			player.curEvent = event
			return event
		} else {
			event := &Event{name: EVENT_PULL, rule: EVENT_RULE_SELF}
			player.curEvent = event
			return event
		}
	case EVENT_PULL2:
		event := &Event{name: EVENT_BUHUA, rule: EVENT_RULE_SELF}
		player.curEvent = event
		return event
	default:
		return nil
	}
}
