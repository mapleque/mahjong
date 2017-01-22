package mahjong

// finit status machine
type FSM struct {
	PlayerList   map[string]*Player `json:"playerList"`   // player list
	MaxPlayer    int                `json:"maxPlayer"`    // max player attend the game
	Rule         string             `json:"rule"`         // game rule
	GameStatus   GameStatus         `json:"gameStatus"`   // game status
	SetStatus    SetStatus          `json:"setStatus"`    // set status
	GameRound    int                `json:"gameRound"`    // game round now
	MaxGameRound int                `json:"maxGameRound"` // max game round

	PlayerChain []string    `json:"playerChain"` // player sequence in set
	RuleGuard   RuleGuard   `json:"ruleGuard"`   // a guard process the set
	EventQueue  *EventQueue `json:"eventQueue"`  // a event stack hold the follow events
	WaitOp      bool        `json:"waitOp"`      // if waiting for user op
	OpPlayer    string      `json:"opPlayer"`    // op playerId
	CurEvent    *Event      `json:"curEvent"`    // current event

	CardPool     []*Card `json:"cardPool"`     // card pool now
	WaitCard     *Card   `json:"waitCard"`     // card wait for op
	DependPlayer string  `json:"dependPlayer"` // last depend playerId

	WinInfo *WinInfo `json:"winInfo"`
}

type Player struct {
	Zhuang bool `json:"zhuang"`

	WaitCard    *Card      `json:"waitCard"`    // 当前牌
	HandCards   []*Card    `json:"handCards"`   // 手牌
	UpCardSet   []*CardSet `json:"upCardSet"`   // 亮着的牌
	DownCardSet []*CardSet `json:"downCardSet"` // 扣着的牌
	OutCards    []*Card    `json:"outCards"`    // 出过的牌
	FlourCards  []*Card    `json:"flourCards"`  // 花牌

	DoEvent EventName `json:"doEvent"`
	DoIndex *[]int    `json:"doIndex"`
}

// 游戏状态
type GameStatus int

const (
	_                   GameStatus = iota
	GAME_STATUS_PREPARE            // 等待开始
	GAME_STATUS_ING                // 进行中
	GAME_STATUS_END                // 结束
)

// 局内状态
type SetStatus int

const (
	_                  SetStatus = iota
	SET_STATUS_PREPARE           // 等待初始化
	SET_STATUS_ING               // 等待用户操作
	SET_STATUS_END               // 等待系统处理
)

func CreateFSM(rule string, maxPlayer int) *FSM {
	fsm := &FSM{
		PlayerList:   make(map[string]*Player),
		MaxPlayer:    maxPlayer,
		Rule:         rule,
		GameStatus:   GAME_STATUS_PREPARE,
		SetStatus:    SET_STATUS_PREPARE,
		GameRound:    0,
		MaxGameRound: 4}
	return fsm
}

func (fsm *FSM) Join(playerId string) bool {
	if len(fsm.PlayerList) >= fsm.MaxPlayer {
		return false
	}
	if fsm.GameStatus != GAME_STATUS_PREPARE {
		return false
	}
	fsm.PlayerList[playerId] = &Player{
		HandCards:   []*Card{},
		UpCardSet:   []*CardSet{},
		DownCardSet: []*CardSet{},
		OutCards:    []*Card{},
		FlourCards:  []*Card{},
		DoIndex:     &[]int{}}
	fsm.process()
	return true
}

func (fsm *FSM) Op(playerId string, op int, indexs []int) bool {
	fsm.OpPlayer = playerId
	fsm.PlayerList[playerId].DoEvent = EventName(op)
	fsm.PlayerList[playerId].DoIndex = &indexs
	ret := fsm.RuleGuard.op(fsm)
	if ret {
		fsm.process()
	}
	return ret
}

func (fsm *FSM) Next() {
	fsm.process()
}

func (fsm *FSM) process() {
	switch fsm.GameStatus {
	case GAME_STATUS_PREPARE:
		if len(fsm.PlayerList) >= fsm.MaxPlayer {
			// TODO shuffle player chain
			// player chain must be an ordered list
			for playerId, _ := range fsm.PlayerList {
				fsm.PlayerChain = append(fsm.PlayerChain, playerId)
			}
			fsm.GameStatus = GAME_STATUS_ING
			fsm.process()
		}
	case GAME_STATUS_ING:
		// deal set
		switch fsm.SetStatus {
		case SET_STATUS_PREPARE:
			// TODO wait for every player ready
			fsm.RuleGuard = createRuleGuard(fsm)
			fsm.SetStatus = SET_STATUS_ING
			fsm.process()
		case SET_STATUS_ING:
			if !fsm.RuleGuard.process(fsm) {
				fsm.SetStatus = SET_STATUS_END
			}
		case SET_STATUS_END:
			if fsm.GameRound >= fsm.MaxGameRound {
				fsm.GameStatus = GAME_STATUS_END
			} else {
				fsm.GameRound++
				fsm.SetStatus = SET_STATUS_PREPARE
			}
			fsm.process()
		}
	case GAME_STATUS_END:
		// do nothing
	}
}
