package model

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/go-xorm/xorm"
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	"github.com/yshd_game/confdata"
	"math"
	"strconv"
	"time"
)

func CheckErr(err error) int {
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	return common.ERR_SUCCESS
}

//用户表结构
type User struct {
	Uid      int    `xorm:"int(11) pk not null autoincr"`  //用户ID
	Account  string `xorm:"varchar(256)  UNIQUE(ACCOUNT)"` //账户
	OpenId   string `xorm:"varchar(50)  UNIQUE(OPENID)"`   //微信OPENID
	Platform int    `xorm:"int(11) default(0)"`            //来源平台
	// SinaAccount      string `xorm:"varchar(255)  UNIQUE(SINAACCOUNT)"`   //sina账户
	// QqAccount        string `xorm:"varchar(255)  UNIQUE(QQACCOUNT)"`     //qq账户
	// WeixinAccount    string `xorm:"varchar(255)  UNIQUE(WEIXInACCOUNT)"` //微信账户
	Tel           string `xorm:"varchar(20)   UNIQUE(TEL) "`    //绑定电话
	NickName      string `xorm:"varchar(128)" UNIQUE(NICKNAME)` //昵称
	Pwd           string `xorm:"varchar(32) not null"`          //密码
	Sex           int    `xorm:"int(11) default(0)"`            //性别
	UserLevel     int    `xorm:"int(11) default(1)"`            //用户等级
	AnchorLevel   int    `xorm:"int(11) default(1)"`            //主播等级
	UserExp       int    `xorm:"int(11) default(0)"`            //用户经验
	AnchorExp     int    `xorm:"int(11) default(0)"`            //主播经验
	Image         string `xorm:"varchar(128)"`                  //头像图片
	Diamond       int    `xorm:"int(11) default(0)"`            //钻石
	Location      string `xorm:"varchar(30)"`                   //位置
	Focus         int    `xorm:"int(11) default(0)"`            //关注人数
	Fans          int    `xorm:"int(11) default(0)"`            //粉丝人数
	Coupons       int    `xorm:"int(11) default(0)"`            //米粒
	FrozenCoupons int    `xorm:"int(11) default(0)"`            //冻结的米粒
	ThawTime      string `xorm:"varchar(64)"`                   //冻结时间
	Score         int64  `xorm:" not null default(0)"`          //游戏分
	Moon          int64  `xorm:" not null default(0)"`

	Push       bool
	Signature  string `xorm:"varchar(255)"` //签名
	Token      string `xorm:"varchar(128)  UNIQUE(TOKEN) " `
	ExpireTime int64  `xorm:"int(11) default(0)" `
	//Cover            string `xorm:"varchar(128)"`       //封面
	//Live             string `xorm:"varchar(30)"`        //预直播流
	Robot bool `xorm:"int(11) default(0)"` //是否机器人
	//Statue           int    `xorm:"int(11) default(0)"` //直播状态
	//RoomId           string `xorm:"varchar(128)"`       //所在房间ID

	WatchId    int64 `xorm:"default(0)"` //当前记录的观看记录ID
	NewPay     bool  //是否是获得首冲奖励充值用户
	NewPayTime int64 `xorm:" default(0)" `
	GroupId    int   `xorm:"int(11) default(0)"` //用户组别ID
	Forbid     bool  //封直播开播权限
	ForbidTime int64 //封直播开播权限到期时间
	//CacheRealImage   string `xorm:"varchar(180)"` //实名认证图片
	AuthRealInfo     bool   //是否已经实名
	RegisterTime     int64  //注册时间
	AccountType      int    `xorm:"int(11) default(1)"` //账户类型
	ForbidPowers     int    `xorm:"int(11) default(0)"` //限制所有帐户权限
	ForbidPowersTime int64  `xorm:"bigint default(0)"`  //限制所有帐户权限到期时间
	RegisterChannel  string `xorm:"varchar(128)"`       //注册渠道
	Device           string `xorm:"varchar(128)"`       //注册设备唯一码
	RegisterFrom     int    `xorm:"int(11) default(0)"` // 注册平台（1-IOS,2-Android）
	CanLinkMic       int    `xorm:"int(11) default(0)"` // 是否可以连麦
	//UnionId          string `xorm:"varchar(60)"`
	AdminId int `xorm:"int(11) default(0)"` //外键关联admin表所属的管理
	//CommossionId int `xorm:"int(11) default(0)"` //对应分配规则
	Imei    string `xorm:"varchar(128)"` //注册imie码
	Version int    `xorm:"version" default(1)`
}

//赠送礼物记录
type GiftRecord struct {
	Id       int64
	SendUser int `xorm:"int(11) not null "`  //赠送者
	RevUser  int `xorm:"int(11) not null "`  //接收者
	GiftId   int `xorm:"int(11) not null"`   //礼物ID
	Num      int `xorm:"int(11) default(0)"` //数量
	Value    int `xorm:"int(11) default(0)"` //总价值
	//	CreateTime time.Time //赠送时间
	RecordTime int64
	AdminId    int    `xorm:"int(11)  default(0) not null "` //所属管理ID
	MoneyType  int    `xorm:"int(11)  default(0) not null "` //金钱类型
	RoomId     string `xorm:"varchar(80) `
}

//礼物赠送详细分配记录
type GiftAssignedDetail struct {
	Id           int64
	GiftRecordId int64
	Identity     int   `xorm:"int(11) not null "` //身份
	MoneyType    int   `xorm:"int(11) not null "` //金钱类型
	Num          int64 `xorm:"int(11) not null "` //金钱数量
}

type GagRecord struct {
	Id         int64
	Owner      int       `xorm:"int(11) not null "`                       //房主
	RoomId     string    `xorm:"varchar(255) index(INDEX_ROOM_USER)"`     //房间ID
	Uid        int       `xorm:"int(11) not null index(INDEX_ROOM_USER)"` //被禁言人
	CreateTime time.Time //禁言时间
}

//提现记录
type CashRecord struct {
	Id               int64
	OwnerId          int       `xorm:"int(11) not null "`  //提现用户
	CouponsBefore    int       `xorm:"int(11) default(0)"` // 本笔订单开始前米粒余额
	CouponsAfter     int       `xorm:"int(11) default(0)"` // 本笔订单开始后米粒余额
	Rice             int       `xorm:"int(11) default(0)"` // 消耗米粒
	Money            int       `xorm:"int(11) default(0)"` // 将要获得RMB
	Statue           int       `xorm:"int(11) default(0)"` // 提现状态
	CardNo           string    `xorm:"varchar(40)"`        //卡号
	Bank             string    `xorm:"varchar(30)"`        // 银行
	RealName         string    `xorm:"varchar(30)"`        //姓名
	CreateTime       time.Time //提现创建时间
	CashType         int       `xorm:"int(2) not null "` //提现类型（0为银行，1为微信）
	LastOprationTime time.Time //最后操作时间
	FinishTime       time.Time //完成时间
	CashNum          string    `xorm:"varchar(30)"` //提现单号
	ErrCodeDes       string    `xorm:"varchar(40)"` //提现错误描述
	Version          int       `xorm:"version" default(1)`
}

//月亮商城订单信息
type ItemRecord struct {
	Id             int64  //订单号
	OwnerId        int    `xorm:"int(11) not null "`      //发起者ID
	Moon           int    `xorm:"int(11) not null "`      //消耗的游戏币
	ItemName       string `xorm:"int(11) not null "`      //商品名字
	ItemId         int    `xorm:"varchar(128) not null "` //实物ID
	Status         int    `xorm:"int(11) not null "`      //状态
	Name           string `xorm:"varchar(40)"`            //收货人
	Tel            string `xorm:"varchar(40)"`            //电话
	Addr           string `xorm:"varchar(128)"`           //地址信息
	CreateTime     time.Time
	Title          string    `xorm:"varchar(128)"` //商品标题
	Icon           string    `xorm:"varchar(128)"` //商品图片
	Logistics      string    `xorm:"varchar(128)"` //物流公司
	TrackingNumber string    `xorm:"varchar(128)"` //订单号
	DeliveryTime   time.Time //发货时间
	Money          int       `xorm:"int(11) not null default(0)`
}

//月亮商城提现记录
type MoonCashRecord struct {
	Id           int64
	OwnerId      int       `xorm:"int(11) not null "`  //提现用户
	Moon         int64     `xorm:"int(11) default(0)"` // 消耗米粒
	Money        int       `xorm:"int(11) default(0)"` // 将要获得RMB
	Statue       int       `xorm:"int(11) default(0)"` // 提现状态
	CreateTime   time.Time //提现创建时间
	FinishTime   time.Time //完成时间
	CashNum      string    `xorm:"varchar(30)"`       //提现单号
	ErrCodeDes   string    `xorm:"varchar(40)"`       //提现错误描述
	OperatorId   int       `xorm:"int(11) not null "` //提现用户
	ItemRecordId int       `xorm:"int(20) not null "` //提现用户
}

type ExchangeToRiceRecord struct {
	Id         int64
	Uid        int   `xorm:"int(11) not null "` //用户
	Diamond    int   `xorm:" not null "`
	Score      int64 `xorm:" not null "`
	CreateTime int64
}

type NicknameResetRecord struct {
	Id            int `xorm:"int(11) pk not null autoincr"`
	Uid           int `xorm:"int(11) "`           //用户ID
	Score         int `xorm:"int(11) default(0)"` //
	OperationTime int `xorm:"int(11) default(0)"` //来源平台
}

func GetGagByUidAndRoomID(uid int, roomid string) *GagRecord {
	gag := &GagRecord{}
	has, err := orm.Where("uid=? and room_id=?", uid, roomid).Get(gag)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return nil
	}
	if has {
		return gag
	}
	return nil
}

func GetUserByUid(uid int) (*User, int) {
	user := &User{}
	has, err := orm.Where("uid=?", uid).Get(user)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return user, common.ERR_UNKNOWN
	}

	if has {
		return user, common.ERR_SUCCESS
	}
	return user, common.ERR_ACCOUNT_EXIST
}

/*
func GetUserByUid(uid int) (*User, bool) {
	user := &User{}
	has, err := orm.Where("uid=?", uid).Get(user)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())

	}
	return user, has
}
*/

func GetUserByUidStr(uid string) (*User, int) {
	uid_, _ := strconv.Atoi(uid)
	return GetUserByUid(uid_)
}

/*
func GetUserByUidStr(uid string) (*User, bool) {
	uid_, _ := strconv.Atoi(uid)
	return GetUserByUid(uid_)
}
*/
/*
func GetUserByAccount(acc string) (*User, bool) {
	user := &User{}
	has, err := orm.Where("account=?", acc).Get(user)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
	}
	return user, has
}
*/

func GetUserByAccountAndPlatfrom(acc string, platform int) (*User, bool) {
	user := &User{}
	has, err := orm.Where("account=? and platform=?", acc, platform).Get(user)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
	}
	return user, has
}

func GetUserByNickName(nickname string) (*User, bool) {
	user := &User{}
	has, err := orm.Where("nick_name=?", nickname).Get(user)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
	}
	return user, has
}

//func GetUserByNickName(nickname string) (*User, bool) {
func CheckGetUserByNickName(nickname string) bool {
	u := &User{}
	has, err := orm.Where("nick_name=?", nickname).Get(u)
	if err != nil {

		common.Log.Errf("orm err is %s", err.Error())
		return false
	}
	return has

}

func GetUserByTel(tel string) (*User, bool) {
	user := &User{}
	has, err := orm.Where("tel=?", tel).Get(user)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
	}
	return user, has
}

func GetUserByToken(stoken string) (*User, bool) {
	user := &User{}
	has, err := orm.Where("token=?", stoken).Get(user)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
	}
	return user, has
}

func (self *User) GenNewNickName() error {
	random_id := common.RadnomRange(1, 100)
	nick := fmt.Sprintf("大米_%d_%d", time.Now().Unix(), random_id)
	_, has := GetUserByNickName(nick)
	if has {
		return self.GenNewNickName()
	}
	_, err := self.SetNick(nick)
	return err
}

/*
func (self *User) Update( ) (err error) {
	_, err = orm.Where("uid=?", self.Uid).MustCols("sex").Update(self)
	if err != nil {
		common.Log.Errf("db error:", err.Error())
		return
	}
	return
}

func (self *User) UpdateFront(filed string) (err error) {
	_, err = orm.Where("uid=?", self.Uid).MustCols(filed).Update(self)
	if err != nil {
		common.Log.Errf("db error:", err.Error())
		return
	}
	return
}


func (self *User) UpdateByMustColS(filed string)  (err error,aff int64)  {
	_, err = orm.Where("uid=?", self.Uid).MustCols(filed).Update(self)
	if err != nil {
		common.Log.Errf("db error:", err.Error())
		return
	}
	return
}
*/
func (self *User) UpdateByColSWithSession(session *xorm.Session, filed ...string) (aff int64, err error) {
	aff, err = session.Where("uid=?", self.Uid).Cols(filed...).Update(self)
	if err != nil {
		common.Log.Errf("db error:", err.Error())
		return
	}
	return
}
func (self *User) UpdateByColS(filed ...string) (aff int64, err error) {

	aff, err = orm.Where("uid=?", self.Uid).Cols(filed...).Update(self)
	if err != nil {
		common.Log.Errf("db error:", err.Error())
		if len(filed) >= 0 {
			common.Log.Infof("user info filed %d,%s", len(filed), filed[0])
		}
		/*
			_,err=orm.Where("uid=?",self.Uid).Update(self)
			if err!=nil {
				return
			}
		*/
		return
	}
	return
}

//增加用户经验(连续逻辑过程必须用事务来处理)
func (self *User) AddUserExp(session *xorm.Session, exp int, use_session bool) (aff int64, err error) {
	self.UserExp += exp
	self.LevelUpUser()
	if use_session {
		aff, err = self.UpdateByColSWithSession(session, "user_exp", "user_level")
	} else {
		aff, err = self.UpdateByColS("user_exp", "user_level")
	}
	return
}

/*
//增加用户经验单独增加
func (self *User) AddUserExp(exp int,session xorm.Session) (aff int64,err error) {
	self.UserExp += exp
	self.LevelUpUser()
	aff,err=self.UpdateByColS("user_exp","user_level")
	if err != nil {
		common.Log.Errf("orm is err %s",err.Error())
		return
	}
	return
}
*/

//增加主播经验(事务来处理连续逻辑)
func (self *User) AddAnchorExpBySession(session *xorm.Session, exp int) (aff int64, err error) {
	self.AnchorExp += exp
	self.LevelUpAnchor()
	aff, err = self.UpdateByColSWithSession(session, "anchor_exp", "anchor_level")
	return
}

/*
//增加主播经验(单独一次处理调用没有连续逻辑过程)
func (self *User) AddAnchorExp(exp int) (aff int64,err error) {
	self.AnchorExp += exp
	self.LevelUpAnchor()
	aff,err= self.UpdateByColS("anchor_exp","anchor_level")
	/*
		if err=self.UpdateFront("anchor_exp");err!=nil{
			common.Log.Err(err)
			return err
		}

	_, err = orm.Where("uid=?", self.Uid).Incr("anchor_exp", exp).Update(self)
	if err != nil {
		common.Log.Err(err)
		return
	}

		for self.UserExp >= GetLevelExp(self.AnchorLevel) {
			self.LevelUpAnchor()
		}


	return
}
*/
//升级用户等级
func (self *User) LevelUpUser() {
	c, has := GetUserExpByLevel(self.UserLevel)
	if has {
		if self.UserExp >= c.Exp {
			self.UserLevel++
			self.LevelUpUser()
		}
	}
}

//升级主播等级
func (self *User) LevelUpAnchor() {
	c, has := GetAnchorById(self.AnchorLevel)
	if has {
		if self.AnchorExp >= c.Exp {
			self.AnchorLevel++
			self.LevelUpAnchor()
		}
	}
}

//修改昵称
func (self *User) SetNick(nickname string) (int64, error) {
	self.NickName = nickname
	return self.UpdateByColS("nick_name")
}

//修改性别
func (self *User) SetSex(sex int) (int64, error) {
	self.Sex = sex
	return self.UpdateByColS("sex")
}

//修改位置
func (self *User) SetLocation(location string) (int64, error) {
	self.Location = location
	return self.UpdateByColS("location")
}

//修改签名
func (self *User) SetSignature(signature string) (int64, error) {
	self.Signature = signature
	return self.UpdateByColS("signature")
}

//设置推送
func (self *User) SetPush(push bool) (int64, error) {
	self.Push = push
	return self.UpdateByColS("push")
}

//关注别人
func (self *User) FocusOtherPublic(other int) int {
	other_entry, ret := GetUserByUid(other)
	if ret != common.ERR_SUCCESS {
		return ret
	}

	ret = FocusOther(self.Uid, other)
	if ret == common.ERR_SUCCESS {
		sess := GetUserSessByUid(self.Uid)
		if sess != nil {
			res := ResponseFocus{
				MType:         common.MESSAGE_TYPE_FOCUS,
				Uid:           self.Uid,
				Uid2:          other,
				NickNameSelf:  self.NickName,
				NickNameOther: other_entry.NickName,
			}
			SendMsgToRoom(sess.Roomid, res)
		}
		ReportActionDate(other, "subscribe", self.Uid)

	}
	return ret
}

//取消关注
func (self *User) CancleFocusPublic(oid int) int {
	return CancleFocusByOid(self.Uid, oid)
}

//增加钻石
func (self *User) addDiamondV2(session *xorm.Session, num int, use_session bool) int {
	if common.MAX_MONEY_LIMIT-self.Diamond < num {
		return common.ERR_MONEY_LIMIT
	}
	var aff int64
	var err error
	if use_session {
		aff, err = session.Where("uid=?", self.Uid).Incr("diamond", num).Cols("diamond").Update(self)
	} else {
		aff, err = orm.Where("uid=?", self.Uid).Incr("diamond", num).Cols("diamond").Update(self)
	}

	if err != nil {
		common.Log.Errf("orm is err %s", err.Error())
		return common.ERR_UNKNOWN
	}

	if aff == 0 {
		return common.ERR_DB_UPDATE
	}
	self.Diamond += num
	return common.ERR_SUCCESS
	/*
		count, err := orm.Where("uid=?", self.Uid).Incr("diamond", num).Update(self)
		if err != nil {
			common.Log.Errf("db error:", err.Error())
			return common.ERR_UNKNOWN
		}
		if count == 1 {
			common.Log.Infof("user addDiamond() 成功了 uid=%d,time=%d", self.Uid, time.Now().Unix())
		} else {
			common.Log.Infof("user addDiamond() 出错了 uid=%d,time=%d", self.Uid, time.Now().Unix())
		}

		sql := fmt.Sprintf("SELECT uid,diamond FROM go_user WHERE uid = %d", self.Uid)
		rowArray, _ := orm.Query(sql)
		for _, row := range rowArray {
			ss := make(map[string]string)
			for colName, colValue := range row {
				value := common.BytesToString(colValue)
				ss[colName] = value

				if colName == "diamond" {
					common.Log.Infof("uid=%d当前diamond余额为:%s", self.Uid, colValue)
				}
			}
		}

		self.Diamond += num
		return common.ERR_SUCCESS
	*/
}

//减少钻石
func (self *User) subDiamondV2(session *xorm.Session, num int, use_session bool) int {
	var aff int64
	var err error
	if self.Diamond >= num {
		if use_session {
			aff, err = session.Where("uid=?", self.Uid).Decr("diamond", num).Cols("diamond").Update(self)
		} else {
			aff, err = orm.Where("uid=?", self.Uid).Decr("diamond", num).Cols("diamond").Update(self)
		}

		if err != nil {
			common.Log.Errf("db error:", err.Error())
			return common.ERR_UNKNOWN
		}
		if aff == 0 {
			return common.ERR_DB_UPDATE
		}
		self.Diamond -= num
		return common.ERR_SUCCESS
	}
	return common.ERR_DIAMOND
}

//增加米粒
func (self *User) addRiceV2(session *xorm.Session, num int, use_session bool) int {
	if common.MAX_MONEY_LIMIT-self.Coupons < num {
		return common.ERR_MONEY_MONEY_MORE
	}
	var aff int64
	var err error
	nowtime := common.GetFormartTime()
	if nowtime == self.ThawTime {
		self.FrozenCoupons += num
	} else {
		self.FrozenCoupons = num
		self.ThawTime = nowtime
	}

	if use_session {
		aff, err = session.Where("uid=?", self.Uid).Incr("coupons", num).Cols("coupons", "frozen_coupons", "thaw_time").Update(self)
	} else {
		aff, err = orm.Where("uid=?", self.Uid).Incr("coupons", num).Cols("coupons", "frozen_coupons", "thaw_time").Update(self)
	}

	if err != nil {
		common.Log.Errf("db error:", err.Error())
		return common.ERR_UNKNOWN
	}

	if aff == 0 {
		return common.ERR_DB_UPDATE
	}
	self.Coupons += num
	return common.ERR_SUCCESS
}

//减少米粒
func (self *User) subRiceV2(session *xorm.Session, num int, use_session bool) int {
	if self.Coupons < num {
		return common.ERR_MONEY_MONEY_LESS
	}
	var aff int64
	var err error
	/*
		res,err:=orm.Exec("update go_user set coupons=coupons-? where uid=?",num,self.Uid)
		if err != nil {
			common.Log.Errf("db error:", err.Error())
			return common.ERR_UNKNOWN
		}
		aff,err:=res.RowsAffected()
		if err != nil {
			common.Log.Errf("db error:", err.Error())
			return common.ERR_UNKNOWN
		}

		if aff==0 {
			return common.ERR_DB_UPDATE
		}
	*/
	if use_session {
		aff, err = session.Where("uid=?", self.Uid).Decr("coupons", num).Cols("coupons").Update(self)
	} else {
		aff, err = orm.Where("uid=?", self.Uid).Decr("coupons", num).Cols("coupons").Update(self)
	}

	if err != nil {
		common.Log.Errf("db error:", err.Error())
		return common.ERR_UNKNOWN
	}

	if aff == 0 {
		return common.ERR_DB_UPDATE
	}
	self.Coupons -= num
	return common.ERR_SUCCESS
}

//增加游戏币
func (self *User) addScoreV2(session *xorm.Session, num int64, use_session bool) int {
	var aff int64
	var err error
	if use_session {
		aff, err = session.Where("uid=?", self.Uid).Incr("score", num).Cols("score").Update(self)
	} else {
		aff, err = orm.Where("uid=?", self.Uid).Incr("score", num).Cols("score").Update(self)
	}

	if err != nil {
		common.Log.Errf("db error:", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff == 0 {
		return common.ERR_DB_UPDATE
	}
	self.Score += num
	return common.ERR_SUCCESS
}

//减少游戏币
func (self *User) subScoreV2(session *xorm.Session, num int64, use_session bool) int {
	var aff int64
	var err error
	if self.Score >= num {
		if use_session {
			aff, err = session.Where("uid=?", self.Uid).Decr("score", num).Cols("score").Update(self)
		} else {
			aff, err = orm.Where("uid=?", self.Uid).Decr("score", num).Cols("score").Update(self)
		}

		if err != nil {
			common.Log.Errf("db error:", err.Error())
			return common.ERR_UNKNOWN
		}
		if aff == 0 {
			return common.ERR_DB_UPDATE
		}
		self.Score -= num
		return common.ERR_SUCCESS
	}
	return common.ERR_MONEY_MONEY_LESS
}

//增加月亮
func (self *User) addMoon(num int64) int {
	aff, err := orm.Where("uid=?", self.Uid).Incr("moon", num).Cols("moon").Update(self)
	if err != nil {
		common.Log.Errf("db error:", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff == 0 {
		return common.ERR_DB_UPDATE
	}
	self.Moon += num
	return common.ERR_SUCCESS
}

//增加月亮
func (self *User) addMoonV2(session *xorm.Session, num int64, use_session bool) int {
	var aff int64
	var err error
	if use_session {
		aff, err = session.Where("uid=?", self.Uid).Incr("moon", num).Cols("moon").Update(self)

	} else {
		aff, err = orm.Where("uid=?", self.Uid).Incr("moon", num).Cols("moon").Update(self)
	}
	if err != nil {
		common.Log.Errf("db error:", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff == 0 {
		return common.ERR_DB_UPDATE
	}
	self.Moon += num
	return common.ERR_SUCCESS
}

//减少月亮
func (self *User) subMoonV2(session *xorm.Session, num int64, use_session bool) int {
	var aff int64
	var err error
	if self.Moon >= num {
		if use_session {
			aff, err = session.Where("uid=?", self.Uid).Decr("moon", num).Cols("moon").Update(self)
			if err != nil {
				common.Log.Errf("db error:", err.Error())
				return common.ERR_UNKNOWN
			}

		} else {
			aff, err = orm.Where("uid=?", self.Uid).Decr("moon", num).Cols("moon").Update(self)
			if err != nil {
				common.Log.Errf("db error:", err.Error())
				return common.ERR_UNKNOWN
			}
		}
		if aff == 0 {
			return common.ERR_DB_UPDATE
		}
		self.Moon -= num
		return common.ERR_SUCCESS
	}
	return common.ERR_MONEY_MOON_LESS
}

func (self *User) CheckMoney(mtype int, num int64) int {
	if num < 0 {
		return common.ERR_PARAM
	}

	if num == 0 {
		return common.ERR_SUCCESS
	}

	switch mtype {
	case common.MONEY_TYPE_DIAMOND:
		if self.Diamond < int(num) {
			return common.ERR_DIAMOND
		}
	case common.MONEY_TYPE_RICE:
		if self.Coupons < int(num) {
			return common.ERR_RICE
		}
	case common.MONEY_TYPE_SCORE:
		if self.Score < num {
			return common.ERR_MONEY_MONEY_LESS
		}
	case common.MONEY_TYPE_MOON:
		if self.Moon < num {
			return common.ERR_MONEY_MONEY_LESS
		}
	default:
		common.Log.Errf("err money type %d", mtype)
		return common.ERR_PARAM
	}
	return common.ERR_SUCCESS
}

//连续逻辑必须用事务减钱
func (self *User) DelMoney(session *xorm.Session, mtype int32, num int64, use_session bool) int {
	if num < 0 {
		return common.ERR_PARAM
	}
	if num == 0 {
		return common.ERR_SUCCESS
	}

	common.Log.Infof("user del money uid=%d,moneytype=%d,num=%d,time=%d", self.Uid, mtype, num, time.Now().Unix())

	switch mtype {
	case common.MONEY_TYPE_DIAMOND:
		return self.subDiamondV2(session, int(num), use_session)
	case common.MONEY_TYPE_RICE:
		return self.subRiceV2(session, int(num), use_session)
	case common.MONEY_TYPE_SCORE:
		return self.subScoreV2(session, num, use_session)
	case common.MONEY_TYPE_MOON:
		return self.subMoonV2(session, num, use_session)
	default:
		common.Log.Errf("err money type %d", mtype)
	}
	return common.ERR_PARAM
}

//增加金钱连续逻辑必须用事务
func (self *User) AddMoney(session *xorm.Session, mtype int32, num int64, use_session bool) int {
	if num < 0 {
		return common.ERR_PARAM
	}
	if num == 0 {
		return common.ERR_SUCCESS
	}

	common.Log.Infof("user add money uid=%d,moneytype=%d,num=%d,time=%d", self.Uid, mtype, num, time.Now().Unix())
	switch mtype {
	case common.MONEY_TYPE_DIAMOND:
		return self.addDiamondV2(session, int(num), use_session)
	case common.MONEY_TYPE_RICE:
		return self.addRiceV2(session, int(num), use_session)
	case common.MONEY_TYPE_SCORE:
		return self.addScoreV2(session, num, use_session)
	case common.MONEY_TYPE_MOON:
		return self.addMoonV2(session, num, use_session)
	default:
		common.Log.Errf("err money type %d", mtype)

	}
	return common.ERR_PARAM
}

//更新expire 时间
func (self *User) SetExpireTime() {
	self.ExpireTime = time.Now().Unix() + common.TOKEN_EXPIRE_TIME
	self.UpdateByColS("expire_time")
}

//增加米粒sky专用
func (self *User) addCoupons(num int) int {
	if common.MAX_MONEY_LIMIT-self.Coupons < num {
		return common.ERR_MONEY_MONEY_MORE
	}

	self.Coupons += num
	count, err := orm.Where("uid=?", self.Uid).MustCols("coupons").Update(self)
	if err != nil {
		common.Log.Errf("db error:", err.Error())
		return common.ERR_UNKNOWN
	}
	if count == 1 {
		common.Log.Infof("user addCoupons() 增加星星成功了 uid=%d,time=%d", self.Uid, time.Now().Unix())
	} else {
		common.Log.Infof("user addCoupons() 增加星星出错了 uid=%d,time=%d", self.Uid, time.Now().Unix())
	}

	sql := fmt.Sprintf("SELECT uid,coupons FROM go_user WHERE uid = %d", self.Uid)
	rowArray, _ := orm.Query(sql)
	for _, row := range rowArray {
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			ss[colName] = value

			if colName == "coupons" {
				common.Log.Infof("uid=%d当前coupons余额为:%s", self.Uid, colValue)
			}
		}
	}

	return common.ERR_SUCCESS
}

func (self *User) TipGiftV2(gid, num, anchorId int) int {
	gift, ok := GetGiftById(gid)
	if !ok {
		return common.ERR_GIFT_EXIST
	}
	if self.Uid == anchorId {
		return common.ERR_SEND_GIFT_SELF
	}
	other, ret := GetUserByUid(anchorId)
	if ret != common.ERR_SUCCESS {
		return ret
	}
	othersess := GetUserSessByUid(anchorId)
	if othersess == nil {
		return common.ERR_USER_OFFLINE
	}
	selfsess := GetUserSessByUid(self.Uid)
	if selfsess == nil {
		return common.ERR_USER_OFFLINE
	}

	if othersess.Roomid != selfsess.Roomid {

		return common.ERR_STAY_IN_SAME_CHAT
	}

	if common.AccountAuthSwitch == true {
		if self.AccountType == 1 && other.AccountType == 0 {
			return common.ERR_SEND_GIFT_ROBOT
		}

		if self.AccountType == 0 && other.AccountType == 1 {
			return common.ERR_SEND_GIFT_ROBOT
		}
	}

	allnum := gift.Price * 1
	ret, bind_user_id, bind_group_id := GetGroupId(other.AdminId)
	if ret != common.ERR_SUCCESS {
		return ret
	}
	commossion, ok := GetFamilyPercent(anchorId, gift.Category)
	if ok == false {
		return common.ERR_CONFGI_ITEM
	}
	var send_num, bind_num, sys_num float32

	send_num = float32(allnum) * commossion.OwnerPercent
	sys_num = float32(allnum) * commossion.SystemPercent

	auth := other.CheckAuthReal()
	if auth == false {
		return common.ERR_AUTH_REAL
	} else {
		if ret := self.CheckMoney(common.MONEY_TYPE_SCORE, int64(allnum)); ret != common.ERR_SUCCESS {
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

		ret := self.DelMoney(session, common.MONEY_TYPE_SCORE, int64(allnum), true)
		if ret != common.ERR_SUCCESS {
			session.Rollback()
			return ret
		}

		if self.AccountType == 0 || self.AccountType == 1 && other.AccountType == 1 {
			ret := other.AddMoney(session, common.MONEY_TYPE_MOON, int64(send_num), true)
			if ret != common.ERR_SUCCESS {
				session.Rollback()
				return common.ERR_UNKNOWN
			}

			aff, err := other.AddAnchorExpBySession(session, int(allnum))
			if err != nil || aff == 0 {
				session.Rollback()
				return common.ERR_UNKNOWN
			}

			if bind_group_id != 11 { //官方主播没有管理者分成
				bind_num = float32(allnum) * commossion.LeaderPercent
				bind_user, _ := GetUserByUid(bind_user_id)
				if bind_user == nil {
					session.Rollback()
					return common.ERR_BIND_USER
				}
				if bind_user_id == self.Uid {
					bind_user = self
				}

				ret := bind_user.AddMoney(session, common.MONEY_TYPE_MOON, int64(bind_num), true)
				if ret != common.ERR_SUCCESS {
					session.Rollback()
					return common.ERR_UNKNOWN
				}
			}
		}

		err = session.Commit()
		if err != nil {
			common.Log.Errf("err is %s", err.Error())
			return common.ERR_UNKNOWN
		}

		rl, err := orm.Exec("insert into `go_gift_record` (`send_user`,`rev_user`,`gift_id`,`num`,`value`,`record_time`,`admin_id`,`money_type`,`room_id`) values (?,?,?,1,?,?,?,?,?)", self.Uid, anchorId, gid, allnum, time.Now().Unix(), bind_user_id, common.MONEY_TYPE_SCORE, selfsess.Roomid)
		if err != nil {
			common.Log.Errf("orm is error:  %s", err.Error())
			return common.ERR_UNKNOWN
		}
		resId, err := rl.LastInsertId()
		if err != nil {
			common.Log.Errf("orm is error:  %s", err.Error())
			return common.ERR_UNKNOWN
		}

		r1 := &GiftAssignedDetail{GiftRecordId: resId, Identity: common.LIVE_IDENTITY_USER, MoneyType: common.MONEY_TYPE_MOON, Num: int64(send_num)}

		r2 := &GiftAssignedDetail{GiftRecordId: resId, Identity: common.LIVE_IDENTITY_SYS, MoneyType: common.MONEY_TYPE_MOON, Num: int64(sys_num)}

		r3 := &GiftAssignedDetail{GiftRecordId: resId, Identity: common.LIVE_IDENTITY_ADMIN, MoneyType: common.MONEY_TYPE_MOON, Num: int64(bind_num)}

		aff, err := orm.Insert(r1, r2, r3)
		if err != nil || aff == 0 {
			if err != nil {
				common.Log.Errf("err is %s", err.Error())
			}
			return common.ERR_UNKNOWN
		}

	}

	chat := GetChatRoom(selfsess.Roomid)

	if chat != nil {
		mutex_chat_guardv2.Lock()
		if self.AccountType == 1 && other.AccountType == 1 {
			chat.AddMoon(int(send_num))
		} else if self.AccountType == 0 {
			chat.AddMoon(int(send_num))
		}
		mutex_chat_guardv2.Unlock()
	}

	IncrAudience(selfsess.Roomid, self.Uid, gift.Price)

	para := &ResponseGift{
		MType:       common.MESSAGE_TYPE_GIFT,
		SendId:      self.Uid,
		SendName:    self.NickName,
		SendImage:   self.Image,
		RevId:       anchorId,
		RevName:     other.NickName,
		GiftId:      gid,
		GiftNum:     num,
		SendLevel:   self.UserLevel,
		GiftDynamic: gift.Dynamic,
		GiftName:    gift.Name,
	}
	GiftSay(selfsess.Roomid, para, self.Uid, self.Token)

	TriggerTask(confdata.TargetType_gift, self.Uid, 1)
	InsertLog(common.ACTION_TYPE_LOG_SEND_GIFT, self.Uid, "")
	ReportActionDate(other.Uid, "gift", self.Uid)
	return common.ERR_SUCCESS
}

//打赏游戏礼物
/*
func (self *User) TipGift(gid, num, anchorId int) int {
	gift, ok := GetGiftById(gid)
	if !ok {
		return common.ERR_GIFT_EXIST
	}
	if self.Uid == anchorId {
		return common.ERR_SEND_GIFT_SELF
	}
	other, ret := GetUserByUid(anchorId)
	if ret != common.ERR_SUCCESS {
		return ret
	}
	othersess := GetUserSessByUid(anchorId)
	if othersess == nil {
		return common.ERR_USER_OFFLINE
	}
	selfsess := GetUserSessByUid(self.Uid)
	if selfsess == nil {
		return common.ERR_USER_OFFLINE
	}

	if othersess.Roomid != selfsess.Roomid {

		return common.ERR_STAY_IN_SAME_CHAT
	}

	if common.AccountAuthSwitch == true {
		if self.AccountType == 1 && other.AccountType == 0 {
			return common.ERR_SEND_GIFT_ROBOT
		}

		if self.AccountType == 0 && other.AccountType == 1 {
			return common.ERR_SEND_GIFT_ROBOT
		}
	}

	allnum := gift.Price * 1

	var bind_user_id int

	session := orm.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		common.Log.Errf("orm is error:  %s", err.Error())
		session.Rollback()
		return common.ERR_UNKNOWN
	}

	auth := other.CheckAuthReal()
	if auth == false {
		return common.ERR_UNKNOWN
		ret := self.subScore(int64(allnum))
		if ret != common.ERR_SUCCESS {
			return common.ERR_MONEY_MONEY_LESS
		}
		other.addMoon(int64(allnum))
		orm.Insert(&GiftRecord{
			SendUser:   self.Uid,
			RevUser:    anchorId,
			GiftId:     gid,
			Num:        1,
			RecordTime: time.Now().Unix(),
			Value:      allnum,
			AdminId:    bind_user_id,
			MoneyType:  common.MONEY_TYPE_SCORE,
		})

		chat := GetChatRoom(selfsess.Roomid)
		if chat != nil {
			chat.AddMoon(allnum)
		}
	} else {
		ret, bind_user_id, bind_group_id := GetGroupId(other.AdminId)
		if ret != common.ERR_SUCCESS {
			return ret
		}

		commossion, ok := GetFamilyPercent(anchorId, gift.Category)
		if ok == false {
			return common.ERR_CONFGI_ITEM
		}

		var owner_num, bind_num float32
		if bind_group_id == 11 {
			ret = self.DelMoney(common.MONEY_TYPE_SCORE, int64(allnum))
			if ret != common.ERR_SUCCESS {
				session.Rollback()
				return ret
			}

			owner_num = float32(allnum) * commossion.OwnerPercent
			//测试帐号送礼不会收到钱
			if self.AccountType == 1 && other.AccountType == 1 {
				ret := other.AddMoney(common.MONEY_TYPE_MOON, int64(owner_num))
				if ret != common.ERR_SUCCESS {
					session.Rollback()
					return ret
				}
				err := other.AddAnchorExp(int(owner_num))
				if err != nil {
					session.Rollback()
					return common.ERR_UNKNOWN
				}
			} else if self.AccountType == 0 {
				ret := other.AddMoney(common.MONEY_TYPE_MOON, int64(owner_num))
				if ret != common.ERR_SUCCESS {
					session.Rollback()
					return ret
				}
				err := other.AddAnchorExp(int(owner_num))
				if err != nil {
					session.Rollback()
					return common.ERR_UNKNOWN
				}
			}

		} else {

			bind_user, _ := GetUserByUid(bind_user_id)

			if bind_user == nil {
				return common.ERR_BIND_USER
			}

			if bind_user_id == self.Uid {
				bind_user = self
			}

			ret = self.DelMoney(common.MONEY_TYPE_SCORE, int64(allnum))
			if ret != common.ERR_SUCCESS {
				session.Rollback()
				return ret
			}

			owner_num = float32(allnum) * commossion.OwnerPercent

			//测试帐号送礼不会收到钱
			if self.AccountType == 1 && other.AccountType == 1 {
				ret := other.AddMoney(common.MONEY_TYPE_MOON, int64(owner_num))
				if ret != common.ERR_SUCCESS {
					session.Rollback()
					return ret
				}
				err := other.AddAnchorExp(int(owner_num))
				if err != nil {
					session.Rollback()
					return common.ERR_UNKNOWN
				}
				bind_num = float32(allnum) * commossion.LeaderPercent
				ret = bind_user.AddMoney(common.MONEY_TYPE_MOON, int64(bind_num))
				if ret != common.ERR_SUCCESS {
					session.Rollback()
					return ret
				}
			} else if self.AccountType == 0 {
				ret := other.AddMoney(common.MONEY_TYPE_MOON, int64(owner_num))
				if ret != common.ERR_SUCCESS {
					session.Rollback()
					return ret
				}
				err := other.AddAnchorExp(int(owner_num))
				if err != nil {
					session.Rollback()
					return common.ERR_UNKNOWN
				}
				bind_num = float32(allnum) * commossion.LeaderPercent
				ret = bind_user.AddMoney(common.MONEY_TYPE_MOON, int64(bind_num))
				if ret != common.ERR_SUCCESS {
					session.Rollback()
					return ret
				}
			}

		}

		sys_num := float32(allnum) * commossion.SystemPercent

		rl, err := orm.Exec("insert into `go_gift_record` (`send_user`,`rev_user`,`gift_id`,`num`,`value`,`record_time`,`admin_id`,`money_type`,`room_id`) values (?,?,?,1,?,?,?,?,?)", self.Uid, anchorId, gid, allnum, time.Now().Unix(), bind_user_id, common.MONEY_TYPE_SCORE, selfsess.Roomid)

		if err != nil {
			session.Rollback()
			common.Log.Errf("orm is error:  %s", err.Error())
			return common.ERR_UNKNOWN
		}
		resId, err := rl.LastInsertId()
		if err != nil {
			session.Rollback()
			common.Log.Errf("orm is error:  %s", err.Error())
			return common.ERR_UNKNOWN
		}

		r1 := &GiftAssignedDetail{GiftRecordId: resId, Identity: common.LIVE_IDENTITY_USER, MoneyType: common.MONEY_TYPE_MOON, Num: int64(owner_num)}

		r2 := &GiftAssignedDetail{GiftRecordId: resId, Identity: common.LIVE_IDENTITY_SYS, MoneyType: common.MONEY_TYPE_MOON, Num: int64(sys_num)}

		r3 := &GiftAssignedDetail{GiftRecordId: resId, Identity: common.LIVE_IDENTITY_ADMIN, MoneyType: common.MONEY_TYPE_MOON, Num: int64(bind_num)}

		_, err = orm.Insert(r1, r2, r3)
		if err != nil {
			session.Rollback()
			common.Log.Errf("orm is error:  %s", err.Error())
			return common.ERR_UNKNOWN
		}

		chat := GetChatRoom(selfsess.Roomid)
		if chat != nil {
			if self.AccountType == 1 && other.AccountType == 1 {
				chat.AddMoon(int(owner_num))
			} else if self.AccountType == 0 {
				chat.AddMoon(int(owner_num))
			}
		}

	}

	IncrAudience(selfsess.Roomid, self.Uid, gift.Price)

	para := &ResponseGift{
		MType:       common.MESSAGE_TYPE_GIFT,
		SendId:      self.Uid,
		SendName:    self.NickName,
		SendImage:   self.Image,
		RevId:       anchorId,
		RevName:     other.NickName,
		GiftId:      gid,
		GiftNum:     num,
		SendLevel:   self.UserLevel,
		GiftDynamic: gift.Dynamic,
		GiftName:    gift.Name,
	}
	GiftSay(selfsess.Roomid, para, self.Uid, self.Token)

	TriggerTask(confdata.TargetType_gift, self.Uid, 1)
	InsertLog(common.ACTION_TYPE_LOG_SEND_GIFT, self.Uid, "")

	err = session.Commit()
	if err != nil {
		return common.ERR_UNKNOWN
	}
	return common.ERR_SUCCESS
}
*/
func (self *User) Test2() int {
	t := time.Now()
	cur_tm := t.AddDate(0, 0, -1)
	s := make([]RoomList, 0)
	err := orm.Where("create_time<? and statue=? and finish_time='0001-01-01 00:00:00'", cur_tm, common.ROOM_RESTART).Find(&s)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return common.ERR_UNKNOWN
	}

	for _, v := range s {
		timeNow := v.CreateTime.Format("2006-01-02 15:04:05")
		limit_time := v.CreateTime.Add(3600 * 3 * time.Second)
		u2 := limit_time.Format("2006-01-02 15:04:05")
		res, err := orm.Query("select * from go_room_list where owner_id=? and create_time>? and create_time<=? and statue=? order by create_time limit 1", v.OwnerId, timeNow, u2, common.ROOM_FINISH)
		if err != nil {
			common.Log.Errf("db err %s", err.Error())
			return common.ERR_UNKNOWN
		}
		if len(res) == 1 {
			b, _ := res[0]["create_time"]
			ctime := common.BytesToString(b)
			t, _ = time.Parse("2006-01-02 15:04:05", ctime)
			_, err := orm.Exec("update go_room_list set finish_time=? where room_id=?", t, v.RoomId)
			if err != nil {
				common.Log.Errf("db err %s", err.Error())
				return common.ERR_UNKNOWN
			}

		} else {
			_, err := orm.Exec("update go_room_list set finish_time=? where room_id=?", u2, v.RoomId)
			if err != nil {
				common.Log.Errf("db err %s", err.Error())
				return common.ERR_UNKNOWN
			}
		}
	}
	return 0
}

func (self *User) Test() int {

	/*
		session := orm.NewSession()
		defer session.Close()
		err := session.Begin()
		if err != nil {
			common.Log.Errf("orm is error:  %s", err.Error())
			session.Rollback()
			return common.ERR_UNKNOWN
		}

		num := 10
		aff, err := orm.Where("uid=?", self.Uid).Incr("coupons", num).Cols("coupons").Update(self)
		if err != nil {
			common.Log.Errf("db error:", err.Error())
			return common.ERR_UNKNOWN
		}
		if aff == 0 {
			return -1
		} else {
			self.Sex = 0
			aff, err := orm.Where("uid=?", self.Uid).Cols("sex").Update(self)
			if err != nil {
				common.Log.Errf("db error:", err.Error())
				return common.ERR_UNKNOWN
			}
			if aff == 0 {
				return -2
			}
			return 0
		}

	*/

	/*
		o, ret := GetUserByUid(3)
		if ret != common.ERR_SUCCESS {
			return ret
		}

		self.AccountType = 1
		o.AccountType = 0
		if self.AccountType == 0 || self.AccountType == 1 && o.AccountType == 1 {
			return 100
		}
		session := orm.NewSession()
		defer session.Close()
		err := session.Begin()
		allnum := 10
		if err != nil {
			common.Log.Errf("orm is error:  %s", err.Error())
			session.Rollback()
			return common.ERR_UNKNOWN
		}

		ret = o.DelMoney(session, common.MONEY_TYPE_RICE, 20, true)
		if ret != common.ERR_SUCCESS {
			session.Rollback()
			return ret
		}

		ret = self.AddMoney(session, common.MONEY_TYPE_RICE, int64(allnum), true)
		if ret != common.ERR_SUCCESS {
			session.Rollback()
			return ret
		}

		ret = self.AddMoney(session, common.MONEY_TYPE_RICE, int64(15), true)
		if ret != common.ERR_SUCCESS {
			session.Rollback()
			return ret
		}

		si, _ := GetUserByUid(4)
		si.Sex = 0
		aff, err := session.Where("uid=?", si.Uid).Cols("sex").Update(self)
		//aff,err:=self.SetSex(0)
		if err != nil || aff == 0 {
			if err != nil {
				common.Log.Errf("err is %s", err.Error())
			}
			session.Rollback()
			return common.ERR_UNKNOWN
		}

		err = session.Commit()
		if err != nil {
			return common.ERR_UNKNOWN
		}


	*/

	/*
		var i int

		i = 4
		for ; i < 7; i++ {
			time.Sleep(1 * time.Microsecond)
			go func(i int) {
				u, ret := GetUserByUid(i)
				if ret == common.ERR_SUCCESS {
					o, _ := GetUserByUid(3)
					msg := fmt.Sprintf("befor  uid=%d,time=%v,money=%d,other diamond=%d", i, time.Now().UnixNano()/1e6, u.Diamond, o.Coupons)
					godump.Dump(msg)

					ret := u.SendGiftV2(1, 1, 3)

					u, _ := GetUserByUid(i)
					o2, _ := GetUserByUid(3)
					msg = fmt.Sprintf("after uid=%d, time=%v,money=%d,other diamond=%d,ret:=%d", i, time.Now().UnixNano()/1e6, u.Diamond, o2.Coupons, ret)
					godump.Dump(msg)

				}
			}(i)
		}
		go func() {
			u, _ := GetUserByUid(11)
			msg := fmt.Sprintf("befor  uid=%d,time=%v,money=%d,other diamond=%d", u.Uid, time.Now().Unix(), u.Diamond)
			godump.Dump(msg)
			ret := u.SendGiftV2(11, 1, 3)

			msg = fmt.Sprintf("after uid=%d, time=%v,money=%d,ret:=%d", u.Uid, time.Now().Unix(), u.Diamond, ret)
			godump.Dump(msg)
		}()

	*/
	var s WinXinPayCallBackReq
	s.Transaction_id = "11"
	s.Return_msg = "err"
	s.Return_code = "ss"
	common.Log.Errf("%v", s)

	return common.ERR_SUCCESS

}

func (self *User) SendGiftV2(gid, num, otherid int) int {
	gift, ok := GetGiftById(gid)
	if !ok {
		return common.ERR_GIFT_EXIST
	}
	if self.Uid == otherid {
		return common.ERR_SEND_GIFT_SELF
	}
	other, ret := GetUserByUid(otherid)
	if ret != common.ERR_SUCCESS {
		return ret
	}

	othersess := GetUserSessByUid(otherid)
	if othersess == nil {
		return common.ERR_USER_OFFLINE
	}

	selfsess := GetUserSessByUid(self.Uid)
	if selfsess == nil {
		return common.ERR_USER_OFFLINE
	}
	if othersess.Roomid != selfsess.Roomid {
		return common.ERR_STAY_IN_SAME_CHAT
	}

	//selfsess:=&SessUser{}
	//selfsess.Roomid="1"

	if common.AccountAuthSwitch == true {
		if self.AccountType == 1 && other.AccountType == 0 {
			return common.ERR_SEND_GIFT_ROBOT
		}

		if self.AccountType == 0 && other.AccountType == 1 {
			return common.ERR_SEND_GIFT_ROBOT
		}
	}

	allnum := gift.Price * 1

	auth := other.CheckAuthReal()

	ret, bind_user_id, bind_group_id := GetGroupId(other.AdminId)
	if ret != common.ERR_SUCCESS {
		return ret
	}
	commossion, ok := GetFamilyPercent(otherid, gift.Category)
	if ok == false {
		return common.ERR_CONFGI_ITEM
	}
	var send_num, bind_num, sys_num float32

	send_num = float32(allnum) * commossion.OwnerPercent
	sys_num = float32(allnum) * commossion.SystemPercent

	if auth == false {
		return common.ERR_AUTH_REAL
	} else {

		if ret := self.CheckMoney(common.MONEY_TYPE_DIAMOND, int64(allnum)); ret != common.ERR_SUCCESS {
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

		ret := self.DelMoney(session, common.MONEY_TYPE_DIAMOND, int64(allnum), true)
		if ret != common.ERR_SUCCESS {
			common.Log.Errf("send filed del money uid=%d err=1", self.Uid)
			session.Rollback()
			return ret
		}

		if self.AccountType == 0 || self.AccountType == 1 && other.AccountType == 1 {
			ret := other.AddMoney(session, common.MONEY_TYPE_RICE, int64(send_num), true)
			if ret != common.ERR_SUCCESS {
				common.Log.Errf("send filed add money uid=%d err=2", self.Uid)
				session.Rollback()
				return common.ERR_UNKNOWN
			}

			aff, err := other.AddAnchorExpBySession(session, int(allnum))
			if err != nil || aff == 0 {
				common.Log.Errf("send filed add money uid=%d err=3", self.Uid)
				session.Rollback()
				return common.ERR_UNKNOWN
			}

			if bind_group_id != 11 { //官方主播没有管理者分成
				bind_num = float32(allnum) * commossion.LeaderPercent
				bind_user, _ := GetUserByUid(bind_user_id)
				if bind_user == nil {
					session.Rollback()
					return common.ERR_BIND_USER
				}
				if bind_user_id == self.Uid {
					bind_user = self
				}

				ret := bind_user.AddMoney(session, common.MONEY_TYPE_RICE, int64(bind_num), true)
				if ret != common.ERR_SUCCESS {
					common.Log.Errf("send filed del money uid=%d err=4", self.Uid)
					session.Rollback()
					return common.ERR_UNKNOWN
				}
			}
		}

		err = session.Commit()
		if err != nil {
			return common.ERR_UNKNOWN
		}

		rl, err := orm.Exec("insert into `go_gift_record` (`send_user`,`rev_user`,`gift_id`,`num`,`value`,`record_time`,`admin_id`,`money_type`,`room_id`) values (?,?,?,1,?,?,?,?,?)", self.Uid, otherid, gid, allnum, time.Now().Unix(), bind_user_id, common.MONEY_TYPE_DIAMOND, selfsess.Roomid)
		if err != nil {
			common.Log.Errf("orm is error:  %s", err.Error())
			//session.Rollback()
			return common.ERR_UNKNOWN
		}
		resId, err := rl.LastInsertId()
		if err != nil {
			//session.Rollback()
			common.Log.Errf("orm is error:  %s", err.Error())
			return common.ERR_UNKNOWN
		}

		r1 := &GiftAssignedDetail{GiftRecordId: resId, Identity: common.LIVE_IDENTITY_USER, MoneyType: common.MONEY_TYPE_RICE, Num: int64(send_num)}

		r2 := &GiftAssignedDetail{GiftRecordId: resId, Identity: common.LIVE_IDENTITY_SYS, MoneyType: common.MONEY_TYPE_RICE, Num: int64(sys_num)}

		r3 := &GiftAssignedDetail{GiftRecordId: resId, Identity: common.LIVE_IDENTITY_ADMIN, MoneyType: common.MONEY_TYPE_RICE, Num: int64(bind_num)}

		aff, err := orm.Insert(r1, r2, r3)
		if err != nil || aff == 0 {
			if err != nil {
				common.Log.Errf("err is %s", err.Error())
			}
			//session.Rollback()
			return common.ERR_UNKNOWN
		}

	}

	chat := GetChatRoom(selfsess.Roomid)
	if chat != nil {
		mutex_chat_guardv2.Lock()
		if self.AccountType == 1 && other.AccountType == 1 {
			chat.AddRice(int(send_num))
		} else if self.AccountType == 0 {
			chat.AddRice(int(send_num))
		}
		mutex_chat_guardv2.Unlock()
	}

	extre, has := GetUserExtraByUid(self.Uid)
	if has {
		extre.ConsumerStatistics(allnum)
	}

	IncrAudience(selfsess.Roomid, self.Uid, gift.Price)

	para := &ResponseGift{
		MType:       common.MESSAGE_TYPE_GIFT,
		SendId:      self.Uid,
		SendName:    self.NickName,
		SendImage:   self.Image,
		RevId:       otherid,
		RevName:     other.NickName,
		GiftId:      gid,
		GiftNum:     num,
		SendLevel:   self.UserLevel,
		GiftDynamic: gift.Dynamic,
		GiftName:    gift.Name,
	}
	/*
		m := &GiftRecord{
			SendUser:   self.Uid,
			RevUser:    otherid,
			GiftId:     gift.GiftId,
			Num:        1,
			Value:      gift.Price,
			RecordTime: time.Now().Unix(),
			AdminId:    bind_user_id,
			MoneyType:  common.MONEY_TYPE_DIAMOND,
			RoomId:     selfsess.Roomid,
		}
	*/

	GiftSay(selfsess.Roomid, para, self.Uid, self.Token)
	TriggerTask(confdata.TargetType_gift, self.Uid, 1)
	InsertLog(common.ACTION_TYPE_LOG_SEND_GIFT, self.Uid, "")

	ReportActionDate(otherid, "gift", self.Uid)
	return common.ERR_SUCCESS
}

//赠送礼物
/*
func (self *User) SendGift(gid, num, otherid int) int {
	gift, ok := GetGiftById(gid)
	if !ok {
		return common.ERR_GIFT_EXIST
	}
	if self.Uid == otherid {
		return common.ERR_SEND_GIFT_SELF
	}
	other, ret := GetUserByUid(otherid)
	if ret != common.ERR_SUCCESS {
		return ret
	}

	othersess := GetUserSessByUid(otherid)
	if othersess == nil {
		return common.ERR_USER_OFFLINE
	}

	selfsess := GetUserSessByUid(self.Uid)
	if selfsess == nil {
		return common.ERR_USER_OFFLINE
	}
	if othersess.Roomid != selfsess.Roomid {
		return common.ERR_STAY_IN_SAME_CHAT
	}

	if common.AccountAuthSwitch == true {
		if self.AccountType == 1 && other.AccountType == 0 {
			return common.ERR_SEND_GIFT_ROBOT
		}

		if self.AccountType == 0 && other.AccountType == 1 {
			return common.ERR_SEND_GIFT_ROBOT
		}
	}

	allnum := gift.Price * 1
	//var bind_user_id int
	auth := other.CheckAuthReal()

	session := orm.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		common.Log.Errf("orm is error:  %s", err.Error())
		session.Rollback()
		return common.ERR_UNKNOWN
	}
	if auth == false {
		return common.ERR_AUTH_REAL
	} else {
		ret, bind_user_id, bind_group_id := GetGroupId(other.AdminId)
		if ret != common.ERR_SUCCESS {
			return ret
		}

		commossion, ok := GetFamilyPercent(otherid, gift.Category)
		if ok == false {
			return common.ERR_CONFGI_ITEM
		}
		var owner_num, bind_num float32

		if bind_group_id == 11 {
			ret = self.DelMoney(common.MONEY_TYPE_DIAMOND, int64(allnum))
			if ret != common.ERR_SUCCESS {
				return ret
			}

			if self.AccountType == 1 && other.AccountType == 1 {
				ret := other.AddMoney(common.MONEY_TYPE_RICE, int64(owner_num))

				if ret != common.ERR_SUCCESS {
					session.Rollback()
					return ret
				}

				err := other.AddAnchorExp(int(owner_num))
				if err != nil {
					session.Rollback()
					return common.ERR_UNKNOWN
				}

			} else if self.AccountType == 0 {
				ret := other.AddMoney(common.MONEY_TYPE_RICE, int64(owner_num))

				if ret != common.ERR_SUCCESS {
					session.Rollback()
					return ret
				}

				err := other.AddAnchorExp(int(owner_num))

				if err != nil {
					session.Rollback()
					return common.ERR_UNKNOWN
				}

			} else {
				bind_user, _ := GetUserByUid(bind_user_id)
				if bind_user == nil {
					return common.ERR_BIND_USER
				}
				if bind_user_id == self.Uid {
					bind_user = self
				}

				ret = self.DelMoney(common.MONEY_TYPE_DIAMOND, int64(allnum))
				if ret != common.ERR_SUCCESS {
					session.Rollback()
					return ret
				}

				owner_num = float32(allnum) * commossion.OwnerPercent

				if self.AccountType == 1 && other.AccountType == 1 {
					ret = other.AddMoney(common.MONEY_TYPE_RICE, int64(owner_num))
					if ret != common.ERR_SUCCESS {
						session.Rollback()
						return ret
					}

					err = other.AddAnchorExp(int(owner_num))
					if err != nil {
						session.Rollback()
						return common.ERR_UNKNOWN
					}
					bind_num = float32(allnum) * commossion.LeaderPercent

					ret = bind_user.AddMoney(common.MONEY_TYPE_RICE, int64(bind_num))
					if ret != common.ERR_SUCCESS {
						session.Rollback()
						return ret
					}
				} else if self.AccountType == 0 {
					ret = other.AddMoney(common.MONEY_TYPE_RICE, int64(owner_num))
					if ret != common.ERR_SUCCESS {
						session.Rollback()
						return ret
					}

					err = other.AddAnchorExp(int(owner_num))
					if err != nil {
						session.Rollback()
						return common.ERR_UNKNOWN
					}
					bind_num = float32(allnum) * commossion.LeaderPercent

					ret = bind_user.AddMoney(common.MONEY_TYPE_RICE, int64(bind_num))
					if ret != common.ERR_SUCCESS {
						session.Rollback()
						return ret
					}
				}
			}
		}
		sys_num := float32(allnum) * commossion.SystemPercent

		rl, err := orm.Exec("insert into `go_gift_record` (`send_user`,`rev_user`,`gift_id`,`num`,`value`,`record_time`,`admin_id`,`money_type`,`room_id`) values (?,?,?,1,?,?,?,?,?)", self.Uid, otherid, gid, allnum, time.Now().Unix(), bind_user_id, common.MONEY_TYPE_DIAMOND, selfsess.Roomid)
		if err != nil {
			common.Log.Errf("orm is error:  %s", err.Error())
			session.Rollback()
			return common.ERR_UNKNOWN
		}
		resId, err := rl.LastInsertId()
		if err != nil {
			session.Rollback()
			common.Log.Errf("orm is error:  %s", err.Error())
			return common.ERR_UNKNOWN
		}

		r1 := &GiftAssignedDetail{GiftRecordId: resId, Identity: common.LIVE_IDENTITY_USER, MoneyType: common.MONEY_TYPE_RICE, Num: int64(owner_num)}

		r2 := &GiftAssignedDetail{GiftRecordId: resId, Identity: common.LIVE_IDENTITY_SYS, MoneyType: common.MONEY_TYPE_RICE, Num: int64(sys_num)}

		r3 := &GiftAssignedDetail{GiftRecordId: resId, Identity: common.LIVE_IDENTITY_ADMIN, MoneyType: common.MONEY_TYPE_RICE, Num: int64(bind_num)}

		_, err = orm.Insert(r1, r2, r3)
		if err != nil {
			session.Rollback()
			return common.ERR_UNKNOWN
		}
		err = session.Commit()
		if err != nil {
			return common.ERR_UNKNOWN
		}

		chat := GetChatRoom(selfsess.Roomid)
		if chat != nil {
			if self.AccountType == 1 && other.AccountType == 1 {
				chat.AddRice(int(owner_num))
			} else if self.AccountType == 0 {
				chat.AddRice(int(owner_num))
			}
		}
	}

	extre, has := GetUserExtraByUid(self.Uid)
	if has {
		extre.ConsumerStatistics(allnum)
	}

	IncrAudience(selfsess.Roomid, self.Uid, gift.Price)

	para := &ResponseGift{
		MType:       common.MESSAGE_TYPE_GIFT,
		SendId:      self.Uid,
		SendName:    self.NickName,
		SendImage:   self.Image,
		RevId:       otherid,
		RevName:     other.NickName,
		GiftId:      gid,
		GiftNum:     num,
		SendLevel:   self.UserLevel,
		GiftDynamic: gift.Dynamic,
		GiftName:    gift.Name,
	}

	GiftSay(selfsess.Roomid, para, self.Uid, self.Token)
	TriggerTask(confdata.TargetType_gift, self.Uid, 1)
	InsertLog(common.ACTION_TYPE_LOG_SEND_GIFT, self.Uid, "")

	return common.ERR_SUCCESS
}
*/

//设置token
func (self *User) SetToken(token string) (int64, error) {
	self.Token = token
	self.ExpireTime = time.Now().Unix() + common.TOKEN_EXPIRE_TIME

	return self.UpdateByColS("token", "expire_time")
}

//设置头像
func (self *User) SetFace(face string) (aff int64, err error) {
	if face == "" {
		return
	}
	willdelete := self.Image

	self.Image = face

	aff, err = self.UpdateByColS("image")
	if err != nil {
		return
	}
	/*
		_, err := orm.Where("uid=?", self.Uid).MustCols("Image").Update(self)
		if err != nil {
			common.Log.Err("set image error: , %s", err.Error())
			return
		}
	*/
	if willdelete != "" {
		DelQiNiuFile3(bucket_face, DomainFace, willdelete)
	}
	return
}

//设置新密码
func (self *User) SetPwd(pwd string) (int64, error) {
	self.Pwd = pwd

	return self.UpdateByColS("pwd")
	/*
		_, err := orm.Where("uid=?", self.Uid).MustCols("pwd").Update(self)
		if err != nil {
			common.Log.Err("set pwd error: , %s", err.Error())
		}
	*/
}

//重置密码
func (self *User) ReSetPwd(oldpwd, newpwd string) int {
	if self.Platform != common.PLATFORM_SELF {
		return common.ERR_THIRD_PWD
	}
	if self.Pwd == oldpwd {
		aff, err := self.SetPwd(newpwd)
		if err != nil || aff == 0 {
			return common.ERR_UNKNOWN
		}
		return common.ERR_SUCCESS
	}
	return common.ERR_PWD
}

//进入房间
func (self *User) JoinRoom(rid string, isAnchor bool, room_type int, ip string) {
	ret, wid := JoinChat(self.Uid, rid, room_type, ip)
	if !ret {
		common.Log.Infof("user join chat record uid=%d ,time=%d", self.Uid, time.Now().Unix())
	} else {
		self.WatchId = wid
	}

	e, ok := GetUserExtraByUid(self.Uid)
	if ok {
		e.SetNextOnlineRewardTime()
	}

	ResetGuard(self.Uid)
	aff, err := self.UpdateByColS("watch_id")
	if err != nil || aff == 0 {
		common.Log.Errf("db update valide uid=?,watch_id=? ", self.Uid, wid)
	}
}

//离开房间
func (self *User) LeaveRoom() {
	e, ok := GetUserExtraByUid(self.Uid)
	if ok {
		e.FinishOnlineRewardTime()
	}

	if self.WatchId == 0 {
		common.Log.Infof("user leavel chat record uid=%d", self.Uid)
		return
	}
	ret := LeaveChat(self.WatchId)

	if !ret {
		common.Log.Errf("user leavel chat record watch_id=%d ,time=%d", self.WatchId, time.Now().Unix())
	}

	//orm.Where("uid=?", self.Uid).Update()
	_, err := orm.Exec("update go_user set watch_id=? where uid=?", 0, self.Uid)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
	}

}

//检查账号封禁
func (self *User) CheckPowerAccount() bool {
	if self.ForbidPowers == common.FORBID_POWERS_FOREVER {
		return true
	}

	if self.ForbidPowers == common.FORBID_POWERS_TIME {
		now := time.Now().Unix()
		if now <= self.ForbidPowersTime {
			return true
		}
	}
	return false
}

//检查账号封停状态
func (self *User) CheckAccount() bool {
	if self.Forbid == true {
		now := time.Now().Unix()
		if now <= self.ForbidTime {
			return true
		}
	}
	return false
}

//检查实名状态
func (self *User) CheckAuthReal() bool {
	return self.AuthRealInfo
}

//设置实名状态
func (self *User) SetAuthReal() {
	self.AuthRealInfo = true
	self.UpdateByColS("auth_real_info")
}

//封禁账号
func (self *User) ForbidSelf() (int64, error) {
	if self.CheckAccount() {
		return 0, nil
	}
	self.Forbid = true
	self.ForbidTime = time.Now().Unix() + common.FORBID_ACCOUNT_KEEP_TIME
	return self.UpdateByColS("forbid", "forbid_time")
}

//获取封禁状态
func (self *User) CheckAccountForbid() int {
	switch self.ForbidPowers {
	case common.FOCUS_STATUE_NONE:
		return common.ERR_SUCCESS
	case common.FORBID_POWERS_FOREVER:
		return common.ERR_FORBID_EVER
	case common.FORBID_POWERS_TIME:
		if self.ForbidPowersTime > time.Now().Unix() {
			return common.ERR_FORBID_POWER
		} else {
			return common.ERR_SUCCESS
		}
	default:
		return common.ERR_UNKNOWN
	}
}

//设置首冲标志位
func (self *User) SetNewPay(session *xorm.Session, diamond int) int {
	if self.NewPay == true {
		return common.ERR_SUCCESS
	}
	self.NewPay = true
	self.NewPayTime = time.Now().Unix()

	send_socre := int(float32(diamond) * 0.1)
	ret := self.AddMoney(session, common.MONEY_TYPE_SCORE, int64(send_socre), true)
	if ret != common.ERR_SUCCESS {
		return ret
	}
	aff, err := self.UpdateByColSWithSession(session, "new_pay", "new_pay_time")
	if err != nil || aff == 0 {
		return common.ERR_UNKNOWN
	}
	msg := fmt.Sprintf("亲爱的用户%s:你于%s充值的%d钻石已经到账，系统赠送%d游戏币。祝你玩的愉快。", self.NickName, common.GetFormartTime2(), diamond, send_socre)
	SendLetter(1, self.Uid, msg)
	return common.ERR_SUCCESS
}

/*
//获得关注列表
func (self *User) GetFocusList(index int) ([]map[string]string, int) {
	return GetFocusList(self.Uid, index)
}


//获得正在直播关注列表
func (self *User) GetLiveFocusList(index int) ([]map[string]string, int) {
	return GetLiveFocusList(self.Uid, index)
}
*/

//获得粉丝列表
func (self *User) GetFansList(index int) ([]map[string]string, int) {
	return GetFansList(self.Uid, index)
}

func (self *User) GetFocusInfo(oid int) int {

	var user1 int
	var user2 int
	one_focus := 0

	if self.Uid > oid {
		user1 = self.Uid
		user2 = oid
		one_focus = 1
	} else {
		user1 = oid
		user2 = self.Uid
		one_focus = 2
	}
	focus := &Focus{}
	has, err := orm.Where("user1=? and user2=?", user1, user2).Get(focus)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if has {
		if one_focus == 1 && focus.OneFocus == 1 {
			return common.ERR_SUCCESS
		} else if one_focus == 2 && focus.TwoFocus == 1 {
			return common.ERR_SUCCESS
		}
		return common.ERR_FOCUS_EXIST
	}
	return common.ERR_FOCUS_EXIST
}

func (self *User) IsSuperUser() bool {
	sql := fmt.Sprintf("SELECT * FROM php_identity WHERE uid=%d AND type = 1", self.Uid)
	rowArray, _ := orm.Query(sql)

	if len(rowArray) >= 1 { //超级用户
		return true
	}
	return false
}

func (self *User) GetLoginInfo(info *LoginInfo) {
	info.Uid = self.Uid
	info.Account = self.Account
	info.Tel = self.Tel
	info.NickName = self.NickName
	info.Sex = self.Sex
	info.UserLevel = self.UserLevel
	info.AnchorLevel = self.AnchorLevel
	info.UserExp = self.UserExp
	info.AnchorExp = self.AnchorExp
	info.Image = self.Image
	info.Diamond = self.Diamond

	info.Focus, info.Fans = GetFocusCount(self.Uid)
	//info.Focus = self.Focus
	//info.Fans = self.Fans
	info.Coupons = self.Coupons
	info.Push = self.Push
	info.Signature = self.Signature
	info.Token = self.Token
	info.Location = self.Location
	info.CanLinkMic = self.CanLinkMic
	info.Score = self.Score
	info.Moon = self.Moon
	info.AuthReal = self.AuthRealInfo
	info.IsSuperUser = self.IsSuperUser()
}

func (self *User) GetUserIndex(info *UserIndexInfo) {
	info.NickName = self.NickName
	info.Sex = self.Sex
	info.Signature = self.Signature
	info.Location = self.Location
	info.Image = self.Image
	info.Uid = self.Uid
	info.UserLevel = self.UserLevel
	info.AnchorLevel = self.AnchorLevel
	info.Diamond = self.Diamond
	info.Rice = self.Coupons
	info.Moon = self.Moon
	info.Score = self.Score
	info.IsSuperUser = self.IsSuperUser()
}

func (self *User) GetChatUser(info *UserInfo) {
	info.Chat.Uid = self.Uid
	info.Chat.Sex = self.Sex
	info.Chat.NickName = self.NickName
	info.Chat.UserLevel = self.UserLevel
	info.Chat.AnchorLevel = self.AnchorLevel
	info.Chat.Signature = self.Signature
	info.Chat.Location = self.Location
	info.Chat.Image = self.Image
	info.Chat.IsSuperUser = self.IsSuperUser()
}

//获取送礼总贡献
func GetSendMoneyNum(uid int, money_type int) (sum int, ret int) {
	var sql string
	u, ok := GetUserByUid(uid)
	if ok == common.ERR_SUCCESS {
		if u.AccountType == 1 {
			sql = fmt.Sprintf("SELECT sum(a.value) as count FROM go_gift_record a LEFT JOIN go_user b ON a.send_user=b.uid  WHERE a.rev_user=%d and a.money_type=%d ", uid, money_type)
		} else {
			sql = fmt.Sprintf("SELECT sum(a.value) as count FROM go_gift_record a LEFT JOIN go_user b ON a.send_user=b.uid  WHERE a.rev_user=%d and a.money_type=%d AND b.account_type!=1 ", uid, money_type)
		}
	} else {
		return 0, 0
	}
	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		ret = common.ERR_UNKNOWN
		return
	}
	if len(rowArray) == 1 {
		b, ok := rowArray[0]["count"]
		if ok {
			sum = common.BytesToInt(b)
			ret = common.ERR_SUCCESS
			return
		}
	}

	sum = 0
	ret = common.ERR_SUCCESS
	return
}

func GetOpenGuardMoneyNum(anchor int) (sum int, ret int) {
	var sql string
	u, ok := GetUserByUid(anchor)
	if ok == common.ERR_SUCCESS {
		if u.AccountType == 1 {
			sql = fmt.Sprintf("SELECT sum(a.price) as count FROM go_guard_record a LEFT JOIN go_user b ON a.uid=b.uid  WHERE   a.anchor_id=%d  ", anchor)
		} else {
			sql = fmt.Sprintf("SELECT sum(a.price) as count FROM go_guard_record a LEFT JOIN go_user b ON a.uid=b.uid  WHERE   a.anchor_id=%d AND b.account_type!=1 ", anchor)
		}
	} else {
		return 0, 0
	}
	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		ret = common.ERR_UNKNOWN
		return
	}
	if len(rowArray) == 1 {
		b, ok := rowArray[0]["count"]
		if ok {
			sum = common.BytesToInt(b)
			ret = common.ERR_SUCCESS
			return
		}
	}

	sum = 0
	ret = common.ERR_SUCCESS
	return
}

func GetSendGameGiftRank(uid, index int) ([]map[string]string, int) {
	sql := fmt.Sprintf("select a.rev_user,a.count,b.nick_name,b.user_level as level,b.image,b.sex from (select rev_user,sum(Value) as count from go_gift_record where send_user=%d and money_type=%d group by rev_user) a left join go_user b on a.rev_user=b.uid  where b.account_type!=1 order by a.count  desc limit %d,%d ", uid, common.MONEY_TYPE_SCORE, index*common.SEND_GIFT_PAGE_COUNT, common.SEND_GIFT_PAGE_COUNT)
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

func GetSendDiamonGiftRank(uid, index int) ([]map[string]string, int) {
	sql := fmt.Sprintf("select a.rev_user,a.count,b.nick_name,b.user_level as level,b.image,b.sex from (select rev_user,sum(Value) as count from go_gift_record where send_user=%d and money_type=%d group by rev_user) a left join go_user b on a.rev_user=b.uid  where b.account_type!=1 order by a.count  desc limit %d,%d ", uid, common.MONEY_TYPE_DIAMOND, index*common.SEND_GIFT_PAGE_COUNT, common.SEND_GIFT_PAGE_COUNT)
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

func GetSendDiamonGiftRankWeek(begin_time, end_time int64) ([]map[string]string, int) {
	sql := fmt.Sprintf("select a.send_user,a.count,b.nick_name,b.user_level as level,b.image,b.sex from (select send_user,sum(value) as count from go_gift_record  where money_type=%d and record_time>%d and record_time<%d group by send_user) a left join go_user b on a.send_user=b.uid  where b.account_type!=1  order by a.count  desc limit %d,%d ", common.MONEY_TYPE_DIAMOND, begin_time, end_time, 0, 50)
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

func GetGainGameGiftRank(uid, index int) ([]map[string]string, int) {
	sql := fmt.Sprintf("select a.send_user,a.count,b.nick_name,b.user_level as level,b.image,b.sex from (select send_user,sum(Value) as count from go_gift_record where rev_user=%d and money_type=%d group by send_user) a left join go_user b on a.send_user=b.uid  where b.account_type!=1 order by a.count  desc limit %d,%d ", uid, common.MONEY_TYPE_SCORE, index*common.SEND_GIFT_PAGE_COUNT, common.SEND_GIFT_PAGE_COUNT)
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

func GetGainDiamonGiftRank(uid, index int) ([]map[string]string, int) {
	sql := fmt.Sprintf("select a.send_user,a.count,b.nick_name,b.user_level as level,b.image,b.sex from (select send_user,sum(Value) as count from go_gift_record where rev_user=%d and money_type=%d group by send_user) a left join go_user b on a.send_user=b.uid  where b.account_type!=1 order by a.count  desc limit %d,%d ", uid, common.MONEY_TYPE_DIAMOND, index*common.SEND_GIFT_PAGE_COUNT, common.SEND_GIFT_PAGE_COUNT)
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

func GetGameWinScoreRankWeek(begin_time, end_time int64) ([]map[string]string, int) {
	sql := fmt.Sprintf("select a.uid,a.count,b.nick_name,b.user_level as level,b.image,b.sex  from (select uid,sum(win_score) as count from go_win_score_record  where create_time > %d and create_time < %d group by uid) a left join go_user b on a.uid = b.uid  where b.account_type!=1  order by a.count  desc limit %d,%d", begin_time, end_time, 0, 50)

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

//获取禁言状态
func (self *User) GagStatusUser(rid string) int {
	gag := GetGagByUidAndRoomID(self.Uid, rid)
	if gag != nil {
		return common.ERR_ALREADY_GAG
	}
	return common.ERR_SUCCESS
}

//禁言
func (self *User) GagUser(oid int) int {
	other, ret := GetUserByUid(oid)
	if ret != common.ERR_SUCCESS {
		return ret
	}

	sess := GetUserSessByUid(self.Uid)
	if sess == nil {
		return common.ERR_STAY_IN_SAME_CHAT
	}
	rid := sess.Roomid
	gag := GetGagByUidAndRoomID(other.Uid, rid)
	if gag != nil {
		return common.ERR_ALREADY_GAG
	}

	aff_row, err := orm.Insert(GagRecord{Owner: self.Uid, Uid: oid, RoomId: rid, CreateTime: time.Now()})

	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		return common.ERR_DB_ADD
	}

	sess_other := GetUserSessByUid(oid)
	if sess_other == nil {
		return common.ERR_SUCCESS
	}

	if rid == sess_other.Roomid {
		var res ResponseClose
		res.MType = common.MESSAGE_TYPE_GAG
		SendMsgToUser(oid, res)
	}

	return common.ERR_SUCCESS
}

//取消禁言
func (self *User) CancelGagUser(oid int) int {
	other, ret := GetUserByUid(oid)
	if ret != common.ERR_SUCCESS {
		return ret
	}

	sess := GetUserSessByUid(self.Uid)
	if sess == nil {
		return common.ERR_STAY_IN_SAME_CHAT
	}
	rid := sess.Roomid
	gag := GetGagByUidAndRoomID(other.Uid, rid)
	if gag != nil {
		aff_row, err := orm.Where("room_id=? and uid=?", rid, oid).Delete(gag)
		if err != nil {
			common.Log.Errf("db err %s", err.Error())
			return common.ERR_UNKNOWN
		}

		if aff_row == 0 {
			return common.ERR_DB_DEL
		}
		var res ResponseClose
		res.MType = common.MESSAGE_TYPE_UNGAG
		SendMsgToUser(oid, res)
		return common.ERR_SUCCESS
	}
	return common.ERR_GAG_EXIST
}

//提现扣除米粒
/*
func (self *User) FrozenRiceFunc(rice int) bool {
	if self.Coupons >= rice && rice > 0 {
		self.Coupons -= rice
		self.FrozenRice += rice
		self.Update()
		return true
	}
	return false
}
*/

//提现
func (self *User) ExchangeRiceToBank(money int) int {
	if common.AccountAuthSwitch {
		if self.AccountType == common.ACCOUNT_TYPE_TEST {
			return common.ERR_SEND_GIFT_ROBOT
		}
	}

	rice := self.CashExchange(money)
	if rice <= 0 {
		return common.ERR_OVER_CASH
	}
	if rice < common.MIN_CASH_RICE {
		return common.ERR_CASH_MIN
	}
	extra, has := GetUserExtraByUid(self.Uid)
	if !has {
		return common.ERR_UNKNOWN
	}

	session := orm.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		common.Log.Errf("orm is error:  %s", err.Error())
		session.Rollback()
		return common.ERR_UNKNOWN
	}
	ret := self.DelMoney(session, common.MONEY_TYPE_RICE, int64(rice), true)
	if ret != common.ERR_SUCCESS {
		session.Rollback()
		return ret
	}
	tradeNO := fmt.Sprintf("BK%sUID%d", time.Now().Format("20060102150405"), self.Uid)

	aff_row, err := session.InsertOne(&CashRecord{OwnerId: self.Uid,
		Rice:       rice,
		Money:      money,
		Statue:     common.CASH_BANK_STATUE_VET,
		CreateTime: time.Now(),
		Bank:       extra.Bank,
		CardNo:     extra.CardNo,
		RealName:   extra.RealName,
		CashType:   0,
		CashNum:    tradeNO},
	)

	if err != nil {
		common.Log.Errf("cash bank failed %d,%d,%d,%s,%s", self.Uid, rice, money, extra.Bank, extra.CardNo)
		session.Rollback()
		return common.ERR_UNKNOWN
	}

	if aff_row == 0 {
		session.Rollback()
		return common.ERR_DB_ADD
	}
	err = session.Commit()
	if err != nil {
		return common.ERR_UNKNOWN
	}
	return common.ERR_SUCCESS

}

//兑换物品
func (self *User) ExchangeItem(item_id int, name, tel, addr string) int {
	Moon_Item_mutex.Lock()
	defer Moon_Item_mutex.Unlock()
	item := GetItemById(item_id)
	if item == nil {
		return common.ERR_CONFGI_ITEM
	}

	if item.Stock <= 0 {
		return common.ERR_STOCK_NIL
	}

	session := orm.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		common.Log.Errf("orm is error:  %s", err.Error())
		session.Rollback()
		return common.ERR_UNKNOWN
	}
	ret := self.DelMoney(session, common.MONEY_TYPE_MOON, int64(item.Moon), true)
	if ret != common.ERR_SUCCESS {
		session.Rollback()
		return ret
	}
	ret = item.DelStock(session)
	if ret != common.ERR_SUCCESS {
		session.Rollback()
		return ret
	}
	aff_row, err := session.InsertOne(&ItemRecord{
		OwnerId:    self.Uid,
		Moon:       item.Moon,
		ItemId:     item_id,
		ItemName:   item.Name,
		Status:     common.SCORE_EXCHANGE_STATUS_VET,
		Name:       name,
		Tel:        tel,
		Addr:       addr,
		CreateTime: time.Now(),
		Title:      item.Title,
		Icon:       item.Icon,
		Money:      item.Money,
	})
	if err != nil || aff_row == 0 {
		if err != nil {
			common.Log.Errf("db err %s", err.Error())
		}
		session.Rollback()
		return common.ERR_UNKNOWN
	}

	err = session.Commit()
	if err != nil {
		return common.ERR_UNKNOWN
	}
	return common.ERR_SUCCESS

}

func (self *User) CashQuotaAll() int {
	//return self.Coupons / common.CASH_BASE_VALUE * 3 / 10
	return self.Coupons / common.CASH_PROPORTION
}

func (self *User) CashRice() int {
	nowtime := common.GetFormartTime()
	if nowtime != self.ThawTime {
		self.FrozenCoupons = 0
		self.UpdateByColS("frozen_coupons")
	}
	return self.Coupons - self.FrozenCoupons
}

//可提现金额
func (self *User) CashQuota() int {
	nowtime := common.GetFormartTime()
	if nowtime != self.ThawTime {
		self.FrozenCoupons = 0
	}
	cahs_rice := self.Coupons - self.FrozenCoupons
	if cahs_rice < 0 {
		cahs_rice = 0
	}

	code, _, todayCashed := self.GetTodayCashedRice()
	if code != common.ERR_SUCCESS {
		todayCashed = common.CASH_MAX_VALUE
	} else if todayCashed < 0 {
		todayCashed = 0
	}

	//canCashMoney := cahs_rice / common.CASH_BASE_VALUE * 3 / 10
	canCashMoney := cahs_rice / common.CASH_PROPORTION //当日剩余星星最大可提现金额

	if todayCashed >= common.CASH_MAX_VALUE { //如果当日已提现金额大于最大额（5000）的情况下
		return 0
	} else { //如果当日已提现金额小于最大额（5000）的情况下todayCashed < 5000
		if canCashMoney+todayCashed <= common.CASH_MAX_VALUE { //如果当日已经提取的人民币+当日剩余星星换算的最大可提现金额加起来不到5000，则返回当日剩余星星最大可提现金额
			return canCashMoney
		} else { //否则返回 5000 - 当日已提取的和剩余可提取中的小的
			if common.CASH_MAX_VALUE-todayCashed < canCashMoney {
				return common.CASH_MAX_VALUE - todayCashed
			} else {
				return canCashMoney
			}
		}
	}

	return cahs_rice / common.CASH_PROPORTION
}

//兑换比例金钱
func (self *User) CashExchange(money int) int {
	if self.CashQuota() < money {
		return 0
	}
	ret := math.Ceil(float64(money * common.CASH_PROPORTION))
	return int(ret)
}

func (self *User) DelSessionById(oid int) int {
	return DelSessionById(self.Uid, oid)
}

func HistoryCountGiftNum(uid int, rev int) int {
	res, err := orm.Query("select sum(value) as count from go_gift_record where send_user=? and rev_user=?", uid, rev)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
	}
	if len(res) == 0 {
		return 0
	}
	history, ok := res[0]["count"]
	if !ok {
		return 0
	}
	return common.BytesToInt(history)
}

func (self *User) MoonOrderList(index int) (res []ItemRecord) {
	res = make([]ItemRecord, 0)
	err := orm.Where("owner_id=?", self.Uid).Limit(common.MOON_ORDER_PAGE_COUNT, index*common.MOON_ORDER_PAGE_COUNT).Find(&res)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return
	}

	for k, v := range res {
		v.Icon = fmt.Sprintf("http://h5.17playlive.com/%s", v.Icon)
		res[k] = v
	}
	return
}

func (self *User) ExchangeScore(num int) int {
	c := GetScoreById(num)
	if c == nil {
		return common.ERR_CONFGI_ITEM
	}

	ret := self.CheckMoney(common.MONEY_TYPE_DIAMOND, int64(num))
	if ret != common.ERR_SUCCESS {
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

	ret = self.DelMoney(session, common.MONEY_TYPE_DIAMOND, int64(num), true)
	if ret != common.ERR_SUCCESS {
		session.Rollback()
		return ret
	}

	ret = self.AddMoney(session, common.MONEY_TYPE_SCORE, int64(c.Score), true)
	if ret != common.ERR_SUCCESS {
		session.Rollback()
		return ret
	}

	err = session.Commit()
	if err != nil {
		return common.ERR_UNKNOWN
	}

	m := &ExchangeToRiceRecord{
		Uid:        self.Uid,
		Diamond:    num,
		Score:      c.Score,
		CreateTime: time.Now().Unix(),
	}

	_, err = orm.InsertOne(m)
	if err != nil {
		common.Log.Errf("diamond exchange to rice err is %s", err.Error())
	}
	return common.ERR_SUCCESS
}

//统计用户当天已提现金额
func (self *User) GetTodayCashedRice() (int, int, int) {
	nowtime := time.Now()
	tomorrow := nowtime.AddDate(0, 0, 1)
	tomorrowBegin := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, time.Local)
	tomorrowStr := tomorrowBegin.Format("2006-01-02 15:04:05")
	todayBegin := time.Date(nowtime.Year(), nowtime.Month(), nowtime.Day(), 0, 0, 0, 0, time.Local)
	todayBeginStr := todayBegin.Format("2006-01-02 15:04:05")

	sql := fmt.Sprintf("SELECT owner_id,sum(rice) rice,sum(money) money FROM go_cash_record WHERE owner_id = '%d' AND create_time >= '%s' AND create_time < '%s' AND statue != 2 AND statue != 5 AND statue != 6", self.Uid, todayBeginStr, tomorrowStr)

	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return common.ERR_UNKNOWN, 0, 0
	} else if len(rowArray) == 0 {
		return common.ERR_SUCCESS, 0, 0
	}

	rice := common.BytesToString(rowArray[0]["rice"])
	rice_, _ := strconv.Atoi(rice)

	money := common.BytesToString(rowArray[0]["money"])
	money_, _ := strconv.Atoi(money)

	return common.ERR_SUCCESS, rice_, money_
}

//提现
func (self *User) ExchangeRiceToWeiXin(money int) int {
	couponsBeforeExchange := self.Coupons - self.FrozenCoupons
	if couponsBeforeExchange < 0 {
		couponsBeforeExchange = 0
	}
	if common.AccountAuthSwitch {
		if self.AccountType == common.ACCOUNT_TYPE_TEST {
			return common.ERR_SEND_GIFT_ROBOT
		}
	}

	rice := self.CashExchange(money)
	if rice <= 0 {
		return common.ERR_OVER_CASH
	}
	if rice < common.MIN_CASH_RICE {
		return common.ERR_CASH_MIN
	}
	extra, has := GetUserExtraByUid(self.Uid)
	if !has {
		return common.ERR_UNKNOWN
	}

	session := orm.NewSession()
	defer session.Close()
	err := session.Begin()
	if err != nil {
		common.Log.Errf("orm is error:  %s", err.Error())
		session.Rollback()
		return common.ERR_UNKNOWN
	}

	ret := self.DelMoney(session, common.MONEY_TYPE_RICE, int64(rice), true)
	if ret != common.ERR_SUCCESS {
		session.Rollback()
		return ret
	}

	couponsAfterExchange := self.Coupons - self.FrozenCoupons
	if couponsAfterExchange < 0 {
		couponsAfterExchange = 0
	}

	tradeNO := fmt.Sprintf("WX%sUID%d", time.Now().Format("20060102150405"), self.Uid)
	aff_row, err := orm.InsertOne(&CashRecord{
		OwnerId:       self.Uid,
		CouponsBefore: couponsBeforeExchange, //本笔订单开始前米粒余额
		CouponsAfter:  couponsAfterExchange,  //本笔订单开始后米粒余额
		Rice:          rice,
		Money:         money,
		Statue:        common.CASH_BANK_STATUE_VET,
		CreateTime:    time.Now(),
		Bank:          extra.Bank,
		CardNo:        extra.CardNo,
		RealName:      extra.RealName,
		CashType:      1,
		CashNum:       tradeNO,
	})
	if err != nil {
		common.Log.Errf("cash bank failed %d,%d,%d,%s,%s", self.Uid, rice, money, extra.Bank, extra.CardNo)
		session.Rollback()
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		session.Rollback()
		return common.ERR_DB_ADD
	}
	err = session.Commit()
	if err != nil {
		return common.ERR_UNKNOWN
	}
	return common.ERR_SUCCESS
}

func (self *User) HasDaySettle() bool {
	rowArray, _ := orm.Query("select * from php_anchor  where uid = ? && settle_type = 1", self.Uid)

	if len(rowArray) > 0 {
		return true
	}

	return false
}

func (self *User) GetDayDuration() int {
	stdtime := time.Now()

	t1 := time.Date(stdtime.Year(), stdtime.Month(), stdtime.Day(), 0, 0, 0, 0, time.Local)

	rowArray, _ := orm.Query("select SUM(UNIX_TIMESTAMP(finish_time) - UNIX_TIMESTAMP(create_time)) total_second from go_room_list WHERE owner_id = ? AND finish_time > create_time AND create_time >= ? AND statue = 3", self.Uid, t1.Format("2006-01-02 15:04:05"))

	if len(rowArray) > 0 {
		totalSecond := common.BytesToString(rowArray[0]["total_second"])
		totalSecond_, _ := strconv.Atoi(totalSecond)
		if totalSecond_ <= 0 {
			return 0
		} else {
			return totalSecond_
		}
	}

	return 0
}

func (self *User) HasResetNickName() bool {
	has, err := orm.Where("uid=?", self.Uid).Get(&NicknameResetRecord{})
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return false
	}
	if has { //以前修改过昵称了
		return true
	}
	return false
}

func (self *User) ResetNickName(nickName string) int {
	has, err := orm.Where(" uid=?", self.Uid).Get(&NicknameResetRecord{})
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if has { //以前修改过昵称了
		if self.Score <= int64(NickNameResetMoney) {
			return common.ERR_SCORE_NOTENOUGH
		} else { //有足够的余额去修改昵称
			result := self.DelMoney(nil, common.MONEY_TYPE_SCORE, int64(NickNameResetMoney), false)
			if result == common.ERR_SUCCESS {
				m := &NicknameResetRecord{
					Uid:           self.Uid,
					Score:         NickNameResetMoney,
					OperationTime: int(time.Now().Unix()),
				}

				_, err := self.SetNick(nickName)
				if err != nil {
					return common.ERR_UNKNOWN
				}

				_, err = orm.Insert(m)
				if err != nil {
					common.Log.Errf("orm err is %s", err.Error())
					return common.ERR_UNKNOWN
				}
				return common.ERR_SUCCESS
			} else {
				return result
			}
		}
	} else { //第一次修改昵称
		m := &NicknameResetRecord{
			Uid:           self.Uid,
			Score:         0,
			OperationTime: int(time.Now().Unix()),
		}

		_, err := self.SetNick(nickName)
		if err != nil {
			return common.ERR_UNKNOWN
		}

		_, err = orm.Insert(m)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
		return common.ERR_SUCCESS
	}
}

func TradeNoToWeiXinPay(ID int, desc string) int {

	cash := &CashRecord{}
	has, err := orm.Where("id=? and (statue=3||statue=7||statue=8)", ID).Get(cash)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if has { // 如果单号存在的情况下
		user, ret := GetUserByUid(cash.OwnerId)
		if ret != common.ERR_SUCCESS {
			return ret
		}

		extra, extraHas := GetUserExtraByUid(cash.OwnerId)
		if !extraHas {
			return common.ERR_ACCOUNT_EXIST
		}

		cash.CashNum = fmt.Sprintf("WX%sUID%d", time.Now().Format("20060102150405"), cash.OwnerId)
		orm.Where("id=?", ID).Update(cash)

		if desc == "" {
			desc = "17玩直播微信提现"
		}

		errcode, errCodeDes := Weixinpay(user.OpenId, extra.RealName, desc, cash.CashNum, int64(cash.Money*100))

		if errcode == common.ERR_SUCCESS {
			cash.Statue = common.CASH_PAY_SUCCESS
		} else if errcode == common.ERR_WEIXIN_CASH_NOTENOUGH { //我们公司的微信账户余额不足
			cash.Statue = common.CASH_PAY_NOTENOUGH
		} else if errcode == common.ERR_WEIXIN_CASH_SYSTEMERROR { //腾讯系统错误
			cash.Statue = common.CASH_PAY_TENCENT_SYSTEMERROR
		} else {
			cash.Statue = common.CASH_PAY_FAIL
			cash.ErrCodeDes = errCodeDes

			//打款失败时候自动返还金额给用户
			//user, ret := GetUserByUid(cash.OwnerId)
			//if ret != common.ERR_SUCCESS {
			//	return ret
			//}
			//user.addRice(cash.Rice)
			user.addCoupons(cash.Rice)
		}
		cash.FinishTime = time.Now()
		cash.CouponsAfter = user.Coupons

		orm.Where("id=?", ID).Update(cash)

		return errcode
	} else {
		return common.ERR_WEIXIN_TRADENO_ERROR
	}
}

func TradeNoRejectWeiXinPay(ID int, desc string) int {

	cash := &CashRecord{}
	has, err := orm.Where("id=? and (statue=2 || statue=10)", ID).Get(cash)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if has { // 如果单号存在的情况下
		user, ret := GetUserByUid(cash.OwnerId)
		if ret != common.ERR_SUCCESS {
			return ret
		}

		result := user.addCoupons(cash.Rice)
		if result == common.ERR_SUCCESS {
			cash.Statue = common.CASH_PAY_RICE_RETURN
			cash.FinishTime = time.Now()
			cash.CouponsAfter = user.Coupons
			orm.Where("id=?", ID).Update(cash)
			return common.ERR_SUCCESS
		} else {
			cash.CouponsAfter = user.Coupons
			orm.Where("id=?", ID).Update(cash)
			return common.ERR_UNKNOWN
		}
	} else {
		return common.ERR_WEIXIN_TRADENO_ERROR
	}
}

func MoonToWeiXinPay(ID int, desc string) int {

	cash := &MoonCashRecord{}
	has, err := orm.Where("id=? and (statue=1||statue=7||statue=8)", ID).Get(cash)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if has { // 如果单号存在的情况下
		user, ret := GetUserByUid(cash.OwnerId)
		if ret != common.ERR_SUCCESS {
			return ret
		}

		extra, extraHas := GetUserExtraByUid(cash.OwnerId)
		if !extraHas {
			return common.ERR_ACCOUNT_EXIST
		}

		item := &ItemRecord{}
		itemhas, err := orm.Where("id=?", cash.ItemRecordId).Get(item)
		if err != nil {
			return common.ERR_UNKNOWN
		}

		cash.CashNum = fmt.Sprintf("WX%sUID%dMOON", time.Now().Format("20060102150405"), cash.OwnerId)
		orm.Where("id=?", ID).Update(cash)

		if desc == "" {
			desc = "17玩月亮商城提现"
		}

		errcode, errCodeDes := Weixinpay(user.OpenId, extra.RealName, desc, cash.CashNum, int64(cash.Money*100))

		if errcode == common.ERR_SUCCESS {
			cash.Statue = common.MOON_CASH_STATUS_SUCCESS
			if itemhas { //更新go_item_record表中记录的状态为线上打款成功
				item.Status = common.SCORE_EXCHANGE_STATUS_ONLINE_SUCCESS
				orm.Where("id=?", item.Id).Update(item)
			}
		} else if errcode == common.ERR_WEIXIN_CASH_NOTENOUGH { //我们公司的微信账户余额不足
			cash.Statue = common.CASH_PAY_NOTENOUGH
		} else if errcode == common.ERR_WEIXIN_CASH_SYSTEMERROR { //腾讯系统错误
			cash.Statue = common.CASH_PAY_TENCENT_SYSTEMERROR
		} else {
			cash.Statue = common.MOON_CASH_STATUS_FAILED
			cash.ErrCodeDes = errCodeDes

			// 打款失败时候自动返还金额给用户
			user, ret := GetUserByUid(cash.OwnerId)
			if ret != common.ERR_SUCCESS {
				return ret
			}

			user.addMoon(cash.Moon)

			if itemhas { //更新go_item_record表中记录的状态为线上打款失败
				item.Status = common.SCORE_EXCHANGE_STATUS_ONLINE_FAILED
				orm.Where("id=?", item.Id).Update(item)
			}
		}
		cash.FinishTime = time.Now()

		orm.Where("id=?", ID).Update(cash)

		return errcode
	} else {
		return common.ERR_WEIXIN_TRADENO_ERROR
	}
}

//用户1是否关注了用户2[用户1是否为用户2的粉丝]
func User1IsFocus2(uid1, uid2 int) int {
	sql := fmt.Sprintf("SELECT user_relation(%d,%d) is_focus", uid1, uid2)
	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return 0
	}
	isFocus := common.BytesToInt(rowArray[0]["is_focus"])

	return isFocus
}

type OutUserInfo3 struct {
	Uid         int    `json:"uid"`
	Image       string `json:"image"`
	Signature   string `json:"signature"`
	Sex         int    `json:"sex"`
	Userlevel   int    `json:"user_level"`
	AnchorLevel int    `json:"anchor_level"`
	NickName    string `json:"nick_name"`
	Location    string `json:"location"`
	RoomId      string `json:"room_id"`
	Cover       string `json:"cover"`
	LiveUrl     string `json:"live"`
	Viewer      int    `json:"viewer"`
	GameType    int    `json:"game_type"`
	RoomName    string `json:"room_name"`
	Statue      int    `json:"statue"`
	Score       int64  `json:"score"`
	IsFocus     int    `json:"is_focus"`
	FlvUrl      string `json:"flv_url"`
}

func GetSearchList(uid, index int, keyWords string) (out_users []OutUserInfo3, ret int) {

	sql1 := "SELECT uid"
	str2 := fmt.Sprintf(" limit %d,%d", index*common.FOCUS_LIST_PAGE_COUNT, common.FOCUS_LIST_PAGE_COUNT)
	sql := sql1 + ` FROM go_user WHERE nick_name LIKE '%` + keyWords + `%' || uid ='` + keyWords + `'` + str2

	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return
	}

	all_user := make([]int, 0)
	live_rooms := make(map[int]*CacheUser, 0)
	for _, row := range rowArray {
		uid := row["uid"]
		uid_ := common.BytesToInt(uid)

		all_user = append(all_user, uid_)
		u, err := GetCacheUser(uid_)
		if err == redis.Nil {
			continue
		} else if err != nil {
			return
		} else {
			if u.Status == common.USER_STATUE_LIVE {
				//live_users=append(live_users,uid_)
				live_rooms[uid_] = u
			}
		}
	}

	users := make([]User, 0)
	if len(all_user) == 0 {
		return
	}
	err = orm.In("uid", all_user).Find(&users)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return
	}

	user, _ := GetUserByUid(uid)
	acountType := user.AccountType
	if acountType == 1 { //如果是测试号
		for _, v := range users {
			var out OutUserInfo3
			out.Uid = v.Uid
			out.AnchorLevel = v.AnchorLevel
			out.Cover = v.Image
			out.Image = v.Image
			out.Location = v.Location
			out.NickName = v.NickName
			out.Signature = v.Signature
			out.Score = v.Score
			out.Userlevel = v.UserLevel
			out.AnchorLevel = v.AnchorLevel
			out.IsFocus = User1IsFocus2(uid, v.Uid)
			c, ok := live_rooms[v.Uid]
			if !ok {
				out.Statue = common.USER_STATUE_LEAVE
				out_users = append(out_users, out)
				continue
			} else {
				out.RoomId = c.RoomId
				room := &RoomList{}
				has, err := orm.Where("room_id=?", c.RoomId).Get(room)
				if err != nil {
					common.Log.Errf("db err %s", err.Error())
					return
				}
				if has {
					out.LiveUrl = room.LiveUrl
					out.GameType = room.GameType
					out.RoomName = room.RoomName
					out.Statue = room.Statue
					out.FlvUrl = room.FlvUrl
				} else {
					continue
				}
				out_users = append(out_users, out)
			}
			//
		}
	} else { //如果是正常号[非测试号]
		for _, v := range users {
			if v.AccountType == 0 { //非测试号返回
				var out OutUserInfo3
				out.Uid = v.Uid
				out.AnchorLevel = v.AnchorLevel
				out.Cover = v.Image
				out.Image = v.Image
				out.Location = v.Location
				out.NickName = v.NickName
				out.Signature = v.Signature
				out.Score = v.Score
				out.Userlevel = v.UserLevel
				out.AnchorLevel = v.AnchorLevel
				out.IsFocus = User1IsFocus2(uid, v.Uid)
				c, ok := live_rooms[v.Uid]
				if !ok {
					out.Statue = common.USER_STATUE_LEAVE
					out_users = append(out_users, out)
					continue
				} else {
					out.RoomId = c.RoomId
					room := &RoomList{}
					has, err := orm.Where("room_id=?", c.RoomId).Get(room)
					if err != nil {
						common.Log.Errf("db err %s", err.Error())
						return
					}
					if has {
						out.LiveUrl = room.LiveUrl
						out.GameType = room.GameType
						out.RoomName = room.RoomName
						out.Statue = room.Statue
						out.FlvUrl = room.FlvUrl
					} else {
						continue
					}

					out_users = append(out_users, out)
				}
			}
			//
		}
	}

	return
}

//获取搜索的总长度
func GetKeyWordsSearchLength(keyWords string) (int, int) {

	sql1 := "SELECT uid"
	sql := sql1 + ` FROM go_user WHERE nick_name LIKE '%` + keyWords + `%' || uid ='` + keyWords + `'`

	rowArray, err := orm.Query(sql)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return common.ERR_DB_FIND, 0
	}

	return common.ERR_SUCCESS, len(rowArray)
}
