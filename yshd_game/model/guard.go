package model

import (
	"fmt"
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	"time"
)

type Guard struct {
	Uid      int
	AnchorId int
	//CreateTime int64
	FinishTime int64
	Tip        int
}

type GuardRecord struct {
	Id         int64
	Uid        int `xorm:"int(11) not null "`
	AnchorId   int `xorm:"int(11) not null "`
	CreateTime int64
	Price      int `xorm:"int(11) not null "`
	GuardType  int `xorm:"int(11) not null "`
}

type GuardRecordDetail struct {
	GuardRecordId int64 `xorm:"int(11) not null "`
	Identity      int   `xorm:"int(11) not null "` //身份
	MoneyType     int   `xorm:"int(11) not null "` //金钱类型
	Num           int64 `xorm:"int(11) not null "` //金钱数量
}

func GetGuard(uid, anchor int) *Guard {
	m := &Guard{}
	has, err := orm.Where("uid=? and anchor_id=?", uid, anchor).Get(m)
	if err != nil {
		common.Log.Err(err.Error())
		return nil
	}
	if has == false {
		return nil
	}
	return m
}

func CheckGuard(uid, anchor int) int {
	guard := GetGuard(uid, anchor)
	if guard != nil {
		if guard.FinishTime >= time.Now().Unix() {
			return 1
		}
	}
	return 0
}

func ResetGuard(uid int) {
	m := make([]Guard, 0)
	err := orm.Where("uid=? ", uid).Find(&m)
	if err != nil {
		common.Log.Err(err.Error())
		return
	}
	now_time := time.Now().Unix()
	for _, v := range m {
		if v.FinishTime <= now_time {
			anchor, ret := GetUserByUid(v.AnchorId)
			if ret != common.ERR_SUCCESS {
				continue
			}
			msg := fmt.Sprintf("你对主播：%s的守护已经到期，若要继续守护请及时竞价", anchor.NickName)
			SendLetter(1, uid, msg)
			_, err := orm.Where("uid=? and anchor_id=?", uid, v.AnchorId).Delete(&v)
			if err != nil {
				common.Log.Err(err.Error())
			}

		} else if v.FinishTime <= now_time+24*3600 && v.Tip == 0 {
			anchor, ret := GetUserByUid(v.AnchorId)
			if ret != common.ERR_SUCCESS {
				continue
			}
			msg := fmt.Sprintf("你对主播：%s的守护还有1天就要到期，若要继续守护请及时续费", anchor.NickName)
			SendLetter(1, uid, msg)
			v.Tip = 1
			_, err := orm.Where("uid=? and anchor_id=?", uid, v.AnchorId).Update(v)
			if err != nil {
				common.Log.Err(err.Error())
			}
		}
	}
}

func OpenGuardAnchor(uid, anchor_id int) int {
	now_time := time.Now().Unix()

	var first int
	guard := GetGuard(uid, anchor_id)
	if guard != nil {
		if guard.FinishTime-now_time+common.GUARD_KEEP_TIME >= common.GUARD_KEEP_TIME_MAX {
			return common.ERR_GUARD_LIMIE_TIME
		}
	}

	user, ret := GetUserByUid(uid)
	if ret != common.ERR_SUCCESS {
		return ret
	}

	if user.IsSuperUser() == true {
		return common.ERR_GUARD_SUPER
	}
	anchor, ret := GetUserByUid(anchor_id)
	if ret != common.ERR_SUCCESS {
		return ret
	}

	ret, bind_user_id := GetBindUser(anchor.AdminId)
	if ret != common.ERR_SUCCESS {
		//godump.Dump(ret)
		return ret
	}

	bind_user, ret2 := GetUserByUid(bind_user_id)
	if ret2 != common.ERR_SUCCESS {
		return ret2
	}

	commossion, has := GetGuardPercent(anchor_id, 1)
	if has == false {
		return common.ERR_CONFGI_ITEM
	}

	allnum, has := GetConfigGuardPrice(1)
	if has == false {
		return common.ERR_CONFGI_ITEM
	}

	selfsess := GetUserSessByUid(uid)
	if selfsess == nil {
		return common.ERR_USER_OFFLINE
	}

	chat := GetChatRoom(selfsess.Roomid)

	var send_num, bind_num, sys_num float32
	send_num = float32(allnum) * commossion.OwnerPercent
	sys_num = float32(allnum) * commossion.SystemPercent
	bind_num = float32(allnum) * commossion.LeaderPercent

	if ret := user.CheckMoney(common.MONEY_TYPE_DIAMOND, int64(allnum)); ret != common.ERR_SUCCESS {
		return ret
	}

	session := orm.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		common.Log.Errf("orm is error:  %s", err.Error())
		session.Rollback()
		return common.ERR_UNKNOWN
	}

	ret = user.DelMoney(session, common.MONEY_TYPE_DIAMOND, int64(allnum), true)
	if ret != common.ERR_SUCCESS {
		session.Rollback()
		return ret
	}

	if user.AccountType == 0 || user.AccountType == 1 && anchor.AccountType == 1 {

		ret = anchor.AddMoney(session, common.MONEY_TYPE_RICE, int64(send_num), true)
		if ret != common.ERR_SUCCESS {
			session.Rollback()
			return common.ERR_UNKNOWN
		}
		if ret2 == common.ERR_SUCCESS {
			ret = bind_user.AddMoney(session, common.MONEY_TYPE_RICE, int64(bind_num), true)
			if ret != common.ERR_SUCCESS {
				session.Rollback()
				return common.ERR_UNKNOWN
			}
		}
	}

	if guard != nil {
		guard.FinishTime = guard.FinishTime + common.GUARD_KEEP_TIME
		guard.Tip = 0
		aff, err := session.Where("uid=? and anchor_id=?", uid, anchor_id).Update(guard)
		if err != nil {
			common.Log.Errf("orm is error:  %s", err.Error())
			return common.ERR_UNKNOWN
		}

		if aff == 0 {
			session.Rollback()
			return common.ERR_DB_UPDATE
		}
	} else {
		m := &Guard{
			Uid:        uid,
			AnchorId:   anchor_id,
			FinishTime: time.Now().Unix() + common.GUARD_KEEP_TIME,
		}

		_, err = session.InsertOne(m)
		if err != nil {
			common.Log.Errf("orm is error:  %s", err.Error())
			return common.ERR_UNKNOWN
		}

		first = 1
	}

	err = session.Commit()
	if err != nil {
		common.Log.Errf("err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if first == 1 {

		if chat != nil {
			room := GetChatRoom(chat.room.Rid)
			mutex_chat_guardv2.Lock()
			room.UpdateAudience(uid)
			mutex_chat_guardv2.Unlock()
		}
	}

	rl, err := orm.Exec("insert into `go_guard_record` (`uid`,`anchor_id`,`create_time`,`price`,`guard_type`) values (?,?,?,?,?)", uid, anchor_id, time.Now().Unix(), allnum, 1)
	if err != nil {
		common.Log.Errf("orm is error:  %s", err.Error())
		return common.ERR_UNKNOWN
	}
	resId, err := rl.LastInsertId()
	if err != nil {
		common.Log.Errf("orm is error:  %s", err.Error())
		return common.ERR_UNKNOWN
	}

	r1 := &GuardRecordDetail{GuardRecordId: resId, Identity: common.LIVE_IDENTITY_USER, MoneyType: common.MONEY_TYPE_RICE, Num: int64(send_num)}

	r2 := &GuardRecordDetail{GuardRecordId: resId, Identity: common.LIVE_IDENTITY_SYS, MoneyType: common.MONEY_TYPE_RICE, Num: int64(sys_num)}

	r3 := &GuardRecordDetail{GuardRecordId: resId, Identity: common.LIVE_IDENTITY_ADMIN, MoneyType: common.MONEY_TYPE_RICE, Num: int64(bind_num)}

	aff, err := orm.Insert(r1, r2, r3)
	if err != nil || aff == 0 {
		if err != nil {
			common.Log.Errf("err is %s", err.Error())
		}
		return common.ERR_UNKNOWN
	}

	if chat != nil {
		mutex_chat_guardv2.Lock()
		if user.AccountType == 1 && anchor.AccountType == 1 {
			chat.AddRice(int(send_num))
		} else if user.AccountType == 0 {
			chat.AddRice(int(send_num))
		}
		mutex_chat_guardv2.Unlock()
	}

	sess := GetUserSessByUid(uid)
	if sess != nil {
		msg := ResponseGuardOpen{
			MType:    common.MESSAGE_OPEN_GUARD,
			Uid:      uid,
			NickName: user.NickName,
			AnchorId: anchor_id,
			First:    first,
		}
		SendMsgToRoom(sess.Roomid, msg)
	}

	return common.ERR_SUCCESS
}

func GetGuardCount(anchor int) (res int64) {
	res, err := orm.Where("anchor_id=? and finish_time>UNIX_TIMESTAMP()", anchor).Count(&Guard{})
	if err != nil {
		common.Log.Errf("orm is error:  %s", err.Error())
		return
	}
	return
}
func ListAnchorGuard(anchor, index int) (retMap []map[string]string, ret int) {
	retMap = make([]map[string]string, 0)
	//rowArray, err := orm.Query("SELECT b.uid,b.`nick_name`,b.`image` FROM  (SELECT SUM(VALUE) AS weight,send_user  FROM go_gift_record WHERE  send_user IN (SELECT uid FROM  go_guard WHERE anchor_id=? )  GROUP BY send_user) a  LEFT JOIN go_user b ON a.send_user=b.uid    ORDER BY a.weight limit ?,10", anchor,index*10)
	rowArray, err := orm.Query("SELECT q.uid,p.`nick_name`,p.`image`  FROM (SELECT SUM(b.value) AS weight,a.uid  FROM go_guard a LEFT JOIN go_gift_record b ON a.uid=b.send_user WHERE a.anchor_id=? and a.finish_time>UNIX_TIMESTAMP() GROUP BY a.uid )  q LEFT JOIN go_user p ON q.uid=p.uid limit ?,10", anchor, index*10)
	if err != nil {
		common.Log.Errf("orm is error:  %s", err.Error())
		ret = common.ERR_UNKNOWN
		return
	}

	for _, row := range rowArray {
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			ss[colName] = value
		}
		retMap = append(retMap, ss)
	}
	ret = common.ERR_SUCCESS
	return
}

func ListSelfGuard(anchor, uid, index int) (retMap []map[string]string, ret int) {
	retMap = make([]map[string]string, 0)
	//rowArray, err := orm.Query("SELECT b.uid,b.`nick_name`,b.`image` FROM  (SELECT SUM(VALUE) AS weight,send_user  FROM go_gift_record WHERE  send_user IN (SELECT uid FROM  go_guard WHERE anchor_id=? and uid!=?)  GROUP BY send_user) a  LEFT JOIN go_user b ON a.send_user=b.uid    ORDER BY a.weight limit ?,10", anchor,uid,index*10)

	rowArray, err := orm.Query("SELECT q.uid,p.`nick_name`,p.`image`  FROM (SELECT SUM(b.value) AS weight,a.uid  FROM go_guard a LEFT JOIN go_gift_record b ON a.uid=b.send_user WHERE a.uid!=? AND a.anchor_id=? and a.finish_time>UNIX_TIMESTAMP() GROUP BY a.uid )  q LEFT JOIN go_user p ON q.uid=p.uid limit ?,10", uid, anchor, index*10)
	if err != nil {
		common.Log.Errf("orm is error:  %s", err.Error())
		ret = common.ERR_UNKNOWN
		return
	}

	for _, row := range rowArray {
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			ss[colName] = value
		}
		retMap = append(retMap, ss)
	}
	ret = common.ERR_SUCCESS
	return
}
