package db

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/liudng/godump"
	. "hdt_app_go/appmeta/log"
	proto "hdt_app_go/protcol"
	"time"
)

func NewRedis(addr string, pwd string) *redis.Client {
	options := &redis.Options{
		Addr:     addr,
		Password: pwd,
	}
	c := redis.NewClient(options)
	s := c.Ping()
	if err := s.Err(); err != nil {
		Log.Fatal(err)
	}
	return c
}

var (
	CircleTime                       = "global_circle_time"
	_keyCacheSns                     = "qxsns_"
	_keyToken                        = "token_"
	KeyHourRankingOfHdtDig           = "K_HourRankingOfHdtDig"
	KeyHourRankingOfHdtDigUpdateTime = "K_HourRankingUpdateTime"
)

func (d *Dao) GetIpHashKey(ip string) string {
	/*
		h:=sha1.New()
		h.Write([]byte(ip))
		bs:=h.Sum(nil)
	*/
	godump.Dump(ip)
	return fmt.Sprintf("iplimit:%s", ip)
}

/*
//获取用户当前周期key
func (d *Dao) getUserActionLimitKey(uid uint32, action string) string {
	return fmt.Sprintf("%d:%s:%s", uid, d.getCircleTime(), action)
}

//更新全局当前周期
func (d *Dao) UpdateCircleTime(circle string) {
	d.RedisCli.Set(CircleTime, circle, 7200*time.Second)
}

func (d *Dao) getCircleTime() string {
	return d.RedisCli.Get(CircleTime).String()
}

func (d *Dao) SetUserActionLimit(uid uint32, action string) {
	key := d.getUserActionLimitKey(uid, action)
	d.RedisCli.Incr(key)
	d.RedisCli.Expire(key, 3600*time.Second)
}

func (d *Dao) GetUserActionLimit(uid uint32, action string) (count int, err error) {
	key := d.getUserActionLimitKey(uid, action)
	counts, err := d.RedisCli.Get(key).Int64()
	if err != nil {
		Log.Err(err)
		return
	}
	count = int(counts)
	return
}
*/

func keyCacheSns(tel string) string {
	return fmt.Sprintf("%s%s", _keyCacheSns, tel)
}

func keyToken(tel string) string {
	return fmt.Sprintf("%s%s", _keyToken, tel)
}

func (d *Dao) AddQianXunCode(tel string, code string) int32 {
	err := d.redisCli.Set(keyCacheSns(tel), code, 300*time.Second).Err()
	if err != nil {
		Log.Err(err.Error())
		return proto.ERR_UNKNOWN
	}
	return proto.ERR_OK
}

func (d *Dao) GetQianXunCode(tel string) (ret int32, code string) {
	code, err := d.redisCli.Get(keyCacheSns(tel)).Result()
	if err == redis.Nil {
		ret = proto.ERR_SNS_TIMEOUT
		return
	} else if err != nil {
		Log.Err(err.Error())
		ret = proto.ERR_UNKNOWN
	} else {
		ret = proto.ERR_OK
	}
	return
}

func (d *Dao) DelQianXunCode(tel string) {
	err := d.redisCli.Del(keyCacheSns(tel)).Err()
	if err == redis.Nil {
		return
	} else if err != nil {
		Log.Err(err.Error())
	}
}

func (d *Dao) QianXunSnsVerify(tel string, snsText string) int32 {
	result, code := d.GetQianXunCode(tel)
	if result != proto.ERR_OK {
		return result
	}

	if code == snsText {
		d.DelQianXunCode(tel)
		return proto.ERR_OK
	}
	return proto.ERR_SNS_CORRECT
}

func (d *Dao) SetUserToken(tel string, token string) int32 {
	err := d.redisCli.Set(_keyToken+tel, token, 86400*time.Second).Err()
	if err != nil {
		Log.Err(err.Error())
		return proto.ERR_UNKNOWN
	}

	return proto.ERR_OK
}

func (d *Dao) GetUserToken(tel string) (int32, string) {
	val, err := d.redisCli.Get(keyToken(tel)).Result()
	if err != nil {
		Log.Err(err.Error())
		return proto.ERR_UNKNOWN, val
	}

	return proto.ERR_OK, val
}

//func (d *Dao) RedisSaveRankingList(m map[string]float64,endTime int64) error{
//	d.redisCli.HDel(KeyHdtRankingList) //Del(KeyHdtRankingList)
//
//	//err = d.RedisCli.LPush(KeyHdtRankingList, rebyte).Err()
//	for k, v := range m {
//
//	}
//	return nil
//}
/*
func (d *Dao) GetHourRankingOfHdtDig() []map[string]string {
	d.SetHourRankingOfHdtDig()
	res, err := d.redisCli.LRange(KeyHourRankingOfHdtDig, 0, 49).Result()
	if err != nil {
		Log.Err(err.Error())
		return nil
	}
	retMap := make([]map[string]string, 0)
	for _, v := range res {
		s := make(map[string]string)
		err := json.Unmarshal([]byte(v), &s)
		if err != nil {
			Log.Err(err.Error())
			continue
		}
		retMap = append(retMap, s)
	}
	return retMap
}

func (d *Dao)SetHourRankingOfHdtDig() {
	t, err := d.redisCli.Get(KeyHourRankingOfHdtDigUpdateTime).Int64()
	if err == redis.Nil {

	} else if err != nil {
		Log.Err(err.Error())
		return
	}
	nowTime := time.Now()
	m := nowTime.Month()

	lbt := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), nowTime.Hour(), 0, 0, 0, time.Local)
	LastBalanceTime := lbt.Unix() //上个结算周期时间戳

	if t < LastBalanceTime { //如果记录的时间小于后台上个结算周期时间戳，则表示需要更新排行榜了
		s, _ := d.GetMySqlRankingList(LastBalanceTime)
		err := d.redisCli.Del(KeyHourRankingOfHdtDig).Err()
		if err != nil {
			Log.Err(err.Error())
			return
		}

		for _, v := range s {
			b, err := json.Marshal(v)
			if err != nil {
				Log.Err(err.Error())
				continue
			}
			err = d.redisCli.RPush(KeyHourRankingOfHdtDig, b).Err()
			if err != nil {
				Log.Err(err.Error())
				continue
			}
		}

		err = d.redisCli.Set(KeyHourRankingOfHdtDigUpdateTime, m, 0).Err()
		if err != nil {
			Log.Err(err.Error())
			return
		}
	}
}
*/

type UserHdtInfo struct {
	Tel string
	Hdt string
}

func (d *Dao) GetHourRankingOfHdtDig() map[string]string {
	d.SetHourRankingOfHdtDig()
	res, err := d.redisCli.LRange(KeyHourRankingOfHdtDig, 0, -1).Result()
	if err != nil {
		Log.Err(err.Error())
		return nil
	}
	//retMap := make([]map[string]string, 0)
	retMap := make(map[string]string, 0)
	for _, v := range res {
		dat := &UserHdtInfo{}
		err := json.Unmarshal([]byte(v), &dat)
		if err != nil {
			Log.Err(err.Error())
			continue
		}
		//s := make(map[string]string)
		//s[dat.Tel] = dat.Hdt
		//retMap = append(retMap, s)
		retMap[dat.Tel] = dat.Hdt
	}

	return retMap
}

func (d *Dao) SetHourRankingOfHdtDig() {
	t, err := d.redisCli.Get(KeyHourRankingOfHdtDigUpdateTime).Int64()
	if err == redis.Nil {

	} else if err != nil {
		Log.Err(err.Error())
		return
	}
	nowTime := time.Now()
	//m := nowTime.Month()

	lbt := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), nowTime.Hour(), 0, 0, 0, time.Local)
	LastBalanceTime := lbt.Unix() //上个结算周期时间戳

	if t < LastBalanceTime { //如果记录的时间小于后台上个结算周期时间戳，则表示需要更新排行榜了
		s, _ := d.GetMySqlRankingList(LastBalanceTime)
		err := d.redisCli.Del(KeyHourRankingOfHdtDig).Err()
		if err != nil {
			Log.Err(err.Error())
			return
		}

		for k, v := range s {
			var data UserHdtInfo
			data.Tel = k
			data.Hdt = v
			b, err := json.Marshal(data)
			if err != nil {
				Log.Err(err.Error())
				continue
			}

			err = d.redisCli.LPush(KeyHourRankingOfHdtDig, b).Err()
			if err != nil {
				Log.Err(err.Error())
				continue
			}
		}

		err = d.redisCli.Set(KeyHourRankingOfHdtDigUpdateTime, nowTime.Unix(), 0).Err()
		if err != nil {
			Log.Err(err.Error())
			return
		}
	}
}
