package model

import (
	"container/list"
	"fmt"
	"github.com/yshd_game/common"
	"sort"
	"time"
)

type Letter struct {
	//Id    int64
	SessionId        int `xorm:"int(11) pk not null autoincr"`         //会话ID
	User1            int `xorm:"int(11) not null UNIQUE(LETTER_SESS)"` //用户ID userID必须>user2ID
	User2            int `xorm:"int(11) not null UNIQUE(LETTER_SESS)"` //用户ID
	Unread1          int `xorm:"int(11) not null default(0)"`          //记录user2给user1发送未读消息数量
	Unread2          int `xorm:"int(11) not null default(0)"`          //记录user1给user2发送未读消息数量
	Read             bool
	LetterSysVersion int64
}

//添加外键约束
//ALTER TABLE msg ADD CONSTRAINT c_sess FOREIGN KEY(session_id) REFERENCES letter(session_id);
//ALTER TABLE letter  ADD UNIQUE KEY(user1, user2);
type LetterMsg struct {
	Id          int64     `xorm:"int(11) pk not null autoincr"`
	SessionId   int       `xorm:"int(11) not null"` //会话ID
	IsSend      bool      //是否是主动发送方
	MessageBody string    //消息内容
	DelFlag1    bool      `xorm:"int(11) not null default(0)"` //user1 删除标志位
	DelFlag2    bool      `xorm:"int(11) not null default(0)"` //user2删除标志位
	CreateTime  time.Time //创建时间
}

type LetterMsgSys struct {
	Id               int64
	SendId           int    `xorm:"int(11) not null"`
	MessageBody      string //消息内容
	LetterSysVersion int64  //创建时间同时也是版本号
}

type LetterMsgResq struct {
	Id          int64 `xorm:"int(11) pk not null autoincr"`
	SessionId   int   `xorm:"int(11) not null"`
	IsSend      bool
	MessageBody string
	CreateTime  time.Time
}

func (self *LetterMsgResq) SetByOut(msg *LetterMsg) {
	self.Id = msg.Id
	self.SessionId = msg.SessionId
	self.IsSend = msg.IsSend
	self.MessageBody = msg.MessageBody
	self.CreateTime = msg.CreateTime
}

type LetterManager struct {
	letter list.List
}

func SendLetterToFamliy(send_id int, msg string, family_id int) int {
	res, err := orm.Query("select * from php_family where id=? ", family_id)
	if err != nil {
		common.Log.Errf("mysql error is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if len(res) == 0 {
		return common.ERR_CONFGI_ITEM
	}
	var admin_id int
	b, ok := res[0]["admin_id"]
	if ok {
		admin_id = common.BytesToInt(b)
	}

	msgarr := make([]*LetterMsg, 0)
	res2, err := orm.Query("SELECT a.uid,b.session_id FROM go_user a LEFT JOIN go_letter b ON a.uid=b.user1 WHERE a.admin_id=? AND b.user2=?", admin_id, send_id)
	if err != nil {
		common.Log.Errf("mysql error is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	for _, v := range res2 {
		m := &LetterMsg{}

		var uid int
		b, ok := v["uid"]
		if ok {
			uid = common.BytesToInt(b)
		} else {
			continue
		}

		b, ok = v["session_id"]
		if ok {
			if len(b) != 0 {
				m.SessionId = common.BytesToInt(b)
			} else {
				_, err := orm.InsertOne(&Letter{User1: uid, User2: send_id, Read: false})
				if err != nil {
					common.Log.Errf("mysql error is %s", err.Error())
					return common.ERR_UNKNOWN
				}

				s := &Letter{}
				_, err = orm.Where("user1=? and user2=?", uid, send_id).Get(s)
				if err != nil {
					common.Log.Errf("orm err is %s", err.Error())
					return common.ERR_UNKNOWN
				}
				m.SessionId = s.SessionId
			}
		} else {
			_, err := orm.InsertOne(&Letter{User1: uid, User2: send_id, Read: false})
			if err != nil {
				common.Log.Errf("mysql error is %s", err.Error())
				return common.ERR_UNKNOWN
			}

			s := &Letter{}
			_, err = orm.Where("user1=? and user2=?", uid, send_id).Get(s)
			if err != nil {
				common.Log.Errf("orm err is %s", err.Error())
				return common.ERR_UNKNOWN
			}

			m.SessionId = s.SessionId
		}

		m.IsSend = false
		m.MessageBody = msg
		m.CreateTime = time.Now()
		msgarr = append(msgarr, m)
	}

	aff_row, err := orm.Insert(&msgarr)
	if err != nil {
		common.Log.Errf("mysql error is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		return common.ERR_DB_ADD
	}
	return common.ERR_SUCCESS
}

func SendLetterToAnchor(send_id int, msg string) int {
	msgarr := make([]*LetterMsg, 0)
	res, err := orm.Query("SELECT  a.uid,c.session_id  FROM php_anchor a  LEFT JOIN (SELECT a.uid ,b.session_id ,b.user1   FROM php_anchor a LEFT JOIN go_letter b ON a.uid=b.user1 WHERE  b.user2=1) c ON a.uid=c.user1 ")
	if err != nil {
		common.Log.Errf("mysql error is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	for _, v := range res {
		m := &LetterMsg{}

		var uid int
		b, ok := v["uid"]
		if ok {
			uid = common.BytesToInt(b)
		}

		b, ok = v["session_id"]
		if ok {
			if len(b) != 0 {
				m.SessionId = common.BytesToInt(b)
			} else {
				_, err := orm.InsertOne(&Letter{User1: uid, User2: send_id, Read: false})
				if err != nil {
					common.Log.Errf("mysql error is %s", err.Error())
					return common.ERR_UNKNOWN
				}

				s := &Letter{}
				_, err = orm.Where("user1=? and user2=?", uid, send_id).Get(s)
				if err != nil {
					common.Log.Errf("orm err is %s", err.Error())
					return common.ERR_UNKNOWN
				}

				m.SessionId = s.SessionId
			}
		} else {
			_, err := orm.InsertOne(&Letter{User1: uid, User2: send_id, Read: false})
			if err != nil {
				common.Log.Errf("mysql error is %s", err.Error())
				return common.ERR_UNKNOWN
			}

			s := &Letter{}
			_, err = orm.Where("user1=? and user2=?", uid, send_id).Get(s)
			if err != nil {
				common.Log.Errf("orm err is %s", err.Error())
				return common.ERR_UNKNOWN
			}

			m.SessionId = s.SessionId
		}

		m.IsSend = false
		m.MessageBody = msg
		m.CreateTime = time.Now()
		msgarr = append(msgarr, m)
	}

	aff_row, err := orm.Insert(&msgarr)
	if err != nil {
		common.Log.Errf("mysql error is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		return common.ERR_DB_ADD
	}
	return common.ERR_SUCCESS
}

func SendLetterToAllV2(send_id int, msg string) int {
	m := &LetterMsgSys{
		SendId:           send_id,
		MessageBody:      msg,
		LetterSysVersion: time.Now().Unix(),
	}
	_, err := orm.InsertOne(m)
	if err != nil {
		common.Log.Errf("mysql error is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	return common.ERR_SUCCESS
}

func SendLetterToAll(send_id int, msg string) int {
	msgarr := make([]*LetterMsg, 0)
	res, err := orm.Query("SELECT a.uid,c.session_id,c.user1  FROM go_user a  LEFT JOIN  (SELECT a.uid ,b.session_id ,b.user1   FROM go_user a LEFT JOIN go_letter b ON a.uid=b.user2  WHERE a.uid=1 ) c ON a.uid=c.user1 ")
	if err != nil {
		common.Log.Errf("mysql error is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	for _, v := range res {
		m := &LetterMsg{}

		var uid int
		b, ok := v["uid"]
		if ok {
			uid = common.BytesToInt(b)
		}

		b, ok = v["session_id"]
		if ok {
			if len(b) != 0 {
				m.SessionId = common.BytesToInt(b)
			} else {
				_, err := orm.InsertOne(&Letter{User1: uid, User2: send_id, Read: false})
				if err != nil {
					common.Log.Errf("mysql error is %s", err.Error())
					return common.ERR_UNKNOWN
				}

				s := &Letter{}
				_, err = orm.Where("user1=? and user2=?", uid, send_id).Get(s)
				if err != nil {
					common.Log.Errf("orm err is %s", err.Error())
					return common.ERR_UNKNOWN
				}

				m.SessionId = s.SessionId
			}
		} else {
			_, err := orm.InsertOne(&Letter{User1: uid, User2: send_id, Read: false})
			if err != nil {
				common.Log.Errf("mysql error is %s", err.Error())
				return common.ERR_UNKNOWN
			}

			s := &Letter{}
			_, err = orm.Where("user1=? and user2=?", uid, send_id).Get(s)
			if err != nil {
				common.Log.Errf("orm err is %s", err.Error())
				return common.ERR_UNKNOWN
			}

			m.SessionId = s.SessionId
		}

		m.IsSend = false
		m.MessageBody = msg
		m.CreateTime = time.Now()
		msgarr = append(msgarr, m)
	}

	aff_row, err := orm.Insert(&msgarr)
	if err != nil {
		common.Log.Errf("mysql error is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		return common.ERR_DB_ADD
	}
	return common.ERR_SUCCESS
}

func SendLetter(send_id, rev_id int, msg string) int {
	letter := &Letter{}
	var user1, user2 int
	var csend bool

	if send_id == 0 {
		return common.ERR_PARAM
	}
	if rev_id == 0 {
		return common.ERR_PARAM
	}

	if send_id == rev_id {
		return common.ERR_LETTER_TO_SELF
	}
	has, err := orm.Where("owner_id=? and black_id=?  ", rev_id, send_id).Get(&Black{})
	if err != nil {
		common.Log.Errf("mysql error is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if has {
		return common.ERR_BLACK_IN
	}
	if send_id > rev_id {
		user1 = send_id
		user2 = rev_id
		csend = true
	} else {
		user1 = rev_id
		user2 = send_id
		csend = false
	}
	has, err = orm.Where("user1=? and user2=?", user1, user2).Get(letter)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	if !has {
		aff_row, err := orm.InsertOne(&Letter{User1: user1, User2: user2, Read: false})
		if err != nil {
			common.Log.Errf("mysql error is %s", err.Error())
			return common.ERR_UNKNOWN
		}

		if aff_row == 0 {
			return common.ERR_DB_ADD
		}
		_, err = orm.Where("user1=? and user2=?", user1, user2).Get(letter)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}

	}
	aff_row, err := orm.InsertOne(&LetterMsg{SessionId: letter.SessionId, IsSend: csend, MessageBody: msg, CreateTime: time.Now()})
	if err != nil {
		common.Log.Errf("mysql error is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		return common.ERR_DB_ADD
	}
	return common.ERR_SUCCESS
}

func CheckNewSysMsg(uid int) {
	sys := &LetterMsgSys{}
	has, err := orm.Desc("letter_sys_version").Limit(1, 0).Get(sys)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return
	}
	if has {
		s := &Letter{}
		has, err := orm.Where("user1=? and user2=1", uid).Get(s)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return
		}
		if has == false {
			SendLetter(1, uid, sys.MessageBody)

			has, err := orm.Where("user1=? and user2=1", uid).Get(s)
			if err != nil {
				common.Log.Errf("orm err is %s", err.Error())
				return
			}
			if has {
				s.LetterSysVersion = sys.LetterSysVersion
				_, err = orm.Where("session_id", s.SessionId).Cols("letter_sys_version").Update(s)
				if err != nil {
					common.Log.Errf("orm err is %s", err.Error())
					return
				}
			}

		} else {
			if s.LetterSysVersion < sys.LetterSysVersion {
				arr := make([]LetterMsgSys, 0)
				err := orm.Where("letter_sys_version > ?", s.LetterSysVersion).Find(&arr)

				if err != nil {
					common.Log.Errf("orm err is %s", err.Error())
					return
				}

				letter_arr := make([]LetterMsg, 0)
				for _, v := range arr {
					m := LetterMsg{}
					m.SessionId = s.SessionId
					m.IsSend = false
					m.MessageBody = v.MessageBody

					tm := time.Unix(v.LetterSysVersion, 0)
					m.CreateTime = tm
					letter_arr = append(letter_arr, m)
				}
				_, err = orm.Insert(&letter_arr)
				if err != nil {
					common.Log.Errf("orm err is %s", err.Error())
					return
				}
				s.LetterSysVersion = sys.LetterSysVersion
				_, err = orm.Where("session_id=?", s.SessionId).Cols("letter_sys_version").Update(s)
				if err != nil {
					common.Log.Errf("orm err is %s", err.Error())
					return
				}

			}
		}

	}
}
func ShowLetterList(uid, index int) ([]map[string]string, int) {
	//sql := fmt.Sprintf("SELECT MAX(m.create_time) AS record_time, m.sess AS session_id,m.nick_name,m.image,m.unread,m.uid,  m.message_body   FROM  (SELECT  sess, nick_name,image,uid,unread,create_time,message_body FROM  (SELECT a.session_id AS sess,b.nick_name,b.image,b.uid,a.unread1 AS unread FROM go_letter a  LEFT JOIN   go_user  b  ON  a.user2=b.uid WHERE a.user1= %d ) f INNER JOIN go_letter_msg g ON f.sess=g.session_id  WHERE g.`del_flag1` =0 UNION ALL SELECT sess, nick_name,image,uid,unread,create_time,message_body FROM  (SELECT c.session_id AS sess ,d.nick_name,d.image,d.uid,c.unread2 AS unread FROM go_letter c LEFT JOIN  go_user  d ON c.user1=d.uid WHERE c.user2=%d ) k INNER JOIN go_letter_msg l ON k.sess=l.session_id  WHERE l.`del_flag2` =0  ORDER BY create_time DESC) m GROUP BY session_id  ORDER BY record_time DESC  limit %d,%d", uid, uid, (index)*common.MSG_LIST_PAGER_COUNT, common.MSG_LIST_PAGER_COUNT)

	CheckNewSysMsg(uid)

	rowArray, err := orm.Query("SELECT r.record_time,r.session_id,r.nick_name,r.image, r.unread,r.uid, r.message_body ,t.type FROM   (SELECT MAX(m.create_time) AS record_time, m.sess AS session_id,m.nick_name,m.image,m.unread,m.uid,  m.message_body   FROM  (SELECT  sess, nick_name,image,uid,unread,create_time,message_body FROM  (SELECT a.session_id AS sess,b.nick_name,b.image,b.uid,a.unread1 AS unread FROM go_letter a  LEFT JOIN   go_user  b  ON  a.user2=b.uid WHERE a.user1= ? ) f INNER JOIN go_letter_msg g ON f.sess=g.session_id  WHERE g.`del_flag1` =0 UNION ALL SELECT sess, nick_name,image,uid,unread,create_time,message_body FROM  (SELECT c.session_id AS sess ,d.nick_name,d.image,d.uid,c.unread2 AS unread FROM go_letter c LEFT JOIN  go_user  d ON c.user1=d.uid WHERE c.user2=? ) k INNER JOIN go_letter_msg l ON k.sess=l.session_id  WHERE l.`del_flag2` =0  ORDER BY create_time DESC) m GROUP BY session_id) r LEFT JOIN php_identity t ON r.uid=t.uid  ORDER BY r.record_time DESC  limit ?,?", uid, uid, (index)*common.MSG_LIST_PAGER_COUNT, common.MSG_LIST_PAGER_COUNT)
	//rowArray, err := orm.Query(sql)
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

type LetterWrapper struct {
	letter []LetterMsg
	by     func(p, q *LetterMsg) bool
}

func (pw LetterWrapper) Len() int { // 重写 Len() 方法
	return len(pw.letter)
}
func (pw LetterWrapper) Swap(i, j int) { // 重写 Swap() 方法
	pw.letter[i], pw.letter[j] = pw.letter[j], pw.letter[i]
}
func (pw LetterWrapper) Less(i, j int) bool { // 重写 Less() 方法
	return pw.by(&pw.letter[i], &pw.letter[j])
}

func ShowLetterDeatil(session_id, index, uid, oid int) ([]LetterMsg, int) {
	msg := make([]LetterMsg, 0)
	if uid == 0 {
		return nil, common.ERR_PARAM
	}
	if oid == 0 {
		return nil, common.ERR_PARAM
	}
	min := index * common.MSG_LIST_PAGER_COUNT

	var del string
	var filed string
	if uid > oid {
		del = "del_flag1"
		filed = "unread1"
	} else {
		del = "del_flag2"
		filed = "unread2"
	}

	var err error
	if del == "del_flag1" {

		err = orm.Where("session_id =? and del_flag1=false", session_id).OrderBy("create_time desc").Limit(common.MSG_LIST_PAGER_COUNT, min).Find(&msg)
	} else if del == "del_flag2" {
		err = orm.Where("session_id =? and del_flag2=false", session_id).OrderBy("create_time desc").Limit(common.MSG_LIST_PAGER_COUNT, min).Find(&msg)
	}

	//err := orm.Where("session_id =? and ?=false", session_id,del).Limit(common.MSG_LIST_PAGER_COUNT, min).Find(&msg)
	if err != nil {
		common.Log.Errf("mysql error is %s", err.Error())
		return msg, common.ERR_UNKNOWN
	}
	if len(msg) == 0 {
		return msg, common.ERR_SUCCESS
	}
	if uid < oid {
		for i := 0; i < len(msg); i++ {
			if msg[i].IsSend {
				msg[i].IsSend = false
			} else {
				msg[i].IsSend = true
			}

		}
	}

	l := &Letter{}
	_, err = orm.Where("session_id=?", session_id).Get(l)
	if err != nil {
		common.Log.Errf("mysql error is %s", err.Error())
		return msg, common.ERR_UNKNOWN
	}
	if filed == "unread1" {
		if l.Unread1 != 0 {
			aff_row, err2 := orm.Where("session_id=?", session_id).MustCols(filed).SetExpr(filed, "0").Update(Letter{})
			if err2 != nil {
				common.Log.Errf("mysql error is %s", err2.Error())
				return msg, common.ERR_UNKNOWN
			}
			if aff_row == 0 {
				return msg, common.ERR_DB_UPDATE
			}
		}
	} else {
		if l.Unread2 != 0 {
			aff_row, err2 := orm.Where("session_id=?", session_id).MustCols(filed).SetExpr(filed, "0").Update(Letter{})
			if err2 != nil {
				common.Log.Errf("mysql error is %s", err2.Error())
				return msg, common.ERR_UNKNOWN
			}

			if aff_row == 0 {
				return msg, common.ERR_DB_UPDATE
			}
		}
	}

	sort.Sort(LetterWrapper{msg, func(p, q *LetterMsg) bool {
		return q.Id < p.Id //  递减排序
	}})
	return msg, common.ERR_SUCCESS
}

func GetLetterSessionId(uid, oid int) (int, bool) {
	letter := &Letter{}
	var user1, user2 int
	if uid > oid {
		user1 = uid
		user2 = oid
	} else {
		user1 = oid
		user2 = uid
	}
	has, err := orm.Where("user1=? and user2=?", user1, user2).Get(letter)
	if err != nil {
		common.Log.Errf("mysql error is %s", err.Error())
		return 0, false
	}
	if has {
		return letter.SessionId, true
	} else {
		letter.User1 = user1
		letter.User2 = user2
		letter.Read = false
		_, err2 := orm.InsertOne(letter)
		if err2 != nil {
			common.Log.Errf("mysql error is %s", err2.Error())
			return 0, false
		}
		return letter.SessionId, true
	}
	return 0, false
}

func DelSessionById(uid, oid int) int {

	//var del int
	var delstr string
	if uid > oid {
		delstr = "del_flag1"
	} else {
		delstr = "del_flag2"
	}

	sid, flag := GetLetterSessionId(uid, oid)
	if flag {
		_, err := orm.Where("session_id=?", sid).SetExpr(delstr, "true").Update(LetterMsg{})

		if err != nil {
			common.Log.Errf("mysql error is %s", err.Error())
			return common.ERR_UNKNOWN
		}

		_, err = orm.Where("session_id=?  and del_flag1=true and del_flag2=true", sid).Delete(LetterMsg{})
		if err != nil {
			common.Log.Errf("mysql error is %s", err.Error())
			return common.ERR_UNKNOWN
		}

		return common.ERR_SUCCESS
	}
	return common.ERR_SESSION_EXIST
}

func GetLetterUnreadNum(uid int) (int, int) {
	sql := fmt.Sprintf("select sum(unread1) as count from go_letter where user1=%d", uid)
	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return common.ERR_UNKNOWN, 0
	}
	var count1, count2 int
	if len(rowArray) == 1 {
		count1 = common.BytesToInt(rowArray[0]["count"])
	}

	sql = fmt.Sprintf("select sum(unread2) as count from go_letter where user2=%d", uid)
	rowArray2, err2 := orm.Query(sql)
	if err2 != nil {
		common.Log.Errf("db err %s", err2.Error())
		return common.ERR_UNKNOWN, 0
	}

	if len(rowArray2) == 1 {
		count2 = common.BytesToInt(rowArray2[0]["count"])
	}
	return common.ERR_SUCCESS, count1 + count2
}
