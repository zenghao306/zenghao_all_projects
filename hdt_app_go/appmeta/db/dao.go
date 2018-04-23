package db

import (
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
	. "hdt_app_go/appmeta/conf"
)

type Dao struct {
	redisCli *redis.Client
	mysqlCli *xorm.Engine
}

func NewDao() *Dao {
	addr := Cfg.MustValue("redis", "addr")
	pwd := Cfg.MustValue("redis", "pwd")

	server := Cfg.MustValue("db", "server")
	username := Cfg.MustValue("db", "username")
	password := Cfg.MustValue("db", "password")
	dbName := Cfg.MustValue("db", "db_name")
	dbPort := Cfg.MustValue("db", "db_port")

	return &Dao{
		redisCli: NewRedis(addr, pwd),
		mysqlCli: NewMysql(server, username, password, dbName, dbPort),
	}
}
