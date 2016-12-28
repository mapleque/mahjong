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
	doEvent  *Event
}

// 创建一个游戏
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

// 用户加入一个游戏
func (game *Game) Join(playerId int) {
	player := &Player{
		id:      playerId,
		prepare: false}
	game.playerChain[playerId] = player
}

// 用户退出一个游戏
func (game *Game) Out(playerId int) {
	if _, ok := playerChain[playerId]; ok {
		delete(playerChain, playerId)
	}
}

// 用户准备开始游戏
// 有这一步是给用户反悔时间，如果用户看有不喜欢的人，就可以不开始并退出
func (game *Game) Ready(playerId int) {
	game.playerChain[playerId].prepare = true
	game.process()
}

func (game *Game) process() {
	switch game.status {
	case GAME_WAITING:
		// 如果是等待状态，要先看人齐不齐，再看是不是都准备了
		// 满足上面的条件，就创建牌局
		// 否则继续等
		if len(game.playerChain) == game.maxPlayer && game.allReady() {
			game.createSet()
			return
		}
		break
	case GAME_GAMING:
		// 如果是进行状态，说明已经有牌局了，直接执行牌局的process
		game.curSet.process()
		break
	default:
		break
	}
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
	game.process()
	return set
}

func (game *Game) Op(playerId int, op string, index []int) {
	// TODO op do some thing
	game.curSet.status = SET_OP_DONE
	game.process()
}

func (set *Set) process() {
	switch set.status {
	case SET_NEED_INIT:
		set.init()
		break
	case SET_WAIT_OP:
		set.rule.doEvent(set)
		if set.curPlayer().doEvent != nil {
			set.status = SET_OP_DONE
			set.process()
		}
		break
	case SET_OP_DONE:
		set.rule.nextEvent(set)
		set.curPlayer().doEvent = nil
		set.status = SET_WAIT_OP
		set.process()
		break
	default:
		break
	}
}

func (set *Set) init() {
	// 初始化牌局
	set.rule.initPool(set)
	set.curPlayerIndex = 0
	set.curPlayer().zhuang = true
	set.curEvent = set.rule.getNextEvent(set.curPlayer())
	set.status = SET_WAIT_OP
}

func (set *Set) curPlayer() *Player {
	return set.playerChain[set.curPlayerIndex]
}
