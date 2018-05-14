package model

import (
	"github.com/codegangsta/martini"
	"github.com/garyburd/redigo/redis"
	"github.com/yshd_game/common"
	//	"net/http"
	"time"
)

/*
func InitResid() {
	spec := redis.DefaultSpec().Host("192.168.1.142").Port(6379)

	//spec := redis.DefaultSpec().Host(common.Cfg.MustValue("", "redis_ip", "192.168.1.142")).Port(6379).Password("shang1234")
	client, err := redis.NewSynchClientWithSpec(spec)
	if err != nil {
		common.Log.Panicf("error on connect redis server is %s", err.Error())

	}
	client.Get(common.Cfg.MustValue("", "db_key"))
	if err != nil {
		common.Log.Panicf("error on Get db_key", err.Error())
	}
}
*/
var redigo *redis.Pool

func RedisInit() martini.Handler {
	proto := common.Cfg.MustValue("redis", "redis_proto")
	addr := common.Cfg.MustValue("redis", "redis_addr")
	pwd := common.Cfg.MustValue("redis", "redis_pwd")

	redigo = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 600 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(proto, addr, redis.DialPassword(pwd))
			if err != nil {
				panic(err)
				return nil, err
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	_, err := redigo.Get().Do("PING")
	if err != nil {
		common.Log.Panicln(err.Error())
	}
	println("Init Redis middleware successfully.")
	return redigo
	/*
		return func(res http.ResponseWriter, r *http.Request, c martini.Context) {
			c.MapTo(redigo.Get(), (*redis.Conn)(nil))
		}
	*/
}
