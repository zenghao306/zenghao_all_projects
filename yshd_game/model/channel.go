package model

import (
	"fmt"
	"github.com/yshd_game/common"
)

func ChannelVersion(channel, version string) (int, []map[string]string) {
	retMap := make([]map[string]string, 0)

	sql := fmt.Sprintf("SELECT pay_type,star_cash_switch,moon_cash_switch,ranking_switch,share_switch,banner_switch,select_show_switch,model FROM php_channel_version_config WHERE channel_id = '%s' && version = '%s'", channel, version)

	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return common.ERR_UNKNOWN, retMap
	}

	if len(rowArray) > 0 { // 如果查询到了相关信息情况
		for _, row := range rowArray {
			ss := make(map[string]string)
			for colName, colValue := range row {
				value := common.BytesToString(colValue)
				ss[colName] = value
			}
			retMap = append(retMap, ss)
		}
	} else { //如果没查询到，给予默认值
		ss := make(map[string]string)
		ss["pay_type"] = "2"
		ss["star_cash_switch"] = "1"
		ss["moon_cash_switch"] = "1"
		ss["ranking_switch"] = "1"
		ss["share_switch"] = "1"
		ss["banner_switch"] = "1"
		ss["select_show_switch"] = "1"
		ss["model"] = "0"

		retMap = append(retMap, ss)
	}

	return common.ERR_SUCCESS, retMap
}
