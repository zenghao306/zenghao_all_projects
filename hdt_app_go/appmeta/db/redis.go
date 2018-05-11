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

func (d *Dao) DelUserToken(tel string) {
	err := d.redisCli.Del(_keyToken+tel).Err()
	if err == redis.Nil {
		return
	} else if err != nil {
		Log.Err(err.Error())
	}
}

func (d *Dao) SetUserToken(tel string, token string) int32 {
	d.DelUserToken(tel)
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

type MinePoolTask struct {
	AppName string
	AppIcoPath string
	AppId int64
	Time int64
	Hdt float64
	HdtTaskBalance float64
}

func (d *Dao) GetMinePoolTaskDigInfo()( lis []*proto.MinePoolTaskListRes_MinePoolTask) {
	exit := d.redisCli.Exists(KeyHourRankingOfHdtDig).Val()
	if exit == 0 { //如果缓存数据不存在则重新获取并存储到redis里
		d.SetMinePoolTaskDigInfo()
	}

	res, err := d.redisCli.LRange(KeyHourRankingOfHdtDig, 0, -1).Result() //获取缓存数据
	if err != nil {
		Log.Err(err.Error())
		return nil
	}

	lists := make([]*proto.MinePoolTaskListRes_MinePoolTask, 0)
	for _, v := range res { //依次展示
		dat := &MinePoolTask{}
		err := json.Unmarshal([]byte(v), &dat)
		if err != nil {
			Log.Err(err.Error())
			continue
		}
		temp := &proto.MinePoolTaskListRes_MinePoolTask{}
		temp.HdtTaskBalance = dat.HdtTaskBalance
		temp.Hdt = dat.Hdt
		temp.Time = dat.Time
		temp.AppId = dat.AppId
		temp.AppIcoPath = dat.AppIcoPath
		temp.AppName = dat.AppName
		temp.Style = 2
		lists = append(lists, temp)
	}

	return lists
}

//获取挖矿排行榜数据并缓存到redis
func (d *Dao) SetMinePoolTaskDigInfo() {
		_, list := d.GetMinePoolTaskDigList()
		err := d.redisCli.Del(KeyHourRankingOfHdtDig).Err()
		if err != nil {
			Log.Err(err.Error())
			return
		}

		for _, v := range list {
			var data MinePoolTask
			data.AppName = v.AppName
			data.AppIcoPath = v.AppIcoPath
			data.AppId = v.AppId
			data.Time = v.Time
			data.Hdt = v.Hdt
			data.HdtTaskBalance = v.HdtTaskBalance
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

			//以下一段代码是设置KeyHourRankingOfHdtDig的有效时间【下一个整点的第一秒】
			t := time.Now()
			t1:=time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 1, 0, time.Local)
			after := 3600 - (t.Unix() - t1.Unix())
			tm := time.Unix(time.Now().Unix() + after, 0)
			d.redisCli.ExpireAt(KeyHourRankingOfHdtDig,tm)
		}
}
