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
			"data": Checker{
				"fsmId": Rule("string", STATUS_INVALID_ID, "fsm id")}}},
		joinFilter)
	baseRouter.NewDocRouter(&Doc{
		Path:        "op",
		Description: "do something",
		Input: Checker{
			"token": Rule("string", STATUS_INVALID_TOKEN, "user token"),
			"data": Checker{
				"op": Rule("int{2,3,5,6,7,12}", STATUS_INVALID_OP, "user op"),
				"indexs": []string{
					Rule("int", STATUS_INVALID_OP, "op indexs")}}}},
		opFilter)
	baseRouter.NewDocRouter(&Doc{
		Path:        "next",
		Description: "start a new game",
		Input: Checker{
			"token": Rule("string", STATUS_INVALID_TOKEN, "user token")}},
		nextFilter)
	baseRouter.NewDocRouter(&Doc{
		Path:        "leave",
		Description: "leave a game",
		Input: Checker{
			"token": Rule("string", STATUS_INVALID_TOKEN, "user token")}},
		leaveFilter)
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
	params := Map(context.Params["data"])
	fsmId := String(params["fsmId"])
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
	params := Map(context.Params["data"])
	op := Int(params["op"])
	indexs := Array(params["indexs"])
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
	fsm.Next()
	return true
}
func leaveFilter(context *Context) bool {
	playerId := getPlayerId(context)
	fsmId, ok := gameMap[playerId]
	if !ok {
		return false
	}
	toDel := []string{}
	for oPlayer, oFsmId := range gameMap {
		if oFsmId == fsmId {
			toDel = append(toDel, oPlayer)
		}
	}
	_, ok = gameList[fsmId]
	if !ok {
		return false
	}
	delete(gameList, fsmId)
	for _, oPlayer := range toDel {
		delete(gameMap, oPlayer)
	}
	return true
}
