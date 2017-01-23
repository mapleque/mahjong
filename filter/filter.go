package filter

import (
	. "github.com/coral"
	. "github.com/mahjong/constant"
	"github.com/mahjong/mahjong"
	"strconv"
)

func InitRouter(server *Server) {
	baseRouter := server.NewRouter("/", DefaultFilter)
	baseRouter.NewDocRouter(&Doc{
		Path:        "info",
		Description: "get info of now",
		Input: Checker{
			"token": Rule("string", STATUS_INVALID_TOKEN, "user token")}},
		infoFilter)
	baseRouter.NewDocRouter(&Doc{
		Path:        "create",
		Description: "create a game",
		Input: Checker{
			"token": Rule("string", STATUS_INVALID_TOKEN, "user token")}},
		createFilter)
	baseRouter.NewDocRouter(&Doc{
		Path:        "join",
		Description: "join a game",
		Input: Checker{
			"token": Rule("string", STATUS_INVALID_TOKEN, "user token"),
			"fsmId": Rule("string", STATUS_INVALID_ID, "fsm id")}},
		joinFilter)
	baseRouter.NewDocRouter(&Doc{
		Path:        "op",
		Description: "do something",
		Input: Checker{
			"token":  Rule("string", STATUS_INVALID_TOKEN, "user token"),
			"op":     Rule("int{2,3,5,6,7,12}", STATUS_INVALID_OP, "user op"),
			"indexs": []string{Rule("int", STATUS_INVALID_OP, "op indexs")}}},
		opFilter)
	baseRouter.NewDocRouter(&Doc{
		Path:        "next",
		Description: "start next game",
		Input: Checker{
			"token": Rule("string", STATUS_INVALID_TOKEN, "user token")}},
		nextFilter)
}

func DefaultFilter(context *Context) bool {
	context.Raw = true
	context.Data = `<!doctype html>
<meta charset='utf-8'>
<title>mahjong</title>
<h1>mahjong</h1>
<p>api doc <a href='/doc'>@see</a></p>`
	return true
}

func getPlayerId(context *Context) string {
	token := String(context.Params["token"])
	return token
}

var gameIdInc = 1

// fsmId -> fsm
var gameList = map[string]*mahjong.FSM{}

// playerId -> fsmId
var gameMap = map[string]string{}

func infoFilter(context *Context) bool {
	playerId := getPlayerId(context)
	fsm := gameList[gameMap[playerId]]
	ret := map[string]interface{}{
		"gameList": gameList,
		"gameInfo": fsm}
	context.Data = ret
	return true
}
func createFilter(context *Context) bool {
	playerId := getPlayerId(context)
	if _, ok := gameMap[playerId]; ok {
		return false
	}
	fsmId := strconv.Itoa(gameIdInc)
	gameIdInc++
	fsm := mahjong.CreateFSM("standard", 2)
	fsm.Join(playerId)
	gameList[fsmId] = fsm
	gameMap[playerId] = fsmId
	return true
}
func joinFilter(context *Context) bool {
	playerId := getPlayerId(context)
	fsmId := String(context.Params["fsmId"])
	if _, ok := gameMap[playerId]; ok {
		return false
	}
	fsm, ok := gameList[fsmId]
	if !ok {
		return false
	}
	fsm.Join(playerId)
	gameMap[playerId] = fsmId
	return true
}
func opFilter(context *Context) bool {
	playerId := getPlayerId(context)
	op := Int(context.Params["op"])
	indexs := Array(context.Params["indexs"])
	var opIndexs []int
	for _, index := range indexs {
		opIndexs = append(opIndexs, Int(index))
	}

	fsmId, ok := gameMap[playerId]
	if !ok {
		return false
	}
	fsm, ok := gameList[fsmId]
	if !ok {
		return false
	}
	return fsm.Op(playerId, op, opIndexs)
}
func nextFilter(context *Context) bool {
	playerId := getPlayerId(context)
	fsmId, ok := gameMap[playerId]
	if !ok {
		return false
	}
	fsm, ok := gameList[fsmId]
	if !ok {
		return false
	}
	if fsm.SetStatus != 3 {
		return false
	}
	fsm.Next()
	return true
}
