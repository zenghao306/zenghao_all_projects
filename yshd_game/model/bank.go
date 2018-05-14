package model

import (
	//"github.com/liudng/godump"
	"fmt"
	"github.com/yshd_game/common"
	"strconv"
)

type PresentBank struct {
	Id       int    `xorm:"int(11) not null pk autoincr"`
	BankName string `xorm:"varchar(40) not null"` //银行名
	//Enabled     int    `xorm:"int(4) not null default(1)"` //1表示有效,2表示无效
}

//func GetPresentBankList() ([]map[string]string, int) {
//
//	sql := fmt.Sprintf("select id,bank_name  from present_bank where enabled = 1")
//	rowArray, err := orm.Query(sql)
//	if err != nil {
//		common.Log.Errf("db err %s", err.Error())
//		return nil, common.ERR_UNKNOWN
//	}
//
//	retMap := make([]map[string]string, 0)
//
//	for _, row := range rowArray {
//		ss := make(map[string]string)
//		for colName, colValue := range row {
//			value := common.BytesToString(colValue)
//			ss[colName] = value
//		}
//		retMap = append(retMap, ss)
//	}
//
//	return retMap, common.ERR_SUCCESS
//
//}

func GetPresentBankList() ([]*PresentBank, int) {

	sql := fmt.Sprintf("select id, bank_name  from go_present_bank where enabled = 1")
	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return nil, common.ERR_UNKNOWN
	}

	list := []*PresentBank{}

	for _, row := range rowArray {
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			ss[colName] = value
		}
		id, _ := strconv.Atoi(ss["id"])
		m := &PresentBank{
			Id:       id,
			BankName: ss["bank_name"],
		}

		list = append(list, m)
	}

	return list, common.ERR_SUCCESS

}

type OutAdminRes struct {
	Uid      int    `json:"uid"`
	UserName string `json:"user_name"`
}

func GetAdminList() []OutAdminRes {
	res, err := orm.Query("select a.uid,b.username from php_auth_group_access a left join php_admin b on a.uid=b.id where a.group_id=9 or   a.group_id=10 or a.group_id=11 order by b.id ")
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return nil
	}

	out := make([]OutAdminRes, 0)

	for _, v := range res {
		var m OutAdminRes
		m.Uid = common.BytesToInt(v["uid"])
		m.UserName = common.BytesToString(v["username"])
		out = append(out, m)
	}

	return out
}
