package mahjong

import (
	"testing"
)

func TestIsHu(t *testing.T) {
}

func TestIsGang(t *testing.T) {
	rule := &StandardRule{}
	if !rule.isGang(
		createCard(1),
		createCard(11),
		createCard(21),
		createCard(31)) {
		t.Error("cards should be gang")
	}
	if rule.isGang(
		createCard(1),
		createCard(12),
		createCard(21),
		createCard(31)) {
		t.Error("cards should not be gang")
	}
}

func TestIsKe(t *testing.T) {
	rule := &StandardRule{}
	if !rule.isKe(
		createCard(11),
		createCard(21),
		createCard(31)) {
		t.Error("cards should be ke")
	}
	if rule.isKe(
		createCard(12),
		createCard(21),
		createCard(31)) {
		t.Error("cards should not be ke")
	}
}

func TestIsShun(t *testing.T) {
	rule := &StandardRule{}
	if !rule.isShun(
		createCard(11),
		createCard(22),
		createCard(33)) {
		t.Error("cards should be shun")
	}
	if !rule.isShun(
		createCard(22),
		createCard(11),
		createCard(33)) {
		t.Error("cards should be shun")
	}
	if rule.isShun(
		createCard(121),
		createCard(122),
		createCard(123)) {
		t.Error("cards should not be shun")
	}
	if rule.isShun(
		createCard(1),
		createCard(12),
		createCard(31)) {
		t.Error("cards should not be shun")
	}
}

func TestIsDui(t *testing.T) {
	rule := &StandardRule{}
	if !rule.isDui(
		createCard(21),
		createCard(31)) {
		t.Error("cards should be dui")
	}
	if rule.isDui(
		createCard(12),
		createCard(31)) {
		t.Error("cards should not be dui")
	}
}

func TestCheckBuhua(t *testing.T) {
	rule := &StandardRule{}
	fsm := &FSM{
		PlayerList: make(map[string]*Player)}
	player := &Player{}
	fsm.PlayerList["1"] = player
	player.HandCards = []*Card{createCard(161)}
	if !rule.checkBuhua(fsm, "1") {
		t.Error("unexpect result")
	}
	player.HandCards = []*Card{createCard(1)}
	if rule.checkBuhua(fsm, "1") {
		t.Error("unexpect result")
	}
	player.HandCards = []*Card{createCard(1)}
	player.WaitCard = createCard(168)
	if !rule.checkBuhua(fsm, "1") {
		t.Error("unexpect result")
	}
}

func TestCheckHu(t *testing.T) {
}

func TestCheckGang(t *testing.T) {
	rule := &StandardRule{}
	fsm := &FSM{
		PlayerList: make(map[string]*Player)}
	player := &Player{}
	fsm.PlayerList["1"] = player

	fsm.DependPlayer = ""
	// self down gang
	player.HandCards = []*Card{
		createCard(2),
		createCard(3),
		createCard(1),
		createCard(11),
		createCard(21)}
	player.WaitCard = createCard(31)
	if !rule.checkGang(fsm, "1") {
		t.Error("unexpect result")
	}

	// self up gang
	player.HandCards = []*Card{
		createCard(2),
		createCard(32),
		createCard(23)}
	player.WaitCard = createCard(31)
	player.UpCardSet = []*CardSet{&CardSet{Cards: []*Card{
		createCard(1),
		createCard(11),
		createCard(21)}}}
	if !rule.checkGang(fsm, "1") {
		t.Error("unexpect result")
	}

	// self no gang
	player.HandCards = []*Card{
		createCard(2),
		createCard(32),
		createCard(23)}
	player.WaitCard = createCard(31)
	player.UpCardSet = []*CardSet{&CardSet{Cards: []*Card{
		createCard(3),
		createCard(14),
		createCard(25)}}}
	if rule.checkGang(fsm, "1") {
		t.Error("unexpect result")
	}

	fsm.DependPlayer = "2"
	fsm.WaitCard = createCard(31)
	player.WaitCard = nil
	// other up gang
	player.HandCards = []*Card{
		createCard(1),
		createCard(11),
		createCard(21)}
	if !rule.checkGang(fsm, "1") {
		t.Error("unexpect result")
	}
	// other no gang
	player.HandCards = []*Card{
		createCard(3),
		createCard(11),
		createCard(21)}
	if rule.checkGang(fsm, "1") {
		t.Error("unexpect result")
	}
}

func TestCheckPeng(t *testing.T) {
	rule := &StandardRule{}
	fsm := &FSM{
		PlayerList: make(map[string]*Player)}
	player := &Player{}
	fsm.PlayerList["1"] = player
	fsm.DependPlayer = "2"
	fsm.WaitCard = createCard(31)
	player.WaitCard = nil
	// peng
	player.HandCards = []*Card{
		createCard(2),
		createCard(11),
		createCard(21)}
	if !rule.checkPeng(fsm, "1") {
		t.Error("unexpect result")
	}
	// no peng
	player.HandCards = []*Card{
		createCard(3),
		createCard(12),
		createCard(21)}
	if rule.checkPeng(fsm, "1") {
		t.Error("unexpect result")
	}
}

func TestCheckChi(t *testing.T) {
	rule := &StandardRule{}
	fsm := &FSM{
		PlayerList: make(map[string]*Player)}
	player := &Player{}
	fsm.PlayerList["1"] = player
	fsm.DependPlayer = "2"
	fsm.WaitCard = createCard(31)
	player.WaitCard = nil
	// chi
	player.HandCards = []*Card{
		createCard(2),
		createCard(13),
		createCard(21)}
	if !rule.checkChi(fsm, "1") {
		t.Error("unexpect result")
	}
	// no chi
	player.HandCards = []*Card{
		createCard(1),
		createCard(12),
		createCard(21)}
	if rule.checkChi(fsm, "1") {
		t.Error("unexpect result")
	}
}
