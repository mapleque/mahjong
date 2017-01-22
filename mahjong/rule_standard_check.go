package mahjong

func (rule *StandardRule) checkBuhua(fsm *FSM, playerId string) bool {
	player := fsm.PlayerList[playerId]
	for _, card := range player.HandCards {
		if rule.isHua(card) {
			return true
		}
	}
	return rule.isHua(player.WaitCard)
}
func (rule *StandardRule) checkHu(fsm *FSM, playerId string) bool {
	fans := rule.processFanList(fsm, playerId)
	fanCount := rule.sumFanCount(fans)
	return fanCount >= 8
}
func (rule *StandardRule) checkGang(fsm *FSM, playerId string) bool {
	player := fsm.PlayerList[playerId]
	if fsm.DependPlayer != "" {
		// gang from other card
		cards := []*Card{fsm.WaitCard}
		for _, card := range player.HandCards {
			if rule.isDui(card, fsm.WaitCard) {
				cards = append(cards, card)
			}
			if len(cards) == 4 {
				return true
			}
		}
	} else {
		// gang from self
		cards := []*Card{player.WaitCard}
		// down
		for _, card := range player.HandCards {
			if rule.isDui(card, player.WaitCard) {
				cards = append(cards, card)
			}
			if len(cards) == 4 {
				return true
			}
		}
		// up
		for _, holdCardSet := range player.UpCardSet {
			if rule.isGang(
				holdCardSet.Cards[0],
				holdCardSet.Cards[1],
				holdCardSet.Cards[2],
				player.WaitCard) {

				return true
			}
		}
	}
	return false
}
func (rule *StandardRule) checkPeng(fsm *FSM, playerId string) bool {
	player := fsm.PlayerList[playerId]
	cards := []*Card{fsm.WaitCard}
	for _, card := range player.HandCards {
		if rule.isDui(card, fsm.WaitCard) {
			cards = append(cards, card)
		}
		if len(cards) == 3 {
			return true
		}
	}
	return false
}
func (rule *StandardRule) checkChi(fsm *FSM, playerId string) bool {
	player := fsm.PlayerList[playerId]
	for _, card1 := range player.HandCards {
		for _, card2 := range player.HandCards {
			if rule.isShun(card1, card2, fsm.WaitCard) {
				return true
			}
		}
	}
	return false
}

func (rule *StandardRule) isHua(card *Card) bool {
	return card != nil && card.Index > 160
}
func (rule *StandardRule) isGang(card1, card2, card3, card4 *Card) bool {
	return rule.isDui(card1, card2) &&
		rule.isDui(card1, card3) &&
		rule.isDui(card1, card4)
}
func (rule *StandardRule) isKe(card1, card2, card3 *Card) bool {
	return rule.isDui(card1, card2) &&
		rule.isDui(card1, card3)
}
func (rule *StandardRule) isShun(card1, card2, card3 *Card) bool {
	var min, mid, max *Card
	if card1.Value < card2.Value {
		min = card1
		max = card2
	} else {
		min = card2
		max = card1
	}
	if card3.Value < min.Value {
		mid = min
		min = card3
	} else if card3.Value > max.Value {
		mid = max
		max = card3
	} else {
		mid = card3
	}
	return max.Value < 30 &&
		(max.Value-mid.Value) == 1 &&
		(mid.Value-min.Value) == 1
}
func (rule *StandardRule) isDui(card1, card2 *Card) bool {
	return card1.Value == card2.Value
}
