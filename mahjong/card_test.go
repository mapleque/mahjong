package mahjong

import (
	"testing"
)

type CardTest struct {
	Name  string
	Index int
	Type  CardType
	Value int
}

var list = []CardTest{
	CardTest{"一萬", 1, CARD_TYPE_SEQ, 1},
	CardTest{"一萬", 11, CARD_TYPE_SEQ, 1},
	CardTest{"九萬", 29, CARD_TYPE_SEQ, 9},
	CardTest{"五条", 115, CARD_TYPE_SEQ, 25},
	CardTest{"东风", 121, CARD_TYPE_WORD, 31},
	CardTest{"發财", 156, CARD_TYPE_WORD, 36},
	CardTest{"菊", 168, CARD_TYPE_FLOUR, 48}}

// 测试创建各种牌是否正确
func TestCreateCard(t *testing.T) {
	card := createCard(1)
	if card == nil {
		t.Error("card shuld not be nil")
	}
	for _, data := range list {
		card := createCard(data.Index)
		if data.Name != card.Name {
			t.Error("card Name should be", data.Name, "but", card.Name)
		}
		if data.Type != card.Type {
			t.Error("card Value should be", data.Type, "but", card.Type)
		}
		if data.Value != card.Value {
			t.Error("card Value should be", data.Value, "but", card.Value)
		}
	}
}

// 测试创建完整的一套牌是否正确
func TestCreateAll(t *testing.T) {
	pool := createAllCards()
	if len(pool) != 144 {
		t.Error("card number should be 144 but", len(pool))
	}
	lastIndex := 0
	for _, card := range pool {
		if card.Index < lastIndex {
			t.Error("card index should be in seq but", lastIndex, card.Index)
		}
		lastIndex = card.Index
	}
}

// 测试排序和乱序
func TestDisOrderCards(t *testing.T) {
	pool := createAllCards()
	disOrderCards(&pool)
	if pool[0].Index == 1 && pool[1].Index == 2 {
		t.Error("pool is not disorder", pool[0].Index)
	}
	orderCards(&pool)
	if pool[0].Value != 1 || pool[144-1].Value != 48 {
		t.Error("pool is not in order", pool[0].Index)
	}
}

// 测试随机方法
func TestRandIndex(t *testing.T) {
	if randIndex(1) != 0 {
		t.Error("wrong result on rand 1")
	}
	if randIndex(0) != 0 {
		t.Error("wrong result on rand 0")
	}
	max := 100
	var seed [10000]int
	for i := 0; i < 10000; i++ {
		r := randIndex(max)
		if r < 0 || r > 100 {
			t.Error("wrong result during rand 100", r)
		}
		seed[i] = r
	}
	sum := 0
	for _, r := range seed {
		sum += r
	}
	avg := sum / len(seed)
	sumstd := 0
	for _, r := range seed {
		sumstd += (r - avg) * (r - avg)
	}
	std := sumstd / len(seed)
	t.Log("rand 10000 times from 0 to 100 result", "avg", avg, "std", std)
}

// 测试交换方法
func TestSwapCards(t *testing.T) {
	for i := 1; i < len(list); i++ {
		card_a := createCard(list[i].Index)
		card_b := createCard(list[i-1].Index)
		swapCards(card_a, card_b)
		if card_a.Index != list[i-1].Index || card_b.Index != list[i].Index {
			t.Error("swap card faild, cards not swap", card_a.Index, card_b.Index)
		}
	}
}
