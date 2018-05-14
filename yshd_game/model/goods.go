package model

import (
	"fmt"
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	"strconv"
)

type Goods struct {
	Id          int64
	GoodsName   string  `xorm:"varchar(120) not null "` // 商品名称
	GoodsTypeId int     `xorm:"int(5) not null "`       // 商品类型id
	GoodsImg    string  `xorm:"varchar(255) not null "` // 商品图片
	GoodsPrice  float32 `xorm:"float(10) not null "`    // 商品价格
	GoodsDes    string  `xorm:"varchar(255) not null "` // 商品描述
	GoodsAddr   string  `xorm:"varchar(255) not null "` // 商品地址
	OnTime      int     `xorm:"int(11) not null "`      // 上架时间
	Status      int     `xorm:"int(11) not null "`      // 商品状态
}

func UserSaleGoodsList(uid, index int) (int, []map[string]string) {
	retMap := make([]map[string]string, 0)

	sql := fmt.Sprintf("SELECT us.uid, us.goods_id,us.status, g.goods_name, t.name goods_type,g.goods_img, g.goods_price, g.goods_descripe, g.goods_address,g.on_time FROM php_user_sale_goods us LEFT JOIN php_goods g ON us.goods_id=g.goods_id LEFT JOIN php_goods_type t ON g.goods_type_id=t.id WHERE us.status != 0 and us.uid=%d ORDER BY us.status DESC limit %d,%d", uid, (index)*common.MSG_LIST_PAGER_COUNT, common.MSG_LIST_PAGER_COUNT)

	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return common.ERR_UNKNOWN, retMap
	}

	for _, row := range rowArray {
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			if colName == "goods_img" {
				ss[colName] = "http://admin.shangtv.cn/" + value
			} else {
				ss[colName] = value
			}
		}
		retMap = append(retMap, ss)
	}
	return common.ERR_SUCCESS, retMap
}

func GetUserSaleGoodsNumbers(uid int) (int, int) {
	sql := fmt.Sprintf("SELECT COUNT(*) num FROM user_sale_goods us WHERE us.uid=%d and us.status!=0", uid)

	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return common.ERR_UNKNOWN, 0
	}

	number := common.BytesToString(rowArray[0]["num"])
	number_, _ := strconv.Atoi(number)

	return common.ERR_SUCCESS, number_
}

// 获取精选列表
func GetShowList(index int) (int, []map[string]string) {
	retMap := make([]map[string]string, 0)

	sql := fmt.Sprintf("SELECT ss.id,link_url,ss.show_name,ss.show_sign,ss.show_img,ss.show_hot, ss.minute,ss.second,g.goods_address FROM php_select_show ss LEFT JOIN php_goods g ON ss.goods_id=g.goods_id WHERE ss.show_status=1 ORDER BY ss.show_weight DESC limit %d,%d", (index)*common.MSG_LIST_PAGER_COUNT, common.MSG_LIST_PAGER_COUNT)
	//sql := fmt.Sprintf("SELECT ss.link_url,ss.show_name,ss.show_sign,ss.show_img,ss.show_hot, ss.minute,ss.second,g.goods_address FROM php_select_show ss LEFT JOIN php_goods g ON ss.goods_id=g.goods_id WHERE ss.show_status=1 ORDER BY ss.show_weight DESC limit %d,%d", (index)*common.MSG_LIST_PAGER_COUNT, common.MSG_LIST_PAGER_COUNT)

	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return common.ERR_UNKNOWN, retMap
	}

	for _, row := range rowArray {
		ss := make(map[string]string)

		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			if colName == "show_img" {
				ss[colName] = "http://h5.17playlive.com/" + value
			} else {
				ss[colName] = value
			}
		}
		retMap = append(retMap, ss)
	}
	return common.ERR_SUCCESS, retMap
}

func SelectShowClick(id int) int {
	sql := fmt.Sprintf("update php_select_show set show_hot=show_hot+1 where id = %d", id)

	_, err := orm.Exec(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return common.ERR_UNKNOWN
	}

	return common.ERR_SUCCESS
}
