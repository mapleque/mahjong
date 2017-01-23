package main

import (
	"flag"

	"github.com/coral"
	"github.com/coral/cache"
	"github.com/coral/config"
	"github.com/coral/db"
	"github.com/coral/log"

	. "github.com/mahjong/constant"
	"github.com/mahjong/filter"
)

var conf config.Configer

func initDB() {
	// add default db
	db.AddDB(
		DEF_DB,
		conf.String("db.DB_DSN"),
		conf.Int("db.DB_MAX_CONNECTION"),
		conf.Int("db.DB_MAX_IDLE"))

	// add other db
	// ...
}

func initRedis() {
	// add default cache
	cache.AddRedis(
		DEF_REDIS,
		conf.String("cache.REDIS_SERVER"),
		conf.String("cache.REDIS_AUTH"),
		conf.Int("cache.REDIS_MAX_CONNECTION"),
		conf.Int("cache.REDIS_MAX_IDLE"))

	// add other cache
	// ...
}

func initLog() {
	// add default logger
	log.AddLogger(
		DEF_LOG,
		conf.String("log.LOG_PATH"),
		conf.Int("log.LOG_MAX_NUMBER"),
		conf.Int64("log.LOG_MAX_SIZE"),
		conf.Int("log.LOG_MAX_LEVEL"),
		conf.Int("log.LOG_MIN_LEVEL"))

	// add other logger
	// ...
}

func initRouter(server *coral.Server) {
	filter.InitRouter(server)
}

func main() {
	confFile := flag.String("ini", "./config/config.ini", "your config file")
	flag.Parse()
	if *confFile != "" {
		config.AddConfiger(config.INI, DEF_CONF, *confFile)
		conf = config.Use(DEF_CONF)

		// init log
		initLog()

		// init db
		initDB()

		// init redis
		initRedis()

		// new server
		server := coral.NewServer(conf.Get("server.HOST"))

		// init router
		initRouter(server)

		// start server
		server.Run()
	} else {
		panic("run with -h to find usage")
	}
}
