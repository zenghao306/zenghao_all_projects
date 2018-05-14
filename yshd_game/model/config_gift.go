package model

import (
	"fmt"
	"github.com/yshd_game/common"
	//"strconv"
	"sort"
)

//var gift_config_map map[int]Gift
var gift_config_map map[int]ConfigGift
var gift_config_arr_cache []ConfigGift

var gift_tip_arr []ConfigGift

//gift_id	weight	price	type	name	icon	pic	category	dynamic
//用户表结构
type ConfigGift struct {
	GiftId   int    `xorm:"int(11) pk not null unique"`
	Type     int    `xorm:"int(11) not null"` //0连击礼物 1不是
	Price    int    `xorm:"int(11) not null"` //礼物价格
	Name     string `xorm:"varchar(20)`       //名字
	Icon     string `xorm:"varchar(50)`       //图标
	Pic      string `xorm:"varchar(50)`       //大图
	Category int    `xorm:"int(11) not null"` //礼物分类 0普通  1豪华 2游戏 3游戏豪华
	//	AllocId  int    `xorm:"int(11) not null"`
	Dynamic int `xorm:"int(11) not null default(0)"` //是不是动态礼物
	Weight  int `xorm:"int(11) not null default(0)"`
}

type GiftWrapper struct {
	gift []ConfigGift
	by   func(p, q *ConfigGift) bool
}

type SortBy func(p, q *ConfigGift) bool

func (pw GiftWrapper) Len() int { // 重写 Len() 方法
	return len(pw.gift)
}
func (pw GiftWrapper) Swap(i, j int) { // 重写 Swap() 方法
	pw.gift[i], pw.gift[j] = pw.gift[j], pw.gift[i]
}
func (pw GiftWrapper) Less(i, j int) bool { // 重写 Less() 方法
	return pw.by(&pw.gift[i], &pw.gift[j])
}

func LoadGift() map[int]ConfigGift {
	gift_config_map = make(map[int]ConfigGift)

	gift_config_arr_cache = make([]ConfigGift, 0)
	//gift_config_map_cache = make(map[string]Gift)
	err := orm.Desc("price").Find(&gift_config_map)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return gift_config_map
	}
	//root := common.Cfg.MustValue("path", "root_path")
	//gift := common.Cfg.MustValue("path", "gift")
	for k, v := range gift_config_map {
		v.Icon = DownloadUrl(DomainGift, v.Icon)
		v.Pic = DownloadUrl(DomainGift, v.Pic)
		gift_config_map[k] = v

		gift_config_arr_cache = append(gift_config_arr_cache, v)
	}

	sort.Sort(GiftWrapper{gift_config_arr_cache, func(p, q *ConfigGift) bool {
		return q.Weight > p.Weight //  递减排序
	}})

	LoadTipGift()
	return gift_config_map
}

func GetGiftById(gid int) (ConfigGift, bool) {
	gift, exist := gift_config_map[gid]
	return gift, exist
}

func GetGiftConfig() (int, []ConfigGift) {
	return common.ERR_SUCCESS, gift_config_arr_cache
	/*
		gift := make([]Gift, 0)
		err := orm.Find(&gift)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN, gift
		}
		return common.ERR_SUCCESS, gift
	*/
}

func AllGiftRank(index int) ([]map[string]string, int) {
	sql := fmt.Sprintf("select a.send_user,a.count,b.nick_name,b.user_level as level,b.image,b.sex from (select send_user,sum(value) as count from go_gift_record  group by send_user) a left join go_user b on a.send_user=b.uid  where b.account_type!=1  order by a.count  desc limit %d,%d ", index*common.ROOM_LIST_PAGE_COUNT, common.ROOM_LIST_PAGE_COUNT)

	//sql:=fmt.Sprintf("select a.send_user,a.count,b.nick_name,b.user_level as level,b.image,b.sex from (select send_user,sum(value) as count from go_gift_record  where money_type=%d group by send_user) a left join go_user b on a.send_user=b.uid  where b.account_type!=1  order by a.count  desc limit 0,50 ",common.MONEY_TYPE_DIAMOND)
	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return nil, common.ERR_UNKNOWN
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

	return retMap, common.ERR_SUCCESS
}

func GetMoonRankList(revUid, index int) ([]map[string]string, int) {

	retMap := make([]map[string]string, 0)
	u, ok := GetUserByUid(revUid)
	if ok != common.ERR_SUCCESS {
		return retMap, ok
	}
	var sql string
	if u.AccountType == 1 {
		sql = fmt.Sprintf("SELECT a.send_user,a.count,b.nick_name,b.anchor_level,b.user_level,b.image,b.sex FROM (SELECT send_user,SUM(value) AS COUNT, MIN(create_time) as CREATE_TIME FROM go_gift_record WHERE rev_user=%d and money_type=%d GROUP BY send_user)a LEFT JOIN go_user b ON a.send_user=b.uid   ORDER BY a.count  DESC, CREATE_TIME ASC  limit %d, %d", revUid, common.MONEY_TYPE_SCORE, index*common.SEND_GIFT_PAGE_COUNT, common.SEND_GIFT_PAGE_COUNT)
	} else {
		sql = fmt.Sprintf("SELECT a.send_user,a.count,b.nick_name,b.anchor_level,b.user_level,b.image,b.sex FROM (SELECT send_user,SUM(value) AS COUNT, MIN(create_time) as CREATE_TIME FROM go_gift_record WHERE rev_user=%d and money_type=%d GROUP BY send_user)a LEFT JOIN go_user b ON a.send_user=b.uid   where b.account_type!=1 ORDER BY a.count  DESC, CREATE_TIME ASC  limit %d, %d", revUid, common.MONEY_TYPE_SCORE, index*common.SEND_GIFT_PAGE_COUNT, common.SEND_GIFT_PAGE_COUNT)
	}

	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return nil, common.ERR_UNKNOWN
	}

	for _, row := range rowArray {
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			ss[colName] = value
		}
		retMap = append(retMap, ss)
	}

	return retMap, common.ERR_SUCCESS
}

func GetCouponsRankList(revUid, index int) ([]map[string]string, int) {

	u, ok := GetUserByUid(revUid)
	if ok != common.ERR_SUCCESS {
		return nil, ok
	}
	var sql string
	if u.AccountType == 1 {
		sql = fmt.Sprintf("SELECT a.send_user,a.count,b.nick_name,b.anchor_level,b.user_level,b.image,b.sex FROM (SELECT send_user,SUM(value) AS COUNT, MIN(create_time) as CREATE_TIME FROM go_gift_record WHERE rev_user=%d and money_type=%d GROUP BY send_user)a LEFT JOIN go_user b ON a.send_user=b.uid  ORDER BY a.count  DESC, CREATE_TIME ASC  limit %d, %d", revUid, common.MONEY_TYPE_DIAMOND, index*common.SEND_GIFT_PAGE_COUNT, common.SEND_GIFT_PAGE_COUNT)
	} else {
		sql = fmt.Sprintf("SELECT a.send_user,a.count,b.nick_name,b.anchor_level,b.user_level,b.image,b.sex FROM (SELECT send_user,SUM(value) AS COUNT, MIN(create_time) as CREATE_TIME FROM go_gift_record WHERE rev_user=%d and money_type=%d GROUP BY send_user)a LEFT JOIN go_user b ON a.send_user=b.uid   where b.account_type!=1 ORDER BY a.count  DESC, CREATE_TIME ASC  limit %d, %d", revUid, common.MONEY_TYPE_DIAMOND, index*common.SEND_GIFT_PAGE_COUNT, common.SEND_GIFT_PAGE_COUNT)
	}

	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return nil, common.ERR_UNKNOWN
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

	return retMap, common.ERR_SUCCESS
}

func GetCouponsRankListWithTime(begin int64, end int64) ([]map[string]string, int) {
	rowArray, err := orm.Query("SELECT a.rev_user,SUM(VALUE) AS count,b.nick_name,b.anchor_level,b.image,b.sex FROM go_gift_record a LEFT JOIN go_user b ON a.rev_user=b.uid  WHERE a.money_type=? AND  a.record_time>? AND a.record_time<? and account_type=0 GROUP BY a.rev_user  ORDER BY count desc limit 0,50 ", common.MONEY_TYPE_DIAMOND, begin, end)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return nil, common.ERR_UNKNOWN
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

	return retMap, common.ERR_SUCCESS
}

func LoadTipGift() int {
	gift_tip_arr = make([]ConfigGift, 0)
	err := orm.Where("category=2 ").Limit(3, 0).Asc("weight").Find(&gift_tip_arr)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return common.ERR_UNKNOWN
	}

	for k, v := range gift_tip_arr {
		v.Icon = DownloadUrl(DomainGift, v.Icon)
		if v.Pic != "" {
			v.Pic = DownloadUrl(DomainGift, v.Pic)
		}

		gift_tip_arr[k] = v
	}
	return common.ERR_SUCCESS
}

func GetTipGift() []ConfigGift {
	return gift_tip_arr
}
