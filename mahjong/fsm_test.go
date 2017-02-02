package mahjong

import (
	"encoding/json"
	"testing"
)

func toJson(v interface{}) string {
	ret, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(ret)
}

func toJsonCards(cards []*Card) string {
	ret := "|"
	for _, card := range cards {
		ret += card.Name + "|"
	}
	return ret
}

func logFsm(t *testing.T, fsm *FSM) {
	t.Log()
	t.Log("current player num", len(fsm.PlayerList))
	t.Log("max player num", fsm.MaxPlayer)
	t.Log("palyer chain", fsm.PlayerChain)
	t.Log("rule", fsm.Rule)
	t.Log("game status", fsm.GameStatus)
	t.Log("set status", fsm.SetStatus)
	t.Log("game round", fsm.GameRound)
	t.Log("max round", fsm.MaxGameRound)
}

func logPlayers(t *testing.T, fsm *FSM) {
	t.Log()
	t.Log("********")
	t.Log("is wait for op", fsm.WaitOp)
	t.Log("op player", fsm.OpPlayer)
	t.Log("depend player", fsm.DependPlayer)
	t.Log("win info", fsm.WinInfo)
	t.Log("card pool remain", len(fsm.CardPool))
	t.Log("event queue", toJson(fsm.EventQueue))
	t.Log("current event", toJson(fsm.CurEvent))
	for _, playerId := range fsm.PlayerChain {
		player := fsm.PlayerList[playerId]
		t.Log("========")
		t.Log("player", playerId, player.Zhuang)
		t.Log("wait card", toJson(player.WaitCard))
		t.Log("hand card", len(player.HandCards), toJsonCards(player.HandCards))
		for _, cardSet := range player.UpCardSet {
			t.Log("up card", toJsonCards(cardSet.Cards))
		}
		for _, cardSet := range player.DownCardSet {
			t.Log("down card", toJsonCards(cardSet.Cards))
		}
		t.Log("out card", len(player.OutCards), toJsonCards(player.OutCards))
		t.Log("flour card", toJsonCards(player.FlourCards))
	}
}

func Test2Player(t *testing.T) {
	fsm := CreateFSM("standard", 2)
	fsm.Join("1")
	fsm.Join("2")
	logFsm(t, fsm)
	logPlayers(t, fsm)
	for i := 0; i < 200; i++ {
		curPlayer := fsm.CurEvent.PlayerId
		if fsm.CurEvent.Event != 7 {
			t.Log("player", curPlayer, "do", 12)
			fsm.Op(curPlayer, 12, []int{fsm.PlayerList[curPlayer].HandCards[0].Index})
		} else {
			t.Log("player", curPlayer, "do", 7)
			fsm.Op(curPlayer, 7, []int{fsm.PlayerList[curPlayer].HandCards[0].Index})
		}
		logPlayers(t, fsm)
		if fsm.SetStatus != 2 {
			logFsm(t, fsm)
			return
		}
	}
	t.Error("set not end")
}

func Test4Player(t *testing.T) {
	fsm := CreateFSM("standard", 4)
	fsm.Join("1")
	fsm.Join("2")
	fsm.Join("3")
	fsm.Join("4")
	logFsm(t, fsm)
	logPlayers(t, fsm)
	for i := 0; i < 200; i++ {
		curPlayer := fsm.CurEvent.PlayerId
		if fsm.CurEvent.Event != 7 {
			t.Log("player", curPlayer, "do", 12)
			fsm.Op(curPlayer, 12, []int{fsm.PlayerList[curPlayer].HandCards[0].Index})
		} else {
			t.Log("player", curPlayer, "do", 7)
			fsm.Op(curPlayer, 7, []int{fsm.PlayerList[curPlayer].HandCards[0].Index})
		}
		logPlayers(t, fsm)
		if fsm.SetStatus != 2 {
			logFsm(t, fsm)
			return
		}
	}
	t.Error("set not end")
}

func Test2PlayerMoreRound(t *testing.T) {
	fsm := CreateFSM("standard", 2)
	fsm.Join("1")
	fsm.Join("2")
	logFsm(t, fsm)
	logPlayers(t, fsm)
	for i := 0; i < 200; i++ {
		curPlayer := fsm.CurEvent.PlayerId
		if fsm.CurEvent.Event != 7 {
			t.Log("player", curPlayer, "do", 12)
			fsm.Op(curPlayer, 12, []int{fsm.PlayerList[curPlayer].HandCards[0].Index})
		} else {
			t.Log("player", curPlayer, "do", 7)
			fsm.Op(curPlayer, 7, []int{fsm.PlayerList[curPlayer].HandCards[0].Index})
		}
		logPlayers(t, fsm)
		if fsm.SetStatus != 2 {
			logFsm(t, fsm)
		}
	}
	if fsm.SetStatus != 3 {
		t.Error("set not end")
	}
	fsm.Next()
	logFsm(t, fsm)
	logPlayers(t, fsm)
	for i := 0; i < 200; i++ {
		curPlayer := fsm.CurEvent.PlayerId
		if fsm.CurEvent.Event != 7 {
			t.Log("player", curPlayer, "do", 12)
			fsm.Op(curPlayer, 12, []int{fsm.PlayerList[curPlayer].HandCards[0].Index})
		} else {
			t.Log("player", curPlayer, "do", 7)
			fsm.Op(curPlayer, 7, []int{fsm.PlayerList[curPlayer].HandCards[0].Index})
		}
		logPlayers(t, fsm)
		if fsm.SetStatus != 2 {
			logFsm(t, fsm)
		}
	}
}
