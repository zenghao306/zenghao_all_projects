package model

import (
	"fmt"
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	"github.com/yshd_game/confdata"
	"strconv"
	"strings"
	"time"
)

type UserExtra struct {
	Uid                    int       `xorm:"int(11) pk not null "`        //用户ID
	Bank                   string    `xorm:"varchar(30)"`                 //银行类型
	CardNo                 string    `xorm:"varchar(40) UNIQUE(CARDNO)"`  //卡号
	RealName               string    `xorm:"varchar(30)"`                 //真名
	CashTel                string    `xorm:"varchar(20) UNIQUE(CASHTEL)"` //提现手机
	CashTime               time.Time //提现时间
	IsChangeTel            bool      //提现手机可以重置标志位
	MonthConsume           int       //月消费
	ConsumeLastTime        string    `xorm:"varchar(30)"`       //当前统计的消费月份
	AddConsumeLevel        int       `xorm:"int(11) default(0)` //已经增加过的经验次数对应月消费额对应配置的的cid索引
	DayWatchTime           int       `xorm:"int(11) default(0)` //当天已经观看的时长
	DayAnchorTime          int       `xorm:"int(11) default(0)` //当天已经主播的时长
	LastWatchTime          string    `xorm:"varchar(30)"`       //当前统计的消费天
	AddWacthTimes          int       `xorm:"int(11) default(0)` //已经增加过的经验次数对应观看时长
	AddAnchorTimes         int       `xorm:"int(11) default(0)` //主播一天增加过的经验统计
	PlayBackCount          int       `xorm:"int(11)`
	PlayBackRecommandType  int       `xorm:"int(11)` //0普通1会议室
	PlayBackRecommandRid   string    `xorm:"varchar(48)"`
	DailyRefreshTime       time.Time //每日刷新过时间
	OnlineRewardTimes      int       `xorm:"bigint default(0)`
	OnlineRewardFinishTime int64
	//	OnlineRewardContinue   int64
}

func GetUserExtraByUid(uid int) (*UserExtra, bool) {
	user := &UserExtra{}
	has, err := orm.Where("uid=?", uid).Get(user)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
	}
	return user, has
}

func GetUserExtraByUidStr(uid string) (*UserExtra, bool) {
	uid_, _ := strconv.Atoi(uid)
	return GetUserExtraByUid(uid_)
}

func (self *UserExtra) Update22() bool {
	if self.Uid == 0 {
		common.Log.Err("update extre error")
		return false
	}
	aff_row, err := orm.Where("uid=?", self.Uid).Update(self)
	if err != nil {
		common.Log.Errf("db error:", err.Error())
		return false
	}

	if aff_row == 0 {
		return false
	}
	return true
}

/*
func (self *UserExtra) UpdateFront(filed string) bool {
	aff_row, err := orm.Where("uid=?", self.Uid).MustCols(filed).Update(self)
	if err != nil {
		common.Log.Errf("db error:", err.Error())
		return false
	}

	if aff_row == 0 {
		return false
	}
	return true
}
*/
func (self *UserExtra) UpdateByColS(filed ...string) (aff_row int64, err error) {
	aff_row, err = orm.Where("uid=?", self.Uid).Cols(filed...).Update(self)
	if err != nil {
		common.Log.Errf("db error:", err.Error())
	}
	if aff_row == 0 {
		common.Log.Errf("aff row")
	}
	return
}

func (self *UserExtra) CheckCashTel() bool {
	if self.CashTel != "" {
		return true
	}
	return false
}

//设置提现电话
func (self *UserExtra) SetCashTel(tel string) int {
	has, err := orm.Where("cash_tel=?", tel).Get(&UserExtra{})
	if err != nil {
		common.Log.Errf("db error:", err.Error())
		return common.ERR_UNKNOWN
	}

	if has {
		return common.ERR_TEL_DUPLICATE
	}
	self.CashTel = tel
	if _, err := self.UpdateByColS("cash_tel"); err != nil {
		return common.ERR_SUCCESS
	}
	return common.ERR_UNKNOWN
}

//设置修改标记
func (self *UserExtra) SetChangeFlag(flag bool) {
	self.IsChangeTel = flag
	self.UpdateByColS("is_change_tel")
}

//提现记录
func (self *UserExtra) GetCashRecord(index int) (int, []CashRecord) {
	record := make([]CashRecord, 0)
	record2 := make([]CashRecord, 0)
	err := orm.Where("owner_id=?", self.Uid).Limit(common.CASH_RECORD_PAGE_COUNT, common.CASH_RECORD_PAGE_COUNT*index).Find(&record)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return common.ERR_UNKNOWN, record
	}

	for _, v := range record { //修改状态值给前端，前端没有9和10两个状态
		if v.Statue == 9 {
			fmt.Printf("\n v.Statue=%d", v.Statue)
			v.Statue = 1
		} else if v.Statue == 10 {
			fmt.Printf("\n v.Statue=%d", v.Statue)
			v.Statue = 2
		}
		record2 = append(record2, v)
	}

	return common.ERR_SUCCESS, record2
}

func (self *UserExtra) ConsumerStatistics(diamond int) {
	nowtime := common.GetFormartTime()
	if nowtime == self.ConsumeLastTime {
		self.MonthConsume += diamond
	} else {
		self.MonthConsume = diamond
		self.ConsumeLastTime = nowtime
		self.AddConsumeLevel = 1
	}
	self.UpdateByColS("month_consume", "consume_last_time", "add_consume_level")
}

//检查消费金额满足增加经验条件（废弃）

/*
func (self *UserExtra) CheckConsumerLevel() {
	return
	c, has := GetComsumerById(self.AddConsumeLevel)
	if has {
		if self.MonthConsume >= c.ConsumeNum {
			self.AddConsumeLevel++
			user, _ := GetUserByUid(self.Uid)
			user.AddUserExp(c.Exp)
			self.CheckConsumerLevel()
		}
	}
}
*/
//刷新检测
func (self *UserExtra) Refresh() {
	nowtime := time.Now()
	nowtime_ := nowtime.Format("20060102")
	refresh := self.DailyRefreshTime.Format("20060102")
	if strings.Compare(nowtime_, refresh) == 0 {
		return
	}

	self.RefreshTaskDialy()
	self.RestOnlineBaseData()

	self.DailyRefreshTime = nowtime

	//orm.Exec("update user_extra set online_reward_times=? , online_reward_time=? , task_refresh_time=? where uid=?",self.OnlineRewardTimes,self.OnlineRewardFinishTime,self.DailyRefreshTime,self.Uid)
	self.UpdateByColS("task_refresh_time", "online_reward_times", "online_reward_finish_time", "daily_refresh_time")
}

//刷新日常任务系统
func (self *UserExtra) RefreshTaskDialy() {
	m, _, _, _, ret := ListTask(self.Uid)
	if ret == common.ERR_SUCCESS {
		for _, v := range m {
			entry, ok := confdata.ConfigData.TaskById[v.TaskId]
			if ok {

				if entry.Type == confdata.TaskType_daily {
					DelTask(self.Uid, v.TaskId)
				}
			}
		}
	}

	for _, v := range confdata.ConfigData.TaskById {
		if v.Type == confdata.TaskType_daily {
			AccecptTask(self.Uid, v.Id)
		}
	}
}

//上线设置下次获取奖励时间点
func (self *UserExtra) SetNextOnlineRewardTime() int {
	self.Refresh()
	if self.OnlineRewardTimes >= 6 {
		return common.ERR_DAILY_ONLINE_TIMES
	}

	if self.OnlineRewardFinishTime == 0 {
		self.OnlineRewardFinishTime = time.Now().Unix() + common.ONlINE_TIME_REWARD
		_, err := self.UpdateByColS("online_reward_finish_time")
		if err != nil {
			return common.ERR_UNKNOWN
		}
	}
	return common.ERR_SUCCESS
	/*
		if common.ONlINE_TIME_REWARD > self.OnlineRewardContinue {
			left := common.ONlINE_TIME_REWARD - self.OnlineRewardContinue

			self.OnlineRewardFinishTime += time.Now().Unix() + left
		} else {
			self.OnlineRewardFinishTime = time.Now().Unix()
		}
		self.UpdateFront("online_reward_finish_time")
	*/
}

//下线记录在线时间满足奖励时间点
func (self *UserExtra) FinishOnlineRewardTime() {
	/*
		if time.Now().Unix() <= self.OnlineRewardFinishTime {
			self.OnlineRewardContinue = common.ONlINE_TIME_REWARD - (self.OnlineRewardFinishTime - time.Now().Unix())
		} else {
			self.OnlineRewardContinue = common.ONlINE_TIME_REWARD
		}
		self.UpdateFront("online_reward_continue")
	*/
	//self.OnlineRewardContinue
}

//提交获取奖励
func (self *UserExtra) PostOnlineReward() int {
	nowtime := time.Now()
	nowtime_ := nowtime.Format("20060102")
	refresh := self.DailyRefreshTime.Format("20060102")
	if strings.Compare(nowtime_, refresh) != 0 {
		return common.ERR_PLEASE_GOINGO_CHAT
	}

	if time.Now().Unix() >= self.OnlineRewardFinishTime && self.OnlineRewardFinishTime != 0 {

		if self.OnlineRewardTimes >= 6 {
			return common.ERR_DAILY_ONLINE_TIMES
		}
		u, _ := GetUserByUid(self.Uid)
		ret := u.AddMoney(nil, common.MONEY_TYPE_SCORE, common.ONLINE_TIME_SOCRE, false)
		if ret != common.ERR_SUCCESS {
			return ret
		}
		self.OnlineRewardTimes += 1
		if self.OnlineRewardTimes != 6 {
			self.OnlineRewardFinishTime = time.Now().Unix() + common.ONlINE_TIME_REWARD
		} else {
			self.OnlineRewardFinishTime = 0
		}

		self.UpdateByColS("online_reward_continue", "online_reward_times", "online_reward_finish_time")
		return common.ERR_SUCCESS
	}
	return common.ERR_ONLINE_REWARD_EARLY
}

//重置每日在线奖励
func (self *UserExtra) RestOnlineBaseData() {
	self.OnlineRewardTimes = 0
	self.OnlineRewardFinishTime = 0
}

type LiveInfo struct {
	Rid       string
	Rice      int
	BeginTime int64
	EndTime   int64
	Moon      int64
}

//直播信息列表
func (self *UserExtra) LiveDetailStatistics(begin, next_tm time.Time, index int) (r []LiveInfo, ret int) {
	r = make([]LiveInfo, 0)
	s := make([]RoomList, 0)
	//err := orm.Where("owner_id=? and create_time>=? and create_time<?", self.Uid, begin, next_tm).Limit(common.LIVE_STATISTICS_PAGE_COUNT, index*common.LIVE_STATISTICS_PAGE_COUNT).Find(&s)

	err := orm.Where("owner_id=? and UNIX_TIMESTAMP(create_time)>=? and UNIX_TIMESTAMP(create_time)<? and  finish_time !='0001-01-01 00:00:00' ", self.Uid, begin.Unix(), next_tm.Unix()).Limit(common.LIVE_STATISTICS_PAGE_COUNT, index*common.LIVE_STATISTICS_PAGE_COUNT).Find(&s)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		ret = common.ERR_UNKNOWN
		return
	}

	for _, v := range s {
		var info LiveInfo
		info.Rice = v.Rice
		info.Rid = v.RoomId
		info.Moon = int64(v.Moon)
		info.BeginTime = v.CreateTime.Unix()
		info.EndTime = v.FinishTime.Unix()
		r = append(r, info)
	}

	return
}

//直播统计
func (self *UserExtra) LiveStatistics(begin, next_tm time.Time) (valid_day int, valid_hours int, rice int, moon int, min int) {
	m := make([]RoomList, 0)
	err := orm.Where("owner_id=? and UNIX_TIMESTAMP(create_time)>=? and UNIX_TIMESTAMP(create_time)<? and  finish_time >create_time ", self.Uid, begin.Unix(), next_tm.Unix()).Find(&m)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return
	}
	allsec := 0
	for _, v := range m {
		diff := v.FinishTime.Sub(v.CreateTime)
		if diff <= 0 {
			continue
		}
		sec := diff / time.Second
		allsec += int(sec)
	}
	valid_hours = allsec / 3600

	res, err := orm.Query("select SUM(rice) as all_rice,SUM(moon)  as all_moon from go_room_list where owner_id=? and  finish_time >create_time and create_time>? and create_time<=? ", self.Uid, begin, next_tm)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return
	}
	b, ok := res[0]["all_rice"]
	if ok {
		rice = common.BytesToInt(b)
	}

	b, ok = res[0]["all_moon"]
	if ok {
		moon = common.BytesToInt(b)
	}

	/*
		res, err = orm.Query("select SUM(moon) as all_moon from go_room_list where owner_id=? and  finish_time !='0001-01-01 00:00:00' and create_time>? and finish_time<=?", self.Uid, begin, next_tm)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return
		}
		b, ok = res[0]["all_moon"]
		if ok {
			moon = common.BytesToInt(b)
		}
	*/

	rep, _ := orm.Query("SELECT UNIX_TIMESTAMP(?) as a", begin)
	_ = common.BytesToInt(rep[0]["a"])

	//godump.Dump(uu)

	ref, err := orm.Query("  SELECT   DATE_FORMAT(create_time,'%Y%m%d') AS days,  SUM(TIMESTAMPDIFF(SECOND,create_time,finish_time))  AS allses   FROM go_room_list WHERE owner_id=? AND finish_time >create_time  AND UNIX_TIMESTAMP(create_time) > UNIX_TIMESTAMP(?) AND UNIX_TIMESTAMP(create_time)<=UNIX_TIMESTAMP(?) GROUP BY days;", self.Uid, begin, next_tm)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return
	}
	sec := 0
	for k, _ := range ref {
		g := common.BytesToInt(ref[k]["allses"])
		if g > 3600 {
			valid_day += 1
		}

		_ = common.BytesToString(ref[k]["days"])

		//godump.Dump(d)
		sec += g

	}
	min = sec / 60
	/*
		res,err:=orm.Query("call statistics_live_hours(?,?,?)",self.Uid,begin,next_tm)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return
		}


		b,ok:=res[0]["valid_sec"]
		if ok {
			valid_sec=common.BytesToInt(b)
			valid_hours=valid_sec/3600
		}

		b,ok=res[0]["orice"]
		if ok {
			rice=common.BytesToInt(b)
		}

		res,err=orm.Query("call statistics_live_days(?,?,?)" ,self.Uid,begin,next_tm)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return
		}
		b,ok=res[0]["valid_day"]
		if ok {
			valid_day=common.BytesToInt(b)
		}
		godump.Dump(valid_hours)
		godump.Dump(valid_day)
	*/
	return
}

func (self *UserExtra) FocusStatistics(begin, next_tm time.Time) int {
	/*
		cur_tm, _ := time.Parse("2015-08-05", date)
		begin := time.Date(cur_tm.Year(), cur_tm.Month(), 0, 0, 0, 0, 0, time.Local)

		next_moneth := cur_tm.AddDate(0, 1, 0)
		next_tm := time.Date(next_moneth.Year(), next_moneth.Month(), 1, 0, 0, 0, 0, time.Local)
		last_sec := next_tm.Add(-1)
	*/
	res, err := orm.Query("select count(*) as fans from ( select * from go_focus where user1=? and two_focus=1 and focus_time2>=? and focus_time2<? union all select * from go_focus where user2=? and one_focus=1  and focus_time1>? and focus_time1<=?) c", self.Uid, begin, next_tm, self.Uid, begin, next_tm)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return 0
	}
	count := res[0]["fans"]
	return common.BytesToInt(count)
}
