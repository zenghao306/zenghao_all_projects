package model

import (
	"github.com/yshd_game/common"
	"time"
)

type NiuNiuRecord struct {
	Id         int64
	GameId     string
	Pos        int
	CreateTime int64
	RoomId     string
	BetNum     int
	AnchorId   int
	MoneyNum   int64
}

type WinScoreRecord struct {
	Id         int64
	Uid        int
	GameId     string
	WinScore   int
	CreateTime int64
}

func RecordResult(gameid string, pos int, roomid string, ctime int64, betNum int, anchorId int) int {
	u, ret := GetUserByUid(anchorId)
	if ret != common.ERR_SUCCESS {
		return ret
	}

	m := &NiuNiuRecord{
		GameId:     gameid,
		Pos:        pos,
		CreateTime: ctime,
		RoomId:     roomid,
		BetNum:     betNum,
		AnchorId:   anchorId,
	}
	/*
		ret, _, _ := GetGroupId(u.AdminId)
		if ret != common.ERR_SUCCESS {
			_, err := orm.Insert(m)
			if err != nil {
				common.Log.Errf("orm err is %s", err.Error())
				return common.ERR_UNKNOWN
			}
			return ret
		}
	*/
	//c := GetFamilyPercent(1, bind_group_id)

	dump, ok := GetGamePercent(anchorId)
	if ok {
		m.MoneyNum = int64(float32(betNum) * dump)
		u.AddMoney(nil, common.MONEY_TYPE_MOON, m.MoneyNum, false)
	}

	_, err := orm.Insert(m)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	return common.ERR_SUCCESS
}

func GetLastestRecordByRid(roomid string) []NiuNiuRecord {
	m := make([]NiuNiuRecord, 0)
	yestodayBeginTime := time.Now().Unix() - 6*3600 //只取6小时以内的[直播应该不会连续超过6小时]
	//err := orm.Where("room_id=?", roomid).Desc("create_time").Limit(20, 0).Find(&m)
	err := orm.Where("create_time > ? AND room_id=?", yestodayBeginTime, roomid).Desc("id").Limit(20, 0).Find(&m)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return nil
	}
	return m
}

//用户赢取游戏币记录
func UserWinScoreRecord(gameid string, uid, score int) int {
	if score <= 0 {
		return common.ERR_UNKNOWN
	}
	m := &WinScoreRecord{
		Uid:        uid,
		GameId:     gameid,
		WinScore:   score,
		CreateTime: time.Now().Unix(),
	}

	_, err := orm.Insert(m)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	return common.ERR_SUCCESS
}
