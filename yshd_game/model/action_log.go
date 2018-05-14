package model

import (
	//"github.com/liudng/godump"
	"fmt"
	"github.com/yshd_game/common"
	"strconv"
	"time"
)

type ActionLog struct {
	Id            int64
	ActionType    int `xorm:"int(11) not null "` //操作类型
	ActionReason  string
	Uid           int    `xorm:"int(11) not null "`  //用户ID
	OperationTime int64  `xorm:"int(11) default(0)"` //操作时间
	Description   string `xorm:"varchar(255)`        //描述
	Ip            string `xorm:"varchar(128)`        //IP
	Info          string `xorm:"varchar(255)`
}

type DailyRecord struct {
	Uid           int   `xorm:"int(11) not null "`  //用户ID
	OperationTime int64 `xorm:"int(11) default(0)"` //记录时间
}

func InsertLog(actType, uid int, desc string) {
	operationTime := time.Now().Unix()
	reason := common.GetDesc(actType)
	_, err := orm.Insert(&ActionLog{ActionType: actType, Uid: uid, ActionReason: reason, OperationTime: operationTime, Description: desc})
	if err != nil {
		common.Log.Err("insert db log error: , %s", err.Error())
		return
	}
}

// 根据用户ID（uid）记录用户每日log
func DailyRecordLog(uid int) {
	operation_time := time.Now().Unix()

	stdtime := time.Now()
	// t1是当日凌晨零点时间戳
	t1 := time.Date(stdtime.Year(), stdtime.Month(), stdtime.Day(), 0, 0, 0, 0, time.Local).Unix()

	// 日期上加一天
	tomorrow := stdtime.AddDate(0, 0, 1)
	// 明日凌晨零点（当日12点）时间戳
	t2 := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, time.Local).Unix()

	rowArray, _ := orm.Query("select * from go_daily_record  where uid = ? && operation_time >= ? && operation_time < ?", uid, t1, t2)
	if len(rowArray) > 0 { //查询出来的条目数大于0说明当日已经记录到了
		return
	} else { // 没查询到那就插入一条数据做记录吧。
		_, err := orm.Insert(&DailyRecord{Uid: uid, OperationTime: operation_time})
		if err != nil {
			common.Log.Err("insert db log error: , %s", err.Error())
			return
		}
	}
}

func InsertLogWithIP(actType, uid int, desc string, ip string, info string) {
	operationTime := time.Now().Unix()
	reason := common.GetDesc(actType)
	_, err := orm.Insert(&ActionLog{ActionType: actType, Uid: uid, ActionReason: reason, OperationTime: operationTime, Description: desc, Ip: ip, Info: info})
	if err != nil {
		common.Log.Err("insert db log error: , %s", err.Error())
		return
	}
}

// 根据用户id（uid）返回用户每周领取奖励记录map，并且还返回当天是否已经领取过
func LoginBonusList(uid int) ([]map[string]string, int) {
	retMap := make([]map[string]string, 0)
	todayGet := 0

	//获取到用户本周领取奖励次数
	currentWeek := common.GetCurentWeekFirstDate() //获取本周星期一凌晨零点时间（字符串类型）
	sql := fmt.Sprintf("SELECT * FROM go_user_login_bonus_record WHERE bonus_date >= '%s' AND owner_id = '%d' ORDER BY bonus_date ASC", currentWeek, uid)
	rowArray, _ := orm.Query(sql) //根据构造的sql语句查询
	length := len(rowArray)       //长度（领取次数）

	// 到每日奖励表（go_daily_bonus_type）里查询奖励list
	rowArray2, _ := orm.Query("SELECT days, name, bonus_money FROM go_daily_bonus_type WHERE status = 1")

	i := 0
	for _, row := range rowArray2 {
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			ss[colName] = value
		}

		// 以下这段是标记奖励表里对应的记录是否已领取
		if i < length {
			ss["has_get"] = "1"
		} else {
			ss["has_get"] = "0"
		}

		retMap = append(retMap, ss)
		i++
	}

	// 这段代码是判断当天是否已领取奖励
	strToday := time.Now().Format("2006-01-02")
	for _, row := range rowArray {
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			if colName == "bonus_date" && value == strToday {
				todayGet = 1
			}
		}
	}

	return retMap, todayGet
}

/*
* 函数名
*   GetDailyLoginBonus
*
* 说明
*       每日第一次登陆用户奖励相应金额
*
* 参数说明
*   uID 用户ID
*
* RETURNS
*   空
 */
func GetDailyLoginBonus(uid int) int {
	today := common.GetFormartTime2()
	sql1 := fmt.Sprintf("SELECT owner_id FROM go_user_login_bonus_record WHERE owner_id = %d && bonus_date='%s'", uid, today)
	rowArray, _ := orm.Query(sql1)
	if len(rowArray) > 0 { //查询如果没有对应的奖励任务
		return common.ERR_DAILY_BOUNTS_HAS_GET
	} else {
		curentWeekFistDay := common.GetCurentWeekFirstDate()
		sql2 := fmt.Sprintf("SELECT b.owner_id ,d.days,d.Bonus_money bonus FROM go_daily_bonus_type d ,(SELECT owner_id,COUNT(*) num FROM go_user_login_bonus_record WHERE owner_id = '%d' and bonus_date >='%s') b WHERE d.days = b.num && d.status = 1", uid, curentWeekFistDay)
		rowArray, _ = orm.Query(sql2)

		if len(rowArray) <= 0 { //如果没有对应的奖励任务
			return common.ERR_HAS_NO_DAILY_BUONTS
		} else {
			user, ret := GetUserByUid(uid)
			if ret != common.ERR_SUCCESS {
				return ret
			}

			bonusScore := 0
			for _, row := range rowArray {
				ss := make(map[string]string)
				for colName, colValue := range row {
					value := common.BytesToString(colValue)
					ss[colName] = value
				}
				bonus := ss["bonus"]
				bonus_, _ := strconv.Atoi(bonus)
				bonusScore += bonus_
			}

			timeNow := time.Now().Format("2006-01-02 15:04:05")
			sql := fmt.Sprintf("INSERT INTO go_user_login_bonus_record(`owner_id`,`bonus_type`,`bonus_money`,`bonus_date`,`time`) SELECT '%d' ,d.days,d.Bonus_money,'%s','%s' FROM go_daily_bonus_type d,(SELECT owner_id,COUNT(*) num FROM go_user_login_bonus_record WHERE owner_id = '%d' and bonus_date >='%s') b WHERE d.days = b.num", uid, today, timeNow, uid, curentWeekFistDay)

			_, err := orm.Query(sql)
			if err != nil {
				common.Log.Errf("db err %s", err.Error())
				return common.ERR_UNKNOWN
			} else {
				user.AddMoney(nil, common.MONEY_TYPE_SCORE, int64(bonusScore), false)
				//user.AddScore(int64(bonusScore))
			}
		}
	}
	return common.ERR_SUCCESS
}

// 根据传入的参数时间戳【times】删除用户之前的牛牛押分数据
func DeleteNiuNiuRaiseRecordLast(times1, times2 int64) {
	sql := fmt.Sprintf("SELECT game_id FROM go_niu_niu_record WHERE create_time <= %d && create_time > %d", times1, times2)
	rowArray, _ := orm.Query(sql)

	if len(rowArray) <= 0 { //如果没有对应的数据，返回
		return
	} else {
		for _, row := range rowArray {
			ss := make(map[string]string)
			for colName, colValue := range row {
				value := common.BytesToString(colValue)
				ss[colName] = value
			}
			gameID := ss["game_id"]

			DelUserRaiseRedisInfo(gameID)
		}
	}
}

// 根据传入的时间戳times删掉此时间戳以前的游戏记录
func DeleteNiuNiuRecordBeforeTheTime(times int64) {
	sql := fmt.Sprintf("DELETE FROM go_niu_niu_record WHERE create_time <= %d AND bet_num = 0", times)

	_, err := orm.Exec(sql)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
	}
}
