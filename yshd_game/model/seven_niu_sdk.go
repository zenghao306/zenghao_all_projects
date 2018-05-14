package model

import (
	"fmt"

	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	//"encoding/hex"
	"bytes"
	"github.com/pili-engineering/pili"
	"github.com/qiniu/api.v7/kodo"
	"github.com/qiniu/api.v7/kodocli"
	"github.com/yshd_game/common"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	sevenniu_access_key = "kBKUKjqdHFLqzsKBod8eBFv65d4aftyEE9TgyWkx"
	sevenniu_secret_key = "J73BQ7QKhI3BNG4AnW2xRwV81qPhTpum00mZY7hP"
	callbackurl         = "http://shangtv.cn:8082/sever_niu_notify"
	//callbackurl = "http://shangtv.hicp.net:12456/sever_niu_notify"
	//callbackbody = `{"key":$(key), "hash":$(etag),"filesize":$(fsize)}`
	callbackbody = "key=$(key)&hash=$(etag)&bucket=$(bucket)&uid=$(x:uid)&utoken=$(x:utoken)"
	bucket_face  = "fashionface"
	bucket_cover = "fashioncover"
	bucket_real  = "fashionreal"
	bucket_hub   = "fashion"
	bucket_front = "idfront"
	bucket_back  = "idback"
	bucket_vod   = "fashion-piliLive"

	BFace = ""
	//DomainFace   = "o6xufm374.bkt.clouddn.com"
	DomainFace  = "face.shangtv.cn"
	DomainCover = "cover.shangtv.cn"
	DomainGift  = "gift.shangtv.cn"
	DomainReal  = "real.shangtv.cn"
	DomainVod   = "vod.shangtv.cn"
	DomainFront = "idfront.shangtv.cn"
	DomainBack  = "idback.shangtv.cn"
	//PublishKey  = "c0017047-0cb6-4897-bf09-e959caac2ffb"

	AccessKey = "vKmd-2d_inZREY6ZVy0DRqkaFt-U8youBeIra1k8"

	SecretKey  = "815KzVcsvOUijX1AoPcxAcj-SydOEgz5diKupT4x"
	ExpireTime = 1000
)

const (
	CACHE_PIC_FACE  = 1
	CACHE_PIC_COVER = 2
	CACHE_PIC_REAL  = 3
	CACHE_PIC_BACK  = 4
	CACHE_PIC_FRONT = 5
)

func GetPicDefine(ptype int) string {
	switch ptype {
	case CACHE_PIC_FACE:
		return bucket_face
	case CACHE_PIC_COVER:
		return bucket_cover
	case CACHE_PIC_REAL:
		return bucket_real
	case CACHE_PIC_BACK:
		return bucket_front
	case CACHE_PIC_FRONT:
		return bucket_back
	default:
		return ""
	}
}

func GetBucket() (string, string) {
	return bucket_face, bucket_cover
}

func GetQiNiuFace() string {
	return bucket_face
}
func GetQiNiuCover() string {
	return bucket_cover
}
func GetQiNiuReal() string {
	return bucket_real
}
func GetQiNiuIdBack() string {
	return bucket_back
}

func GetQiNiuIdFront() string {
	return bucket_front
}

func InitQiNiuKey() {
	DomainFace = common.Cfg.MustValue("qiniu", "domain_face")
	DomainCover = common.Cfg.MustValue("qiniu", "domain_cover")
	DomainGift = common.Cfg.MustValue("qiniu", "domain_gift")
	DomainReal = common.Cfg.MustValue("qiniu", "domain_real")
	DomainVod = common.Cfg.MustValue("qiniu", "domain_vod")

	DomainFront = common.Cfg.MustValue("qiniu", "domain_front")
	DomainBack = common.Cfg.MustValue("qiniu", "domain_back")

	AccessKey = common.Cfg.MustValue("qiniu", "access_key")
	SecretKey = common.Cfg.MustValue("qiniu", "secret_key")

	notify := common.Cfg.MustValue("host", "weixin_notify")

	bucket_face = common.Cfg.MustValue("qiniu", "bucket_face")
	bucket_cover = common.Cfg.MustValue("qiniu", "bucket_cover")
	bucket_real = common.Cfg.MustValue("qiniu", "bucket_real")
	bucket_hub = common.Cfg.MustValue("qiniu", "bucket_hub")

	bucket_front = common.Cfg.MustValue("qiniu", "bucket_front")
	bucket_back = common.Cfg.MustValue("qiniu", "bucket_back")

	callbackurl = fmt.Sprintf("http://%s/sever_niu_notify", notify)

	kodo.SetMac(sevenniu_access_key, sevenniu_secret_key)

	InitQiiuPili()
}

func Gen7NiuToken(uptype int, filekey string) string {
	var bucket string
	if uptype == 1 {
		bucket = bucket_face
	} else if uptype == 2 {
		bucket = bucket_cover
	} else if uptype == 3 {
		bucket = bucket_real
	} else if uptype == 4 {
		bucket = bucket_front
	} else if uptype == 5 {
		bucket = bucket_back
	} else {
		bucket = bucket_face
	}

	key := filekey
	c := kodo.New(0, nil)

	//设置上传的策略
	policy := &kodo.PutPolicy{
		Scope:        bucket + ":" + key, // 上传文件的限制条件，这里限制只能上传一个名为 "foo/bar.jpg" 的文件
		Expires:      600,                // 这是限制上传凭证(uptoken)的过期时长，3600 是一小时
		CallbackUrl:  callbackurl,
		CallbackBody: callbackbody,
	}
	//生成一个上传token
	token := c.MakeUptoken(policy)
	return token
}

func SetUserFile(uptype, uid int, filename string) {
	user, _ := GetUserByUid(uid)
	if uptype == 1 {
		user.SetFace(filename)
	} else if uptype == 2 {
		//user.SetCover(filename)
		common.Log.Errf("begin cancel cover is %s", filename)
	}
}

func DownloadUrl(domain, key string) string {
	//调用MakeBaseUrl()方法将domain,key处理成http://domain/key的形式
	baseUrl := kodo.MakeBaseUrl(domain, key)
	//policy := kodo.GetPolicy{}
	//生成一个client对象
	//c := kodo.New(0, nil)
	//调用MakePrivateUrl方法返回url
	//return c.MakePrivateUrl(baseUrl, &policy)
	return baseUrl
}

type PutRet struct {
	Hash     string `json:"hash"`
	Key      string `json:"key"`
	Filesize int    `json:"filesize"`
}

func UploadLocalQiNiu(filepath string, key string, uid int, token string) bool {
	qitoken := Gen7NiuToken(1, key)
	zone := 0
	uploader := kodocli.NewUploader(zone, nil)

	var ret PutRet
	extra := &kodocli.PutExtra{}
	uids := strconv.Itoa(uid)
	extra.Params = make(map[string]string, 2)
	extra.Params["x:uid"] = uids
	extra.Params["x:utoken"] = token
	//调用PutFile方式上传，这里的key需要和上传指定的key一致
	res := uploader.PutFile(nil, &ret, qitoken, key, filepath, extra)
	if res != nil {
		common.Log.Errf("upload file failed %s", res)
		return false
	}
	return true
}

func GetQiNiuFile(key string, bucket string) int {

	c := kodo.New(0, nil)
	p := c.Bucket(bucket)

	_, res := p.Stat(nil, key)

	if res != nil {
		return common.ERR_UNKNOWN
	}
	return common.ERR_SUCCESS
}

func DelQiNiuFile3(bucket string, domain string, file string) int {
	c := kodo.New(0, nil)
	p := c.Bucket(bucket)
	//common.Log.Infof("del qi niu domain is %s file is %s", bucket, file)

	match := fmt.Sprintf("%s", domain)

	name := strings.Index(file, match)
	ss := common.Substr(file, name+len(domain)+1, len(file)-name)

	//调用Delete方法删除文件
	res := p.Delete(nil, ss)
	if res == nil {
		common.Log.Infof("del qi niu domain is %s file is %s", bucket, file)
		return common.ERR_UNKNOWN
	} else {
		common.Log.Errf("del qi niu domain is %s file is %s reson is %s", bucket, file, res)
	}
	return common.ERR_SUCCESS
}

/*
func DelQiNiuFile2(domain string, file string) {
	c := kodo.New(0, nil)
	p := c.Bucket(domain)
	reg := regexp.MustCompile(`[0-9]+`)

	key := reg.FindAllString(file, -1)
	if len(key) == 0 {
		common.Log.Errf("get qi niu key err is %s", file)
		return
	}
	common.Log.Infof("del qi niu domain is %s file is %s", domain, file)
	return
	//调用Delete方法删除文件
	res := p.Delete(nil, key[0])
	if res == nil {
		common.Log.Infof("del qi niu domain is %s file is %s", domain, file)
	} else {
		common.Log.Errf("del qi niu domain is %s file is %s reson is %s", domain, file, res)
	}
}
*/
func GenSercertUrl(roomid string) string {
	expire := time.Now().Unix() + int64(ExpireTime)
	signstr := fmt.Sprintf("/fashion/%s?e=%d", roomid, expire)

	//godump.Dump(signstr)
	key := []byte(SecretKey)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(signstr))
	//sign2 := hex.EncodeToString(mac.Sum(nil))
	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	a := b64.EncodeToString([]byte(mac.Sum(nil)))
	//a := []byte(mac.Sum(nil))
	//godump.Dump(a)
	reg := regexp.MustCompile(`\+`)
	rep := []byte("-")
	str := reg.ReplaceAll([]byte(a), rep)

	reg = regexp.MustCompile(`\/`)
	rep = []byte("_")
	str = reg.ReplaceAll([]byte(str), rep)

	/*
		reg = regexp.MustCompile(`\=`)
		rep = []byte("")
		str = reg.ReplaceAll([]byte(str), rep)
	*/
	//godump.Dump(string(str))
	//s := AccessKey + ":" + string(str)
	//godump.Dump(s)
	return AccessKey + ":" + string(str)
	/*
		expire := time.Now().Unix() + int64(ExpireTime)
		signstr := fmt.Sprintf("/fashion/%s?expire=%d", roomid, expire)
		key := []byte(PublishKey)
		mac := hmac.New(sha1.New, key)
		mac.Write([]byte(signstr))
		sign2 := hex.EncodeToString(mac.Sum(nil))
		b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
		a := b64.EncodeToString([]byte(mac.Sum(nil)))
		return a
	*/
}

type QiNiuM3u8Ret struct {
	Fname string "json:fname"
}
type QiNiuM3u8Err struct {
	Error string "json:error"
}

var hub pili.Hub

func InitQiiuPili() {
	return
	// 直播空间名称
	// 直播空间必须事先存在，可以在 portal.qiniu.com 上创建

	credentials := pili.NewCredentials(AccessKey, SecretKey)
	hub = pili.NewHub(credentials, bucket_hub)
	// 初始化流对象
	//godump.Dump(hub)
}

/*
func SaveM3u8File(roomid string) int {
	hubName := "fashion-piliLive"
	credentials := pili.NewCredentials(AccessKey, SecretKey)
	hub = pili.NewHub(credentials, hubName)

	godump.Dump("beigin svae")
	r, ok := GetRoomById(roomid)
	if ok == false {
		return common.ERR_ROOM_EXIST
	}

	options := pili.OptionalArguments{ // optional
		Title:           r.RoomName, // optional, auto-generated as default
		PublishKey:      "",         // optional, auto-generated as default
		PublishSecurity: "dynamic",  // optional, can be "dynamic" or "static", "dynamic" as default
	}
	godump.Dump(hub)
	stream, err := hub.CreateStream(options)
	if err != nil {
		fmt.Println("Error:", err)
		return 999
	}
	name := fmt.Sprintf("%s.mp4", roomid)

	notify := common.Cfg.MustValue("host", "weixin_notify")
	callbackurl = fmt.Sprintf("http://%s/sever_niu_notify", notify)

	options2 := pili.OptionalArguments{
		NotifyUrl:    callbackurl,
		UserPipeline: "user_pipeline",
	} // optional

	format := "mp4"
	fname, _ := stream.SaveAs(name, format, r.CreateTime.Unix(), time.Now().Unix(), options2)
	godump.Dump(fname)
	r.Playback = fname.TargetUrl
	_, err = orm.Where("room_id=?", roomid).MustCols("playback").Update(r)
	if err != nil {
		common.Log.Errf("db error:", err.Error())
		return common.ERR_UNKNOWN
	}
	return common.ERR_SUCCESS

}
*/

func SaveMutipleM3u8File(roomid string, uid int, recordId int64) int {

	info, ok := GetMutipleRoomRecordByID(recordId)
	if ok == false {
		return common.ERR_ROOM_EXIST
	}

	ret := CheckPlayList(uid)
	if ret != common.ERR_SUCCESS {
		return ret
	}
	s := Base64UrlQiniu(roomid)
	path := fmt.Sprintf("/v2/hubs/%s/streams/%s/saveas", bucket_hub, s)
	data := "POST " + path
	data += "\nHost: pili.qiniuapi.com"
	if "<Content-Type>" != "" {
		data += "\nContent-Type: application/json"
	}
	data += "\n\n"
	// 6. 添加 Body，前提: Content-Length 存在且 Body 不为空，同时 Content-Type 存在且不为空或 "application/octet-stream"

	// 计算 HMAC-SHA1 签名，并对签名结果做 URL 安全的 Base64 编码
	//

	encodedSign := QiNiuGenToken(data)
	QiNiuToken := "Qiniu " + AccessKey + ":" + encodedSign
	//-----------------------------------//

	url2 := fmt.Sprintf("http://pili.qiniuapi.com/v2/hubs/%s/streams/%s/saveas", bucket_hub, s)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", url2, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", QiNiuToken)
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := client.Do(req) //发送
	if err != nil {
		common.Log.Errf("qiniu token check err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		common.Log.Errf("qiniu token check err is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	if resp.StatusCode == 200 {
		var v QiNiuM3u8Ret
		json.Unmarshal(body, &v)
		info.Playback = v.Fname
		return AddMutiplePlayBack(info.OwnerId, roomid, v.Fname, time.Now().Unix(), recordId)
		/*
			_, err = orm.Where("room_id=?", roomid).MustCols("playback").Update(r)
			if err != nil {
				common.Log.Errf("db error:", err.Error())
				return common.ERR_UNKNOWN
			}
		*/
		return common.ERR_SUCCESS
	} else {
		var v QiNiuM3u8Err
		json.Unmarshal(body, &v)
		common.Log.Errf("save m3u8 file err is %s", v.Error)
		return common.ERR_EXIST_STREAM
	}

	return common.ERR_SUCCESS
}

func SaveM3u8File(roomid string) int {
	info, ok := GetRoomById(roomid)
	if ok == false {
		return common.ERR_ROOM_EXIST
	}

	/*
		ret := CheckPlayList(info.OwnerId)
		if ret != common.ERR_SUCCESS {
			return ret
		}
	*/
	s := Base64UrlQiniu(roomid)
	path := fmt.Sprintf("/v2/hubs/%s/streams/%s/saveas", bucket_hub, s)
	data := "POST " + path
	data += "\nHost: pili.qiniuapi.com"
	if "<Content-Type>" != "" {
		data += "\nContent-Type: application/json"
	}
	data += "\n\n"
	// 6. 添加 Body，前提: Content-Length 存在且 Body 不为空，同时 Content-Type 存在且不为空或 "application/octet-stream"

	// 计算 HMAC-SHA1 签名，并对签名结果做 URL 安全的 Base64 编码
	//

	encodedSign := QiNiuGenToken(data)
	QiNiuToken := "Qiniu " + AccessKey + ":" + encodedSign
	//-----------------------------------//

	url2 := fmt.Sprintf("http://pili.qiniuapi.com/v2/hubs/%s/streams/%s/saveas", bucket_hub, s)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", url2, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", QiNiuToken)
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := client.Do(req) //发送
	if err != nil {
		common.Log.Errf("qiniu token check err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		common.Log.Errf("qiniu token check err is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	if resp.StatusCode == 200 {
		var v QiNiuM3u8Ret
		json.Unmarshal(body, &v)
		info.Playback = v.Fname
		return AddNewPlayBack(info.OwnerId, v.Fname, roomid, time.Now().Unix())
		/*
			_, err = orm.Where("room_id=?", roomid).MustCols("playback").Update(r)
			if err != nil {
				common.Log.Errf("db error:", err.Error())
				return common.ERR_UNKNOWN
			}
		*/
		return common.ERR_SUCCESS
	} else {
		var v QiNiuM3u8Err
		json.Unmarshal(body, &v)
		common.Log.Errf("save m3u8 file err is %s", v.Error)
		return common.ERR_EXIST_STREAM
	}

	return common.ERR_SUCCESS
}

func QiNiuGenToken(dst string) string {
	key := []byte(SecretKey)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(dst))
	//sign2 := hex.EncodeToString(mac.Sum(nil))
	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	a := b64.EncodeToString([]byte(mac.Sum(nil)))
	//a := []byte(mac.Sum(nil))
	//godump.Dump(a)
	reg := regexp.MustCompile(`\+`)
	rep := []byte("-")
	str := reg.ReplaceAll([]byte(a), rep)

	reg = regexp.MustCompile(`\/`)
	rep = []byte("_")
	str = reg.ReplaceAll([]byte(str), rep)
	return string(str)
}

func Base64UrlQiniu(dst string) string {
	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	a := b64.EncodeToString([]byte(dst))
	reg := regexp.MustCompile(`\+`)
	rep := []byte("-")
	str := reg.ReplaceAll([]byte(a), rep)

	reg = regexp.MustCompile(`\/`)
	rep = []byte("_")
	str = reg.ReplaceAll([]byte(str), rep)
	return string(str)
}

func GetStream(roomid string) int {
	s := Base64UrlQiniu(roomid)

	path := fmt.Sprintf("/v2/hubs/%s/streams/%s", bucket_hub, s)
	//	godump.Dump(path)
	data := "GET " + path
	data += "\nHost: pili.qiniuapi.com"
	if "<Content-Type>" != "" {
		data += "\nContent-Type: application/x-www-form-urlencoded"
	}
	data += "\n\n"
	encodedSign := QiNiuGenToken(data)
	QiNiuToken := "Qiniu " + AccessKey + ":" + encodedSign
	//godump.Dump(QiNiuToken)
	url2 := fmt.Sprintf("http://pili.qiniuapi.com/v2/hubs/%s/streams/%s", bucket_hub, s)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url2, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", QiNiuToken)
	resp, err := client.Do(req) //发送
	if err != nil {
		common.Log.Errf("qq token check err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		common.Log.Errf("qq token check err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	//godump.Dump(string(body))
	return common.ERR_SUCCESS
}

type QiNiuCreateRoomReq struct {
	OwnerId  string `json:"owner_id"`
	RoomName string `json:"room_name"`
	//UserMax  int    `json:"user_max"`
}

type QiNiuCreateRoomRet struct {
	RoomName string `json:"room_name"`
}

func CreateMultipleRoom(uid int, user_max int) int {
	var yourReq QiNiuCreateRoomReq
	rid := common.GenNewFileName(uid)
	yourReq.OwnerId = strconv.Itoa(uid)
	yourReq.RoomName = rid

	//yourReq.UserMax = user_max
	bytes_req, err := json.Marshal(yourReq)
	if err != nil {
		common.Log.Err("以json形式编码发送错误, 原因:%s", err.Error())
		return common.ERR_INNER_XML_ENCODE
	}

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://rtc.qiniuapi.com/v1/rooms", bytes.NewReader(bytes_req))
	req.Header.Set("Content-Type", "application/json")

	path := fmt.Sprintf("/v1/rooms")
	data := "POST " + path
	data += "\nHost: rtc.qiniuapi.com"
	if req.Header.Get("Content-Type") != "" {
		data += "\nContent-Type: application/json"
	}
	data += "\n\n"
	data += string(bytes_req)

	encodedSign := QiNiuGenToken(data)

	QiNiuToken := "Qiniu " + AccessKey + ":" + encodedSign
	req.Header.Set("Authorization", QiNiuToken)

	response, err := client.Do(req) //发送
	if err != nil {
		common.Log.Errf("qiniu rtc  err is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	body, _ := ioutil.ReadAll(response.Body)
	switch response.StatusCode {
	case 200:
		var m QiNiuCreateRoomRet
		err = json.Unmarshal(body, &m)
		if err != nil {
			common.Log.Errf("error: %s", err.Error())
			return common.ERR_UNKNOWN
		}
		AddMultipleRoomList(uid, rid)
		break
	case 400:
		return common.ERR_MULTIPLE_ROOM_PARM
	case 611:
		return common.ERR_MULTIPLE_ROOM_EXIST
	default:
		return common.ERR_UNKNOWN

	}
	c := NewChatRoomInfo()
	o := c.GetChatInfo()
	o.Rid = rid
	//o.Image = m.Cover
	o.Uid = uid
	c.Statue = common.ROOM_MULTIPLE
	AddChatRoom(c)
	return common.ERR_SUCCESS
}

type QiNiuGetMutipleRet struct {
	RoomName   string `json:"room_name"`
	OwnerId    string `json:"owner_id"`
	RoomStatus int    `json:"room_status"`
	UserMax    int    `json:"user_max"`
}

func ReqGetMutipleRoom(rname string) (int, QiNiuGetMutipleRet) {
	url := fmt.Sprintf("http://rtc.qiniuapi.com/v1/rooms/%s", rname)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	path := fmt.Sprintf("/v1/rooms/%s", rname)
	data := "GET " + path
	data += "\nHost: rtc.qiniuapi.com"
	data += "\n\n"

	encodedSign := QiNiuGenToken(data)

	QiNiuToken := "Qiniu " + AccessKey + ":" + encodedSign
	req.Header.Set("Authorization", QiNiuToken)

	var m QiNiuGetMutipleRet
	response, err := client.Do(req) //发送
	if err != nil {
		common.Log.Errf("qiniu rtc  err is %s", err.Error())
		return common.ERR_UNKNOWN, m
	}

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode == 200 {
		err = json.Unmarshal(body, &m)
		//godump.Dump(m)
		if err != nil {
			common.Log.Errf("error: %s", err.Error())
			return common.ERR_UNKNOWN, m
		}

	} else if response.StatusCode == 602 {
		return common.ERR_MULTIPLE_ROOM_FIN, m
	}
	return common.ERR_SUCCESS, m
}

type QiNiuDElMutipleRet struct {
	ERROR string `json:"error"`
}

func DelMutipleRoom(rname string) int {
	url := fmt.Sprintf("http://rtc.qiniuapi.com/v1/rooms/%s", rname)
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", url, nil)

	path := fmt.Sprintf("/v1/rooms/%s", rname)
	data := "DELETE " + path
	data += "\nHost: rtc.qiniuapi.com"
	data += "\n\n"

	encodedSign := QiNiuGenToken(data)

	QiNiuToken := "Qiniu " + AccessKey + ":" + encodedSign
	req.Header.Set("Authorization", QiNiuToken)
	response, err := client.Do(req) //发送
	if err != nil {
		common.Log.Errf("qiniu del multiple  err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	//var m QiNiuDElMutipleRet

	if response.StatusCode == 200 {
		/*
			body, _ := ioutil.ReadAll(response.Body)
			err = json.Unmarshal(body, &m)
			if err != nil {
				common.Log.Errf("error: %s", err.Error())
				return common.ERR_UNKNOWN
			}
		*/
		return common.ERR_SUCCESS
	} else if response.StatusCode == 612 {
		return common.ERR_MULTIPLE_ROOM_NOT_FOUND
	} else if response.StatusCode == 613 {
		return common.ERR_MULTIPLE_ROOM_IN_USE
	}
	return common.ERR_UNKNOWN
}

type RoomAccess struct {
	RoomName string `json:"room_name"`
	UserId   string `json:"user_id"`
	Perm     string `json:"perm"`
	Expire   int64  `json:"expire_at"`
}

func GenMutipleToken(uid int, rid string) (int, string, string) {
	expire := time.Now().Unix() + int64(10*3600)
	var ret int
	var j *MultipleRoomList
	permission := ""
	if rid == "" {
		user, ret := GetUserByUid(uid)
		if ret != common.ERR_SUCCESS {
			return ret, "", ""
		} else if user.CanLinkMic == 0 { //主播判断是否有连麦的权限
			return common.ERR_MULTIPLE_NO_POWER_LINK_MIC, "", ""
		}

		ret, j = GetFreeMultipleRoom(uid)
		if ret == common.ERR_SUCCESS || ret == common.ERR_MULTIPLE_HAS_MIC {
			if permission == "" {
				permission = "admin"
			}
			rid = j.RoomId
		} else {
			return ret, "", ""
		}
	} else {
		permission = "user"
	}
	var m RoomAccess
	m.RoomName = rid
	m.UserId = strconv.Itoa(uid)
	m.Perm = permission
	m.Expire = expire
	byte_arr, err := json.Marshal(m)
	if err != nil {
		common.Log.Err("以json形式编码错误, 原因:%s", err.Error())
		return common.ERR_UNKNOWN, "", ""
	}
	room_access := string(byte_arr)
	encode_access := Base64UrlQiniu(room_access)
	encode_sign := QiNiuGenToken(encode_access)
	QiNiuToken := AccessKey + ":" + encode_sign + ":" + encode_access
	return ret, QiNiuToken, rid
}

func InviteUser(uid, oid int, rid string) int {
	var res ResponseInvite
	res.MType = common.MESSAGE_TYPE_MULTIPLE_INVITE
	res.User = uid
	res.Other = oid
	res.Rid = rid
	if SendMsgToUser(oid, res) {
		return common.ERR_SUCCESS
	}
	return common.ERR_USER_OFFLINE
}
