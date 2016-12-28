package mahjong

type Game struct {
	maxPlayer   int             // 指定玩家数
	maxSetNum   int             // 指定局数
	curSetNum   int             // 当前局数
	curSet      *Set            // 当前局
	status      GameStatus      // 游戏状态
	playerChain map[int]*Player // 玩家链
	rule        Rule
}

type Set struct {
	playerTotal         int
	curPlayerIndex      int             // 当前操作玩家
	curEvent            *Event          // 当前事件
	curEventPlayerIndex int             // 当前事件操作玩家
	curEventCard        *Card           // 当前事件的牌
	cardPool            []*Card         // 牌池
	players             map[int]*Player // 所有玩家
	playerChain         map[int]*Player // 玩家链
	status              SetStatus
	rule                RuleImpl
}

// 游戏状态
type GameStatus int

const (
	_            GameStatus = iota
	GAME_WAITING            // 等待开始
	GAME_GAMING             // 进行中
	GAME_FINISH             // 结束
)

// 局内状态
type SetStatus int

const (
	_             SetStatus = iota
	SET_NEED_INIT           // 等待初始化
	SET_WAIT_OP             // 等待用户操作
	SET_OP_DONE             // 等待系统处理
)

type Player struct {
	id      int  // 玩家id
	prepare bool // 玩家准备状态
	zhuang  bool

	waitCard    *Card     // 当前牌
	handCards   []*Card   // 手牌
	upCardSet   []CardSet // 亮着的牌
	downCardSet []CardSet // 扣着的牌
	outCards    []*Card   // 出过的牌
	flourCards  []*Card   // 花牌

	curEvent *Event
}

type Event struct {
	name EventName // 事件名称
	rule EventRule // 事件规则
}

// 事件定义
type EventName int

const (
	_           EventName = iota
	EVENT_BUHUA           // 补花
	EVENT_ZIMO            // 自摸
	EVENT_HU              // 胡
	EVENT_GANG            // 杠
	EVENT_PENG            // 碰
	EVENT_CHI             // 吃
	EVENT_PUSH            // 出牌
	EVENT_PULL            // 抓牌
	EVENT_PULL2           // 抓牌
	EVENT_PULL4           // 抓牌
)

// 事件规则定义
type EventRule int

const (
	_                EventRule = iota
	EVENT_RULE_SELF            //自己
	EVENT_RULE_NEXT            //下家
	EVENT_RULE_OTHER           //其他人
	EVENT_RULE_ALL             //所有人
)

func CreateGame(maxPlayer, maxSet int) *Game {
	game := &Game{
		maxPlayer:   maxPlayer,
		maxSetNum:   maxSet,
		curSetNum:   0,
		status:      GAME_WAITING,
		playerChain: make(map[int]*Player),
		rule:        RULE_STANDARD}
	return game
}

func (game *Game) autoContral() {
	switch game.status {
	case GAME_WAITING:
		if len(game.playerChain) == game.maxPlayer && game.allReady() {
			game.createSet()
			return
		}
		break
	case GAME_GAMING:
		game.curSet.autoContral()
		break
	default:
		break
	}
}

func (game *Game) Join(playerId int) {
	player := &Player{
		id:      playerId,
		prepare: false}
	game.playerChain[playerId] = player
}

func (game *Game) Ready(playerId int) {
	game.playerChain[playerId].prepare = true
	game.autoContral()
}

func (game *Game) allReady() bool {
	for _, player := range game.playerChain {
		if !player.prepare {
			return false
		}
	}
	return true
}

func (game *Game) createSet() *Set {
	// TODO 选择规则
	rule := getRule(RULE_STANDARD)
	set := &Set{
		playerTotal: len(game.playerChain),
		players:     make(map[int]*Player),
		playerChain: make(map[int]*Player),
		status:      SET_NEED_INIT,
		rule:        rule}
	// TODO 随机分配seq
	curSeq := 0
	for playerId, player := range game.playerChain {
		set.players[playerId] = player
		set.playerChain[curSeq] = player
		curSeq += 1
	}
	game.curSet = set
	game.status = GAME_GAMING
	game.autoContral()
	return set
}

func (game *Game) Op(playerId int, op string, index []int) {
	game.curSet.status = SET_OP_DONE
	game.autoContral()
}

func (set *Set) autoContral() {
	switch set.status {
	case SET_NEED_INIT:
		set.init()
		break
	case SET_WAIT_OP:
		switch set.curEvent.name {
		case EVENT_PULL4, EVENT_PULL2, EVENT_PULL:
			// 系统自动代理操作
			set.autoProcess()
		default:
			// wait for user op
			return
		}
		break
	case SET_OP_DONE:
		set.process()
		break
	default:
		break
	}
}

func (set *Set) autoProcess() {
	switch set.curEvent.name {
	case EVENT_PULL4:
		set.rule.getCards(4, &set.cardPool, set.curPlayer())
		break
	case EVENT_PULL2:
		set.rule.getCards(2, &set.cardPool, set.curPlayer())
		break
	case EVENT_PULL:
		set.rule.getCards(1, &set.cardPool, set.curPlayer())
		break
	}
	set.status = SET_OP_DONE
	set.autoContral()
}

func (set *Set) init() {
	// 初始化牌局
	set.cardPool = *set.rule.allCards()
	set.curPlayerIndex = 0
	set.curPlayer().zhuang = true
	set.curEvent = set.rule.getNextEvent(set.curPlayer())
	set.status = SET_WAIT_OP
}

func (set *Set) process() {
	// 让牌局进行
	switch set.curEvent.name {
	case EVENT_PULL4, EVENT_PULL2:
		set.curPlayerIndex = (set.curPlayerIndex + 1) % set.playerTotal
		set.curEvent = set.rule.getNextEvent(set.curPlayer())
		break
	case EVENT_PULL:
		set.curEvent = set.rule.getNextEvent(set.curPlayer())
	}
	set.status = SET_WAIT_OP
	set.autoContral()
}

func (set *Set) curPlayer() *Player {
	return set.playerChain[set.curPlayerIndex]
}
