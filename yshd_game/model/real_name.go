package model

import (
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	"time"
)

type AuthRealInfo struct {
	Id   int64
	Uid  int    `xorm:"int(11) not null"`     //用户ID
	Name string `xorm:"varchar(30) not null"` //名字
	Sex  int    `xorm:"int(11) not null"`     //性别
	//IdentifyType   int    `xorm:"int(11) not null "`      //证件类别1.身份证 2.护照 3.军人证
	Identification string `xorm:"varchar(30) not null "` //身份证
	//Birthday       string `xorm:"varchar(30) not null "`  //生日
	RealImage string `xorm:"varchar(256) not null "` //真人头像
	//	School         string `xorm:"varchar(30) not null "`  //学校
	Statues    int //审核状态
	CommitTime int64
	Tel        string `xorm:"varchar(20)  "`
	IdFront    string `xorm:"varchar(256) "`
	IdBack     string `xorm:"varchar(256)  "`
	AdminId    int    `xorm:"int(11) not null"`
}

func AddRealNameInfo(user *User, name string, sex int, identification, tel string, admin_id int) int {
	/*
		if user.CacheRealImage == "" {
			return common.ERR_UPLOAD_NOTIFY
		}
	*/

	real_image := GetCachePic(user.Uid, GetPicDefine(CACHE_PIC_REAL))
	if real_image == "" {
		return common.ERR_UPLOAD_NOTIFY
	}
	auth := &AuthRealInfo{
		Uid:            user.Uid,
		Name:           name,
		Sex:            sex,
		Identification: identification,
		//	Birthday:       birthday,
		RealImage: real_image,
		Statues:   common.REAL_AUTH_VET,
		//	School:         school,
		CommitTime: time.Now().Unix(),
		//IdentifyType: id_type,
		Tel:     tel,
		AdminId: admin_id,
		IdFront: GetCachePic(user.Uid, GetPicDefine(CACHE_PIC_FRONT)),
		IdBack:  GetCachePic(user.Uid, GetPicDefine(CACHE_PIC_BACK)),
	}

	_, err := orm.Insert(auth)
	if err != nil {
		common.Log.Err("insert real auth error: , %s", err.Error())
		return common.ERR_UNKNOWN
	}
	//user.SetAuthReal()

	//user.SetRealImage("")
	ClearCachePic(user.Uid, GetPicDefine(CACHE_PIC_FRONT))
	ClearCachePic(user.Uid, GetPicDefine(CACHE_PIC_BACK))
	ClearCachePic(user.Uid, GetPicDefine(CACHE_PIC_REAL))
	return common.ERR_SUCCESS
}

func CheckAuthRecord(uid int) bool {
	has, err := orm.Where("uid=? and statues=?", uid, common.REAL_AUTH_VET).Get(&AuthRealInfo{})

	if err != nil {
		common.Log.Err("get real auth error: , %s", err.Error())
		return false
	}
	if has {
		return true
	}
	return false
}

func GetRealInfoByUid(uid int) *AuthRealInfo {
	s := &AuthRealInfo{}
	_, err := orm.Where("uid=?", uid).Get(s)
	if err != nil {
		common.Log.Err("get real auth info error: , %s", err.Error())
		return s
	}
	return s
}
