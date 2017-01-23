package mahjong

func (rule *StandardRule) initEventQueue(fsm *FSM) {
	// push init event in event stack
	for i := 0; i < 3; i++ {
		for _, playerId := range fsm.PlayerChain {
			fsm.EventQueue.put(&Event{playerId, EVENT_PULL4})
		}
	}
	for seq, playerId := range fsm.PlayerChain {
		if seq == 0 {
			fsm.EventQueue.put(&Event{playerId, EVENT_PULL2})
		} else {
			fsm.EventQueue.put(&Event{playerId, EVENT_PULL1})
		}
	}
	for _, playerId := range fsm.PlayerChain {
		fsm.EventQueue.put(&Event{playerId, EVENT_INIT_BUHUA})
	}
	fsm.EventQueue.put(&Event{fsm.PlayerChain[0], EVENT_EXPAND_AFTER_PULL})
}

func (rule *StandardRule) dealEvent(fsm *FSM, event *Event) {
	switch event.Event {
	case EVENT_EXPAND_AFTER_PULL:
		rule.afterPull(fsm, event.PlayerId)
	case EVENT_EXPAND_AFTER_PUSH:
		rule.afterPush(fsm, event.PlayerId)
	case EVENT_EXPAND_PUSH_FINISH:
		rule.pushFinish(fsm, event.PlayerId)
	case EVENT_HU, EVENT_GANG, EVENT_PENG, EVENT_CHI, EVENT_PUSH:
		fsm.WaitOp = true
	case EVENT_CI:
		rule.doCi(fsm, event.PlayerId)
	case EVENT_PULL:
		rule.doPull(fsm, event.PlayerId)
	case EVENT_PULL4:
		rule.doInitPull(fsm, event.PlayerId, 4)
	case EVENT_PULL2:
		rule.doInitPull(fsm, event.PlayerId, 2)
	case EVENT_PULL1:
		rule.doInitPull(fsm, event.PlayerId, 1)
	case EVENT_BUHUA:
		rule.doBuhua(fsm, event.PlayerId)
	case EVENT_INIT_BUHUA:
		rule.doInitBuhua(fsm, event.PlayerId)
	}
}

func (rule *StandardRule) afterPull(fsm *FSM, playerId string) {
	if rule.checkBuhua(fsm, playerId) {
		fsm.EventQueue.put(&Event{playerId, EVENT_BUHUA})
	} else {
		if rule.checkHu(fsm, playerId) {
			fsm.EventQueue.put(&Event{playerId, EVENT_HU})
		}
		if rule.checkGang(fsm, playerId) {
			fsm.EventQueue.put(&Event{playerId, EVENT_GANG})
		}
		fsm.EventQueue.put(&Event{playerId, EVENT_PUSH})
	}
}

func (rule *StandardRule) afterPush(fsm *FSM, playerId string) {
	fsm.DependPlayer = playerId
	cur := 0
	for seq, otherPlayerId := range fsm.PlayerChain {
		if playerId == otherPlayerId {
			cur = seq
		}
	}
	totalPlayer := len(fsm.PlayerChain)
	for i := 1; i < totalPlayer; i++ {
		nextPlayerId := fsm.PlayerChain[(cur+i)%totalPlayer]
		if rule.checkHu(fsm, nextPlayerId) {
			fsm.EventQueue.put(&Event{nextPlayerId, EVENT_HU})
		}
	}
	for i := 1; i < totalPlayer; i++ {
		nextPlayerId := fsm.PlayerChain[(cur+i)%totalPlayer]
		if rule.checkGang(fsm, nextPlayerId) {
			fsm.EventQueue.put(&Event{nextPlayerId, EVENT_GANG})
		}
	}
	for i := 1; i < totalPlayer; i++ {
		nextPlayerId := fsm.PlayerChain[(cur+i)%totalPlayer]
		if rule.checkPeng(fsm, nextPlayerId) {
			fsm.EventQueue.put(&Event{nextPlayerId, EVENT_PENG})
		}
	}
	nextPlayerId := fsm.PlayerChain[(cur+1)%totalPlayer]
	if rule.checkChi(fsm, nextPlayerId) {
		fsm.EventQueue.put(&Event{nextPlayerId, EVENT_CHI})
	}
	fsm.EventQueue.put(&Event{playerId, EVENT_EXPAND_PUSH_FINISH})
	fsm.EventQueue.put(&Event{nextPlayerId, EVENT_PULL})
}

func (rule *StandardRule) pushFinish(fsm *FSM, playerId string) {
	fsm.DependPlayer = ""
	fsm.WaitCard = nil
}

func (rule *StandardRule) doCi(fsm *FSM, playerId string) {
	fsm.DependPlayer = ""
	poolSize := len(fsm.CardPool)
	if poolSize == 0 {
		fsm.EventQueue.clear()
		return
	}
	cards := fsm.CardPool[0:1]
	fsm.CardPool = fsm.CardPool[1:]
	fsm.PlayerList[playerId].WaitCard = cards[0]
	fsm.EventQueue.put(&Event{playerId, EVENT_EXPAND_AFTER_PULL})
}

func (rule *StandardRule) doPull(fsm *FSM, playerId string) {
	poolSize := len(fsm.CardPool)
	if poolSize == 0 {
		fsm.EventQueue.clear()
		return
	}
	cards := fsm.CardPool[poolSize-1 : poolSize]
	fsm.CardPool = fsm.CardPool[:poolSize-1]
	fsm.PlayerList[playerId].WaitCard = cards[0]
	fsm.EventQueue.put(&Event{playerId, EVENT_EXPAND_AFTER_PULL})
}

func (rule *StandardRule) doBuhua(fsm *FSM, playerId string) {
	poolSize := len(fsm.CardPool)
	if poolSize == 0 {
		fsm.EventQueue.clear()
		return
	}
	player := fsm.PlayerList[playerId]
	cards := fsm.CardPool[poolSize-1 : poolSize]
	fsm.CardPool = fsm.CardPool[:poolSize-1]
	player.FlourCards = append(player.FlourCards, cards[0])
	swapCards(player.WaitCard, cards[0])
	fsm.EventQueue.put(&Event{playerId, EVENT_EXPAND_AFTER_PULL})
}

func (rule *StandardRule) doPush(fsm *FSM, playerId string) bool {
	player := fsm.PlayerList[playerId]
	index := player.DoIndex
	var pushCard *Card
	// 先看是不是刚抓上来的牌
	if player.WaitCard != nil && player.WaitCard.Index == (*index)[0] {
		// 如果是，直接换指针
		pushCard = player.WaitCard
	} else {
		// 如果不是，在手牌里找
		for i, card := range player.HandCards {
			if card.Index == (*index)[0] {
				if player.WaitCard != nil {
					swapCards(player.HandCards[i], player.WaitCard)
					pushCard = player.WaitCard
				} else {
					pushCard = card
					removeCards(&player.HandCards, []*Card{card})
				}
				orderCards(&player.HandCards)
				break
			}
		}
	}
	if pushCard == nil {
		return false
	}
	player.WaitCard = nil
	fsm.WaitCard = pushCard
	player.OutCards = append(player.OutCards, fsm.WaitCard)
	fsm.EventQueue.put(&Event{playerId, EVENT_EXPAND_AFTER_PUSH})
	return true
}
func (rule *StandardRule) doHu(fsm *FSM, playerId string) bool {
	// calculate fan
	fans := rule.processFanList(fsm, playerId)
	fanCount := rule.sumFanCount(fans)
	// build win info
	losePlayerIds := []string{}
	if fsm.DependPlayer != "" {
		losePlayerIds = append(losePlayerIds, fsm.DependPlayer)
	} else {
		for _, otherPlayerId := range fsm.PlayerChain {
			if playerId != otherPlayerId {
				losePlayerIds = append(losePlayerIds, otherPlayerId)
			}
		}
	}
	fsm.WinInfo = &WinInfo{
		Fans:          fans,
		Count:         fanCount,
		WinPlayerId:   playerId,
		LosePlayerIds: losePlayerIds}
	fsm.EventQueue.clear()
	return true
}

func (rule *StandardRule) doGang(fsm *FSM, playerId string) bool {
	player := fsm.PlayerList[playerId]
	cur := 0
	for seq, otherPlayerId := range fsm.PlayerChain {
		if playerId == otherPlayerId {
			cur = seq
		}
	}
	if fsm.DependPlayer != "" {
		// gang from other card
		cardSet := &CardSet{Cards: []*Card{fsm.WaitCard}}
		toRemove := []*Card{}
		for _, card := range player.HandCards {
			if rule.isDui(card, fsm.WaitCard) {
				cardSet.Cards = append(cardSet.Cards, card)
				toRemove = append(toRemove, card)
			}
			if len(cardSet.Cards) == 4 {
				break
			}
		}
		if len(cardSet.Cards) != 4 {
			return false
		}
		player.UpCardSet = append(player.UpCardSet, cardSet)
		removeCards(&player.HandCards, toRemove)
		removeCards(
			&fsm.PlayerList[fsm.DependPlayer].OutCards,
			[]*Card{fsm.WaitCard})

		fsm.DependPlayer = playerId
		fsm.EventQueue.clear()
		totalPlayer := len(fsm.PlayerChain)
		for i := 1; i < totalPlayer; i++ {
			nextPlayerId := fsm.PlayerChain[(cur+i)%totalPlayer]
			if rule.checkHu(fsm, nextPlayerId) {
				fsm.EventQueue.put(&Event{nextPlayerId, EVENT_HU})
			}
		}
		fsm.EventQueue.put(&Event{playerId, EVENT_CI})
		return true
	}
	// gang from self
	cardSet := &CardSet{Cards: []*Card{player.WaitCard}}
	toRemove := []*Card{}
	for _, card := range player.HandCards {
		if rule.isDui(card, player.WaitCard) {
			cardSet.Cards = append(cardSet.Cards, card)
			toRemove = append(toRemove, card)
		}
		if len(cardSet.Cards) == 4 {
			break
		}
	}
	if len(cardSet.Cards) == 4 {
		// down
		player.DownCardSet = append(player.DownCardSet, cardSet)
		removeCards(&player.HandCards, toRemove)
		player.WaitCard = nil
		fsm.EventQueue.clear()
		fsm.EventQueue.put(&Event{playerId, EVENT_CI})
		return true
	}
	// up
	for _, holdCardSet := range player.UpCardSet {
		if rule.isGang(
			holdCardSet.Cards[0],
			holdCardSet.Cards[1],
			holdCardSet.Cards[2],
			player.WaitCard) {

			holdCardSet.Cards = append(holdCardSet.Cards, player.WaitCard)
			removeCards(&player.HandCards, toRemove)
			fsm.EventQueue.clear()
			fsm.DependPlayer = playerId
			fsm.WaitCard = player.WaitCard
			player.WaitCard = nil
			totalPlayer := len(fsm.PlayerChain)
			for i := 1; i < totalPlayer; i++ {
				nextPlayerId := fsm.PlayerChain[(cur+i)%totalPlayer]
				if rule.checkHu(fsm, nextPlayerId) {
					fsm.EventQueue.put(&Event{nextPlayerId, EVENT_HU})
				}
			}
			fsm.EventQueue.put(&Event{playerId, EVENT_CI})
			return true
		}
	}
	return false
}

func (rule *StandardRule) doPeng(fsm *FSM, playerId string) bool {
	// move peng cards to up cards
	player := fsm.PlayerList[playerId]
	cardSet := &CardSet{Cards: []*Card{fsm.WaitCard}}
	toRemove := []*Card{}
	for _, card := range player.HandCards {
		if rule.isDui(card, fsm.WaitCard) {
			cardSet.Cards = append(cardSet.Cards, card)
			toRemove = append(toRemove, card)
		}
		if len(cardSet.Cards) == 3 {
			break
		}
	}
	if len(cardSet.Cards) != 3 {
		return false
	}
	player.UpCardSet = append(player.UpCardSet, cardSet)
	removeCards(&player.HandCards, toRemove)
	removeCards(
		&fsm.PlayerList[fsm.DependPlayer].OutCards,
		[]*Card{fsm.WaitCard})
	fsm.WaitCard = nil
	fsm.EventQueue.clear()
	fsm.EventQueue.put(&Event{playerId, EVENT_PUSH})
	return true
}

func (rule *StandardRule) doChi(fsm *FSM, playerId string) bool {
	// move chi cards to up cards
	// use player.DoIndex
	player := fsm.PlayerList[playerId]
	indexs := *player.DoIndex
	if len(indexs) != 2 {
		return false
	}
	cardSet := &CardSet{Cards: []*Card{fsm.WaitCard}}
	toRemove := []*Card{}
	for _, card := range player.HandCards {
		for _, index := range indexs {
			if card.Index == index {
				cardSet.Cards = append(cardSet.Cards, card)
				toRemove = append(toRemove, card)
				break
			}
		}
		if len(cardSet.Cards) == 3 {
			break
		}
	}
	player.UpCardSet = append(player.UpCardSet, cardSet)
	removeCards(&player.HandCards, toRemove)
	removeCards(
		&fsm.PlayerList[fsm.DependPlayer].OutCards,
		[]*Card{fsm.WaitCard})
	fsm.WaitCard = nil
	fsm.EventQueue.clear()
	fsm.EventQueue.put(&Event{playerId, EVENT_PUSH})
	return true
}

func (rule *StandardRule) doInitBuhua(fsm *FSM, playerId string) {
	player := fsm.PlayerList[playerId]
	// find all flour cards and add into player.flourCards
	// change by new cards from pool
	// check and add recreaction if need
	for _, card := range player.HandCards {
		if rule.isHua(card) {
			for {
				poolSize := len(fsm.CardPool)
				cards := fsm.CardPool[poolSize-1 : poolSize]
				fsm.CardPool = fsm.CardPool[:poolSize-1]
				player.FlourCards = append(player.FlourCards, cards[0])
				if !rule.isHua(cards[0]) {
					swapCards(card, cards[0])
					break
				}
			}
		}
	}
	if player.WaitCard != nil && rule.isHua(player.WaitCard) {
		for {
			poolSize := len(fsm.CardPool)
			cards := fsm.CardPool[poolSize-1 : poolSize]
			fsm.CardPool = fsm.CardPool[:poolSize-1]
			player.FlourCards = append(player.FlourCards, cards[0])
			if !rule.isHua(cards[0]) {
				swapCards(player.WaitCard, cards[0])
				break
			}
		}
	}
	orderCards(&player.HandCards)
}

func (rule *StandardRule) doInitPull(fsm *FSM, playerId string, num int) {
	poolSize := len(fsm.CardPool)
	cards := fsm.CardPool[poolSize-num : poolSize]
	fsm.CardPool = fsm.CardPool[:poolSize-num]
	player := fsm.PlayerList[playerId]
	for _, card := range cards {
		player.HandCards = append(player.HandCards, card)
	}
	orderCards(&player.HandCards)
	if num == 2 && len(player.HandCards) == 14 {
		player.WaitCard = player.HandCards[13]
		player.HandCards = player.HandCards[:13]
	}
}
