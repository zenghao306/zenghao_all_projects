package model

import (
	"fmt"
	"github.com/yshd_game/common"
)

var (
	EveryDayLowestCashTime = common.EVERY_DAY_LOWEST_CASH_TIME
	NiuNiuPem              = 0.5
	NickNameResetMoney     = common.NICKNAMERESETMONEY
)

func BannerList() (int, []map[string]string) {
	retMap := make([]map[string]string, 0)

	sql := fmt.Sprintf("SELECT title,link_url,banner_img,weight FROM php_app_banner WHERE `status` = 1 ORDER BY weight DESC")

	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return common.ERR_UNKNOWN, retMap
	}

	for _, row := range rowArray {
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			if colName == "banner_img" {
				ss[colName] = "http://h5.17playlive.com/" + value
			} else {
				ss[colName] = value
			}
		}
		retMap = append(retMap, ss)
	}

	return common.ERR_SUCCESS, retMap
}

type ConfigSystemVariable struct {
	VariableName string `xorm:"varchar(20) pk not null "` //变量名
	ValueInt     int
	ValueFloat   float64
}

//系统相关变量初始化或者重设定
func SystemVariableInitOrReset() int {
	lowestCashTime := &ConfigSystemVariable{} //直播有效天对应最低在线时长
	has, err := orm.Where("variable_name=?", "every_day_lowest_cash_time").Get(lowestCashTime)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if has {
		EveryDayLowestCashTime = lowestCashTime.ValueInt
	} else {
		EveryDayLowestCashTime = common.EVERY_DAY_LOWEST_CASH_TIME
	}

	niuniu_pem := &ConfigSystemVariable{}
	has2, err := orm.Where("variable_name=?", "niuniu_pem").Get(niuniu_pem)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if has2 {
		NiuNiuPem = niuniu_pem.ValueFloat
	} else {
		NiuNiuPem = 0.5
	}

	nickNameResetMoney := &ConfigSystemVariable{} //获取重设昵称所需游戏币
	has3, err := orm.Where("variable_name=?", "nick_name_reset_money").Get(nickNameResetMoney)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if has3 {
		NickNameResetMoney = nickNameResetMoney.ValueInt
	} else {
		NickNameResetMoney = common.NICKNAMERESETMONEY
	}

	return common.ERR_SUCCESS
}
