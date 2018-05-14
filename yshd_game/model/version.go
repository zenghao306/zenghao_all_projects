package model

import (
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	"strconv"
	"strings"
	//"github.com/liudng/godump"
	//"github.com/liudng/godump"
)

type VersionManager struct {
	Id         int    `xorm:"int(11) not null pk autoincr"`
	Version    string `xorm:"varchar(50) not null"`        //大版本号
	Title      string `xorm:"varchar(50) not null "`       //标题
	Content    string `xorm:"varchar(255) not null "`      //更新内容
	MinVersion string `xorm:"varchar(50) not null "`       //最低版本
	PubTime    int    `xorm:"int(11) not null "`           //版本记录时间
	Platform   int    `xorm:"int(11) not null default(1)"` //平台类型1 android ，2 ios
	Url        string `xorm:"varchar(255) not null "`      //最新下载地址
	Status     int    `xorm:"int(11) not null default(1)"` //1表示最近可用,2表示不可用
	ChannelId  string `xorm:"varchar(125) not null"`
}

type Version_C struct {
	Version    string `json:"version"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Url        string `json:"url"`
	MinVersion string `json:"min_version"`
}

func CheckVersion(vp string, channelID string) int {

	_, cache_android_version, cache_min_android_version, ret := GetAndroidVersion(channelID)
	if ret != common.ERR_SUCCESS {
		return ret
	}
	sli := strings.Split(vp, ".")

	if len(sli) != 3 {
		return common.ERR_PARAM
	}

	for i := 0; i < 3; i++ {
		value, err := strconv.Atoi(sli[i])
		if err != nil {
			return common.ERR_PARAM
		}

		if value == cache_min_android_version[i] {
			continue
		} else if value < cache_min_android_version[i] {
			return common.ERR_VERSION_MUST_UPDATE
		} else {
			break
		}
	}

	for i := 0; i < 3; i++ {
		value, err := strconv.Atoi(sli[i])
		if err != nil {
			return common.ERR_PARAM
		}

		if value == cache_android_version[i] {
			continue
		} else if value < cache_android_version[i] {
			return common.ERR_VERSION_UPDATE
		} else {
			return common.ERR_SUCCESS
		}
	}
	return common.ERR_SUCCESS
}

func GetAndroidVersion(channelID string) (androidVersion Version_C, cacheAndroidVersion []int, cacheMinAndroidVersion []int, ret int) {
	cacheAndroidVersion = make([]int, 3)

	cacheMinAndroidVersion = make([]int, 3)

	t := VersionManager{}
	has, err2 := orm.Where("platform=? and status=? and channel_id=?", common.OS_TYPE_ANDROID, 1, channelID).Desc("pub_time").Get(&t)
	if err2 != nil {
		common.Log.Errf("orm err is %s", err2.Error())
		ret = common.ERR_UNKNOWN
		return
	}
	if has {
		androidVersion.Version = t.Version
		androidVersion.Title = t.Title
		androidVersion.Content = t.Content
		androidVersion.Url = t.Url
		androidVersion.MinVersion = t.MinVersion

		sli := strings.Split(t.Version, ".")
		if len(sli) != 3 {
			common.Log.Panicln("android version length panic")
			return
		}
		var err3 error
		for i := 0; i < 3; i++ {
			cacheAndroidVersion[i], err3 = strconv.Atoi(sli[i])
			if err3 != nil {
				common.Log.Panicln("android version type panic")
				return
			}
		}

		sli = strings.Split(t.MinVersion, ".")
		if len(sli) != 3 {
			common.Log.Panicln("android version length panic")
			return
		}
		for i := 0; i < 3; i++ {
			cacheMinAndroidVersion[i], err3 = strconv.Atoi(sli[i])
			if err3 != nil {
				common.Log.Panicln("android min version type panic")
				return
			}
		}
	} else {
		ret = common.ERR_CONFGI_ITEM
	}
	return
}
