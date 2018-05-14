package controller

import (
	//"github.com/liudng/godump"
	//"github.com/liudng/godump"
	"github.com/martini-contrib/render"
	"github.com/yshd_game/common"
	"github.com/yshd_game/model"
	"net/http"
)

type WinXinLoginReq struct {
	AccessToken  string `form:"access_token" binding:"required"`
	FreashToken  string `form:"refresh_token"  binding:"required"`
	OpneId       string `form:"openid"  binding:"required"`
	ChannelId    string `form:"channel_id" `
	Deviec       string `form:"device" `
	RegisterFrom int    `form:"register_from"  `
	Imei         string `form:"imei"`
}

type QQLoginReq struct {
	AccessToken string `form:"access_token" binding:"required"`
	//OpneId      string `form:"openid"  binding:"required"`
	Type         int    `form:"type"  `
	ChannelId    string `form:"channel_id" `
	Deviec       string `form:"device" `
	RegisterFrom int    `form:"register_from"  `
}

type SinaLoginReq struct {
	AccessToken  string `form:"access_token" binding:"required"`
	Uid          string `form:"uid"  binding:"required"`
	Type         int    `form:"type"  `
	ChannelId    string `form:"channel_id" `
	Deviec       string `form:"device" `
	RegisterFrom int    `form:"register_from"  `
}

type FaceBookLoginReq struct {
	OpenID       string `form:"openid"  binding:"required"`
	City         string `form:"city" `
	HeadImg      string `form:"head_img" `
	NickName     string `form:"nick_name" `
	ChannelId    string `form:"channel_id" `
	Deviec       string `form:"device" `
	Sex          int    `form:"sex"  `
	RegisterFrom int    `form:"register_from"  `
}

type TwitterLoginReq struct {
	OpenID       string `form:"openid"  binding:"required"`
	City         string `form:"city" `
	HeadImg      string `form:"head_img" `
	NickName     string `form:"nick_name" `
	ChannelId    string `form:"channel_id" `
	Deviec       string `form:"device" `
	Sex          int    `form:"sex"  `
	RegisterFrom int    `form:"register_from"  `
}

type HuaWeiLoginReq struct {
	OpenID       string `form:"openid"  binding:"required"`
	HeadImg      string `form:"head_img" `
	NickName     string `form:"nick_name" `
	ChannelId    string `form:"channel_id" `
	Deviec       string `form:"device" `
	Sex          int    `form:"sex"  `
	RegisterFrom int    `form:"register_from"  `
}

type ChannelReq struct {
	ChannelName string `form:"channel_name" binding:"required"`
	Device      string `form:"device"  binding:"required"`
}

type UvReq struct {
	Imei      string `form:"imei" binding:"required"`
	Device    string `form:"device"  binding:"required"`
	ChannelId string `form:"channel_id"  binding:"required"`
}

type UIDReq struct {
	Uid int `form:"uid" binding:"required"`
}

func ChecOtherkLoginIn(uid int) {
	defer common.PrintPanicStack()
	sess_user := model.GetUserSessByUid(uid)
	if sess_user != nil {
		var data model.ResponseErr

		data.MType = common.MESSAGE_TYPE_ERR

		u, ret := model.GetUserByUid(uid)
		if ret == common.ERR_SUCCESS {
			if u.Token == sess_user.Token {
				data.Err = common.ERR_SELF_LOGIN
			} else {
				data.Err = common.ERR_OTHER_LOGIN
			}
		}

		if sess_user.Sess.IsClosed() {
			model.DelUserSession(uid)
			common.Log.Errf("sess exist err direct del %d", sess_user.Uid)
			return
		}

		err := sess_user.Sess.CloseWithMsgAndJson(data)
		if err != nil {
			common.Log.Errf("sess exist err is %s", err.Error())
		}
	}
}

func AdminCloseSession(r render.Render, d UIDReq) {
	ret_value := make(map[string]interface{})
	user, ret := model.GetUserByUid(d.Uid)
	if ret == common.ERR_SUCCESS {
		ret_value[ServerTag] = model.ChartRoomClose(user.Uid)
		ChecOtherkLoginIn(user.Uid)
	} else {
		ret_value[ServerTag] = ret
	}
	r.JSON(http.StatusOK, ret_value)
}

//curl -d 'access_token=1&refresh_token=962c482aa556a59d06392c32a93d00d3&openid=5' 'http://192.168.1.12:3000/third_login/weixin_login'
func WeiXinLogin(req *http.Request, r render.Render, d WinXinLoginReq) {
	ret_value := make(map[string]interface{})
	if err := model.CheckWeixinToekn(d.AccessToken, d.OpneId); err == 0 {
		err_no, union_id, register := model.GetWeiXinUserinfo(d.AccessToken, d.OpneId, d.ChannelId, d.Deviec, d.RegisterFrom)
		ret_value[ServerTag] = err_no
		ret_value["register"] = register
		if err_no == 0 {
			//user, _ := model.GetUserByAccount(d.OpneId)
			user, has := model.GetUserByAccountAndPlatfrom(union_id, common.PLATFORM_WEIXIN)
			if has {
				info := &model.LoginInfo{}
				token := common.GenUserToken(user.Uid)
				user.SetToken(token)
				user.GetLoginInfo(info)
				ret_value["user"] = info

				model.DailyRecordLog(user.Uid)

				go ChecOtherkLoginIn(user.Uid)
			}

		}
	} else {
		ret_value[ServerTag] = common.ERR_WEIXIN_TOKEN_FAILED
	}
	r.JSON(http.StatusOK, ret_value)
}

func FaceBookLogin(r render.Render, d FaceBookLoginReq) {

	ret_value := make(map[string]interface{})

	err_no := model.GetThirdPartyUserInfo(d.OpenID, d.City, d.HeadImg, d.NickName, d.ChannelId, d.Deviec, d.Sex, common.PLATFORM_FACE_BOOK, d.RegisterFrom)
	ret_value[ServerTag] = err_no

	if err_no == 0 {
		user, _ := model.GetUserByAccountAndPlatfrom(d.OpenID, common.PLATFORM_FACE_BOOK)
		info := &model.LoginInfo{}
		token := common.GenUserToken(user.Uid)
		user.SetToken(token)
		user.GetLoginInfo(info)
		ret_value["user"] = info

		model.DailyRecordLog(user.Uid)

		go ChecOtherkLoginIn(user.Uid)
	}
	r.JSON(http.StatusOK, ret_value)
}

func FaceBookLoginGet(req *http.Request, r render.Render, d FaceBookLoginReq) {
	ret_value := make(map[string]interface{})

	err_no := model.GetThirdPartyUserInfo(d.OpenID, d.City, d.HeadImg, d.NickName, d.ChannelId, d.Deviec, d.Sex, common.PLATFORM_FACE_BOOK, d.RegisterFrom)
	ret_value[ServerTag] = err_no

	if err_no == 0 {
		user, _ := model.GetUserByAccountAndPlatfrom(d.OpenID, common.PLATFORM_FACE_BOOK)
		info := &model.LoginInfo{}
		token := common.GenUserToken(user.Uid)
		user.SetToken(token)
		user.GetLoginInfo(info)
		ret_value["user"] = info

		model.DailyRecordLog(user.Uid)

		go ChecOtherkLoginIn(user.Uid)
	}
	r.JSON(http.StatusOK, ret_value)
}

func TwitterLogin(r render.Render, d TwitterLoginReq) {
	ret_value := make(map[string]interface{})

	err_no := model.GetThirdPartyUserInfo(d.OpenID, d.City, d.HeadImg, d.NickName, d.ChannelId, d.Deviec, d.Sex, common.PLATFORM_TWITTER, d.RegisterFrom)
	ret_value[ServerTag] = err_no
	if err_no == 0 {
		user, _ := model.GetUserByAccountAndPlatfrom(d.OpenID, common.PLATFORM_TWITTER)
		info := &model.LoginInfo{}
		token := common.GenUserToken(user.Uid)
		user.SetToken(token)
		user.GetLoginInfo(info)
		ret_value["user"] = info

		model.DailyRecordLog(user.Uid)

		go ChecOtherkLoginIn(user.Uid)
	}

	r.JSON(http.StatusOK, ret_value)
}

func HuaWeiLogin(r render.Render, d HuaWeiLoginReq) {
	ret_value := make(map[string]interface{})
	err_no := model.GetThirdPartyUserInfo(d.OpenID, "huoxing", d.HeadImg, d.NickName, d.ChannelId, d.Deviec, d.Sex, common.PLATFORM_HUAWEI, d.RegisterFrom)
	ret_value[ServerTag] = err_no
	if err_no == 0 {
		user, _ := model.GetUserByAccountAndPlatfrom(d.OpenID, common.PLATFORM_HUAWEI)
		info := &model.LoginInfo{}
		token := common.GenUserToken(user.Uid)
		user.SetToken(token)
		user.GetLoginInfo(info)
		ret_value["user"] = info

		model.DailyRecordLog(user.Uid)

		go ChecOtherkLoginIn(user.Uid)
	}

	r.JSON(http.StatusOK, ret_value)
}

//curl -d 'access_token=BAA789BB3217521DEF2303BFD19A43FB&type=2&channel_id=self&register_from=1' 'http://192.168.1.12:3003/third_login/qq_login'
//curl -d 'access_token=BAA789BB3217521DEF2303BFD19A43FB&type=2&channel_id=self&register_from=1' 'http://120.76.156.177:3003/third_login/qq_login'
func QQLogin(r render.Render, d QQLoginReq) {
	ret_value := make(map[string]interface{})
	unionid, err_no, register := model.QQLogin(d.AccessToken, d.Type, d.ChannelId, d.Deviec, d.RegisterFrom)
	if err_no == common.ERR_SUCCESS {

		ret_value["register"] = register
		info := &model.LoginInfo{}
		//user, _ := model.GetUserByAccount(d.OpneId)
		user, _ := model.GetUserByAccountAndPlatfrom(unionid, common.PLATFORM_QQ)
		token := common.GenUserToken(user.Uid)
		user.SetToken(token)
		user.GetLoginInfo(info)
		ret_value["user"] = info

		model.DailyRecordLog(user.Uid)

		go ChecOtherkLoginIn(user.Uid)
	}
	ret_value[ServerTag] = err_no
	r.JSON(http.StatusOK, ret_value)
}

/*
func QQH5Login(r render.Render, d QQLoginReq) {

	ret_value := make(map[string]interface{})
	unionid, err_no := model.QQLogin(d.AccessToken, d.Type, d.ChannelId, d.Deviec, d.RegisterFrom,d.Imei)
	if err_no == common.ERR_SUCCESS {
		info := &model.LoginInfo{}
		//user, _ := model.GetUserByAccount(d.OpneId)
		user, _ := model.GetUserByAccountAndPlatfrom(unionid, common.PLATFORM_QQ)
		token := common.GenUserToken(user.Uid)
		user.SetToken(token)
		user.GetLoginInfo(info)
		ret_value["user"] = info

		model.DailyRecordLog(user.Uid)

		go ChecOtherkLoginIn(user.Uid)
	}
	ret_value[ServerTag] = err_no
	r.JSON(http.StatusOK, ret_value)

}
*/
func SinaLogin(r render.Render, d SinaLoginReq) {
	ret_value := make(map[string]interface{})
	err_no := model.SinaSDKLogin(d.AccessToken, d.Uid, d.ChannelId, d.Deviec, d.RegisterFrom)
	if err_no == common.ERR_SUCCESS {
		info := &model.LoginInfo{}
		//user, _ := model.GetUserByAccountWithSinaPlatform(d.AccessToken)
		user, _ := model.GetUserByAccountAndPlatfrom(d.AccessToken, common.PLATFORM_SINA)
		token := common.GenUserToken(user.Uid)
		user.SetToken(token)
		user.GetLoginInfo(info)
		ret_value["user"] = info

		model.DailyRecordLog(user.Uid)

		go ChecOtherkLoginIn(user.Uid)
	}
	ret_value[ServerTag] = err_no
	r.JSON(http.StatusOK, ret_value)
}

func RegisterDevice(r render.Render, d ChannelReq) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag] = model.RegisterNewChannel(d.ChannelName, d.Device)
	r.JSON(http.StatusOK, ret_value)
}

//curl -d 'imei=a10501986d2d4ef7a9ccc6b121eb38c9&device=iphone9&channel_id=5926707307fe65561d000dd9'  '192.168.1.12:3003/uv'
//curl -d 'imei=323fb5d800ae45b6b4c5304a985f5232&device=iphone9&channel_id=5926707307fe65561d000dd9' 'http://shangtv.cn:3003/uv'
func AddUv(r render.Render, d UvReq) {
	//godump.Dump(d)
	ret_value := make(map[string]interface{})
	model.GetUvInfo(d.Imei, d.Device, d.ChannelId)
	ret_value[ServerTag] = model.AddImei(d.Imei, d.Device, d.ChannelId)

	r.JSON(http.StatusOK, ret_value)
}
