package model

import (
	"fmt"
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	"strconv"
)

// type GameBonusCommission struct {
// 	Id         int64
// 	GameName   string  `xorm:"varchar(30) not null "`  // 游戏名称
// 	Desc       string  `xorm:"varchar(200) not null "` // 游戏描述
// 	BonusTimes int     `xorm:"int(11) not null "`      // 奖励倍数
// 	Commission float32 `xorm:"float(11) not null "`    // 系统抽成
// }

// func AddRecommendGoods(name, goodsImg, linkUrl string, price float32, ownerID int) int {
// 	if GetRecommendGoodsCount(ownerID) == 0 { //如果是第一个推荐的商品，则默认IsGroom为1
// 		_, err := orm.InsertOne(&RecommendGoods{Name: name, GoodsImg: goodsImg, Price: price, OwnerId: ownerID, LinkUrl: linkUrl, IsGroom: 1})
// 		if err != nil {
// 			common.Log.Errf("mysql error is %s", err.Error())
// 			return common.ERR_UNKNOWN
// 		}
// 	} else {
// 		_, err := orm.InsertOne(&RecommendGoods{Name: name, GoodsImg: goodsImg, Price: price, OwnerId: ownerID, LinkUrl: linkUrl, IsGroom: 0})
// 		if err != nil {
// 			common.Log.Errf("mysql error is %s", err.Error())
// 			return common.ERR_UNKNOWN
// 		}
// 	}

// 	return common.ERR_SUCCESS
// }

// func DelRecommendGoods(uid, goodsId int) int {
// 	_, err := orm.Where("id=? and owner_id=?", uid, goodsId).Delete(RecommendGoods{})
// 	if err != nil {
// 		common.Log.Errf("mysql error is %s", err.Error())
// 		return common.ERR_UNKNOWN
// 	}
// 	return common.ERR_SUCCESS
// }

// // 获取推荐商品个数
// func GetRecommendGoodsCount(uid int) int {
// 	sql := fmt.Sprintf("select count(*) num FROM recommend_goods WHERE owner_id = %d", uid)

// 	rowArray, err := orm.Query(sql)
// 	if err != nil {
// 		common.Log.Errf("db err %s", err.Error())
// 		return 0
// 	}
// 	number := common.BytesToString(rowArray[0]["num"])
// 	number_, _ := strconv.Atoi(number)

// 	return number_
// }

// func GoodsList(goodsTypeID, index int) (int, []map[string]string) {
// 	var sql string
// 	retMap := make([]map[string]string, 0)

// 	if goodsTypeID == 0 { //全部类别
// 		sql := fmt.Sprintf("SELECT g.goods_id, g.goods_name, t.name goods_type,g.goods_img, g.goods_price, g.goods_descripe, g.goods_address,g.on_time FROM goods g LEFT JOIN goods_type t ON g.goods_type_id=t.id WHERE g.status=1 limit %d,%d", (index)*common.MSG_LIST_PAGER_COUNT, common.MSG_LIST_PAGER_COUNT)
// 	} else { //根据商品类别的ID查询
// 		sql := fmt.Sprintf("SELECT g.goods_id, g.goods_name, t.name goods_type,g.goods_img, g.goods_price, g.goods_descripe, g.goods_address,g.on_time FROM goods g LEFT JOIN goods_type t ON g.goods_type_id=t.id WHERE g.status=1 and g.goods_type_id=%d limit %d,%d", goodsTypeID, (index)*common.MSG_LIST_PAGER_COUNT, common.MSG_LIST_PAGER_COUNT)
// 	}

// 	godump.Dump(sql)

// 	rowArray, err := orm.Query(sql)
// 	if err != nil {
// 		common.Log.Errf("db err %s", err.Error())
// 		return common.ERR_UNKNOWN, retMap
// 	}

// 	for _, row := range rowArray {
// 		ss := make(map[string]string)
// 		for colName, colValue := range row {
// 			value := common.BytesToString(colValue)
// 			ss[colName] = value
// 		}
// 		retMap = append(retMap, ss)
// 	}
// 	return common.ERR_SUCCESS, retMap
// }

var GameNiuNiuBonustimes, GameTexasPokBonustimes int
var GameNiuNiuCommission, GameTexasPokCommission float64

func GetGameConfigureInit() int {
	var bonustimes int
	var commission float64
	var isGameNiuNiu bool = false
	var isGameTexas bool = false
	fmt.Println("GetGameConfigureInit()")
	sql := fmt.Sprintf("SELECT * From go_game_bonus_comm")

	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return common.ERR_UNKNOWN
	}

	for _, row := range rowArray {
		isGameNiuNiu = false
		isGameTexas = false
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			if colName == "character" {
				if value == "NIUNIU" {
					isGameNiuNiu = true
				} else if value == "DEZHOU" {
					isGameTexas = true
				}
			} else if colName == "bonus_times" {
				bonustimes, err = strconv.Atoi(value)
				if err != nil {
					common.Log.Errf("db err %s", err.Error())
					return common.ERR_UNKNOWN
				}
			} else if colName == "commission" {
				commission, err = strconv.ParseFloat(value, 64)
				if err != nil {
					common.Log.Errf("db err %s", err.Error())
					return common.ERR_UNKNOWN
				}
			}
		}
		if isGameNiuNiu {
			GameNiuNiuBonustimes = bonustimes
			GameNiuNiuCommission = commission
		} else if isGameTexas {
			GameTexasPokBonustimes = bonustimes
			GameTexasPokCommission = commission
		}
	}
	//fmt.Printf("GameNiuNiuBonustimes=%d,GameNiuNiuCommission=%f", GameNiuNiuBonustimes, GameNiuNiuCommission)
	//fmt.Printf("GameTexasPokBonustimes=%d,GameTexasPokCommission=%f", GameTexasPokBonustimes, GameTexasPokCommission)
	return common.ERR_SUCCESS

}

func CheckUserBetReward(uid, timeStart int) (int, int, int) {
	sql := fmt.Sprintf("SELECT game_id,pos From go_niu_niu_record WHERE bet_num >0 AND create_time >= %d ", timeStart)
	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
	}
	userGain := 0
	userNotGainRaise := 0
	userTotalRaise := 0

	for _, row := range rowArray {
		ss := make(map[string]string)
		var gameId string
		var pos int
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			ss[colName] = value

			if colName == "game_id" {
				gameId = value
			}

			if colName == "pos" {
				pos, _ = strconv.Atoi(value)
			}
		}
		win, winno, totalRaise := CheckUserBetRewardFromRedis(gameId, uid, pos)
		userGain += win
		userNotGainRaise += winno
		userTotalRaise += totalRaise
	}
	return userGain, userNotGainRaise, userTotalRaise
}

func BakUserRaisePosInfoWithGameID(tStart, tEnd int64) {
	sql := fmt.Sprintf("SELECT game_id,pos From go_niu_niu_record WHERE bet_num >0 AND create_time >= %d AND create_time < %d ", tStart, tEnd)

	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
	}

	for _, row := range rowArray {
		ss := make(map[string]string)
		var gameId string
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			ss[colName] = value

			if colName == "game_id" {
				gameId = value
			}

		}
		mRaise := GetUserRaisePosInfoWithGameId(gameId)
		for key, row := range mRaise {
			sql := fmt.Sprintf("INSERT INTO go_game_raise_info(`game_id`,`uid`,`raise0_num`,`raise1_num`,`raise2_num`) VALUES ('%s',%d,%d,%d,%d)", gameId, key, row.pos0, row.pos1, row.pos2)

			_, err := orm.Query(sql)
			if err != nil {
				common.Log.Errf("db err %s", err.Error())
			}
		}
	}

}

//根据传入的参数计算是否押中
func GetRaiseSuccByPercent(f float64) bool {
	var value int64

	value = int64(f * 10000)
	random := common.RandInt64(1, 10000)
	if random <= value {
		return true
	} else {
		return false
	}
}

func GetMinIndexOfThree(left, middle, right int) int {
	if left <= middle && left <= right {
		return 0
	} else if middle <= left && middle <= right {
		return 1
	} else {
		return 2
	}
}
