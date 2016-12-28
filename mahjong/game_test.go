package mahjong

import (
	"testing"
)

func TestGameSingle(t *testing.T) {
	game := CreateGame(1, 1)
	game.Join(1)
	if game.status != GAME_WAITING {
		t.Error("game status should be ", GAME_WAITING, "but ", game.status)
	}
	game.Ready(1)
	if game.status != GAME_GAMING {
		t.Error("game status should be ", GAME_GAMING, "but ", game.status)
	}
	if game.curSet.status != SET_WAIT_OP {
		t.Error("set status should be ", SET_WAIT_OP, "but ", game.curSet.status)
	}
	game.Op(1, "op", nil)
	if game.curSet.status != SET_WAIT_OP {
		t.Error("set status should be ", SET_WAIT_OP, "but ", game.curSet.status)
	}
}

func TestGameMulty(t *testing.T) {
	game := CreateGame(2, 2)
	game.Join(1)
	if game.status != GAME_WAITING {
		t.Error("game status should be ", GAME_WAITING, "but ", game.status)
	}
	game.Ready(1)
	if game.status != GAME_WAITING {
		t.Error("game status should be ", GAME_WAITING, "but ", game.status)
	}
	game.Join(2)
	if game.status != GAME_WAITING {
		t.Error("game status should be ", GAME_WAITING, "but ", game.status)
	}
	game.Ready(2)
	if game.status != GAME_GAMING {
		t.Error("game status should be ", GAME_GAMING, "but ", game.status)
	}
	if game.curSet.status != SET_WAIT_OP {
		t.Error("set status should be ", SET_WAIT_OP, "but ", game.curSet.status)
	}
	game.Op(1, "op", nil)
	if game.curSet.status != SET_WAIT_OP {
		t.Error("set status should be ", SET_WAIT_OP, "but ", game.curSet.status)
	}
}
