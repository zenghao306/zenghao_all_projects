package model

import (
	"github.com/yshd_game/common"
	"time"
)

type ConfigAndroidTradeItem struct {
	ItemId   string `xorm:"varchar(20) pk not null "`
	Describe string `xorm:"varchar(20) "`
	Money    int    `xorm:"int(11) default(0)"` //单位分
	Diamond  int    `xorm:"int(11) default(0)"`
}

/*
type AndroidTradeItemCache struct {
	ItemId   string  `xorm:"varchar(20) pk not null "`
	Describe string  `xorm:"varchar(20) "`
	Money    float32 `xorm:"int(11) default(0)"` //单位分
	Diamond  int     `xorm:"int(11) default(0)"`
}
*/

type ConfigIosTradeItem struct {
	ItemId   string `xorm:"varchar(20) pk not null "`
	Describe string `xorm:"varchar(20) "`
	Money    int    `xorm:"int(11) default(0)"` //单位元
	Diamond  int    `xorm:"int(11) default(0)"`
	Category int    `xorm:"int(11) default(0)"`
	//ChannelId  string    `xorm:"varchar(40) default(0)"`
}

type ConfigGoogleTradeItem struct {
	ItemId   string `xorm:"varchar(20) pk not null "`
	Describe string `xorm:"varchar(20) "`
	Money    int    `xorm:"int(11) default(0)"` //单位分
	Diamond  int    `xorm:"int(11) default(0)"`
}

type GoogleTradeItemCache struct {
	ItemId   string  `xorm:"varchar(20) pk not null "`
	Describe string  `xorm:"varchar(20) "`
	Money    float32 `xorm:"float default(0)"` //单位分
	Diamond  int     `xorm:"int(11) default(0)"`
}

type IOSConfigV2 struct {
	ItemId   string `json:"ItemId""`
	Describe string `json:"Describe"`
	Money    int    `json:"Money"` //单位元
	Diamond  int    `json:"Diamond"`
	Category int    `json:"Category"`
	ExtraNum int    `json:"ExtraNum"`
}

type AndroidTradeItemCache struct {
	ItemId   string  `json:"ItemId""`
	Describe string  `json:"Describe"`
	Money    float32 `json:"Money"` //单位元
	Diamond  int     `json:"Diamond"`
	ExtraNum int     `json:"ExtraNum"`
}

//android load config
var android_trade_map map[string]ConfigAndroidTradeItem
var android_cache_arr []AndroidTradeItemCache

func LoadAndroidPay() map[string]ConfigAndroidTradeItem {
	android_trade_map = make(map[string]ConfigAndroidTradeItem)
	android_cache_arr = make([]AndroidTradeItemCache, 0)
	err := orm.Find(&android_trade_map)
	if err != nil {
		common.Log.Panic("orm err is %s", err.Error())
	}

	android_cache_arr_temp := make([]ConfigAndroidTradeItem, 0)

	err = orm.OrderBy("money").Find(&android_cache_arr_temp)
	if err != nil {
		common.Log.Panic("orm err is %s", err.Error())
	}
	now_time := time.Now().Unix()
	for _, v := range android_cache_arr_temp {
		var a AndroidTradeItemCache
		a.Describe = v.Describe
		a.Diamond = v.Diamond
		a.ItemId = v.ItemId
		a.Money = float32(float32(v.Money) / 100)

		active := GetChargeActive(v.ItemId)
		if active != nil {
			if active.Status == 1 && active.BeginTime < now_time && active.FinishTime > now_time {
				a.ExtraNum = int(active.ExtraNum)
			}
		}
		android_cache_arr = append(android_cache_arr, a)
	}

	/*
		err = orm.OrderBy("money").Find(&android_cache_arr)
		if err != nil {
			common.Log.Panic("orm err is %s", err.Error())
		}

		for k, v := range android_cache_arr {
			v.Money = v.Money / 100
			android_cache_arr[k] = v
		}
	*/
	return android_trade_map
}

func GetAndroidItem(id string) (*ConfigAndroidTradeItem, bool) {
	if v, ok := android_trade_map[id]; ok {
		return &v, true
	} else {
		return &v, false
	}
}

func GetAndroidItemConfig() (int, []AndroidTradeItemCache) {
	return common.ERR_SUCCESS, android_cache_arr
}

//ios  load  config
var ios_trade_map map[string]ConfigIosTradeItem
var ios_cache_arr []ConfigIosTradeItem

func LoadIOSPay() map[string]ConfigIosTradeItem {
	ios_trade_map = make(map[string]ConfigIosTradeItem)
	ios_cache_arr = make([]ConfigIosTradeItem, 0)
	err := orm.Find(ios_trade_map)
	if err != nil {
		common.Log.Panic("orm err is %s", err.Error())
	}

	err = orm.OrderBy("money").Find(&ios_cache_arr)
	if err != nil {
		common.Log.Panic("orm err is %s", err.Error())
	}
	return ios_trade_map
}

func GetIOSItem(id string) (*ConfigIosTradeItem, bool) {
	if v, ok := ios_trade_map[id]; ok {
		return &v, true
	} else {
		return &v, false
	}
}

func GetIOSItemConfig() (int, []ConfigIosTradeItem) {
	return common.ERR_SUCCESS, ios_cache_arr
}

//
//func GetIOSItemConfigByChannel(category int) (ret int, res []ConfigIosTradeItem) {
func GetIOSItemConfigByChannel(category int) (ret int, res []IOSConfigV2) {
	//	res = make([]ConfigIosTradeItem, 0)
	//	err := orm.Where("category=?", category).OrderBy("money").Find(&res)
	res = make([]IOSConfigV2, 0)
	err := orm.Sql("SELECT * FROM go_config_ios_trade_item a LEFT JOIN (SELECT * FROM go_charge_active WHERE STATUS=1 AND begin_time<UNIX_TIMESTAMP() AND finish_time>UNIX_TIMESTAMP()) b ON a.item_id=b.item_id WHERE a.category=?  ORDER BY a.money ASC", category).Find(&res)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		ret = common.ERR_UNKNOWN
		return

	}
	if len(res) == 0 {
		ret = common.ERR_CONFGI_ITEM
		return
	}
	ret = common.ERR_SUCCESS
	return
}

//google  load  config
var google_trade_map map[string]ConfigGoogleTradeItem
var google_cache_arr []GoogleTradeItemCache

func LoadGooglePay() map[string]ConfigGoogleTradeItem {
	google_trade_map = make(map[string]ConfigGoogleTradeItem)
	google_cache_arr = make([]GoogleTradeItemCache, 0)
	err := orm.Find(&google_trade_map)
	if err != nil {
		common.Log.Panic("orm err is %s", err.Error())
	}

	android_cache_arr_temp := make([]GoogleTradeItemCache, 0)
	err = orm.OrderBy("money").Find(&android_cache_arr_temp)
	if err != nil {
		common.Log.Panic("orm err is %s", err.Error())
	}

	for _, v := range android_cache_arr_temp {
		var a AndroidTradeItemCache
		a.Describe = v.Describe
		a.Diamond = v.Diamond
		a.ItemId = v.ItemId
		a.Money = float32(float32(v.Money) / 100)
		android_cache_arr = append(android_cache_arr, a)
	}
	return google_trade_map
}

func GetGoogleItem(id string) (*ConfigGoogleTradeItem, bool) {
	if v, ok := google_trade_map[id]; ok {
		return &v, true
	} else {
		return &v, false
	}
}

func GetGoogleItemConfig() (int, []GoogleTradeItemCache) {
	return common.ERR_SUCCESS, google_cache_arr
}
