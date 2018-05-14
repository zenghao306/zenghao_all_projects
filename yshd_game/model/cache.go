package model

import (
	//	"github.com/yshd/common"
	//"github.com/liudng/godump"

	"github.com/yshd_game/common"
	"strconv"
	"time"
)

//推荐列表缓存
type CacheRecommandList struct {
	timeout         int64
	recommand       []map[string]string //少
	test_recommand  []map[string]string //全部
	recommand2      []map[string]string //不带后续游戏【3】的
	test_recommand2 []map[string]string //不带后续游戏【3】的
	multiple        []map[string]string
	test_multiple   []map[string]string
	mutipleTest     []map[string]string
}

var RecommandCache CacheRecommandList

func InitCahce() {
	RecommandCache.recommand = make([]map[string]string, 0)
	RecommandCache.test_recommand = make([]map[string]string, 0)
	RecommandCache.recommand2 = make([]map[string]string, 0)
	RecommandCache.test_recommand2 = make([]map[string]string, 0)
	RecommandCache.multiple = make([]map[string]string, 0)
	RecommandCache.test_multiple = make([]map[string]string, 0)
	RecommandCache.mutipleTest = make([]map[string]string, 0)

	RoomListCache.room_list = make([]map[string]string, 0)
	RoomListCache.room_test_list = make([]map[string]string, 0)

	CreatePlayBackRoom()
}

func GetMutipleTest() []map[string]string {
	return RecommandCache.mutipleTest
}
func GetRecommandCache(rtype int, uid int) ([]map[string]string, []map[string]string) {
	ReFreshRecommand(uid)
	if rtype == common.ROOM_TYPE_NOMARL {
		return RecommandCache.recommand, RecommandCache.multiple
	}
	return RecommandCache.test_recommand, RecommandCache.test_multiple
}

func ReFreshRecommand(uid int) {
	nowtime := time.Now().Unix()
	if RecommandCache.timeout <= nowtime {
		rids := GetDaGuanRecommand(uid)
		RecommandCache.recommand, RecommandCache.test_recommand, _ = GetRecommandListV3(rids) //GetRecommandWithPlayUrl(0)

		//RecommandCache.res,RecommandCache.test=GetRecommandListV2()

		RecommandCache.multiple, RecommandCache.test_multiple, _ = GetMultipleRoomList(0)
		RecommandCache.timeout = nowtime + int64(5)
	}
}

func GetRecommandCache2(rtype int, uid int) ([]map[string]string, []map[string]string) {
	ReFreshRecommand2(uid)
	if rtype == common.ROOM_TYPE_NOMARL {
		return RecommandCache.recommand2, RecommandCache.multiple
	}
	return RecommandCache.test_recommand2, RecommandCache.test_multiple
}

func ReFreshRecommand2(uid int) {
	nowtime := time.Now().Unix()
	if RecommandCache.timeout <= nowtime {
		rids := GetDaGuanRecommand(uid)
		RecommandCache.recommand2, RecommandCache.test_recommand2, _ = GetRecommandListV4(rids) //GetRecommandWithPlayUrl(0)

		//RecommandCache.res,RecommandCache.test=GetRecommandListV2()

		RecommandCache.multiple, RecommandCache.test_multiple, _ = GetMultipleRoomList(0)
		RecommandCache.timeout = nowtime + int64(5)
	}
}

//主页房间缓存
type CacheMainList struct {
	timeout        int64
	room_list      []map[string]string
	room_test_list []map[string]string
}

var RoomListCache CacheMainList

func GetRoomCache() []map[string]string {
	ReFreshRoom()
	return RoomListCache.room_list
}

func ReFreshRoom() {
	nowtime := time.Now().Unix()
	if RoomListCache.timeout <= nowtime {
		RoomListCache.room_list, RoomListCache.room_test_list, _, _, _ = GetRecommandList(0)
		RoomListCache.timeout = nowtime + int64(5)
	}
}

//回拨缓存
//
type CachePlayBackList struct {
	timeout   int64
	recommand []map[string]string
}

var PlayBackCache CachePlayBackList

func CreatePlayBackRoom() {
	r, _ := GetRecommandWithPlayUrl(0)
	for _, m := range r {
		c := NewChatRoomInfo()
		o := c.GetChatInfo()
		o.Rid = m["room_id"]
		o.Uid, _ = strconv.Atoi(m["owner_id"])
		o.Image = m["cover"]
		c.Statue = common.ROOM_PLAYBACK
		AddChatRoom(c)
	}

	k, _ := GetRecommandMutipleWithPlayUrl(0)
	for _, m := range k {
		c := NewChatRoomInfo()
		o := c.GetChatInfo()
		o.Rid = m["room_id"]
		o.Uid, _ = strconv.Atoi(m["owner_id"])
		o.Image = m["cover"]
		c.Statue = common.ROOM_PLAYBACK
		c.IsMultiple = true
		AddChatRoom(c)
	}
}

//
func InitMutipleRoom() {
	r := make([]MultipleRoomList, 0)
	err := orm.Find(&r)
	if err != nil {
		common.Log.Errf("db err %s", err.Error())
		return
	}
	for _, m := range r {
		c := NewChatRoomInfo()
		o := c.GetChatInfo()
		o.Rid = m.RoomName
		o.Image = m.Cover
		o.Uid = 0
		c.Statue = common.ROOM_MULTIPLE
		AddChatRoom(c)
	}
}

//{"cover":"http://face.17playlive.com/100029149621090444032","image":"http://face.17playlive.com/100029149621090444032",
// "live_url":"http://vod.shangtv.cn/recordings/z1.fashion.100029149647737417649/0_1496477735.m3u8","location":"重庆市","nick_name":"苏黎世",
// "owner_id":"100029","room_id":"100029149647737417649","room_name":"房间1496477375","sex":"0","viewer":"69"}]
type RecommandRoomRes struct {
	Cover     string `json:"cover"`
	Image     string `json:"image"`
	LiveUrl   string `json:"live_url"`
	Location  string `json:"location"`
	NickName  string `json:"nick_name"`
	OwnerId   string `json:"owner_id"`
	RoomId    string `json:"room_id"`
	RoomNname string `json:"room_name"`
	Sex       string `json:"sex"`
	Viewer    string `json:"viewer"`
	FlvUrl    string `json:"flv_url"`
	Statue    string `json:"statue"`
	GameType  string `json:"game_type"`
}

//func GetRecommandListV2()  (res []RecommandRoomRes,test []RecommandRoomRes)  {

func GetRecommandListV2() (res []map[string]string, test []map[string]string) {
	rids := make([]string, 0)
	uids := make([]int, 0)

	//	res=make([]RecommandRoomRes,0)
	//	test=make([]RecommandRoomRes,0)

	for k, v := range chat_room_manager {
		if v.Statue == common.ROOM_ONLIVE {
			rids = append(rids, k)
			uids = append(uids, v.room.Uid)
		}
	}

	if len(rids) == 0 {
		return
	}

	//rooms := make([]RoomList, 0)
	/*
		err:=orm.In("room_id",rids).Find(&rooms)
		if err != nil {
			common.Log.Errf("db err %s", err.Error())
			return
		}
	*/

	/*
		var str string
		for i:=0;i<len(rids);i++  {
			str+=rids[i]
			str+=","
		}

		l:=str[0:len(str)-1]
		rowArray,err:=orm.Query("select a.room_id,a.room_name,a.owner_id,a.location,a.cover,a.live_url ,a.flv_url ,b.nick_name,b.sex,a.statue,b.image,a.roomtype,a.game_type from go_room_list a left join go_user where a.room_id in (?) order by a.weight ,b.coupons,a.create_time desc",l)
		if err != nil {
			common.Log.Errf("db err %s", err.Error())
			return
		}

		for k,v:=range rowArray  {
			ss := make(map[string]string)
			for colName, colValue := range v {
				if colName=="room_id" {
					value := common.BytesToString(colValue)
					chat := GetChatRoom(value)
					if chat==nil {

					}
				}
			}
		}

	*/
	/*
		users := make(map[int]User, 0)
		err=orm.In("uid",uids).Find(&users)
		if err != nil {
			common.Log.Errf("db err %s", err.Error())
			return
		}
	*/

	/*
		for _,v:=range rooms  {
			u,ok:=users[v.OwnerId]
			if !ok {
				continue
			}

			chat := GetChatRoom(v.RoomId)
			var view string
			if  chat==nil{
				view="0"
			}else{
				view=strconv.Itoa(chat.GetCount() + chat.GetVRobotCount())
			}

			m:=RecommandRoomRes{
				Cover:v.Cover,
				Image:u.Image,
				LiveUrl:v.LiveUrl,
				Location:v.Location,
				NickName:u.NickName,
				OwnerId:  strconv.Itoa( v.OwnerId),
				RoomId:v.RoomId,
				RoomNname:v.RoomName,
				Sex:strconv.Itoa( u.Sex),
				Viewer:view,
				Statue: strconv.Itoa(common.ROOM_ONLIVE),
				GameType:strconv.Itoa(v.GameType),
			}

			if v.Roomtype==0 {
				res=append(res,m)
				test=append(test,m)
			}else if v.Roomtype==1 {
				test=append(test,m)
			}

		}
	*/

	return
}
