package model

import (
	"github.com/go-xorm/xorm"
	"github.com/yshd_game/common"
	"sync"
)

//月亮商城
type ConfigItem struct {
	ItemId  int    `xorm:"int(11) pk not null"`
	Name    string `xorm:"varchar(20)`                   //名字
	Title   string `xorm:"varchar(128)`                  //副标题
	Detail  string `xorm:"varchar(128)`                  //说明
	Icon    string `xorm:"varchar(128)`                  //图标
	Moon    int    `xorm:"int(11)  not null"`            //花费月亮
	Stock   int    `xorm:"int(11)  not null"`            //库存数量
	Shelves int    `xorm:"int(11)  not null default(0)"` //0正常 1下架
	Money   int    `xorm:"int(11)  not null default(0)"` //货品对应人民币
	Status  int    `xorm:"int(11)  not null default(0)"` //0启动 1删除
}

var configItemMgr map[int]ConfigItem
var Moon_Item_mutex sync.Mutex

var configItmeCache []ConfigItem

func LoadConfigItem() map[int]ConfigItem {
	configItemMgr = make(map[int]ConfigItem)
	configItmeCache = make([]ConfigItem, 0)
	err := orm.Find(&configItemMgr)
	if err != nil {
		common.Log.Panicf("orm is err %s", err.Error())
		return configItemMgr
	}

	err = orm.Find(&configItmeCache)
	if err != nil {
		common.Log.Panicf("orm is err %s", err.Error())
		return configItemMgr
	}
	return configItemMgr
}

func GetItemById(item_id int) *ConfigItem {
	var s ConfigItem
	_, err := orm.Where("item_id=?", item_id).Get(&s)
	if err != nil {
		common.Log.Errf("orm is err %s", err.Error())
		return nil
	}
	return &s
	/*
		m, ok := configItemMgr[item_id]
		if ok {
			return &m
		}
		return nil
	*/
}

func GetItemConfig() []ConfigItem {
	configItmes := make([]ConfigItem, 0)
	err := orm.Where("status=0").Find(&configItmes)
	if err != nil {
		common.Log.Errf("orm is err %s", err.Error())
		return nil
	}

	for k, v := range configItmes {
		v.Icon = "http://h5.17playlive.com/" + v.Icon
		configItmes[k] = v
	}
	return configItmes
	//return configItmeCache
}

//减少库存
func (self *ConfigItem) DelStock(session *xorm.Session) int {
	if self.Stock > 0 {
		aff, err := session.Where("item_id=?", self.ItemId).Decr("stock", 1).Update(&ConfigItem{})
		if err != nil {
			return common.ERR_UNKNOWN
		}

		if aff == 0 {
			return common.ERR_DB_UPDATE
		}
		self.Stock--
		return common.ERR_SUCCESS
	}

	/*
		_, err := session.Where("item_id=?", self.ItemId).Cols("stock").Update(self)
		if err != nil {
			common.Log.Errf("orm is err %s", err.Error())
			return common.ERR_UNKNOWN
		}
	*/
	return common.ERR_SUCCESS
}
