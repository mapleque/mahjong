package mahjong

import (
	"testing"
)

func logCards(cards []*Card) string {
	str := ""
	for _, card := range cards {
		str += " " + logCard(card)
	}
	return str
}
func logCard(card *Card) string {
	return card.info()
}

func TestGetCard(t *testing.T) {
	player := &Player{}
	pool := createAllCards()
	total := len(pool)
	disOrderCards(&pool)
	getCards(4, &pool, player)
	if len(player.handCards) != 4 {
		t.Error("get card faild", logCards(player.handCards))
	}
	t.Log(len(pool))
	if len(pool) != total-4 {
		t.Error("get card faild", len(pool), total, "pool not change")
	}
	getCards(4, &pool, player)
	t.Log(logCards(player.handCards))
	getCards(4, &pool, player)
	t.Log(logCards(player.handCards))
	getCards(2, &pool, player)
	t.Log(logCards(player.handCards))
	if len(player.handCards) != 13 {
		t.Error("get card faild", logCards(player.handCards))
	}
	if player.waitCard == nil {
		t.Error("get card faild", logCards(player.handCards))
	}
}

func TestPushCard(t *testing.T) {
	pool := createAllCards()
	disOrderCards(&pool)
	player := &Player{}
	getCards(4, &pool, player)
	getCards(1, &pool, player)
	pushIndex := player.handCards[0].Index
	player.doIndex = &[]int{pushIndex}
	card := pushCard(player)
	if card.Index != pushIndex {
		t.Error("push card increct",
			logCards(player.handCards),
			logCard(player.waitCard),
			logCard(card), pushIndex)
	}
	if len(player.handCards) != 4 {
		t.Error("push card faild", len(player.handCards), "wrong hand card num")
	}
	for _, curCard := range player.handCards {
		if curCard.Index == pushIndex {
			t.Error("push card faild", pushIndex, "is still in hand")
		}
	}
}

func TestIsGang(t *testing.T) {
	handCards := []*Card{
		createCard(1),
		createCard(11),
		createCard(31)}
	waitCard := createCard(21)
	if !isGang(handCards, []*CardSet{}, waitCard) {
		t.Error("cards should be gang", logCards(handCards), logCard(waitCard))
	}
	waitCardWrong := createCard(22)
	if isGang(handCards, []*CardSet{}, waitCardWrong) {
		t.Error("cards should not be gang",
			logCards(handCards), logCard(waitCardWrong))
	}
}

func TestIsPen(t *testing.T) {
	handCards := []*Card{
		createCard(11),
		createCard(31)}
	waitCard := createCard(21)
	if !isPeng(handCards, waitCard) {
		t.Error("cards should be peng", logCards(handCards), logCard(waitCard))
	}
	waitCardWrong := createCard(22)
	if isPeng(handCards, waitCard) {
		t.Error("cards should not be peng",
			logCards(handCards), logCard(waitCardWrong))
	}
}

func TestIsChi(t *testing.T) {
	handCards := []*Card{
		createCard(11),
		createCard(32)}
	waitCard := createCard(23)
	if !isChi(handCards, waitCard) {
		t.Error("cards should be chi", logCards(handCards), logCard(waitCard))
	}
	waitCardWrong := createCard(22)
	if isChi(handCards, waitCard) {
		t.Error("cards should not be chi",
			logCards(handCards), logCard(waitCardWrong))
	}
}
