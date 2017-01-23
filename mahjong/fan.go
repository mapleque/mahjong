package mahjong

type WinInfo struct {
	Fans          []Fan    `json:"fans"`
	Count         int      `json:"count"`
	WinPlayerId   string   `json:"winPlayerId"`
	LosePlayerIds []string `json:"losePlayerIds"`
}

type Fan struct {
	Index int    // 序号
	Name  string // 名称
	Count int    // 番数

	without []int                   // 不记番序号
	check   func(*FSM, string) bool // 校验函数
}
