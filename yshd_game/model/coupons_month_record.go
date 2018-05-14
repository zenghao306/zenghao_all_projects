package model

import (
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	"time"
)

type CouponsMonthRecord struct {
	Uid          int    `xorm:"int(11) not null UNIQUE(FOCUSE_USER)"`
	YearMonths   string `xorm:"varchar(40) not null"` //产生年月
	BeginBalance int    `xorm:"int(11) default(0)"`
	Income       int    `xorm:"int(11) default(0)"`
	Rice         int    `xorm:"int(11) default(0)"`
	Money        int    `xorm:"int(11) default(0)"`
	EndBalance   int    `xorm:"int(11) default(0)"`
}

// 每月初所有用户米粒数备份。
func ConponsBkAtEveryMonthBegin() {
	sql := "insert into go_coupons_month_record(uid,year_months,begin_balance) select uid,'" + common.GetCurentYearMonthString() + "',coupons from user"
	_, err := orm.Exec(sql)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
	}
}

// 所有用户(上月末)米粒数备份。
func ConponsBkAtEveryMonthEnd() {
	lastYearMonth := common.GetLastYearMonthString()
	sql := "CALL statistics_conpons_bk_everymonth_end('" + lastYearMonth + "')"
	_, err := orm.Exec(sql)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
	}
}

// 删掉上月所有体现金额为零的记录（因为体现为零不必做记录）
func DeleteCashIsZeroRecordAtLastMonth() {
	lastYearMonth := common.GetLastYearMonthString()
	sql := "DELETE FROM go_coupons_month_record WHERE money = 0 AND year_months='" + lastYearMonth + "'"
	_, err := orm.Exec(sql)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
	}
}

// 上月所有用户收入（米粒）记录到数据库表
func AllUsersIncomeRiceMoneyAtLastMonth() {
	lastYear, lastMonth := common.GetLastYearMonth()

	t1 := time.Date(lastYear, time.Month(lastMonth), 1, 0, 0, 0, 0, time.Local)

	stdtime := time.Now()

	t2 := time.Date(stdtime.Year(), stdtime.Month(), 1, 0, 0, 0, 0, time.Local)

	lastYearMonth := common.GetLastYearMonthString()

	// 以下是统计上月所有用户收入（米粒）记录到数据库表
	sql := "CALL statistics_income('" + t1.Format("2006-01-02 15:04:05") + "','" + t2.Format("2006-01-02 15:04:05") + "','" + lastYearMonth + "')"
	_, err := orm.Exec(sql)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
	}

	// 以下是统计上月所有用户提现米粒|提现金额统计到数据库表中
	sql = "CALL statistics_rice_money('" + t1.Format("2006-01-02 15:04:05") + "','" + t2.Format("2006-01-02 15:04:05") + "','" + lastYearMonth + "')"
	_, err = orm.Exec(sql)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
	}
}
