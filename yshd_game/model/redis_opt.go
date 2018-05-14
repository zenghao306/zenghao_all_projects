package model

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	"github.com/yshd_game/confdata"
	"time"
	//"sort"
	"errors"

	"strconv"
	"sync"
	//"github.com/liudng/godump"
	//"github.com/liudng/godump"
)

const (
	ERR_REDIS_STR = "err over bet num"
)

var redisCli *redis.Client

var (
	_keyGameID = "mpn_"
	//_keyRouID  = "rou_"
	//_keyCachePicTable="pic_user"
	_keyCachePic   = "pic_"
	_keyCahcehUser = "user_"
	_keyCacheRoom  = "audience_"
	_keyCacheSns   = "qxsns_"
	_keyCacheUV    = "uvinfo_set"
)

type MBetInfo struct {
	Uid    int
	Num    int
	GameId string
}

var MBet map[int]MBetInfo
var MBetGuard sync.RWMutex

func init() {
	MBet = make(map[int]MBetInfo, 0)
}

/*
func keyUser(uid int) string {
	return fmt.Sprintf("key:%d", uid)
}
*/
func keyGameID(gameID string) string {
	return _keyGameID + gameID
}

func keyRoomAudience(rid string) string {
	return _keyCacheRoom + rid
}

/*
func keyUserRoundID(userID int) string {
	nowtime := time.Now()
	date := nowtime.Format("20060102")
	return fmt.Sprintf("%s%d_%s", _keyRouID, userID, date)
}
*/
func keyCachePic(userID int) string {
	return fmt.Sprintf("%s%d", _keyCachePic, userID)
}

func keyCacheUser(uid int) string {
	return fmt.Sprintf("%s%d", _keyCahcehUser, uid)
}

func keyCacheSns(tel string) string {
	return fmt.Sprintf("%s%s", _keyCacheSns, tel)
}

func keyCacheUv() string {
	return fmt.Sprintf("%s", _keyCacheUV)
}

func keyCacheSendGift(date_type string) string {
	switch date_type {
	case "week":
		return fmt.Sprintf("%s", "week_send_gift_rank")
	case "month":
		return fmt.Sprintf("%s", "month_send_gift_rank")
	default:
		return ""
	}
}

func keyCacheRevGift(date_type string) string {
	switch date_type {
	case "week":
		return fmt.Sprintf("%s", "week_rev_gift_rank")
	case "month":
		return fmt.Sprintf("%s", "month_rev_gift_rank")
	case "all":
		return fmt.Sprintf("%s", "all_rev_gift_rank")
	default:
		return ""
	}
}

func keyCacheGameRank(date_type string) string {
	switch date_type {
	case "week":
		return fmt.Sprintf("%s", "game_winner_week_rank")
	case "month":
		return fmt.Sprintf("%s", "game_winner_month_rank")
	case "all":
		return fmt.Sprintf("%s", "game_winner_all_rank")
	default:
		return ""
	}
}

type BetInfo struct {
	Uid int
	Num int
	Pos int
}

type BetInfo2 struct {
	Uid int
	Num int
}

type UserBetRecord struct {
	Id     int64
	GameId string `xorm:"varchar(128)"`
	Uid    int
	Num    int
	Pos    int
}

/*
func JoinStr(gameId string, uid int) string {
	return fmt.Sprintf("%s_%d", gameId, uid)
}
*/
func InitRedis() {
	//return
	addr := common.Cfg.MustValue("redis", "redis_addr")

	pwd := common.Cfg.MustValue("redis", "redis_pwd")

	common.Log.Debugf("redis info %s", addr)

	redisCli = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd, // no password set
		DB:       0,   // use default DB
	})

	_, err := redisCli.Ping().Result()
	if err != nil {
		common.Log.Panicf("connent redis err is %v", err)
	}
}

//加注
func AddGameBetByID(gameid string, uid, num, pos int) error {
	/*
		var val BetInfo

		res,err:=redisCli.HGet(keyGameID(gameid),keyUser(uid)).Result()
		if err==redis.Nil {
			val.Uid = uid
			val.Pos = pos
			val.Num = num
		}else if err!=nil {
			common.Log.Errf(err.Error())
			return err
		}else{
			err = json.Unmarshal([]byte(res), &val)
			if err != nil {
				common.Log.Err(err)
				return err
			}

			if val.Num+num>100000 {
				var serr error = errors.New(ERR_REDIS_STR)
				return serr
			}
			val.Num+=num
		}
		rebyte, err := json.Marshal(val)
		if err != nil {
			common.Log.Errf(err.Error())
			return err
		}
		err=redisCli.HSet(keyGameID(gameid),keyUser(uid),rebyte).Err()
		if err != nil {
			common.Log.Errf(err.Error())
			return err
		}
	*/
	MBetGuard.Lock()
	defer MBetGuard.Unlock()
	u, ok := MBet[uid]
	if ok {
		if u.GameId == gameid {
			if u.Num+num > 100000 {
				var serr error = errors.New(ERR_REDIS_STR)
				return serr
			}
			u.Num += num
			MBet[uid] = u
		} else {
			var m MBetInfo
			m.Uid = uid
			m.Num = num
			m.GameId = gameid
			MBet[uid] = m
		}
	} else {
		var m MBetInfo
		m.Uid = uid
		m.Num = num
		m.GameId = gameid
		MBet[uid] = m
	}

	var val BetInfo
	val.Uid = uid
	val.Pos = pos
	val.Num = num
	rebyte, err := json.Marshal(&val)
	if err != nil {
		common.Log.Errf(err.Error())
		return err
	}

	//redisCli.HSet(keyGameID(gameid),keyUser(uid),rebyte)
	err = redisCli.LPush(keyGameID(gameid), rebyte).Err()

	if err != nil {
		common.Log.Errf("redis err is %s", err.Error())
		return err
	}

	if num >= 1000 {
		r := &UserBetRecord{
			Uid:    uid,
			Pos:    pos,
			Num:    num,
			GameId: gameid,
		}
		_, err = orm.InsertOne(r)
		if err != nil {
			common.Log.Errf("redis err is %s", err.Error())
			return err
		}
	}
	return nil
}

func GetBetTableByID(gameid string) (interface{}, error) {
	return redisCli.HGetAll(keyGameID(gameid)).Result()
}

//pos=-1退回

//每局结束调用，每个用户发奖励
//pos是位置（左中右）,pos为-1时候退回
//bet倍率
func ResultBetReward(gameid string, pos int, bet int, roomid string, ctime int64, anchorId int) map[int]int {
	uGainRecord := make(map[int]int)

	s, err := redisCli.LRange(keyGameID(gameid), 0, -1).Result() //遍历

	if err != nil {
		common.Log.Errf("redis err is %s", err.Error())
		return uGainRecord
	}

	/*
		var allBet int
		keys, _,err:=redisCli.HScan(keyGameID(gameid),0,"key:*",0).Result()
		if err != nil {
			godump.Dump(err)
			common.Log.Errf("redis err is %s",err.Error())
			return uGainRecord
		}
		godump.Dump(keys)
		for _,v:=range keys  {
			godump.Dump(v)
			res,err:=redisCli.HGet(keyGameID(gameid),v).Result()
			if err==redis.Nil {
				common.Log.Errf("redis is err %s",string(res))
				continue
			}else if err!=nil {
				common.Log.Errf("redis is err %s",err.Error())
				continue
			}else{
				dat := &BetInfo{}
				err := json.Unmarshal([]byte(res), &dat)
				if err != nil {
					common.Log.Err(err)
					continue
				}
				godump.Dump(dat)
				user, ok := GetUserByUid(dat.Uid)
				if !ok {
					continue
				}
				if  dat.Pos==pos{



					allBet+=dat.Num
					godump.Dump(allBet)
					uGainRecord[dat.Uid] += dat.Num
				} else if dat.Pos == -1 {
					user.AddMoney(common.MONEY_TYPE_SCORE, int64(dat.Num))
				}else{
					uGainRecord[dat.Uid] += 0
				}

			}
		}


		for k, v := range uGainRecord {
			user, _ := GetUserByUid(k) //根据用户ID获取user对象

			dump_num := int(float32(v*bet) - float32(v) * 0.1)
			UserWinnerScore(k, int(user.Score), dump_num)

			user.AddMoney(common.MONEY_TYPE_SCORE, int64(dump_num))
			if v > 0 {
				TriggerTask(confdata.TargetType_win, k, 1)
			}
		}

		RecordResult(gameid, pos, roomid, ctime, allBet, anchorId)

		return uGainRecord
	*/
	var allBet int
	for _, v := range s {

		dat := &BetInfo{}
		err := json.Unmarshal([]byte(v), &dat)
		if err != nil {
			common.Log.Err(err)
			continue
		}
		user, ret := GetUserByUid(int(dat.Uid))
		if ret != common.ERR_SUCCESS {
			continue
		}
		MBetGuard.Lock()

		_, ok := MBet[dat.Uid]
		if ok {
			delete(MBet, dat.Uid)
		}
		MBetGuard.Unlock()
		if dat.Pos == pos {
			//user.AddMoney(common.MONEY_TYPE_SCORE, int64(dat.Num*bet))
			//user.AddScore(int64(dat.Num * bet))
			uGainRecord[dat.Uid] += int(float32(dat.Num*bet) - float32(dat.Num)*0.1) //dat.Num
			//uGainRecord[dat.Uid]+=dat.Num
			//ms = append(ms, BetInfo2{dat.Uid, dat.Num}) //追加到排序slice里
		} else if dat.Pos == -1 {
			user.AddMoney(nil, common.MONEY_TYPE_SCORE, int64(dat.Num), false)
			//user.AddScore(int64(dat.Num))
		} else if dat.Pos >= 0 && pos >= 0 && dat.Pos != pos {
			uGainRecord[dat.Uid] += 0
		}
		allBet += dat.Num
	}

	for k, v := range uGainRecord {
		user, _ := GetUserByUid(k) //根据用户ID获取user对象

		//dump_num := int(float32(v*bet) - float32(v) * 0.1)
		//dump_num := 0
		user.AddMoney(nil, common.MONEY_TYPE_SCORE, int64(v), false)
		//user.AddMoney(common.MONEY_TYPE_SCORE, int64(dump_num))
		//////////////
		sess := GetUserSessByUid(user.Uid)
		if sess != nil && sess.Roomid == roomid {
			UserWinnerScore(k, int(user.Score), v)
		}
		/////////
		//UserWinnerScore(k, int(user.Score), v)

		//[记录用户每局实际赢取的游戏币] added by zenghao 2017-09-13
		UserWinScoreRecord(gameid, k, v)
		//added end

		if v > 0 {
			TriggerTask(confdata.TargetType_win, k, 1)
		}

	}

	ret := RecordResult(gameid, pos, roomid, ctime, allBet, anchorId)
	if ret != common.ERR_SUCCESS {
		common.Log.Errf("record result is err %d", ret)
	}
	return uGainRecord
}

//func getUserNumberOfWins(uid int) int {
//	key := keyUserRoundID(uid)
//	value, err := redisCli.Get(key).Int64()
//	if err != nil {
//		common.Log.Errf("redis error:", err.Error())
//		return -1
//	}
//	return int(value)
//}

//每局结束调用，每个用户发奖励
//pos是位置（左中右）,pos为-1时候退回
//bet倍率
func GetUserRaiseInfo(gameid string, uid int) (int, int, int) {
	posAScore := 0 //第一家押的分数
	posBScore := 0 //第二家押的分数
	posCScore := 0 //第三家押的分数

	s, err := redisCli.LRange(keyGameID(gameid), 0, -1).Result() //遍历
	if err != nil {
		common.Log.Err(err)
		return posAScore, posBScore, posCScore
	}

	for _, v := range s {

		dat := &BetInfo{}
		err := json.Unmarshal([]byte(v), &dat)
		if err != nil {
			common.Log.Err(err)
			continue
		}

		//根据位置值取出相应的分数并赋值给对应的变量
		if dat.Uid == uid {
			if dat.Pos == 0 {
				posAScore += dat.Num
			} else if dat.Pos == 1 {
				posBScore += dat.Num
			} else if dat.Pos == 2 {
				posCScore += dat.Num
			}
		}
	}

	return posAScore, posBScore, posCScore
}

// 根据游戏ID删除redis里押分的数据
// gameID，游戏ID
func DelUserRaiseRedisInfo(gameID string) {
	redisCli.Del(keyGameID(gameID))
}

//上传图片设置回调的缓存图片
func SetCachePic(uid int, pic string, pic_type string) {
	err := redisCli.HSet(keyCachePic(uid), pic_type, pic).Err()
	if err != nil {
		common.Log.Err(err)
	}
}

//获取缓存图片
func GetCachePic(uid int, pic_type string) string {
	res, err := redisCli.HGet(keyCachePic(uid), pic_type).Result()
	if err == redis.Nil {
		common.Log.Err("key2 does not exists")
		return ""
	} else if err != nil {
		common.Log.Err(err)
		return ""
	}
	return res
}

//清空缓存图片
func ClearCachePic(uid int, pic_type string) {
	err := redisCli.HDel(keyCachePic(uid), pic_type).Err()
	if err != nil {
		common.Log.Err(err)
	}
}

type CacheUser struct {
	Uid       int
	PreRoomId string
	Status    int    //是否在直播
	RoomId    string //所在房间ID
	WatchId   int    //当前记录的观看记录ID

}

func SetCacheUser(uid int, u *CacheUser) (err error) {
	b, err := json.Marshal(u)
	if err != nil {
		common.Log.Err(err)
		return
	}

	err = redisCli.Set(keyCacheUser(uid), string(b), 360*time.Second).Err()

	if err != nil {
		common.Log.Err(err)
		return
	}
	return
}

func GetCacheUser(uid int) (u *CacheUser, err error) {
	b, err := redisCli.Get(keyCacheUser(uid)).Bytes()
	if err == redis.Nil {
		return
	} else if err != nil {
		common.Log.Err(err)
		return
	}

	err = json.Unmarshal(b, &u)
	if err != nil {
		common.Log.Err(err)
		return
	}
	return
}

//新增观众
func AddAudience(rid string, uid int, score int) (err error) {
	//f, err := strconv.ParseFloat(score, 32)
	z := redis.Z{
		Score:  float64(score),
		Member: strconv.Itoa(uid),
	}
	err = redisCli.ZAdd(keyRoomAudience(rid), z).Err()
	if err != nil {
		common.Log.Err(err)
		return
	}
	return
}

//给主播增加权重
func IncrAudience(rid string, uid int, incr_score int) (err error) {

	err = redisCli.ZIncrBy(keyRoomAudience(rid), float64(incr_score), strconv.Itoa(uid)).Err()
	if err != nil {
		common.Log.Err(err)
		return
	}
	return
}

//获取观众列表
func GetAudience(rid string, min, max int64) (users []int, err error) {
	res, err := redisCli.ZRevRange(keyRoomAudience(rid), min, max).Result()
	if err != nil {
		common.Log.Err(err)
		return
	}
	users = make([]int, 0)
	for _, v := range res {
		var uid int
		uid, err = strconv.Atoi(v)
		if err != nil {
			return
		}
		users = append(users, uid)
	}
	return
}

//删除指定观众
func DelAudience(rid string, uid int) (err error) {
	err = redisCli.ZRem(keyRoomAudience(rid), strconv.Itoa(uid)).Err()
	if err != nil {
		common.Log.Err(err)
		return
	}

	return
}

//删除观众列表
func DelAudienceKey(rid string) (err error) {
	err = redisCli.Del(keyRoomAudience(rid)).Err()
	if err != nil {
		common.Log.Err(err)
		return
	}

	return
}

func AddQianXunCode(tel string, code string) {
	err := redisCli.Set(keyCacheSns(tel), code, 300*time.Second).Err()
	if err != nil {
		common.Log.Err(err)
		return
	}
	return
}

func GetQianXunCode(tel string) (ret int, code string) {
	code, err := redisCli.Get(keyCacheSns(tel)).Result()
	if err == redis.Nil {
		ret = common.ERR_SNS_TIMEOUT
		return
	} else if err != nil {
		ret = common.ERR_UNKNOWN
	} else {
		ret = common.ERR_SUCCESS
	}
	return
}
func DelQianXunCode(tel string) {
	err := redisCli.Del(keyCacheSns(tel)).Err()
	if err == redis.Nil {
		return
	} else if err != nil {
		common.Log.Err(err.Error())
	}
}

func GetUvInfo(uv string, device string, channel_id string) {
	dat := &UvInfo{
		Imei:       uv,
		Device:     device,
		ChannelId:  channel_id,
		CreateTime: time.Now().Unix(),
	}
	b, err := json.Marshal(dat)
	if err != nil {
		common.Log.Err(err)
		return
	}
	err = redisCli.HSet(keyCacheUv(), uv, b).Err()
	if err != nil {
		common.Log.Err(err)
		return
	}
}

func CheckUserBetRewardFromRedis(gameId string, uid, pos int) (int, int, int) {
	userGain := 0
	userNotGainRaise := 0
	userTotalRaise := 0

	s, err := redisCli.LRange(keyGameID(gameId), 0, -1).Result() //遍历
	if err != nil {
		common.Log.Errf("redis err is %s", err.Error())
		return userGain, userNotGainRaise, userTotalRaise
	}

	for _, v := range s {
		dat := &BetInfo{}
		err := json.Unmarshal([]byte(v), &dat)
		if err != nil {
			common.Log.Err(err)
			continue
		}

		if dat.Uid == uid {
			if dat.Pos == pos {
				userGain += int(float32(dat.Num*common.GAME_BONUS_TIMES) - float32(dat.Num)*0.1) //dat.Num
			} else {
				userNotGainRaise += dat.Num
			}
			userTotalRaise += dat.Num
		}
	}

	return userGain, userNotGainRaise, userTotalRaise
}

type UserRaisePosInfo struct {
	pos0 int
	pos1 int
	pos2 int
}

func GetUserRaisePosInfoWithGameId(gameId string) map[int]*UserRaisePosInfo {
	mRaise := make(map[int]*UserRaisePosInfo)

	s, err := redisCli.LRange(keyGameID(gameId), 0, -1).Result() //遍历
	if err != nil {
		common.Log.Errf("redis err is %s", err.Error())
		return mRaise
	}

	for _, v := range s {
		dat := &BetInfo{}
		err := json.Unmarshal([]byte(v), &dat)
		if err != nil {
			common.Log.Err(err)
			continue
		}
		_, ok := mRaise[dat.Uid]
		if !ok {
			m := &UserRaisePosInfo{}
			mRaise[dat.Uid] = m
		}
		c2, _ := mRaise[dat.Uid]

		if dat.Pos == 0 {
			c2.pos0 = dat.Num
		} else if dat.Pos == 1 {
			c2.pos1 = dat.Num
		} else if dat.Pos == 2 {
			c2.pos2 = dat.Num
		}
	}

	return mRaise
}

func SetWeekSendRank() {
	t, err := redisCli.Get("week_send_update_time").Int64()
	if err == redis.Nil {

	} else if err != nil {
		common.Log.Err(err.Error())
		return
	}
	now_time := time.Now().Unix()
	if now_time-t > int64(7*3600*24) {

		t := time.Now()

		end := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, time.Local)
		begin := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
		weekday := t.Weekday()
		if int(weekday) == 0 {
			weekday = 7
		}
		diff := int(weekday)
		diff2 := int(weekday) + 7

		end_time := end.AddDate(0, 0, -diff)

		begin_tm2 := begin.AddDate(0, 0, -diff2)

		s, _ := GetSendDiamonGiftRankWeek(begin_tm2.Unix(), end_time.Unix())
		err := redisCli.Del(keyCacheSendGift("week")).Err()
		if err != nil {
			common.Log.Err(err.Error())
			return
		}

		for _, v := range s {
			b, err := json.Marshal(v)
			if err != nil {
				common.Log.Err(err.Error())
				continue
			}
			err = redisCli.RPush(keyCacheSendGift("week"), b).Err()
			if err != nil {
				common.Log.Err(err.Error())
				continue
			}
		}

		err = redisCli.Set("week_send_update_time", end_time.Unix(), 0).Err()
		if err != nil {
			common.Log.Err(err.Error())
			return
		}
	}
}

func GetAllSendRank() []map[string]string {
	//sql := fmt.Sprintf("select a.send_user,a.count,b.nick_name,b.user_level as level,b.image,b.sex from (select send_user,sum(value) as count from go_gift_record  group by send_user) a left join go_user b on a.send_user=b.uid  where b.account_type!=1  order by a.count  desc limit %d,%d ", index*common.ROOM_LIST_PAGE_COUNT, common.ROOM_LIST_PAGE_COUNT)

	sql := fmt.Sprintf("select a.send_user,a.count,b.nick_name,b.user_level as level,b.image,b.sex from (select send_user,sum(value) as count from go_gift_record  where money_type=%d group by send_user) a left join go_user b on a.send_user=b.uid  where b.account_type!=1  order by a.count  desc limit 0,50 ", common.MONEY_TYPE_DIAMOND)
	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return nil
	}

	retMap := make([]map[string]string, 0)

	for _, row := range rowArray {
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			ss[colName] = value
		}
		retMap = append(retMap, ss)
	}

	return retMap

}

func GetWeekSendRank() []map[string]string {
	SetWeekSendRank()

	res, err := redisCli.LRange(keyCacheSendGift("week"), 0, 49).Result()
	if err != nil {
		common.Log.Err(err.Error())
		return nil
	}
	retMap := make([]map[string]string, 0)
	for _, v := range res {
		s := make(map[string]string)
		err := json.Unmarshal([]byte(v), &s)
		if err != nil {
			common.Log.Err(err.Error())
			continue
		}
		retMap = append(retMap, s)
	}
	return retMap
}

func GetMonthSendRank() []map[string]string {

	SetMonthSendRank()
	res, err := redisCli.LRange(keyCacheSendGift("month"), 0, 49).Result()
	if err != nil {
		common.Log.Err(err.Error())
		return nil
	}
	retMap := make([]map[string]string, 0)
	for _, v := range res {
		s := make(map[string]string)
		err := json.Unmarshal([]byte(v), &s)
		if err != nil {
			common.Log.Err(err.Error())
			continue
		}
		retMap = append(retMap, s)
	}
	return retMap
}

func SetMonthSendRank() {
	t, err := redisCli.Get("month_send_update_time").Int64()
	if err == redis.Nil {
		t = 0
	} else if err != nil {
		common.Log.Err(err.Error())
		return
	}
	now_time := time.Now()
	m := now_time.Month()

	if int64(m) != t {
		t := time.Now()

		cur := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
		end_time := cur.Unix()

		begin_tm2 := cur.AddDate(0, -1, 0)

		s, _ := GetSendDiamonGiftRankWeek(begin_tm2.Unix(), end_time)
		err := redisCli.Del(keyCacheSendGift("month")).Err()
		if err != nil {
			common.Log.Err(err.Error())
			return
		}

		for _, v := range s {
			b, err := json.Marshal(v)
			if err != nil {
				common.Log.Err(err.Error())
				continue
			}
			err = redisCli.RPush(keyCacheSendGift("month"), b).Err()
			if err != nil {
				common.Log.Err(err.Error())
				continue
			}
		}

		err = redisCli.Set("month_send_update_time", int(m), 0).Err()
		if err != nil {
			common.Log.Err(err.Error())
			return
		}

	}
}

func SetWeekRevRank() {
	t, err := redisCli.Get("week_rev_update_time").Int64()
	if err == redis.Nil {

	} else if err != nil {
		common.Log.Err(err.Error())
		return
	}
	now_time := time.Now().Unix()
	if now_time-t > int64(7*3600*24) {

		t := time.Now()
		end := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, time.Local)
		begin := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
		weekday := t.Weekday()
		if int(weekday) == 0 {
			weekday = 7
		}
		diff := int(weekday)
		diff2 := int(weekday) + 7

		end_time := end.AddDate(0, 0, -diff)

		begin_tm2 := begin.AddDate(0, 0, -diff2)

		s, _ := GetCouponsRankListWithTime(begin_tm2.Unix(), end_time.Unix())
		err := redisCli.Del(keyCacheRevGift("week")).Err()
		if err != nil {
			//godump.Dump(err.Error())
			common.Log.Err(err.Error())
			return
		}

		for _, v := range s {

			b, err := json.Marshal(v)
			if err != nil {
				common.Log.Err(err.Error())
				continue
			}
			err = redisCli.RPush(keyCacheRevGift("week"), b).Err()
			if err != nil {
				common.Log.Err(err.Error())
				continue
			}
		}

		err = redisCli.Set("week_rev_update_time", end_time.Unix(), 0).Err()
		if err != nil {
			common.Log.Err(err.Error())
			return
		}
	}
}

func GetWeekRevRank() []map[string]string {
	SetWeekRevRank()

	res, err := redisCli.LRange(keyCacheRevGift("week"), 0, 49).Result()
	if err != nil {
		common.Log.Err(err.Error())
		return nil
	}
	retMap := make([]map[string]string, 0)
	for _, v := range res {
		s := make(map[string]string)
		err := json.Unmarshal([]byte(v), &s)
		if err != nil {
			common.Log.Err(err.Error())
			continue
		}
		retMap = append(retMap, s)
	}
	return retMap
}

func GetMonthRevRank() []map[string]string {
	SetMonthRevRank()
	res, err := redisCli.LRange(keyCacheRevGift("month"), 0, 49).Result()
	if err != nil {
		common.Log.Err(err.Error())
		return nil
	}
	retMap := make([]map[string]string, 0)
	for _, v := range res {
		s := make(map[string]string)
		err := json.Unmarshal([]byte(v), &s)
		if err != nil {
			common.Log.Err(err.Error())
			continue
		}
		retMap = append(retMap, s)
	}
	return retMap
}

func GetAllRevRank() []map[string]string {
	SetAllRevRank()

	res, err := redisCli.LRange(keyCacheRevGift("all"), 0, 49).Result()
	if err != nil {
		common.Log.Err(err.Error())
		return nil
	}
	retMap := make([]map[string]string, 0)
	for _, v := range res {
		s := make(map[string]string)
		err := json.Unmarshal([]byte(v), &s)
		if err != nil {
			common.Log.Err(err.Error())
			continue
		}
		retMap = append(retMap, s)
	}
	return retMap
}

func SetAllRevRank() {
	t, err := redisCli.Get("all_rev_update_time").Int64()
	if err == redis.Nil {

	} else if err != nil {
		common.Log.Err(err.Error())
		return
	}
	now_time := time.Now().Unix()
	if now_time-t > int64(3600) {
		s, _ := GetCouponsRankListWithTime(0, now_time)
		err := redisCli.Del(keyCacheRevGift("all")).Err()
		if err != nil {

			common.Log.Err(err.Error())
			return
		}

		for _, v := range s {
			//	godump.Dump(v)
			b, err := json.Marshal(v)
			if err != nil {
				common.Log.Err(err.Error())
				continue
			}
			err = redisCli.RPush(keyCacheRevGift("all"), b).Err()
			if err != nil {
				common.Log.Err(err.Error())
				continue
			}
		}

		err = redisCli.Set("all_rev_update_time", now_time, 0).Err()
		if err != nil {
			common.Log.Err(err.Error())
			return
		}
	}
}

func SetMonthRevRank() {
	t, err := redisCli.Get("month_rev_update_time").Int64()
	if err == redis.Nil {

	} else if err != nil {
		common.Log.Err(err.Error())
		return
	}
	now_time := time.Now()
	m := now_time.Month()
	if int64(m) != t {
		t := time.Now()

		cur := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
		end_time := cur.Unix()

		begin_tm2 := cur.AddDate(0, -1, 0)

		s, _ := GetCouponsRankListWithTime(begin_tm2.Unix(), end_time)
		err := redisCli.Del(keyCacheRevGift("month")).Err()
		if err != nil {
			common.Log.Err(err.Error())
			return
		}

		for _, v := range s {
			b, err := json.Marshal(v)
			if err != nil {
				common.Log.Err(err.Error())
				continue
			}
			err = redisCli.RPush(keyCacheRevGift("month"), b).Err()
			if err != nil {
				common.Log.Err(err.Error())
				continue
			}
		}

		err = redisCli.Set("month_rev_update_time", m, 0).Err()
		if err != nil {
			common.Log.Err(err.Error())
			return
		}
	}
}

func SetGameWeekRank() {
	t, err := redisCli.Get("game_rank_week_update_time").Int64()
	if err == redis.Nil {

	} else if err != nil {
		common.Log.Err(err.Error())
		return
	}
	now_time := time.Now().Unix()
	if now_time-t > int64(7*3600*24) {

		t := time.Now()
		end := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, time.Local)
		begin := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
		weekday := t.Weekday()
		if int(weekday) == 0 {
			weekday = 7
		}
		diff := int(weekday)
		diff2 := int(weekday) + 7

		end_time := end.AddDate(0, 0, -diff)
		begin_tm2 := begin.AddDate(0, 0, -diff2)

		s, _ := GetGameWinScoreRankWeek(begin_tm2.Unix(), end_time.Unix())
		err := redisCli.Del(keyCacheGameRank("week")).Err()
		if err != nil {
			common.Log.Err(err.Error())
			return
		}

		for _, v := range s {
			b, err := json.Marshal(v)
			if err != nil {
				common.Log.Err(err.Error())
				continue
			}
			err = redisCli.RPush(keyCacheGameRank("week"), b).Err()
			if err != nil {
				common.Log.Err(err.Error())
				continue
			}
		}

		err = redisCli.Set("game_rank_week_update_time", end_time.Unix(), 0).Err()
		if err != nil {
			common.Log.Err(err.Error())
			return
		}
	}
}

func GetGameWeekRank() []map[string]string {
	SetGameWeekRank()

	res, err := redisCli.LRange(keyCacheGameRank("week"), 0, 49).Result()
	if err != nil {
		common.Log.Err(err.Error())
		return nil
	}
	retMap := make([]map[string]string, 0)
	for _, v := range res {
		s := make(map[string]string)
		err := json.Unmarshal([]byte(v), &s)
		if err != nil {
			common.Log.Err(err.Error())
			continue
		}
		retMap = append(retMap, s)
	}
	return retMap
}

func GetMonthRankOfGame() []map[string]string {
	SetMonthRankOfGame()
	res, err := redisCli.LRange(keyCacheGameRank("month"), 0, 49).Result()
	if err != nil {
		common.Log.Err(err.Error())
		return nil
	}
	retMap := make([]map[string]string, 0)
	for _, v := range res {
		s := make(map[string]string)
		err := json.Unmarshal([]byte(v), &s)
		if err != nil {
			common.Log.Err(err.Error())
			continue
		}
		retMap = append(retMap, s)
	}
	return retMap
}

func SetMonthRankOfGame() {
	t, err := redisCli.Get("game_rank_month_update_time").Int64()
	if err == redis.Nil {

	} else if err != nil {
		common.Log.Err(err.Error())
		return
	}
	now_time := time.Now()
	m := now_time.Month()
	if int64(m) != t {
		t := time.Now()

		cur := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
		end_time := cur.Unix()

		begin_tm2 := cur.AddDate(0, -1, 0)

		s, _ := GetGameWinScoreRankWeek(begin_tm2.Unix(), end_time)
		err := redisCli.Del(keyCacheGameRank("month")).Err()
		if err != nil {
			common.Log.Err(err.Error())
			return
		}

		for _, v := range s {
			b, err := json.Marshal(v)
			if err != nil {
				common.Log.Err(err.Error())
				continue
			}
			err = redisCli.RPush(keyCacheGameRank("month"), b).Err()
			if err != nil {
				common.Log.Err(err.Error())
				continue
			}
		}

		err = redisCli.Set("game_rank_month_update_time", m, 0).Err()
		if err != nil {
			common.Log.Err(err.Error())
			return
		}
	}
}

func GetAllRankOfGame() []map[string]string {
	SetAllRankOfGame()

	res, err := redisCli.LRange(keyCacheGameRank("all"), 0, 49).Result()
	if err != nil {
		common.Log.Err(err.Error())
		return nil
	}
	retMap := make([]map[string]string, 0)
	for _, v := range res {
		s := make(map[string]string)
		err := json.Unmarshal([]byte(v), &s)
		if err != nil {
			common.Log.Err(err.Error())
			continue
		}
		retMap = append(retMap, s)
	}
	return retMap
}

func SetAllRankOfGame() {
	t, err := redisCli.Get("all_game_rank_winner_update_time").Int64()
	if err == redis.Nil {

	} else if err != nil {
		common.Log.Err(err.Error())
		return
	}
	now_time := time.Now().Unix()
	if now_time-t > int64(3600) {
		s, _ := GetGameWinScoreRankWeek(0, now_time)
		err := redisCli.Del(keyCacheGameRank("all")).Err()
		if err != nil {

			common.Log.Err(err.Error())
			return
		}

		for _, v := range s {
			b, err := json.Marshal(v)
			if err != nil {
				common.Log.Err(err.Error())
				continue
			}
			err = redisCli.RPush(keyCacheGameRank("all"), b).Err()
			if err != nil {
				common.Log.Err(err.Error())
				continue
			}
		}

		err = redisCli.Set("all_game_rank_winner_update_time", now_time, 0).Err()
		if err != nil {
			common.Log.Err(err.Error())
			return
		}
	}
}

func ResetRank() int {
	err := redisCli.Del("month_rev_update_time").Err()
	if err == redis.Nil {

	} else if err != nil {
		common.Log.Err(err.Error())
		return common.ERR_UNKNOWN
	}
	err = redisCli.Del("all_rev_update_time").Err()
	if err == redis.Nil {

	} else if err != nil {
		common.Log.Err(err.Error())
		return common.ERR_UNKNOWN
	}
	err = redisCli.Del("week_rev_update_time").Err()
	if err == redis.Nil {

	} else if err != nil {
		common.Log.Err(err.Error())
		return common.ERR_UNKNOWN
	}
	err = redisCli.Del("month_send_update_time").Err()
	if err == redis.Nil {

	} else if err != nil {
		common.Log.Err(err.Error())
		return common.ERR_UNKNOWN
	}
	err = redisCli.Del("week_send_update_time").Err()
	if err == redis.Nil {

	} else if err != nil {
		common.Log.Err(err.Error())
		return common.ERR_UNKNOWN
	}
	return common.ERR_SUCCESS
}
