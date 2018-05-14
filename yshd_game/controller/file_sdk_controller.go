package controller

import (
	"github.com/martini-contrib/render"
	"github.com/yshd_game/common"
	"github.com/yshd_game/model"
	"net/http"
	//"github.com/liudng/godump"
)

type UpLoadTokenReq struct {
	Uid    int    `form:"uid" binding:"required"`
	Token  string `form:"token" binding:"required"`
	UpType int    `form:"uptype" binding:"required"` //1 or 2
}

type UpLoadFileReq struct {
	Uid      int    `form:"uid" binding:"required"`
	Token    string `form:"token" binding:"required"`
	UpType   int    `form:"uptype" binding:"required"` //1 or 2
	FileName string `form:"filename" binding:"required"`
}

type QiNiuNotifyReq struct {
	Key string `form:"key" binding:"required"`
	//Hash string `json:"hash" binding:"required"`
	//FileSize int    `json:"filesize" binding:"required"` //1 or 2
	Uid    int    `form:"uid" binding:"required"`
	Token  string `form:"token" binding:"required"`
	Bucket string `form:"bucket" binding:"required"`
}

func GetSevenNiuToken(r render.Render, d UpLoadTokenReq) {
	common.Log.Info("GetSevenNiuToken() called@@@@@@")

	//godump.Dump("online sevrer token")
	//godump.Dump(d)
	ret_value := make(map[string]interface{})
	ret_value["ErrCode"] = common.ERR_SUCCESS
	key := common.GetOnlyId(d.Uid)
	ret_value["token"] = model.Gen7NiuToken(d.UpType, key)
	ret_value["key"] = key
	//godump.Dump(ret_value)
	r.JSON(http.StatusOK, ret_value)
}

func Set7NiuFileName(r render.Render, d UpLoadFileReq) {
	common.Log.Info("Set7NiuFileName() called@@@@@@")

	ret_value := make(map[string]interface{})
	model.SetUserFile(d.UpType, d.Uid, d.FileName)
	ret_value["ErrCode"] = common.ERR_SUCCESS
}

//func QiNiuNotify(req *http.Request, r render.Render, d QiNiuNotifyReq) {
func QiNiuNotify(req *http.Request, r render.Render) {
	common.Log.Info("QiNiuNotify() called@@@@@@")

	ret_value := make(map[string]interface{})
	key := req.FormValue("key")
	uid := req.FormValue("uid")
	token := req.FormValue("utoken")
	buc := req.FormValue("bucket")

	user, ret := model.GetUserByUidStr(uid)
	if ret == common.ERR_SUCCESS {
		if user.Token == token {
			if buc == model.GetQiNiuFace() {

				ret := model.GetQiNiuFile(key, model.GetPicDefine(model.CACHE_PIC_FACE))
				if ret == common.ERR_SUCCESS {
					user.SetFace(model.DownloadUrl(model.DomainFace, key))
				}

				ret_value["ErrCode"] = ret
				ret_value["image"] = user.Image
			} else if buc == model.GetQiNiuCover() {
				//user.SetCover(model.DownloadUrl(model.DomainCover, key))
				ret_value["ErrCode"] = common.ERR_UNKNOWN
			} else if buc == model.GetQiNiuReal() {
				model.SetCachePic(user.Uid, model.DownloadUrl(model.DomainReal, key), model.GetPicDefine(model.CACHE_PIC_REAL))
				ret_value["ErrCode"] = common.ERR_SUCCESS
			} else if buc == model.GetQiNiuIdFront() {
				model.SetCachePic(user.Uid, model.DownloadUrl(model.DomainFront, key), model.GetPicDefine(model.CACHE_PIC_FRONT))
				ret_value["ErrCode"] = common.ERR_SUCCESS
			} else if buc == model.GetQiNiuIdBack() {
				model.SetCachePic(user.Uid, model.DownloadUrl(model.DomainBack, key), model.GetPicDefine(model.CACHE_PIC_BACK))
				//user.SetRealImage(model.DownloadUrl(model.DomainBack, key))
				ret_value["ErrCode"] = common.ERR_SUCCESS
			} else {
				ret_value["ErrCode"] = common.ERR_PARAM
			}
		}
	} else {
		ret_value["ErrCode"] = common.ERR_PARAM
	}
	r.JSON(http.StatusOK, ret_value)
}
