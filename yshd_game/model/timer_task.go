package model

import (
	//	"github.com/liudng/godump"
	"github.com/yshd_game/common"
	"github.com/yshd_game/timer"
	"os"
	"time"
	//"syscall"
)

func TimerTask() {
	defer common.PrintPanicStack()
	for {
		d := timer.NewDispatcher(1)
		stdtime := time.Now()
		tomorrow := stdtime.AddDate(0, 0, 1)
		t := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, time.Local)
		du := time.Duration(t.Unix()-stdtime.Unix()) * time.Second

		d.AfterFunc(du, func() {
			nowtime := time.Now()
			filename := nowtime.Format("20060102")
			finalname := common.Logpath + filename
			f, _ := os.OpenFile(finalname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
			common.Log.SetNewOutPutFile(f)
			//syscall.Dup2(int(f.Fd()), 2)
			//nowtime2 := time.Now()
			//DeleteNiuNiuRaiseRecordLast(nowtime2.Unix()-24*3600, nowtime2.Unix()-2*24*3600) //删除前一天的redis数据表里游戏押分数据
			//DeleteNiuNiuRecordBeforeTheTime(nowtime2.Unix() - 3*24*3600)                    //删掉2天前押注为0的游戏记录
		})
		(<-d.ChanTimer).Cb()
	}
}

//每日凌晨5点备份游戏押分数据
func TimerTaskGameRaiseBk() {
	defer common.PrintPanicStack()
	for {
		d := timer.NewDispatcher(1)
		stdtime := time.Now()
		tomorrow := stdtime.AddDate(0, 0, 1)
		t := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 5, 0, 0, 0, time.Local)
		du := time.Duration(t.Unix()-stdtime.Unix()) * time.Second
		d.AfterFunc(du, func() {
			stdtime := time.Now()
			t2 := time.Date(stdtime.Year(), stdtime.Month(), stdtime.Day(), 0, 0, 0, 0, time.Local)
			BakUserRaisePosInfoWithGameID(t2.Unix()-24*3600, t2.Unix())
		})
		(<-d.ChanTimer).Cb()
	}
}

// 米粒和体现按月记录
func TimerTaskCouponsMonthRecord() {
	defer common.PrintPanicStack()
	for {
		//LastYearMonth := common.GetLastYearMonthString()
		//godump.Dump(LastYearMonth)
		d := timer.NewDispatcher(1)
		stdtime := time.Now()
		nextMonth := stdtime.AddDate(0, 1, 0)
		t := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, time.Local)
		du := time.Duration(t.Unix()-stdtime.Unix()) * time.Second

		d.AfterFunc(du, func() {
			ConponsBkAtEveryMonthBegin()         // 米粒按月记录【每月的1号0点0分记录本月初的米粒和上月末的米粒】
			AllUsersIncomeRiceMoneyAtLastMonth() // 上月米粒提现统计到数据库表记录。
			ConponsBkAtEveryMonthEnd()           //上月末米粒提现。
			DeleteCashIsZeroRecordAtLastMonth()  // 删掉上月体现记录为零的记录。
		})
		(<-d.ChanTimer).Cb()
	}
}

func TimerTaskRobotSay() {
	defer common.PrintPanicStack()
	for {
		d := timer.NewDispatcher(1)

		d.AfterFunc(1*time.Second, func() {
			CheckChatRobotTimer()
		})
		(<-d.ChanTimer).Cb()
	}

}

func TimerTaskMulitple() {
	defer common.PrintPanicStack()
	CheckLockMultiple()
}

func TimerTaskNiuNiu() {
	defer common.PrintPanicStack()
	{
		d := timer.NewDispatcher(1)

		d.AfterFunc(3*time.Second, func() {

			//str := "appId=149683301163651&channelOrderNo=4000362001201706084831069650&deviceType=01&funcode=N001&mhtCharset=UTF-8&mhtCurrencyType=156&mhtOrderAmt=1&mhtOrderName=17%E7%8E%A9%E7%9B%B4%E6%92%AD-%E7%8E%B0%E5%9C%A8%E6%94%AF%E4%BB%98%E6%B8%A0%E9%81%93%E5%85%85%E5%80%BC&mhtOrderNo=JH56_1496906691&mhtOrderStartTime=20170608152451&mhtOrderTimeOut=300&mhtOrderType=01&nowPayOrderNo=201001201706081525180441073&payChannelType=13&payConsumerId=o0kRqwCCd4iSs0moaDrwCIX2slq4&signType=MD5&signature=e564d11afba3d6bc310b59a6ba42e89c&tradeStatus=A001"
			//mhtOrderNo, nowPayOrderNo, tradeStatus, payConsumerId := URLDecode(str)
			////fmt.Printf("返回结果appId：%s",strMap["appId"])
			//fmt.Printf("返回结果appId：%s,%s,%s,%s",mhtOrderNo,nowPayOrderNo,tradeStatus,payConsumerId)
			//n0, n1, n2 := GetUserRaiseInfo("104374_1498454608", 104374)
			//fmt.Printf("押分为：%d,%d,%d",n0,n1,n2)
			//fmt.Printf("\n调用TestSignDefaultClient()")
			//common.TestSignDefaultClient()
			//fmt.Printf("start@")
			//BakUserRaisePosInfoWithGameID()
			//fmt.Printf("押分为：finish@")
			//min := GetMinIndexOfThree(1000,10000,10000)
			//
		})
		(<-d.ChanTimer).Cb()
	}

}
