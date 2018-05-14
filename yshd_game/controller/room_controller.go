package controller

import (
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/yshd_game/common"
	"github.com/yshd_game/model"
	"mime/multipart"
	"net/http"
	"time"
)

type UploadCoverReq struct {
	Uid   int                   `form:"uid" `
	Token string                `form:"token" binding:"required"`
	File  *multipart.FileHeader `form:"file"`
}

type RecommandListReq struct {
	Index      int    `form:"index" `
	PlayIndex  int    `form:"play_index" `
	Uid        int    `form:"uid"`
	AppVersion string `form:"app_version"`
	Os         int    `form:"os"`
}
type CreateRoomReq struct {
	Uid       int    `form:"uid"`
	Token     string `form:"token" binding:"required"`
	RoomName  string `form:"room_name" `
	Location  string `form:"location"`
	Itype     int    `form:"itype"`
	Save      int    `form:"save"`
	GameType  int    `form:"game_type"`
	Device    string `form:"device"`
	UserAgent string `form:"user_agent"`
}

type CloseRoomReq struct {
	Uid      int    `form:"uid" `
	Token    string `form:"token" binding:"required"`
	RoomId   string `form:"roomid" binding:"required" `
	Multiple int    `form:"multiple"`
}

type PreLiveReq struct {
	Uid   int    `form:"uid"  binding:"required"`
	Token string `form:"token" binding:"required"`
	Line  int    `form:"line" binding:"required"`
	Rid   string `form:"line"`
}

type GetStreamReq struct {
	Roomid string `form:"rid"`
}

type RecommandPlayBackReq struct {
	Uid       int    `form:"uid"`
	Token     string `form:"token" binding:"required"`
	Rid       string `form:"rid" binding:"required"`
	IsMutiple int    `form:"multiple"`
}

//创建房间
//curl -d 'token=ef4929e43770ff59a633bc4bf0b10084&room_name=大米111&location=huoxing' 'http://shangtv.cn:3003/room/create_room'
//curl -d 'token=d4cf6821af1ad5a5fa989c499a55ca81&room_name=eee&location=huoxing&uid=100016' 'http://t1.shangtv.cn:3003/room/create_room'
//curl -d 'token=894235f80c32c1323e57fcd550345db2&room_name=eee&location=huoxing&uid=3' 'http://192.168.1.12:3003/room/create_room'
func CreatRoom(req *http.Request, r render.Render, s sessions.Session, d CreateRoomReq) {
	ret_value := make(map[string]interface{})

	var user *model.User
	if d.Uid == 0 {
		user, _ = model.GetUserByToken(d.Token)
	} else {
		user, _ = model.GetUserByUid(d.Uid)
	}

	//user, _ := model.GetUserByToken(d.Token)
	var cover string
	cover = user.Image
	/*
		if d.Itype == 1 {
			if user.Cover == "" {
				cover = user.Image
			} else {
				cover = user.Cover
			}
		} else if d.Itype == 2 {
			cover = user.Image
		}
	*/

	if d.GameType < 0 || d.GameType > common.GAME_TYPE_MAX { //游戏类型的值做下拦截吧
		d.GameType = 0
	}

	if !common.GameRunningSwitch { //游戏直播间开关关闭情况下
		d.GameType = 0
	}

	weight := model.GetWeightByGroupId2(user.GroupId)
	roomid, err := model.CreateRoom(user.Uid, d.RoomName, cover, d.Location, weight, user.AccountType, d.GameType, d.Device, d.UserAgent)

	if err == common.ERR_SUCCESS {

		c := model.NewChatRoomInfo()

		o := c.GetChatInfo()
		o.Image = user.Image

		o.Rid = roomid
		c.Save = d.Save
		c.RoomType = user.AccountType
		if d.GameType < 0 || d.GameType > common.GAME_TYPE_MAX { //不在定义值值范围内
			d.GameType = 0
		}
		c.GameType = d.GameType

		o.Uid = user.Uid
		c.Statue = common.ROOM_PRE_V2

		o.Image = user.Image

		model.AddChatRoom(c)

		if c.GameType == common.GAME_TYPE_NIUNIU { // 如果客户端是要牛牛添加相应的结构并放到索引里去吧
			r := model.NewRoomInfoNiuNiu()
			r.Rid = roomid
			r.Uid = user.Uid
			model.AddRoomInfoNiuNiu(r)
		} else if c.GameType == common.GAME_TYPE_TEXAS { //如果客户端是要德州扑克添加相应的结构并放到索引里去吧，O(∩_∩)O哈哈~
			r := model.NewRoomInfoTexas()
			r.Rid = roomid
			r.Uid = user.Uid
			model.AddRoomInfoTexas(r)
		} else if c.GameType == common.GAME_TYPE_GOLDEN_FLOWER { // 如果客户端是要砸金花添加相应的结构并放到索引里去吧
			r := model.NewRoomInfoGoldenFlower()
			r.Rid = roomid
			r.Uid = user.Uid
			model.AddRoomInfoGoldenFlower(r)
		}
	}
	ret_value["ErrCode"] = err
	ret_value["roomid"] = roomid

	r.JSON(http.StatusOK, ret_value)
}

//关闭房间
//curl -d 'token=eb459ac5bc9d311eae7c90eee970905c&roomid=89149240812651525' 'http://shangtv.cn:3000/room/close_room'
func CloseRoom(req *http.Request, r render.Render, d CloseRoomReq) {
	ret_value := make(map[string]interface{})
	//user, _ := model.GetUserByToken(d.Token)

	var user *model.User
	if d.Uid == 0 {
		user, _ = model.GetUserByToken(d.Token)
	} else {
		user, _ = model.GetUserByUid(d.Uid)
	}
	chat := model.GetChatRoom(d.RoomId)
	room, has := model.GetRoomById(d.RoomId)
	if !has {
		ret_value[ServerTag] = common.ERR_ROOM_EXIST
		r.JSON(http.StatusOK, ret_value)
		return
	}

	//model.AnchorMgr.SetCloseRoom(user.Uid)
	var dur time.Duration
	if room.Statue == common.ROOM_FINISH {
		dur = room.FinishTime.Sub(room.CreateTime)
	} else {
		dur = time.Now().Sub(room.CreateTime)
	}
	ret_value["time"] = dur / time.Second

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

	ret_value["count"] = chat.GetCount() + chat.GetVRobotCount()
	ret_value["rice"] = chat.GetRice()
	ret_value["moon"] = chat.GetMoon() + model.GetDump(room.RoomId)

	ret_value["dump"] = model.GetDump(room.RoomId)
	// added by zenghao 2017.05.03[解决主播主动关闭房间后还能在列表里看到房间信息的问题]

	sess := model.GetUserSessByUid(user.Uid)
	if sess != nil {
		s := sess.Sess
		model.DirectCloseRoom(user.Uid, d.RoomId)
		s.Close()
	} else {
		model.DirectCloseRoom(user.Uid, d.RoomId)
	}

	r.JSON(http.StatusOK, ret_value)
}

//房间列表
//curl 'http://120.76.96.73:3000/list_room?&index=0'
func ListRoom(req *http.Request, r render.Render) {

	ret_value := make(map[string]interface{})

	//ret_value["RoomList"], ret_value["ErrCode"] = model.GetRoomList(index)
	ret_value["ErrCode"] = common.ERR_SUCCESS
	ret_value["RoomList"] = model.GetRoomCache()
	r.JSON(http.StatusOK, ret_value)
}

func ListRoomRealUserCount(req *http.Request, r render.Render) {

	ret_value := make(map[string]interface{})

	//ret_value["RoomList"], ret_value["ErrCode"] = model.GetRoomList(index)
	//ret_value["ErrCode"] = common.ERR_SUCCESS
	ret_value["RoomList"], ret_value["total_num"], ret_value["ErrCode"] = model.GetRoomRealUserCountList()
	r.JSON(http.StatusOK, ret_value)
}

/*
//上传封面
func UploadCover(w http.ResponseWriter, req *http.Request, r render.Render, d UploadCoverReq) {
	godump.Dump(d)
	ret_value := make(map[string]interface{})
	user, _ := model.GetUserByToken(d.Token)

	path := common.Cfg.MustValue("path", "cover")

	if d.File == nil {
	//	user.SetCover(user.Image)
		ret_value["ErrCode"] = common.ERR_SUCCESS
		r.JSON(http.StatusOK, ret_value)

		return
	}

	filename, ret := common.UploadBinding(d.File, path, user.Uid)
	if !ret {
		http.Error(w, "bad upload", 400)
		return
	}
	//user.SetCover(path + filename)
	ret_value["ErrCode"] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}
*/

//开始直播预创建
//curl -d 'token=4434da173909053d062993e64eef2995&uid=8&line=1'  'http://shangtv.cn:3003/room/live_create'
//curl -d 'token=894235f80c32c1323e57fcd550345db2&uid=3&line=1' 'http://192.168.1.12:3003/room/live_create'
//curl -d 'token=d4cf6821af1ad5a5fa989c499a55ca81&uid=100016&line=1' 'http://t1.shangtv.cn:3003/room/live_create'
func LiveCreate(req *http.Request, r render.Render, s sessions.Session, d PreLiveReq) {
	ret_value := make(map[string]interface{})

	var user *model.User
	if d.Uid == 0 {
		user, _ = model.GetUserByToken(d.Token)
	} else {
		user, _ = model.GetUserByUid(d.Uid)
	}

	if common.RealAuthSwitch {
		if user.CheckAuthReal() == false {
			if model.CheckAuthRecord(d.Uid) {
				ret_value["ErrCode"] = common.ERR_AUTH_VET
				r.JSON(http.StatusOK, ret_value)
				return
			}

			ret_value["ErrCode"] = common.ERR_AUTH_REAL
			r.JSON(http.StatusOK, ret_value)
			return
		}
	}

	if ret := model.CheckLeader(user.Uid); ret != common.ERR_SUCCESS {
		ret_value["ErrCode"] = ret
		r.JSON(http.StatusOK, ret_value)
		return
	}

	if sess := model.GetUserSessByUid(user.Uid); sess != nil {
		if sess.Sess.IsClosed() {
			model.DelUserSession(user.Uid)
			ret_value["ErrCode"] = common.ERR_USER_SESS
			r.JSON(http.StatusOK, ret_value)
			return
		}
		err := sess.Sess.Close()
		if err != nil {
			common.Log.Errf("sess exist %s", err.Error())
		}
		ret_value["ErrCode"] = common.ERR_USER_SESS
		r.JSON(http.StatusOK, ret_value)
		return
	} else {

		if r_status := model.CheckRoomStatus(d.Uid); r_status != common.ERR_SUCCESS {
			ret_value["ErrCode"] = r_status
			r.JSON(http.StatusOK, ret_value)
			return
		}
	}

	if ret := model.AnchorMgr.SetCloseRoom(user.Uid); ret == true {
		ret_value["ErrCode"] = common.ERR_USER_RECONNECT
		r.JSON(http.StatusOK, ret_value)
		return
	}

	if forbid := user.CheckAccount(); forbid == true {
		ret_value["ErrCode"] = common.ERR_FORBID
		r.JSON(http.StatusOK, ret_value)
		return
	}

	live, err := model.PreCreateLiveUrl(user, d.Line)

	ret_value["ErrCode"] = err
	ret_value["url"] = live

	r.JSON(http.StatusOK, ret_value)
}

/*
//取消上传封面
//curl -d 'token=35fc728908fac142b965d4ecdd17ff3c' 'http://192.168.1.12:3000/room/exit_create'
func CancleUploadCover(w http.ResponseWriter, req *http.Request, r render.Render) {
	common.Log.Info("CancleUploadCover () called@@@@@@")
	ret_value := make(map[string]interface{})
	token := req.FormValue("token")
	user, _ := model.GetUserByToken(token)
	//path := common.Cfg.MustValue("path", "cover")
	//save_path := common.StaticPath + path
	if user.Cover == "" {
		ret_value[ServerTag] = common.ERR_UPLOAD_EMPTY
		r.JSON(http.StatusOK, ret_value)
		return
	}
	//_, cover := model.GetBucket()
	//model.DelQiNiuFile(cover, user.Cover)

	ret_value["ErrCode"] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}
*/
/*
func JoinRoom(w http.ResponseWriter, req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	//token := req.FormValue("token")
	//user, _ := model.GetUserByToken(token)

	rid := req.FormValue("rid")
	rid_, _ := strconv.Atoi(rid)
	room := model.GetChatRoom(rid_)
	if room != nil {
		ret_value["room"] = room
		ret_value["ErrCode"] = common.ERR_SUCCESS
	} else {
		ret_value["ErrCode"] = common.ERR_UNKNOWN
	}
	r.JSON(http.StatusOK, ret_value)
}
*/
//curl 'http://192.168.1.12:3000/room_info?uid=8'
func GetRoomInfo(w http.ResponseWriter, req *http.Request, r render.Render, d CommonReqOnlyUid) {
	ret_value := make(map[string]interface{})
	//	rid := req.FormValue("rid")
	//rid_, _ := strconv.Atoi(rid)
	ret_value["ErrCode"] = common.ERR_SUCCESS

	//ret_value["info"], _ = model.GetRoomById(rid_)
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

	//room, has := chat_room_manager[rid]
	//if chat.Session
	room1 := model.GetChatRoom(sess.Roomid)
	if room1.IsMultiple { //如果是会议室
		has, room := model.GetMultipleRoomByRid(sess.Roomid)
		if !has {
			ret_value["ErrCode"] = common.ERR_MULTIPLE_RID
			r.JSON(http.StatusOK, ret_value)
			return
		}
		ret_value["is_multiple"] = true
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
	} else { //如果是普通直播
		room, _ := model.GetRoomById(sess.Roomid)
		if room == nil {
			ret_value["ErrCode"] = common.ERR_ROOM_EXIST
			r.JSON(http.StatusOK, ret_value)
			return
		}

		ret_value["is_multiple"] = false
		ret_value["live"] = room.LiveUrl
		ret_value["status"] = room.Statue

		//path := common.Cfg.MustValue("video", "video_addr")
		//port := common.Cfg.MustValue("video", "video_port")
		//ret_value["mobile"] = fmt.Sprintf("http://%s:%s/hls/%s.m3u8", path, port, rid)
		ret_value["mobile"] = room.MobileUrl
		ret_value["cover"] = room.Cover
		ret_value["rid"] = sess.Roomid
		ret_value["face"] = user.Image
		ret_value["sex"] = user.Sex
		ret_value["nick_name"] = user.NickName
		ret_value["location"] = user.Location
		ret_value["game_type"] = room.GameType

		chat := model.GetChatRoom(sess.Roomid)

		ret_value["count"] = chat.GetCount()
		r.JSON(http.StatusOK, ret_value)
	}

}

func GetRoomInfoByRid(w http.ResponseWriter, req *http.Request, r render.Render, d CommonReqOnlyRid) {
	ret_value := make(map[string]interface{})
	v, has := model.GetPlayBack(d.Rid)
	if has == false {
		ret_value["ErrCode"] = common.ERR_ROOM_EXIST
		r.JSON(http.StatusOK, ret_value)
		return
	}

	user, _ := model.GetUserByUid(v.Uid)
	if user == nil {
		ret_value["ErrCode"] = common.ERR_ACCOUNT_EXIST
		r.JSON(http.StatusOK, ret_value)
		return
	}
	room, _ := model.GetRoomById(d.Rid)

	ret_value["live"] = v.PlayUrl
	ret_value["status"] = common.USER_STATUE_RELIVE
	ret_value["cover"] = room.Cover
	ret_value["sex"] = user.Sex
	ret_value["nick_name"] = user.NickName
	ret_value["location"] = user.Location
	ret_value["rid"] = d.Rid
	r.JSON(http.StatusOK, ret_value)
}

//curl 'http://192.168.1.12:3000/recommand_list?uid=100010'
func RecommandListRoom(r render.Render, d RecommandListReq) {
	ret_value := make(map[string]interface{})

	//ret_value["RoomList"], ret_value["ErrCode"] = model.GetRecommandList(d.Index)
	rtype := 0
	if d.Uid != 0 {
		user, ret := model.GetUserByUid(d.Uid)
		if ret == common.ERR_SUCCESS {
			rtype = user.AccountType
			ret_value["ErrCode"] = common.ERR_SUCCESS
		} else {
			ret_value["ErrCode"] = ret
			r.JSON(http.StatusOK, ret_value)
			return
		}
	} else {
		rtype = common.ROOM_TYPE_NOMARL
	}
	//rtype=1+rtype
	/*
		res:=make([]model.RecommandRoomRes,0)
		test:=make([]model.RecommandRoomRes,0)

		res,test=model.GetRecommandListV2()

		if rtype==1 {
			ret_value["RoomList"]=test
		}else{
			ret_value["RoomList"]=res
		}
	*/
	//_, ret_value["MultipleList"] = model.GetRecommandCache(rtype)
	/*
		if d.AppVersion=="" {
			res,ret:=model.GetRecommandListWithGameType(rtype,1,2)
			if ret==common.ERR_SUCCESS {
				ret_value["RoomList"]=res

				ret_value["MultipleList"]=model.GetMutipleTest()
			}
		}else{
			ret_value["RoomList"], ret_value["MultipleList"] = model.GetRecommandCache(rtype)
		}
	*/
	// 如果APP版本字符串传入的不为空，则表明是最新的版本[带上游戏3]
	if d.AppVersion != "" {
		ret_value["RoomList"], ret_value["MultipleList"] = model.GetRecommandCache(rtype, d.Uid)
	} else { //不带游戏3
		ret_value["RoomList"], ret_value["MultipleList"] = model.GetRecommandCache2(rtype, d.Uid)
	}
	//ret_value["RoomList"], ret_value["MultipleList"] = model.GetRecommandCache(rtype)

	ret_value["PlayBackList"], _ = model.GetRecommandWithPlayUrl(d.PlayIndex)
	ret_value["PlayBackMultipleList"], _ = model.GetRecommandMutipleWithPlayUrl(d.PlayIndex)
	r.JSON(http.StatusOK, ret_value)
}

func RecommandListRoomV2(r render.Render, d RecommandListReq) {
	ret_value := make(map[string]interface{})
	/*
		res:=make([]model.RecommandRoomRes,0)
		test:=make([]model.RecommandRoomRes,0)
		res,test=	model.GetRecommandListV2()
		if len(res)==0 {
			res2:=make([]model.RecommandRoomRes,0)
			ret_value["RoomList"]=res2
		}else{
			ret_value["RoomList"]=res
		}

		if len(test)==0 {
			test2:=make([]model.RecommandRoomRes,0)
			ret_value["RoomList"]=test2
		}else{
			ret_value["TestRoomList"]=test
		}

	*/
	r.JSON(http.StatusOK, ret_value)
}

//curl -d 'token=15311111608&uid=9&line=1&rid=5525' 'http://192.168.1.12:3000/room/live_create'
func PullAddr(r render.Render, d PreLiveReq) {
	common.Log.Info("PullAddr () called@@@@@@")
	ret_value := make(map[string]interface{})

	ret_value["ErrCode"], ret_value["push_url"] = model.GenPullAddr(d.Uid, d.Line, d.Rid)
	r.JSON(http.StatusOK, ret_value)

}

func SaveM8u3(r render.Render, f GetStreamReq) {
	ret_value := make(map[string]interface{})
	ret_value["ErrCode"] = model.GetStream(f.Roomid)
	r.JSON(http.StatusOK, ret_value)
}

//http://192.168.1.12:3000/room/play_list?uid=8&token=92332cf955cce0bf03de5b3afbe51273
func PlayBackList(r render.Render, d CommonReq) {
	ret_value := make(map[string]interface{})
	ret_value["ErrCode"] = common.ERR_SUCCESS
	ret_value["plist"] = model.GetPlayBackList(d.Uid)
	ret_value["multiple_plist"] = model.GetMultiplePlayBackList(d.Uid)
	extra, has := model.GetUserExtraByUid(d.Uid)
	if has {
		ret_value["recommand"] = extra.PlayBackRecommandRid
		ret_value["recommand_type"] = extra.PlayBackRecommandType
	}

	r.JSON(http.StatusOK, ret_value)
}

//curl -d 'uid=8&token=92332cf955cce0bf03de5b3afbe51273&roomid=12' 'http://192.168.1.12:3000/room/del_play'
func DelPlayBack(r render.Render, d CloseRoomReq) {
	ret_value := make(map[string]interface{})

	if d.Multiple == 0 {
		ret_value[ServerTag] = model.DelPlayBackList(d.Uid, d.RoomId)
	} else {
		ret_value[ServerTag] = model.DelMutiplePlayBack(d.Uid, d.RoomId)
	}

	r.JSON(http.StatusOK, ret_value)
}

func CheckPlayBack(r render.Render, d CommonReq) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag] = model.CheckPlayList(d.Uid)
	r.JSON(http.StatusOK, ret_value)
}

//curl -d 'uid=8&token=92332cf955cce0bf03de5b3afbe51273&rid=13' 'http://192.168.1.12:3000/room/recommand'
func RecommandPlayBack(r render.Render, d RecommandPlayBackReq) {
	ret_value := make(map[string]interface{})
	if d.IsMutiple == 1 {
		ret_value[ServerTag] = model.UpdateMutipleRecommandFlag(d.Uid, d.Rid)
	} else {
		ret_value[ServerTag] = model.UpdateRecommandFlag(d.Uid, d.Rid)
	}
	r.JSON(http.StatusOK, ret_value)
}

func CancelRecommand(r render.Render, d CommonReq) {
	ret_value := make(map[string]interface{})
	ret_value[ServerTag] = model.CancelRecommandFlag(d.Uid)
	r.JSON(http.StatusOK, ret_value)
}

func FreshPlayBack(r render.Render, f GetStreamReq) {
	ret_value := make(map[string]interface{})
	model.CreatePlayBackRoom()
	ret_value["ErrCode"] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

func MarkRecommandPlayBack(r render.Render, f GetStreamReq) {
	ret_value := make(map[string]interface{})

	ret_value["ErrCode"] = model.HiddenPlayBack(f.Roomid)
	r.JSON(http.StatusOK, ret_value)
}
func CheckAuthRealController(req *http.Request, r render.Render, d CommonReq) {
	ret_value := make(map[string]interface{})

	user, _ := model.GetUserByUid(d.Uid)
	if common.RealAuthSwitch {
		if user.CheckAuthReal() == false {

			if model.CheckAuthRecord(d.Uid) {
				ret_value["ErrCode"] = common.ERR_AUTH_VET
				r.JSON(http.StatusOK, ret_value)
				return
			}

			ret_value["ErrCode"] = common.ERR_AUTH_REAL
			r.JSON(http.StatusOK, ret_value)
			return
		}
	}

	if ret := model.CheckLeader(user.Uid); ret != common.ERR_SUCCESS {
		ret_value["ErrCode"] = ret
		r.JSON(http.StatusOK, ret_value)
		return
	}

	ret_value["ErrCode"] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}
