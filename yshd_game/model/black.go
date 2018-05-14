package model

import (
	"fmt"
	"github.com/yshd_game/common"
)

type Black struct {
	Id      int64
	OwnerId int `xorm:"int(11) not null UNIQUE(BLACK)"` //主动者ID
	BlackId int `xorm:"int(11) not null UNIQUE(BLACK)"` //黑名单者ID
}

//curl -d 'token=1f59dd2bce2c73d041556ae8f85f9341&uid=7&blackid=5' 'http://192.168.1.12:3000/black/add_black'
func AddBlack(uid, black_id int) int {
	aff_row, err := orm.InsertOne(&Black{OwnerId: uid, BlackId: black_id})
	if err != nil {
		common.Log.Errf("mysql error is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		return common.ERR_DB_ADD
	}
	return common.ERR_SUCCESS

}

//curl -d 'token=1f59dd2bce2c73d041556ae8f85f9341&uid=7&id=1' 'http://192.168.1.12:3000/black/del_black'
/*
func DelBlack(uid, id int) int {
	_, err := orm.Where("id=? and owner_id=?", id, uid).Delete(Black{})
	if err != nil {
		common.Log.Errf("mysql error is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	return common.ERR_SUCCESS
}
*/
func DelBlack(uid, blackid int) int {
	aff_row, err := orm.Where("black_id=? and owner_id=?", blackid, uid).Delete(Black{})
	if err != nil {
		common.Log.Errf("mysql error is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		return common.ERR_DB_DEL
	}
	return common.ERR_SUCCESS

}

//curl  'http://192.168.1.12:3000/black/black_list?index=0&uid=1&token=cc4c1d4636e75a002782ff4e7e267829'
func BlackList(uid, index int) (int, []map[string]string) {
	/*
		black := make([]Black, 0)
		err := orm.Where("owner_id=?", uid).Limit(common.BLACK_LIST_PAGE_COUNT, index*common.BLACK_LIST_PAGE_COUNT).Find(&black)
		if err != nil {
			common.Log.Errf("mysql error is %s", err.Error())
			return common.ERR_UNKNOWN, black
		}
		return common.ERR_SUCCESS, black
	*/
	retMap := make([]map[string]string, 0)
	sql := fmt.Sprintf("SELECT b.image ,b.uid,b.signature,b.location,b.nick_name,b.user_level as level FROM go_black a LEFT JOIN user b ON a.black_id=b.uid WHERE a.owner_id=%d limit %d,%d", uid, (index)*common.MSG_LIST_PAGER_COUNT, common.MSG_LIST_PAGER_COUNT)
	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return common.ERR_UNKNOWN, retMap
	}

	for _, row := range rowArray {
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			ss[colName] = value
		}
		retMap = append(retMap, ss)
	}
	return common.ERR_SUCCESS, retMap
}
