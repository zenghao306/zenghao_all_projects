package controller

import (
	"github.com/martini-contrib/render"
	//"github.com/martini-contrib/sessions"
	//	"github.com/liudng/godump"
	"github.com/yshd_game/common"
	"github.com/yshd_game/model"
	"net/http"
	"time"
)

type VersionReq struct {
	Platform  int    `form:"platform" binding:"required"`
	Version   string `form:"version" binding:"required"`
	ChannelId string `form:"channel_id" binding:"required"`
}

type TaskVersionReq struct {
	Platform int    `form:"platform" binding:"required"`
	Version  string `form:"version" binding:"required"`
}

type LangeReq struct {
	Lang string `form:"lang" binding:"required"`
}

var (
	cash_explain = "温馨提示：\r\n1、每日累计兑换现金不得超过5000元\r\n" +
		"2、日单笔兑换现金必须大于100元，才能开启兑换。\r\n3、当日收到的星星不可当日兑换现金。\r\n4、如果您在兑换过程中遇到问题，请及时联系官方客服。\r\n官方QQ:800179986"
	//cash_notice  = "每月1-5号为系统结算日，提现功能将关闭"

	englisg_cash_explain = "You will be deducted 3% fee based on law,the actual amount will be less."
	englisg_cash_notice  = "Cash function will be closed on the 1-5 monthly"
)

//shangtv.cn:3003/version?platform=1&version=1.0.13&channel_id=guanfang
//获取版本
func GetVersionInfo(req *http.Request, r render.Render, d VersionReq) {
	ret_value := make(map[string]interface{})
	ret := model.CheckVersion(d.Version, d.ChannelId)
	ret_value[ServerTag] = ret
	if ret != common.ERR_SUCCESS {
		if ret != common.ERR_CONFGI_ITEM {
			ret_value["version"], _, _, _ = model.GetAndroidVersion(d.ChannelId)
		}
	}
	r.JSON(http.StatusOK, ret_value)
}

func CashSayController(req *http.Request, r render.Render, d LangeReq) {

	ret_value := make(map[string]interface{})
	ret_value[ServerTag] = common.ERR_SUCCESS

	ret_value["cash_explain"] = cash_explain

	r.JSON(http.StatusOK, ret_value)
}

const ServerVersion = "20170929_V1.0.8" //服务器版本号
var ServerRestartTime string            //服务器重启时间

// 初始化服务器重启时间变量
func InitServerRestartTimeSet() {
	ServerRestartTime = time.Now().Format("2006-01-02 15:04:05")
}

func ServerController(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	ret_value["server_version"] = ServerVersion          //服务器版本号
	ret_value["server_restart_time"] = ServerRestartTime //服务器最近一次重启时间
	r.JSON(http.StatusOK, ret_value)
}
