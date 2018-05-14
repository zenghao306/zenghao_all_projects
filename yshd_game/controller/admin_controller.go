package controller

import (
	"github.com/martini-contrib/render"
	"github.com/yshd_game/common"
	//	"github.com/yshd_game/confdata"
	//"encoding/json"
	"errors"
	"fmt"
	"github.com/yshd_game/confdata"
	"github.com/yshd_game/model"
	"github.com/yshd_game/sensitive"
	"net/http"

	"time"
)

type ChatCloseReq struct {
	Uid       int    `form:"uid"`
	Token     string `form:"token" `
	AnchorId  int    `form:"anchorid" binding:"required"`
	CloseType int    `form:"close_type"`
}

type AddDiamondReq struct {
	Uid     int `form:"uid"  binding:"required"`
	Diamond int `form:"diamond" binding:"required"`
	//Intro string `form:"intro" binding:"required"`
}

type SwitchReq struct {
	Switch int `form:"switch"`
}

type WeiXinPayTadeNoReq struct {
	ID   int    `form:"id" binding:"required"`
	Desc string `form:"desc"`
}

type ReportUserReq struct {
	Uid  int    `form:"uid"  binding:"required"`
	Oid  int    `form:"oid" binding:"required"`
	Desc string `form:"desc" `
}

type NoticeToAllReq struct {
	Msg string `form:"msg" binding:"required"`
}

type ForbidUserReq struct {
	Oid        int `form:"oid" binding:"required"`
	Forbid     int `form:"forbid" `
	Forbidtime int `form:"forbid_time" `
	//Intro string `form:"intro" binding:"required"`
}

type SelfReq struct {
	Str string `form:"str" `
}

type RoomReq struct {
	Rid string `form:"str" binding:"required"`
}

func errrr() error {
	var serr error = errors.New(model.ERR_REDIS_STR)
	return serr
}

func SelfTest2(req *http.Request, r render.Render) {
	model.ReportAnchorDate("add", "100000149723667060183", 100000, 0, "ss", time.Now().Unix())
	//model.ReportActionDate("1","view",3)
	model.ReportAnchorDate("add", "100000149723675477091", 100000, 0, "ssdd", time.Now().Unix())
}

func SelfTest3(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	model.GetDaGuanRecommand(185511)
	r.JSON(http.StatusOK, ret_value)
}

//shangtv.cn:3000/admin/SELF?str=sdee
//192.168.1.12:3003/admin/SELF?str=安局豪华
func SelfTest(req *http.Request, r render.Render, d SelfReq) {
	ret_value := make(map[string]interface{})
	/*
		allnum := 20 * 1
		owner_num = float32(allnum) * commossion.OwnerPercent
		other.AddMoney(common.MONEY_TYPE_MOON, int64(owner_num))
		ret_value["ss"]=model.AddNewPlayBack(7,"ss","123456",123456)
	*/
	//ret_value[ServerTag] = model.DelAudienceKey(d.Str)

	//godump.Dump(model.GetPicDefine(model.CACHE_PIC_REAL))

	//model.SetCachePic(100, model.DownloadUrl(model.DomainReal, "ssss"), model.GetPicDefine(model.CACHE_PIC_REAL))

	//model.SetCachePic(100, model.DownloadUrl(model.DomainReal, "ssss"), model.GetPicDefine(model.CACHE_PIC_FRONT))

	//model.SetCachePic(100, model.DownloadUrl(model.DomainReal, "ssss"), model.GetPicDefine(model.CACHE_PIC_BACK))

	/*
		go func() {
			u2, ok := model.GetUserByUid(2)
			for i := 0; i < 15; i++ {
				if ok == common.ERR_SUCCESS {
					s := fmt.Sprintf("user before add money num=%d ,version=%d", u2.Coupons, u2.Version)
					godump.Dump(s)
					ret := u2.AddMoney(common.MONEY_TYPE_RICE, 10)
					time.Sleep(1 * time.Microsecond)
					s = fmt.Sprintf("user  after add money num=%d ,ret=%d,version=%d", u2.Coupons, ret, u2.Version)
					godump.Dump(s)
				}
			}

		}()

		go func() {
			u, ok := model.GetUserByUid(2)
			for i := 0; i < 15; i++ {
				if ok == common.ERR_SUCCESS {
					s := fmt.Sprintf("user before del money num=%d,version=%d", u.Coupons, u.Version)
					godump.Dump(s)
					ret := u.DelMoney(common.MONEY_TYPE_RICE, 10)
					time.Sleep(1 * time.Microsecond)
					s = fmt.Sprintf("user after del money num=%d ,ret=%d,version=%d", u.Coupons, ret, u.Version)
					godump.Dump(s)
				}
			}
		}()
	*/
	//fmt.Print(strings.Trim("   nnnnss  "," "))
	//ret_value[ServerTag]=strings.Trim("   nnnnss  "," ")

	//u, _ := model.GetUserByUid(2)
	//ret_value[ServerTag] = u.Test()

	//var s chan int
	//s=make(chan int,10)
	//go func() {
	//	for v:=range s{
	//		godump.Dump(v)
	//	}
	//}()

	//close(s)
	//s<-1
	if d.Str == "开始补时间" {
		model.AddFinishTime()
	}
	model.InitConsistData()

	r.JSON(http.StatusOK, ret_value)
}

func AudienceCount(req *http.Request, r render.Render, d SelfReq) {
	ret_value := make(map[string]interface{})
	ret_value["count"] = model.GetUserSessCount()
	r.JSON(http.StatusOK, ret_value)
}

func GenNameAdmin(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	ret_value["ErrCode"] = model.GenName()
	r.JSON(http.StatusOK, ret_value)
}

func FlushNotice(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	model.ResetNotice()
	ret_value[ServerTag] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

func AuthReal(req *http.Request, r render.Render, d SwitchReq) {
	ret_value := make(map[string]interface{})

	if d.Switch == 1 {
		common.RealAuthSwitch = true
	} else {
		common.RealAuthSwitch = false
	}

	ret_value[ServerTag] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

func AuthRealStatus(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})

	ret_value["status"] = common.RealAuthSwitch

	ret_value[ServerTag] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

//curl -d
func CloseChat(req *http.Request, r render.Render, d ChatCloseReq) {
	ret_value := make(map[string]interface{})
	info := fmt.Sprintf("anchorid=%d", d.AnchorId)
	model.InsertLogWithIP(common.ACTION_TYPE_LOG_CLOSE_CHAT, 0, common.GetDesc(common.ACTION_TYPE_LOG_CLOSE_CHAT), common.GetRemoteIp(req), info)
	ret_value["ErrCode"] = model.MonitorClose(d.AnchorId, d.CloseType)
	r.JSON(http.StatusOK, ret_value)
}

func SuperUserCloseChat(req *http.Request, r render.Render, d ChatCloseReq) {
	ret_value := make(map[string]interface{})
	user, _ := model.GetUserByUid(d.Uid)
	if !user.IsSuperUser() {
		ret_value["ErrCode"] = common.ERR_NO_POWER_TO_CLOSE_ROOM
		r.JSON(http.StatusOK, ret_value)
	} else {
		info := fmt.Sprintf("anchorid=%d", d.AnchorId)
		model.InsertLogWithIP(common.ACTION_TYPE_LOG_CLOSE_CHAT, d.Uid, common.GetDesc(common.ACTION_TYPE_LOG_CLOSE_CHAT), common.GetRemoteIp(req), info)
		ret_value["ErrCode"] = model.MonitorClose(d.AnchorId, d.CloseType)
		r.JSON(http.StatusOK, ret_value)
	}
}

//192.168.1.12:3003/admin/ReloadAllConfig
func ReloadAllConfig(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	model.InitAdMgr()
	//model.InitVersionData()
	sensitive.InitkeyWord()
	model.InitSayWord()
	model.LoadAllRobotNickName()

	model.LoadGift()
	model.LoadAndroidPay()
	model.LoadIOSPay()
	model.LoadConfigUserExp()
	model.LoadAnchorExp()
	model.LoadScoreExchange()
	model.LoadConfigItem()
	model.SystemVariableInitOrReset()
	model.LoadToyInfo()

	confdata.InitConfData()
	model.LoadConfigTask()
	ret_value[ServerTag] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

func TEST(req *http.Request, r render.Render, d SwitchReq) {

	ret_value := make(map[string]interface{})
	model.ADDTEST()
	if d.Switch == 1 {
		common.RefreshRecommndSwitch = true
	} else {
		common.RefreshRecommndSwitch = false
	}

	ret_value[ServerTag] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

func TestAddDiamond(req *http.Request, r render.Render, d AddDiamondReq) {
	ret_value := make(map[string]interface{})
	user, ret := model.GetUserByUid(d.Uid)
	if ret == common.ERR_SUCCESS {
		//user.addDiamond(d.Diamond)
		user.AddMoney(nil, common.MONEY_TYPE_DIAMOND, int64(d.Diamond), false)
		user.AccountType = 1
		user.UpdateByColS("account_type")
		ret_value[ServerTag] = common.ERR_SUCCESS
	} else {
		ret_value[ServerTag] = common.ERR_UNKNOWN
	}
	r.JSON(http.StatusOK, ret_value)
}

func AuthAccountType(req *http.Request, r render.Render, d SwitchReq) {
	ret_value := make(map[string]interface{})
	if d.Switch == 1 {
		common.AccountAuthSwitch = true
	} else {
		common.AccountAuthSwitch = false
	}
	ret_value[ServerTag] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

//curl -d 'uid=1&oid=10&desc=ssss' 't1.shangtv.cn:3003/report'
func ReportUser(req *http.Request, r render.Render, d ReportUserReq) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag] = model.AddReport(d.Uid, d.Oid, d.Desc)
	r.JSON(http.StatusOK, ret_value)
}

func SetSwitchCashBank(req *http.Request, r render.Render, d SwitchReq) {
	ret_value := make(map[string]interface{})
	if d.Switch == 1 {
		common.CashBankSwitch = true
	} else {
		common.CashBankSwitch = false
	}
	ret_value[ServerTag] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}
func GetSwitchCashBank(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	ret_value["cash_switch"] = common.CashBankSwitch
	ret_value[ServerTag] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

func SetSwitchOnlineShopping(req *http.Request, r render.Render, d SwitchReq) {
	ret_value := make(map[string]interface{})
	if d.Switch == 1 {
		common.OnlineShoppingSwitch = true
	} else {
		common.OnlineShoppingSwitch = false
	}
	ret_value[ServerTag] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}
func GetSwitchOnlineShopping(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	ret_value["online_shopping_switch"] = common.OnlineShoppingSwitch
	ret_value[ServerTag] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

func SetGameRunningSwitch(req *http.Request, r render.Render, d SwitchReq) {
	ret_value := make(map[string]interface{})
	if d.Switch == 1 {
		common.GameRunningSwitch = true
	} else {
		common.GameRunningSwitch = false
		model.CloseAllGameRoom()
	}
	ret_value[ServerTag] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}
func GetSwitchGameRunningSwitch(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	ret_value["online_shopping_switch"] = common.GameRunningSwitch
	ret_value[ServerTag] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

func SendNotice(req *http.Request, r render.Render, d NoticeToAllReq) {
	ret_value := make(map[string]interface{})
	var sys model.ResponseSys
	sys.MType = common.MESSAGE_TYPE_SYS
	sys.Notice = d.Msg
	model.AdminSysToAll(sys)
	ret_value[ServerTag] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

//192.168.1.12:3003/admin/monitor
func MonitorRoom(req *http.Request, r render.Render) {
	ret_value := model.GetMonitorBaseRoom()
	r.JSON(http.StatusOK, ret_value)
}

func RefreshRoomStatus(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})

	ret_value[ServerTag] = model.RefreshRoomStatus()
	r.JSON(http.StatusOK, ret_value)
}

//curl -d 'oid=8&forbid=1'  'http://shangtc.cn:3000/admin/forbid_user_power'
func ForbidAccount(req *http.Request, r render.Render, d ForbidUserReq) {
	ret_value := make(map[string]interface{})

	//model.InsertLogWithIP(common.ACTION_TYPE_LOG_CLOSE_CHAT, 0, common.GetDesc(common.ACTION_TYPE_LOG_CLOSE_CHAT), common.GetRemoteIp(req))
	ret_value["ErrCode"] = model.ForbidUserPower(d.Oid, d.Forbid, d.Forbidtime)
	r.JSON(http.StatusOK, ret_value)
}

func DelMutipleRoom(r render.Render) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag] = model.DelMutipleRoom("Xixi")
	r.JSON(http.StatusOK, ret_value)
}

func WeiXinPay(r render.Render, d WeiXinPayTadeNoReq) {
	ret_value := make(map[string]interface{})

	ret_value[ServerTag] = common.ERR_SUCCESS
	ret_value["ErrCode"] = model.TradeNoToWeiXinPay(d.ID, d.Desc)
	r.JSON(http.StatusOK, ret_value)
}

func RejectWeiXinPay(r render.Render, d WeiXinPayTadeNoReq) {
	ret_value := make(map[string]interface{})

	ret_value[ServerTag] = common.ERR_SUCCESS
	ret_value["ErrCode"] = model.TradeNoRejectWeiXinPay(d.ID, d.Desc)
	r.JSON(http.StatusOK, ret_value)
}

func MoonWeiXinPay(r render.Render, d WeiXinPayTadeNoReq) {
	ret_value := make(map[string]interface{})

	ret_value[ServerTag] = common.ERR_SUCCESS
	ret_value["ErrCode"] = model.MoonToWeiXinPay(d.ID, d.Desc)
	r.JSON(http.StatusOK, ret_value)
}

func GetRoomInfoAdmin(r render.Render, d RoomReq) {
	ret_value := make(map[string]interface{})
	room := model.GetChatRoom(d.Rid)
	if room != nil {
		ret_value["rid"] = room
	}
	r.JSON(http.StatusOK, ret_value)
}

func GetUserInfoAdmin(r render.Render, d CommonReq) {
	ret_value := make(map[string]interface{})
	user := model.GetUserSessByUid(d.Uid)
	if user != nil {
		ret_value["rid"] = user.Roomid
		ret_value["uid"] = user.Uid

		/*
			err := user.Sess.Close()
			if err != nil {
				ret_value["close"] = err.Error()
			}
		*/
	}
	r.JSON(http.StatusOK, ret_value)
}

func CloseUserInfoAdmin(r render.Render, d CommonReq) {
	ret_value := make(map[string]interface{})
	user := model.GetUserSessByUid(d.Uid)
	if user != nil {
		err := user.Sess.Close()
		if err != nil {
			ret_value["close"] = err.Error()
		}
	}
}

type ChannelSwitchReq struct {
	ChannelId string `form:"channel_id" binding:"required"`
}

func ChannelSwitchController(r render.Render, d ChannelSwitchReq) {
	ret_value := make(map[string]interface{})
	ret_value["switch"], ret_value[ServerTag] = model.GetChannelSwitch(d.ChannelId)
	r.JSON(http.StatusOK, ret_value)

}

type WriteLetterUserReq struct {
	Uid    int    `form:"uid"  binding:"required"`
	Oid    int    `form:"oid"`
	Msg    string `form:"msg" `
	Type   int    `form:"type" `
	Family int    `form:"family"`
}

//120.76.156.177:3003/admin/send_letter?uid=1&type=2&msg=nihhhh
////192.168.1.12:3003/admin/send_letter?uid=1&type=1&msg=nihhhh
func AdminWriterLetterController(r render.Render, d WriteLetterUserReq) {
	ret_value := make(map[string]interface{})

	_, ret := model.GetUserByUid(d.Uid)
	if ret == common.ERR_SUCCESS {

		switch d.Type {
		case 1:
			ret = model.SendLetterToAllV2(d.Uid, d.Msg)
		case 2:
			ret = model.SendLetterToAnchor(d.Uid, d.Msg)
		case 3:
			ret = model.SendLetterToFamliy(d.Uid, d.Msg, d.Family)

		}

		//ret := model.SendLetter(user.Uid, d.Oid, d.Msg)
		if ret == common.ERR_SUCCESS {
			sess := model.GetUserSessByUid(d.Oid)
			if sess != nil {
				m := &model.ResponseLetterUnread{}
				m.MType = common.MESSAGE_TYPE_LETTER_UNREAD
				ret, m.Num = model.GetLetterUnreadNum(d.Oid)
				if ret == common.ERR_SUCCESS {
					sess.Sess.SendMsg(m)
				}

			}
			ret_value[ServerTag] = ret
		}

	}

	r.JSON(http.StatusOK, ret_value)
}

func ResetKey(r render.Render) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag] = model.ResetRank()
	r.JSON(http.StatusOK, ret_value)
}
