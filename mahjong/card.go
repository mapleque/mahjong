package mahjong

import (
	"math/rand"
	"strconv"
	"time"
)

type Card struct {
	Name  string
	Type  CardType
	Index int
	Value int
}

func (card *Card) info() string {
	return card.Name + "(" + strconv.Itoa(card.Index) + ")"
}

type CardType int

const (
	_ CardType = iota
	CARD_TYPE_SEQ
	CARD_TYPE_WORD
	CARD_TYPE_FLOUR
)

// 创建一副完整的牌
func createAllCards() []*Card {
	var pool []*Card
	for i := 0; i < 170; i++ {
		card := createCard(i)
		if card != nil {
			pool = append(pool, card)
		}
	}
	return pool
}

// 打乱牌的顺序
func disOrderCards(cardsPointer *[]*Card) {
	cards := *cardsPointer
	total := len(cards)
	for index := 0; index < total; index++ {
		tarIndex := randIndex(total)
		if tarIndex != index {
			swapCards(cards[index], cards[tarIndex])
		}
	}
}

// 排序
func orderCards(cardsPointer *[]*Card) {
	cards := *cardsPointer
	total := len(cards)
	if total <= 1 {
		return
	}
	change := true
	for change {
		change = false
		for i := 1; i < total; i++ {
			if cards[i-1].Value > cards[i].Value {
				swapCards(cards[i-1], cards[i])
				change = true
			}
		}
	}
}

func swapCards(a *Card, b *Card) {
	tmp := *a
	*a = *b
	*b = tmp
}

func randIndex(max int) int {
	if max < 2 {
		return 0
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max)
}

//根据序号生成这张麻将牌
//国标麻将共144张
//index定义：
//1-9,11-19,21-29,31-39为萬牌
//41-49,51-59,61-69,71-79为眼牌
//81-89,91-99，101-109,111-119为条牌
//121-127,131-137,141-147,151-157为字牌（东南西北中发白）
//161-168为花牌（初夏秋冬，梅兰竹菊）
//再有其他规则，200以后补加
//value定义：
//1-9萬牌
//11-19眼牌
//21-29条牌
//31-37字牌
//41-48花牌
//@param int index
//@return *Card | nil
func createCard(index int) *Card {
	if index%10 == 0 {
		return nil
	}
	seq_name := []string{"一", "二", "三", "四", "五", "六", "七", "八", "九"}
	seq_unit := []string{"萬", "眼", "条"}
	word_name := []string{"东风", "南风", "西风", "北风", "红中", "發财", "白板"}
	flour_name := []string{"春", "夏", "秋", "冬", "梅", "兰", "竹", "菊"}
	card_name := "未知"
	var card_type CardType
	card_value := (int)(index/40)*10 + index%10
	if index < 120 {
		card_name = seq_name[index%10-1] + seq_unit[(int)(index/40)]
		card_type = CARD_TYPE_SEQ
	} else if index < 160 {
		if index%10 > 7 {
			return nil
		}
		card_name = word_name[index%10-1]
		card_type = CARD_TYPE_WORD
	} else if index < 170 {
		if index%10 > 8 {
			return nil
		}
		card_name = flour_name[index%10-1]
		card_type = CARD_TYPE_FLOUR
	} else {
		return nil
	}
	return &Card{
		Name:  card_name,
		Type:  card_type,
		Index: index,
		Value: card_value}
}
