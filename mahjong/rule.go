package mahjong

type Rule int

type RuleGuard interface {
	init(*FSM)
	//@return bool continue
	process(*FSM) bool
	//@return bool success
	op(*FSM) bool
}

func createRuleGuard(fsm *FSM) RuleGuard {
	switch fsm.Rule {
	case "standard":
		guard := &StandardRule{}
		guard.init(fsm)
		return guard
	default:
		return nil
	}
}
