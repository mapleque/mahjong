package mahjong

type StandardRule struct{}

func (rule *StandardRule) init(fsm *FSM) {
	//card pool
	pool := createAllCards()
	disOrderCards(&pool)
	fsm.CardPool = pool
	// zhuang
	for index, playerId := range fsm.PlayerChain {
		fsm.PlayerList[playerId].Zhuang =
			(index+fsm.GameRound)%len(fsm.PlayerChain) == 0
	}
	// add player event
	fsm.EventQueue = &EventQueue{}
	fsm.EventQueue.clear()
	rule.initEventQueue(fsm)
}

func (rule *StandardRule) process(fsm *FSM) bool {
	// pop event
	event, empty := fsm.EventQueue.pop()
	if empty {
		fsm.CurEvent = nil
		// finish set
		return false
	}
	// do it or wait for op
	fsm.CurEvent = event
	rule.dealEvent(fsm, event)
	if !fsm.WaitOp {
		return rule.process(fsm)
	}
	return true
}

func (rule *StandardRule) op(fsm *FSM) bool {
	doPlayer := fsm.PlayerList[fsm.OpPlayer]
	if doPlayer.DoEvent != fsm.CurEvent.Event &&
		doPlayer.DoEvent != EVENT_PASS {
		return false
	}
	switch doPlayer.DoEvent {
	case EVENT_HU:
		if !rule.doHu(fsm, fsm.OpPlayer) {
			return false
		}
	case EVENT_GANG:
		if !rule.doGang(fsm, fsm.OpPlayer) {
			return false
		}
	case EVENT_PENG:
		if !rule.doPeng(fsm, fsm.OpPlayer) {
			return false
		}
	case EVENT_CHI:
		if !rule.doChi(fsm, fsm.OpPlayer) {
			return false
		}
	case EVENT_PUSH:
		if !rule.doPush(fsm, fsm.OpPlayer) {
			return false
		}
	case EVENT_PASS:
		// do nothing
	}
	fsm.WaitOp = false
	return true
}
