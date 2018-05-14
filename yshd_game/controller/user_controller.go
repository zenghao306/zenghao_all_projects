package controller

import (
	"fmt"
	//"github.com/liudng/godump"
	"github.com/martini-contrib/render"
	"github.com/yshd_game/common"
	"github.com/yshd_game/model"
	"github.com/yshd_game/sensitive"
	"github.com/yshd_game/wrap"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

var (
	ServerTag = "ErrCode"
)

/*
func UnserizeToUTF8(rs string) {
	regexp.ReplaceAllString()
}
*/
type SendGiftReq struct {
	Uid    int    `form:"uid" `
	Token  string `form:"token" binding:"required"`
	GiftId int    `form:"giftid" binding:"required"`
	Num    int    `form:"num" `
	Revid  int    `form:"revid" binding:"required"`
	Times  int    `form:"times" `
}

type ModifyCharReq struct {
	Uid      int    `form:"uid" `
	Token    string `form:"token" binding:"required"`
	NickName string `form:"nickname" binding:"required"`
	Sex      int    `form:"sex" `
}

type ModifyPwdReq struct {
	Uid    int    `form:"uid" `
	Token  string `form:"token" binding:"required"`
	OldPwd string `form:"oldpwd" binding:"required"`
	NewPwd string `form:"newpwd" binding:"required"`
}

type ModifyInfoNickReq struct {
	Uid      int    `form:"uid" `
	Token    string `form:"token" binding:"required"`
	NickName string `form:"nick_name" binding:"required"`
}

type GetResetNickReq struct {
	Uid   int    `form:"uid" `
	Token string `form:"token" binding:"required"`
}

type ModifyInfoSexReq struct {
	Uid   int    `form:"uid" `
	Token string `form:"token" binding:"required"`
	Sex   int    `form:"sex" `
}

type ModifyInfoLocationReq struct {
	Uid      int    `form:"uid" `
	Token    string `form:"token" binding:"required"`
	Location string `form:"location" binding:"required"`
}

type ModifyInfoSignReq struct {
	Uid   int    `form:"uid" `
	Token string `form:"token" binding:"required"`
	Sign  string `form:"sign" binding:"required"`
}

type SetPushReq struct {
	Uid   int    `form:"uid" `
	Token string `form:"token" binding:"required"`
	Push  int    `form:"push" binding:"required"`
}

type UploadFaceReq struct {
	Uid   int                   `form:"uid" `
	Token string                `form:"token" binding:"required"`
	File  *multipart.FileHeader `form:"file" binding:"required"`
}
type CommonReqOnlyUid struct {
	Uid int `form:"uid"  binding:"required"`
}
type CommonReqOnlyRid struct {
	Rid string `form:"rid"  binding:"required"`
}

type CommonReqWithRid struct {
	Uid   int    `form:"uid" binding:"required"`
	Token string `form:"token" binding:"required"`
	Rid   string `form:"rid"  binding:"required"`
}

type LiveStatisticsReq struct {
	Uid   int    `form:"uid" `
	Token string `form:"token" binding:"required"`
	//Date  string `form:"date"  binding:"required"`
	MonthIndex int `form:"month_index"  `
	Index      int `form:"index"  `
}

/*
type Call struct {
	Render render.Render
	Req    SendGiftReq
	Reply  map[string]interface{}
	Done   chan *Call
}
*/

/**/
var (
	msg_proc chan *wrap.Call
)

func init() {
	msg_proc = make(chan *wrap.Call, 10)
	go func() {
		for {
			select {

			case msg, ok := <-msg_proc:
				if ok {
					model.SendGift(msg)
				}
			}
		}
	}()
}

//curl "t1.shangtv.cn:3003/gift/send_gift?uid=4&token=887c4366a2ac754ff520fc58187f3849&giftid=1&num=1&revid=7"
//curl "192.168.1.12:3003/gift/send_gift?uid=4&token=887c4366a2ac754ff520fc58187f3849&giftid=1&num=1&revid=7"
func WrapSendGift(q render.Render, d SendGiftReq) {
	//common.Log.Infof("befor call back,time %v", time.Now().UnixNano()/1e6)
	call := CallBackSendGift(q, d, nil)
	//result :=
	if call == nil {
		return
	}
	<-call.Done
	//s := result.Reply

	//godump.Dump(s)
	//common.Log.Infof("after call back,time %v", time.Now().UnixNano()/1e6)
	return
}

func CallBackSendGift(q render.Render, d SendGiftReq, done chan *wrap.Call) *wrap.Call {
	if done == nil {
		done = make(chan *wrap.Call, 50)
	} else {
		if cap(done) == 0 {
			fmt.Println("chan容量为0,无法返回结果,退出此次计算!")
			return nil
		}
	}
	r := make(map[string]interface{})
	r["uid"] = d.Uid
	r["token="] = d.Token
	r["gift_id"] = d.GiftId
	r["num"] = d.Num
	r["rev_id"] = d.Revid
	call := &wrap.Call{
		Render:  q,
		Request: r,
		Reply:   make(map[string]interface{}),
		Done:    done,
		Uid:     d.Uid,
		Token:   d.Token,
		GiftID:  d.GiftId,
		Num:     d.Num,
		RevId:   d.Revid,
	}
	var user *model.User
	var ok bool
	var ret int

	ret_value := make(map[string]interface{})

	if d.Uid == 0 {
		user, ok = model.GetUserByToken(d.Token)
		if ok == false {
			ret_value["ErrCode"] = common.ERR_TOKEN_VALID
			q.JSON(http.StatusOK, ret_value)
			return nil
		}
	} else {
		user, ret = model.GetUserByUid(d.Uid)
		if ret != common.ERR_SUCCESS {
			ret_value["ErrCode"] = common.ERR_TOKEN_VALID
			q.JSON(http.StatusOK, ret_value)
			return nil
		}
	}
	sess := model.GetUserSessByUid(user.Uid)
	if sess != nil {
		room := model.GetChatRoom(sess.Roomid)
		room.MsgProc <- call
	} else {
		ret_value["ErrCode"] = common.ERR_USER_OFFLINE
		q.JSON(http.StatusOK, ret_value)
		return nil
	}
	return call
}

//送礼物
//curl -d 'token=ed3656ae994739f4ace04eef1a1ca0e9&giftid=27&num=1&revid=3' 'http://120.76.156.177:3003/gift/send_gift'
//func SendGift( r render.Render, d SendGiftReq) {

/*
func SendGift(call *Call) {
	r := call.Render
	d := call.Req
	ret_value := call.Reply
	//render.Render, d SendGiftReq

	//ret_value := make(map[string]interface{})

	var user *model.User
	var ok bool
	var ret int
	if d.Uid == 0 {
		user, ok = model.GetUserByToken(d.Token)
		if ok == false {
			ret_value["ErrCode"] = common.TOKEN_EXPIRE_TIME
			r.JSON(http.StatusOK, ret_value)
			call.done()
			return
		}
	} else {
		user, ret = model.GetUserByUid(d.Uid)
		if ret != common.ERR_SUCCESS {
			ret_value["ErrCode"] = common.ERR_ACCOUNT_EXIST
			r.JSON(http.StatusOK, ret_value)
			call.done()
			return
		}
	}
	gift, ok := model.GetGiftById(d.GiftId)
	if !ok {
		ret_value[ServerTag] = common.ERR_GIFT_EXIST
		r.JSON(http.StatusOK, ret_value)
		call.done()
		return
	}

	if gift.Category == common.GIFT_CATEGORY_EXTRAVAGANT || gift.Category == common.GIFT_CATEGORY_HOT {
		ret = user.SendGiftV2(d.GiftId, d.Num, d.Revid)
	} else {
		ret = user.TipGiftV2(d.GiftId, d.Num, d.Revid)
	}

	//ret := user.SendGift(d.GiftId, d.Num, d.Revid)
	ret_value["ErrCode"] = ret

	user_new, _ := model.GetUserByToken(d.Token)
	ret_value["diamond"] = user_new.Diamond
	ret_value["score"] = user_new.Score

	r.JSON(http.StatusOK, ret_value)
	call.done()
}
*/
//修改角色属性
//curl -d 'token=b4ef22b54387fbcddcbf429718dc473a&nickname=大米111&sex=1' 'http://192.168.1.12:3000/user/modify_char'
func ModifyChar(req *http.Request, r render.Render, d ModifyCharReq) {
	ret_value := make(map[string]interface{})
	var user *model.User
	if d.Uid == 0 {
		user, _ = model.GetUserByToken(d.Token)
	} else {
		user, _ = model.GetUserByUid(d.Uid)
	}

	aff, err := user.SetNick(d.NickName)
	if err != nil || aff == 0 {
		ret_value["ErrCode"] = common.ERR_UNKNOWN
		r.JSON(http.StatusOK, ret_value)
		return
	}
	aff, err = user.SetSex(d.Sex)

	if err != nil || aff == 0 {
		ret_value["ErrCode"] = common.ERR_UNKNOWN
		r.JSON(http.StatusOK, ret_value)
		return
	}

	ret_value["ErrCode"] = common.ERR_SUCCESS

	/*
		if ret == nil && ret2 == nil {
			ret_value["ErrCode"] = common.ERR_SUCCESS
		} else {
			ret_value["ErrCode"] = common.ERR_UNKNOWN
		}
	*/
	r.JSON(http.StatusOK, ret_value)
}

//修改角色密码
//curl -d 'oldpwd=12345167&newpwd=22222222222222&token=cbb55501b739a633b924d952caf3a26c' 'http://192.168.1.12:3000/user/modify_pwd'
func ModifyPwd(req *http.Request, r render.Render, d ModifyPwdReq) {
	ret_value := make(map[string]interface{})
	var user *model.User
	if d.Uid == 0 {
		user, _ = model.GetUserByToken(d.Token)
	} else {
		user, _ = model.GetUserByUid(d.Uid)
	}

	ret_value["ErrCode"] = user.ReSetPwd(d.OldPwd, d.NewPwd)
	r.JSON(http.StatusOK, ret_value)
}

//修改个人信息
//curl -d 'image=aa.pic&nick_name=cart&sex=1&location=beijing&sign=hahahahah&token=f13b06d5c77c648e9c7d1e40bd3c5d5a' 'http://192.168.1.12:3000/user/modify_info'
/*
func ModifyInfo(req *http.Request, r render.Render, d ModifyInfoNickReq) {
	ret_value := make(map[string]interface{})
	token := req.FormValue("token")
	user, _ := model.GetUserByToken(token)

	image := req.FormValue("image")
	nick_name := req.FormValue("nick_name")
	sex := req.FormValue("sex")
	sex_, _ := strconv.Atoi(sex)

	location := req.FormValue("location")
	sign := req.FormValue("sign")

	user.Image = image
	user.NickName = nick_name
	user.Sex = sex_
	user.Location = location
	user.Signature = sign

	if user.Update() {
		ret_value["ErrCode"] = common.ERR_SUCCESS
	} else {
		ret_value["ErrCode"] = common.ERR_UNKNOWN
	}

	r.JSON(http.StatusOK, ret_value)
}
*/
//curl -d 'token=15311111605&nick_name=\ue0000_中午' 'http://192.168.1.12:3000/user/modify_info_nick'
func ModifyInfoNick(req *http.Request, r render.Render, d ModifyInfoNickReq) {
	ret_value := make(map[string]interface{})
	var user *model.User
	if d.Uid == 0 {
		user, _ = model.GetUserByToken(d.Token)
	} else {
		user, _ = model.GetUserByUid(d.Uid)
	}

	if sensitive.CheckExistSensitive(d.NickName) {
		ret_value["ErrCode"] = common.ERR_CONTAIN_SENSETIVE
		r.JSON(http.StatusOK, ret_value)
		return
	}
	has := model.CheckGetUserByNickName(d.NickName)
	if has {
		ret_value["ErrCode"] = common.ERR_REPEAT_NICKNAME
		r.JSON(http.StatusOK, ret_value)
		return
	}
	aff, err := user.SetNick(d.NickName)
	if err != nil || aff == 0 {
		ret_value["ErrCode"] = common.ERR_UNKNOWN
		r.JSON(http.StatusOK, ret_value)
		return
	}
	ret_value["ErrCode"] = common.ERR_SUCCESS

	r.JSON(http.StatusOK, ret_value)
}

//func ModifyInfoNick(req *http.Request, r render.Render, d ModifyInfoNickReq) {
//	ret_value := make(map[string]interface{})
//
//	ret_value["ErrCode"] = common.ERR_UNKNOWN
//	r.JSON(http.StatusOK, ret_value)
//
//}

// 修改昵称需要游戏币
func ModifyNickNameWithScore(req *http.Request, r render.Render, d ModifyInfoNickReq) {
	ret_value := make(map[string]interface{})
	var user *model.User
	if d.Uid == 0 {
		user, _ = model.GetUserByToken(d.Token)
	} else {
		user, _ = model.GetUserByUid(d.Uid)
	}

	if sensitive.CheckExistSensitive(d.NickName) {
		ret_value["ErrCode"] = common.ERR_CONTAIN_SENSETIVE
		r.JSON(http.StatusOK, ret_value)
		return
	}
	has := model.CheckGetUserByNickName(d.NickName)
	if has {
		ret_value["ErrCode"] = common.ERR_REPEAT_NICKNAME
		r.JSON(http.StatusOK, ret_value)
		return
	}

	ret_value["ErrCode"] = user.ResetNickName(d.NickName)
	r.JSON(http.StatusOK, ret_value)
}

func GetResetedInfo(req *http.Request, r render.Render, d GetResetNickReq) {
	ret_value := make(map[string]interface{})
	var user *model.User
	ret_value["nickname_reset_score"] = model.NickNameResetMoney
	ret_value["ErrCode"] = common.ERR_SUCCESS
	if d.Uid == 0 {
		user, _ = model.GetUserByToken(d.Token)
	} else {
		user, _ = model.GetUserByUid(d.Uid)
	}
	ret_value["has_reset_nickname"] = user.HasResetNickName()
	r.JSON(http.StatusOK, ret_value)
}

//curl -d 'token=15441420512&sex=1' 'http://192.168.1.12:3000/user/modify_info_sex'
func ModifyInfoSex(req *http.Request, r render.Render, d ModifyInfoSexReq) {
	ret_value := make(map[string]interface{})
	var user *model.User
	if d.Uid == 0 {
		user, _ = model.GetUserByToken(d.Token)
	} else {
		user, _ = model.GetUserByUid(d.Uid)
	}

	if _, err := user.SetSex(d.Sex); err != nil {
		ret_value["ErrCode"] = common.ERR_UNKNOWN
		return
	}
	ret_value["ErrCode"] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}
func ModifyInfoLocation(req *http.Request, r render.Render, d ModifyInfoLocationReq) {
	ret_value := make(map[string]interface{})
	var user *model.User
	if d.Uid == 0 {
		user, _ = model.GetUserByToken(d.Token)
	} else {
		user, _ = model.GetUserByUid(d.Uid)
	}

	if _, err := user.SetLocation(d.Location); err != nil {
		ret_value["ErrCode"] = common.ERR_UNKNOWN
		return
	}
	ret_value["ErrCode"] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

//curl -d 'token=15311111626&sign=办理本科啊啊啊' 'http://192.168.1.11:3000/user/modify_info_sign'
func ModifyInfoSign(req *http.Request, r render.Render, d ModifyInfoSignReq) {
	ret_value := make(map[string]interface{})

	newkey := sensitive.GetSensitiveWord(d.Sign)
	//godump.Dump(newkey)
	var user *model.User
	if d.Uid == 0 {
		user, _ = model.GetUserByToken(d.Token)
	} else {
		user, _ = model.GetUserByUid(d.Uid)
	}

	if _, err := user.SetSignature(newkey); err != nil {
		ret_value["ErrCode"] = common.ERR_UNKNOWN
		return
	}
	ret_value["ErrCode"] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

/*
func UploadFace(w http.ResponseWriter, req *http.Request, r render.Render, d UploadFaceReq) {
	common.Log.Info("UploadFace() called@@@@@@")
	ret_value := make(map[string]interface{})
	user, _ := model.GetUserByToken(d.Token)
	path := common.Cfg.MustValue("path", "face")
	//filename, err := common.Upload(req, path, user.Uid)

	filename, ret := common.UploadBinding(d.File, path, user.Uid)
	if !ret {
		http.Error(w, "bad upload", 600)
		return
	}

	user.SetFace(path + filename)
	ret_value["ErrCode"] = common.ERR_SUCCESS
	ret_value["file"] = filename
	r.JSON(http.StatusOK, ret_value)
}
*/
/*
func Upload(w http.ResponseWriter, req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	file, handler, err := req.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), 500)
		godump.Dump(err.Error())
		return
	}
	defer file.Close()

	filename := fmt.Sprintf("face_11_%d_%s", time.Now().Unix(), handler.Filename)

	f, err := os.Create(filename)
	defer f.Close()
	io.Copy(f, file)
	ret_value["ErrCode"] = 0
	ret_value["file"] = filename
	r.JSON(http.StatusOK, ret_value)
}
*/
//
//curl -d 'push=1&token=962c482aa556a59d06392c32a93d00d3' 'http://192.168.1.12:3000/user/set_push'
func SetPush(req *http.Request, r render.Render, d SetPushReq) {
	ret_value := make(map[string]interface{})
	user, _ := model.GetUserByToken(d.Token)
	var pushret bool
	if d.Push == 1 {
		pushret = true
	} else {
		pushret = false
	}
	if _, err := user.SetPush(pushret); err != nil {
		ret_value["ErrCode"] = common.ERR_UNKNOWN
		return
	}
	ret_value["ErrCode"] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

func UploadHtml(w http.ResponseWriter, req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	file, handler, err := req.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), 500)
		//godump.Dump(err.Error())
		return
	}
	defer file.Close()

	f, err := os.Create("./wangye/" + handler.Filename)
	defer f.Close()
	io.Copy(f, file)
	ret_value["ErrCode"] = 0
	ret_value["file"] = handler.Filename
	r.JSON(http.StatusOK, ret_value)
}

func RefreshUserInfo(r render.Render, d CommonReq) {
	ret_value := make(map[string]interface{})
	user, _ := model.GetUserByUid(d.Uid)
	info := &model.LoginInfo{}
	user.GetLoginInfo(info)
	ret_value["user"] = info
	ret_value["ErrCode"] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

func UploadSensitive(w http.ResponseWriter, req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	file, handler, err := req.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()

	f, err := os.Create("./config/sensitive.txt")
	defer f.Close()
	io.Copy(f, file)
	ret_value["ErrCode"] = 0
	ret_value["file"] = handler.Filename
	r.JSON(http.StatusOK, ret_value)
}

func UploadRobot(w http.ResponseWriter, req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	file, handler, err := req.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()

	f, err := os.Create("./config/robot.txt")
	defer f.Close()
	io.Copy(f, file)
	ret_value["ErrCode"] = 0
	ret_value["file"] = handler.Filename
	r.JSON(http.StatusOK, ret_value)
}

func TipToAnchor(req *http.Request, r render.Render, d SendGiftReq) {
	ret_value := make(map[string]interface{})

	u, ret := model.GetUserByUid(d.Uid)
	if ret == common.ERR_SUCCESS {
		if u.Token == d.Token {
			ret_value[ServerTag] = u.TipGiftV2(d.GiftId, d.Num, d.Revid)
		}
	} else {
		ret_value[ServerTag] = ret
	}

	ret_value["score"] = u.Score
	ret_value["diamond"] = u.Diamond
	r.JSON(http.StatusOK, ret_value)
}

//curl '192.168.1.12:3003/live_statistics/live_info?uid=146167&token=6385d4bd6a351d8ec8252bb43bb91db7&index=0&month_index=0'
//curl 'shangtv.cn:3003/live_statistics/live_info?uid=146167&token=ca168e387f4bc7a764ac32d5b26ad55e&index=0&month_index=0'
func LiveStatisticsController(req *http.Request, r render.Render, d LiveStatisticsReq) {
	ret_value := make(map[string]interface{})
	if d.MonthIndex < 0 || d.MonthIndex > 2 {
		ret_value[ServerTag] = common.ERR_PARAM
		r.JSON(http.StatusOK, ret_value)
		return
	}
	if d.MonthIndex == 0 {
		d.MonthIndex = 0
	}

	var cur_tm time.Time
	u, ok := model.GetUserExtraByUid(d.Uid)
	if ok {

		for i := 0; i <= 2; i++ {
			if d.MonthIndex == i {
				t := time.Now()
				cur := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
				cur_tm = cur.AddDate(0, -i, 0)
			}

		}
		next_moneth := cur_tm.AddDate(0, 1, 0)
		next_tm := time.Date(next_moneth.Year(), next_moneth.Month(), 1, 0, 0, 0, 0, time.Local)

		ret_value["live_list"], ret_value[ServerTag] = u.LiveDetailStatistics(cur_tm, next_tm, d.Index)
		ret_value["sum_day"], ret_value["sum_hours"], ret_value["sum_rice"], ret_value["sum_moon"], ret_value["sum_min"] = u.LiveStatistics(cur_tm, next_tm)
		ret_value["fans"] = u.FocusStatistics(cur_tm, next_tm)
	}
	r.JSON(http.StatusOK, ret_value)
}

type MonethIndexRes struct {
	Data  string
	Index int
}

//http://192.168.1.12:3003/live_statistics/moneth_index
func MonthStatisticsController(req *http.Request, r render.Render) {
	ret_value := make(map[string]interface{})
	all := make([]MonethIndexRes, 0)
	for i := 0; i <= 2; i++ {
		t := time.Now()
		cur := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)

		cur_tm := cur.AddDate(0, -i, 0)
		//begin := time.Date(cur_tm.Year(), cur_tm.Month(), 1, 0, 0, 0, 0, time.Local)
		var m MonethIndexRes

		m.Data = fmt.Sprintf("%d-%d", cur_tm.Year(), cur_tm.Month())
		m.Index = i
		all = append(all, m)
	}
	ret_value["month_index"] = all
	ret_value[ServerTag] = common.ERR_SUCCESS
	r.JSON(http.StatusOK, ret_value)
}

type UserSearchReq struct {
	Uid      int    `form:"uid" binding:"required"`
	Token    string `form:"token" binding:"required"`
	Index    int    `form:"index" `
	KeyWords string `form:"Key_words" binding:"required"`
}

func UserSearch(req *http.Request, r render.Render, d UserSearchReq) {
	ret_value := make(map[string]interface{})

	var user *model.User
	if d.Uid == 0 {
		user, _ = model.GetUserByToken(d.Token)
	} else {
		user, _ = model.GetUserByUid(d.Uid)
	}

	v := make([]model.OutUserInfo3, 0)

	v, ret_value[ServerTag] = model.GetSearchList(user.Uid, d.Index, d.KeyWords)
	if len(v) == 0 {
		ret_value["focus"] = make([]model.OutUserInfo2, 0)
	} else {
		ret_value["focus"] = v
	}

	//user.GetFocusList(d.Index)
	_, ret_value["count"] = model.GetKeyWordsSearchLength(d.KeyWords)
	r.JSON(http.StatusOK, ret_value)
}
