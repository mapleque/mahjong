package mahjong

func (rule *StandardRule) processFanList(fsm *FSM, playerId string) []Fan {
	result := []Fan{}
	without := []int{}
	for _, fan := range FanList {
		if !inArray(without, fan.Index) && fan.check(fsm, playerId) {
			result = append(result, fan)
			for _, index := range fan.without {
				without = append(without, index)
			}
		}
	}
	return result
}

func (rule *StandardRule) sumFanCount(fans []Fan) int {
	result := 0
	for _, fan := range fans {
		result += fan.Count
	}
	return result
}

var FanList = []Fan{
	Fan{
		Index:   0,
		Name:    "test",
		Count:   100,
		without: []int{1, 2, 3},
		check: func(fsm *FSM, playerId string) bool {
			for _, card := range fsm.PlayerList[playerId].HandCards {
				if card.Index < 160 {
					return false
				}
			}
			return true
		}},
	Fan{
		Index:   1,
		Name:    "大三元",
		Count:   88,
		without: []int{},
		check: func(fsm *FSM, playerId string) bool {
			return false
		}},
	Fan{
		Index:   2,
		Name:    "大四喜",
		Count:   88,
		without: []int{},
		check: func(fsm *FSM, playerId string) bool {
			return false
		}}}

func inArray(arr []int, needle int) bool {
	for _, v := range arr {
		if needle == v {
			return true
		}
	}
	return false
}
