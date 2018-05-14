package controller

import (
	"fmt"
	//"github.com/liudng/godump"
	"github.com/martini-contrib/render"
	"github.com/yshd_game/common"
	"github.com/yshd_game/model"
	"net/http"
	"time"
)

type CreateMultipleReq struct {
	Uid   int    `form:"uid"  binding:"required"`
	Token string `form:"token"  binding:"required"`
}

type GetMultipleReq struct {
	Uid   int    `form:"uid"`
	Token string `form:"token" `
	Rid   string `form:"rid"`
}

type DelMultipleReq struct {
	Uid   int    `form:"uid" binding:"required"`
	Token string `form:"token" binding:"required"`
	Rid   string `form:"rid" binding:"required"`
}

type GenMultpleTokenReq struct {
	Uid    int    `form:"uid"`
	Token  string `form:"token" `
	Rid    string `form:"rid" `
	Perm   string `form:"perm"`
	Expire int64  `form:"expire"`
}

type InviteReq struct {
	Uid   int    `form:"uid"`
	Token string `form:"token" `
	Oid   int    `form:"oid"`
	Rid   string `form:"rid" binding:"required"`
}

type CloseMultpleReq struct {
	Uid   int    `form:"uid"`
	Token string `form:"token" binding:"required"`
	//RoomName string `form:"room_name"`
	Rid string `form:"rid" binding:"required"`
}

type AddInfoMultpleReq struct {
	Uid      int    `form:"uid"  binding:"required"`
	Token    string `form:"token"  binding:"required"`
	Rid      string `form:"rid"  binding:"required"`
	RoomName string `form:"room_name"`
	Location string `form:"location"`
	Save     int    `form:"save"`
}
type MultipleLiveReq struct {
	Uid   int    `form:"uid"`
	Token string `form:"token" `
	Rid   string `form:"rid" binding:"required"`
	Line  int    `form:"line" `
}

type ReadyRoomReq struct {
	Uid      int    `form:"uid" binding:"required"`
	Token    string `form:"token" binding:"required" `
	RoomName string `form:"room_name"`
	Location string `form:"location"`
	Rid      string `form:"rid" `
	Save     int    `form:"save" `
	GameType int    `form:"game_type"`
}

//http://t1.shangtv.cn:3003/multiple/create_room?uid=1&token=f8774b3ff6de892115312d8c77bfa79f
func CreateMultiple(req *http.Request, r render.Render, d CreateMultipleReq) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag] = model.CreateMultipleRoom(d.Uid, 3)
	r.JSON(http.StatusOK, ret_value)
}

//http://shangtv.cn:3000/multiple/get_room?uid=13&token=9e15fef13e644dd8af1ad2902b98b842&rid=1222148479455599569
func GetMutipleRoom(req *http.Request, r render.Render, d GetMultipleReq) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag], ret_value["multiple"] = model.ReqGetMutipleRoom(d.Rid)
	r.JSON(http.StatusOK, ret_value)
}

func DelMultopleRoom(req *http.Request, r render.Render, d DelMultipleReq) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag] = model.DelMutipleRoom(d.Rid)
	r.JSON(http.StatusOK, ret_value)
}

//http://192.168.1.12:3000/multiple/gen_token?uid=16&token=92f220be3f60bca4d13e634ed94113d6
//http://shangtv.cn:3000/multiple/gen_token?uid=16&token=92f220be3f60bca4d13e634ed94113d6
//curl 'http://shangtv.cn:3000/multiple/gen_token?uid=16&token=92f220be3f60bca4d13e634ed94113d6'
func GenQiNiuTokenController(req *http.Request, r render.Render, d GenMultpleTokenReq) {
	ret_value := make(map[string]interface{})

	ret_value[ServerTag], ret_value["token"], ret_value["rid"] = model.GenMutipleToken(d.Uid, d.Rid)
	r.JSON(http.StatusOK, ret_value)
}

func InviteController(req *http.Request, r render.Render, d InviteReq) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag] = model.InviteUser(d.Uid, d.Oid, d.Rid)
	r.JSON(http.StatusOK, ret_value)
}

//  shangtv.cn:3000/multiple/exit?uid=89&token=eb459ac5bc9d311eae7c90eee970905c&rid=89149240812651525
//192.168.1.12:3000/multiple/exit?uid=2&token=16909146765fab9ee50593d6ed567330&room_name=1148126770346818
func CloseMultiple(req *http.Request, r render.Render, d CloseMultpleReq) {
	ret_value := make(map[string]interface{})

	user, _ := model.GetUserByUid(d.Uid)

	has, room := model.GetMultipleRoomByRid(d.Rid)
	if !has {
		ret_value[ServerTag] = common.ERR_ROOM_EXIST
		r.JSON(http.StatusOK, ret_value)
		return
	}

	if room.OwnerId != d.Uid {
		common.Log.Errf("close multiple uid err  ownerid =? and param=?", room.OwnerId, d.Uid)
		ret_value[ServerTag] = common.ERR_MULTIPLE_OWNER
		r.JSON(http.StatusOK, ret_value)
		return
	}

	/*
		if room.Statue != common.MULTIPLE_ROOM_BUSY {
			common.Log.Errf("close multiple statue err ownerid =? , param=?  ,room.Statue=?", room.OwnerId, d.Uid, room.Statue)
			ret_value[ServerTag] = common.ERR_MULTIPLE_STATUES
			r.JSON(http.StatusOK, ret_value)
			return
		}
	*/

	common.Log.Debugf("now time ? begin close multiple start room=?,uid=?", time.Now().Unix(), d.Rid, d.Uid)
	model.AnchorMgr.SetCloseRoom(user.Uid)

	var dur time.Duration
	if room.Statue == common.MULTIPLE_ROOM_FIN {
		dur = room.FinishTime.Sub(room.CreateTime)
	} else {
		dur = time.Now().Sub(room.CreateTime)
	}
	ret_value["time"] = dur / time.Second

	chat := model.GetChatRoom(d.Rid)
	if chat == nil {
		ret_value["time"] = dur / time.Second
		ret_value["count"] = room.Count
		ret_value["rice"] = room.Rice
		ret_value["moon"] = room.Moon
		ret_value[ServerTag] = common.ERR_SUCCESS
		r.JSON(http.StatusOK, ret_value)
		return
	}
	if chat.GetChatInfo().Uid != user.Uid {
		ret_value[ServerTag] = common.ERR_NOT_HOUSE_MANAGER
		r.JSON(http.StatusOK, ret_value)
		return
	}
	ret_value["count"] = chat.GetCount()
	ret_value["rice"] = chat.GetRice()
	ret_value["moon"] = chat.GetMoon()
	//user.SetStatue(common.USER_STATUE_LEAVE)

	//ret_value[ServerTag] = model.CloseMultiple(d.Uid, d.Rid)
	ret_value[ServerTag] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

func ComfirmMultiple(req *http.Request, r render.Render, d AddInfoMultpleReq) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag] = model.ComfirmAddMutipleChat(d.Uid, d.RoomName, d.Location, d.Rid)
	r.JSON(http.StatusOK, ret_value)
}

func GenMultipleLive(req *http.Request, r render.Render, d MultipleLiveReq) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag], ret_value["live_url"] = model.GenMultiplePull(d.Uid, d.Line, d.Rid)
	r.JSON(http.StatusOK, ret_value)
}

func GetMultipleRoomInfo(w http.ResponseWriter, req *http.Request, r render.Render, d CommonReqOnlyUid) {
	ret_value := make(map[string]interface{})

	ret_value["ErrCode"] = common.ERR_SUCCESS

	sess := model.GetUserSessByUid(d.Uid)
	if sess == nil {
		ret_value["ErrCode"] = common.ERR_ROOM_EXIST
		r.JSON(http.StatusOK, ret_value)
		return
	}

	user, _ := model.GetUserByUid(d.Uid)
	if user == nil {
		ret_value["ErrCode"] = common.ERR_ACCOUNT_EXIST
		r.JSON(http.StatusOK, ret_value)
		return
	}
	has, room := model.GetMultipleRoomByRid(sess.Roomid)
	if !has {
		ret_value["ErrCode"] = common.ERR_MULTIPLE_RID
		r.JSON(http.StatusOK, ret_value)
		return
	}
	// type MultipleRoomList struct {
	// 	RoomName   string    `xorm:"varchar(255)"  `            //房间名字
	// 	OwnerId    int       `xorm:"not null "`                 //主播ID
	// 	CreateTime time.Time //创建时间
	// 	FinishTime time.Time //结束时间
	// 	Location   string    `xorm:"varchar(128)"  ` //定位

	// 	Rice       int       //收到的米粒
	// 	Count      int       //人数
	// 	Weight     int       `xorm:"not null default(0)"` //排序权重
	// 	LockTime   int64
	// }
	ret_value["live"] = room.LiveUrl  //直播流
	ret_value["status"] = room.Statue //房间状态

	ret_value["mobile"] = room.MobileUrl //移动流
	ret_value["cover"] = room.Cover      //封面图片
	ret_value["rid"] = sess.Roomid       //房间ID

	ret_value["face"] = user.Image
	ret_value["sex"] = user.Sex
	ret_value["nick_name"] = user.NickName
	ret_value["location"] = user.Location

	chat := model.GetChatRoom(sess.Roomid)

	ret_value["count"] = chat.GetCount()
	r.JSON(http.StatusOK, ret_value)
}

//http://t1.shangtv.cn:3003/multiple/ready_room?uid=9&token=d08b80a85541a706f5cc4ec77228c29d
func ReadyMutipleController(req *http.Request, r render.Render, d ReadyRoomReq) {
	ret_value := make(map[string]interface{})
	var temp int
	var rid string
	temp, ret_value["token"], rid = model.GenMutipleToken(d.Uid, "")

	fmt.Printf("CreateMultiple() d.GameType=%d", d.GameType)
	if d.GameType < 0 || d.GameType > common.GAME_TYPE_MAX { //游戏类型的值做下拦截吧
		d.GameType = 0
	}

	if temp == common.ERR_SUCCESS {
		ret_value["rid"] = rid
		temp, ret_value["live_pull"] = model.GenMultiplePull(d.Uid, 1, rid)
		if temp == common.ERR_SUCCESS {

			ret_value[ServerTag] = model.ReadyMutipleChat(d.Uid, d.RoomName, d.Location, rid, d.Save, d.GameType)
			//godump.Dump(ret_value)
			r.JSON(http.StatusOK, ret_value)
			return
		}
	} else if temp == common.ERR_MULTIPLE_HAS_MIC {
		ret_value["rid"] = rid
		_, ret_value["live_pull"] = model.GenMultiplePull(d.Uid, 1, rid)
	}
	ret_value[ServerTag] = temp
	r.JSON(http.StatusOK, ret_value)
}
