package model

import (
	"fmt"
	//"github.com/liudng/godump"
	"github.com/yshd_game/common"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type LoginInfo struct {
	Uid         int //`json:"uid"`
	Account     string
	Tel         string
	NickName    string
	Sex         int
	UserLevel   int
	AnchorLevel int
	UserExp     int
	AnchorExp   int
	Image       string
	Diamond     int
	Focus       int
	Fans        int
	Coupons     int
	Push        bool
	Signature   string
	Token       string
	Location    string
	CanLinkMic  int
	Score       int64
	Moon        int64
	AuthReal    bool
	IsSuperUser bool
}

//验证登陆账号
func Auth(account, pwd string, platform int) (int, *User) {
	/*
		sql := fmt.Sprintf("select  * from user where account = %s and platform =%d  ", account, platform)
		result, err := orm.Query(sql)
		if err != nil {
			common.Log.Errf("err is %s", err.Error())
		}
		if len(result) > 0 {
			real_pwd := common.BytesToString(result[0]["pwd"])
			if real_pwd == pwd {
				return common.ERR_SUCCESS
			}
			return common.ERR_PWD
		}
		return common.ERR_EXIST
	*/
	//user, has := GetUserByAccount(account)
	user, has := GetUserByAccountAndPlatfrom(account, platform)
	if has {
		if user.Pwd == pwd {
			ret := user.CheckAccountForbid()
			if ret != common.ERR_SUCCESS {
				return ret, user
			}

			InsertLog(common.ACTION_TYPE_LOG_LOGIN, user.Uid, "")

			return common.ERR_SUCCESS, user
		}
		return common.ERR_PWD, user
	} else {
		return common.ERR_ACCOUNT_EXIST, user
	}

}

//通过手机号创建新账号
func CreateAccountByTel(pwd, tel, channel_id, device string, registerFrom int) int {
	user := &User{}
	has, err := orm.Where("account=?", tel).Get(user)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	if has == false {
		//new_user := &User{Account: tel, Tel: tel, Pwd: pwd, WeixinAccount: "", SinaAccount: "", QqAccount: "", Platform: common.PLATFORM_SELF, Token: tel, UserLevel: 1, AnchorLevel: 1, UserExp: 0, AnchorExp: 0, RegisterTime: time.Now().Unix()}

		sql := fmt.Sprintf("insert into `go_user` (`account`, `tel`, `pwd`,`platform`,`user_level`,`anchor_level`,`register_time`,`register_channel`,`device`,`register_from`) values ('%s','%s','%s',%d,1,1,'%d','%s','%s',%d)", tel, tel, pwd, common.PLATFORM_SELF, time.Now().Unix(), channel_id, device, registerFrom)

		res, err := orm.Exec(sql)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
		aff_row, err := res.RowsAffected()
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
		if aff_row == 0 {
			return common.ERR_DB_ADD
		}
		/*
			_, err := orm.InsertOne(new_user)
			if err != nil {
				go common.Log.Errf("orm err is %s", err.Error())
				godump.Dump("rrrrrr")
				return common.ERR_UNKNOWN
			}
		*/
		//sql := fmt.Sprintf("insert into user_extra (uid) values (%d)", new_user.Uid)

		user, has := GetUserByTel(tel)
		if !has {
			return common.ERR_UNKNOWN
		}
		u := &UserExtra{}
		has, err = orm.Where("uid=?", user.Uid).Get(u)
		if !has {
			res, err = orm.Exec("insert into go_user_extra (uid) values (?)", user.Uid)
			if err != nil {
				common.Log.Errf("orm err is %s", err.Error())
				return common.ERR_UNKNOWN
			}
			if aff_row, err := res.RowsAffected(); aff_row == 0 {
				if err != nil {
					common.Log.Errf("orm err is %s", err.Error())
					return common.ERR_UNKNOWN
				}
				return common.ERR_DB_ADD
			}
		} else {
			//godump.Dump("user err")
		}

		sql = fmt.Sprintf("insert into go_coupons_month_record (uid,year_months) values (%d,%s)", user.Uid, common.GetCurentYearMonthString())
		res, err = orm.Exec(sql)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
		aff_row, err = res.RowsAffected()
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
		if aff_row == 0 {
			return common.ERR_DB_ADD
		}
		InsertLog(common.ACTION_TYPE_LOG_REGISTER, user.Uid, "")
		return common.ERR_SUCCESS
	}
	return common.ERR_EXIST
}

func CreateAccountByWeiXin(account, wxOpenID, pwd, city, headimgurl, nickname string, platform, sex, registerFrom int, times int, channel_id, device string) int {
	_, has := GetUserByNickName(nickname)
	if has {
		if times >= 100 {
			return common.ERR_REPEAT_THIRD_NICKNAME
		}
		num := common.RadnomRange(1, 1000000)
		nickname = fmt.Sprintf("%s_%d", nickname, num)
		times++
		return CreateAccountByWeiXin(account, wxOpenID, pwd, city, headimgurl, nickname, platform, sex, registerFrom, times, channel_id, device)
	}

	//sql := fmt.Sprintf("insert into `go_user` (`account`, `open_id`,`nick_name`, `pwd`,`location`,`image`,`sex`,`user_level`,`anchor_level`,`register_time`,`register_channel`,`device`,`register_from`,`platform`) values ('%s','%s','%s','%s',\"%s\",'%s',%d,1,1,%d,'%s','%s',%d,%d)", account, wxOpenID, nickname, pwd, city, headimgurl, sex, time.Now().Unix(), channel_id, device, registerFrom, platform)
	//common.Log.Info(sql)
	s := &User{
		Account:         account,
		OpenId:          wxOpenID,
		NickName:        nickname,
		Pwd:             pwd,
		Location:        city,
		Image:           headimgurl,
		Sex:             sex,
		RegisterTime:    time.Now().Unix(),
		RegisterChannel: channel_id,
		Device:          device,
		RegisterFrom:    registerFrom,
		Platform:        platform,
		UserLevel:       1,
		AnchorLevel:     1,
	}
	aff_row, err := orm.Cols("account", "platform", "nick_name", "pwd", "location", "image", "sex", "register_time", "register_channel", "user_level", "anchor_level", "device", "register_from", "opend_id").InsertOne(s)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		return common.ERR_DB_ADD
	}
	/*
		res, err := orm.Exec(sql)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}

		aff_row, err := res.RowsAffected()
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
		if aff_row == 0 {
			return common.ERR_DB_ADD
		}
	*/
	//user, has := GetUserByAccount(account)
	user, has := GetUserByAccountAndPlatfrom(account, platform)
	if !has {
		return common.ERR_UNKNOWN
	}

	u := &UserExtra{}
	has, err = orm.Where("uid=?", user.Uid).Get(u)
	if !has {
		res, err := orm.Exec("insert into go_user_extra (uid) values (?)", user.Uid)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
		if aff_row, err := res.RowsAffected(); aff_row == 0 {
			if err != nil {
				common.Log.Errf("orm err is %s", err.Error())
				return common.ERR_UNKNOWN
			}
			return common.ERR_DB_ADD
		}
	} else {
		//godump.Dump("user err")
	}

	/*
		res, err := orm.Exec("insert into go_user_extra (uid) values (?)", user.Uid)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
		aff_row, err = res.RowsAffected()
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
		if aff_row == 0 {
			return common.ERR_DB_ADD
		}
	*/
	sql := fmt.Sprintf("insert into go_coupons_month_record (uid,year_months) values (%d,%s)", user.Uid, common.GetCurentYearMonthString())
	res, err := orm.Exec(sql)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	aff_row, err = res.RowsAffected()
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		return common.ERR_DB_ADD
	}
	imgResponse, err := http.Get(headimgurl)
	if err != nil {
		common.Log.Errf("download pic err is %s", err.Error())
		token := common.GenUserToken(user.Uid)
		user.SetToken(token)
		return common.ERR_SUCCESS
	} else {
		defer imgResponse.Body.Close()
		imgByte, _ := ioutil.ReadAll(imgResponse.Body)
		path := common.Cfg.MustValue("path", "face")
		filename := common.GenNewFileName(user.Uid)
		filepath := common.StaticPath + path + filename
		fh, errimage := os.Create(filepath)
		defer fh.Close()
		if errimage != nil {
			common.Log.Err("weixin get image err is ")
			return common.ERR_GET_IMAGE_ERR
		}
		fh.Write(imgByte)

		UploadLocalQiNiu(filepath, filename, user.Uid, user.Token)
	}
	token := common.GenUserToken(user.Uid)
	user.SetToken(token)
	return common.ERR_SUCCESS
}

//创建微信新账号
func CreateAccountByWeiXin2(account, pwd, city, headimgurl, nickname string, platform, sex int, times int) int {
	_, has := GetUserByNickName(nickname)
	if has {
		if times >= 100 {
			return common.ERR_REPEAT_THIRD_NICKNAME
		}
		num := common.RadnomRange(1, 10000)
		nickname = fmt.Sprintf("%s_%d", nickname, num)
		times++
		return CreateAccountByWeiXin2(account, pwd, city, headimgurl, nickname, platform, sex, times)
	}

	sql := fmt.Sprintf("insert into `go_user` (`account`, `nick_name`, `pwd`,`location`,`image`,`sex`,`user_level`,`anchor_exp`,`register_time`,`platform`) values ('%s','%s','%s','%s','%s',%d,1,1,%d)", account, nickname, pwd, city, headimgurl, sex, time.Now().Unix(), platform)
	res, err := orm.Exec(sql)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	aff_row, err := res.RowsAffected()
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		return common.ERR_DB_ADD
	}
	//user, has := GetUserByAccountWithWeixinPlatform(account)
	user, has := GetUserByAccountAndPlatfrom(account, common.PLATFORM_WEIXIN)
	if !has {
		return common.ERR_UNKNOWN
	}
	_, err = orm.Exec("insert into go_user_extra (uid) values (?)", user.Uid)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	sql = fmt.Sprintf("insert into go_coupons_month_record (uid,year_months) values (%d,%s)", user.Uid, common.GetCurentYearMonthString())
	res, err = orm.Exec(sql)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	aff_row, err = res.RowsAffected()
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		return common.ERR_DB_ADD
	}
	imgResponse, err := http.Get(headimgurl)
	if err != nil {
		common.Log.Errf("download pic err is %s", err.Error())
		token := common.GenUserToken(user.Uid)
		user.SetToken(token)
		return common.ERR_SUCCESS
	}
	defer imgResponse.Body.Close()
	imgByte, _ := ioutil.ReadAll(imgResponse.Body)
	path := common.Cfg.MustValue("path", "face")
	filename := common.GenNewFileName(user.Uid)
	filepath := common.StaticPath + path + filename
	fh, errimage := os.Create(filepath)
	defer fh.Close()
	if errimage != nil {
		common.Log.Err("weixin get image err is ")
		return common.ERR_GET_IMAGE_ERR
	}
	fh.Write(imgByte)
	token := common.GenUserToken(user.Uid)
	user.SetToken(token)

	UploadLocalQiNiu(filepath, filename, user.Uid, user.Token)
	return common.ERR_SUCCESS
}

//创建qq新账号
func CreateAccountByQQ(unionid, pwd, city, headimgurl, nickname string, platform, sex, registerFrom int, times int, channel_id, device string, openid string) int {

	_, has := GetUserByNickName(nickname)
	if has {
		if times >= 100 {
			return common.ERR_REPEAT_THIRD_NICKNAME
		}
		num := common.RadnomRange(1, 10000000)
		nickname = fmt.Sprintf("%s_%d", nickname, num)
		times++
		return CreateAccountByQQ(unionid, pwd, city, headimgurl, nickname, platform, sex, registerFrom, times, channel_id, device, openid)
	}
	/*
		sql := fmt.Sprintf("insert into `go_user` (`account`, `platform`, `nick_name`, `pwd`,`location`,`image`,`sex`,`user_level`,`anchor_exp`,`register_time`,`register_channel`,`device`,`register_from`,`open_id`,`moon`) values ('%s',%d,'%s','%s','%s','%s',%d,1,1,%d,'%s','%s',%d,        '%s',      0)", unionid, platform, nickname, pwd, city, headimgurl, sex, time.Now().Unix(), channel_id, device, registerFrom, openid)
	*/
	//sql := fmt.Sprintf("insert into `go_user` (`account`, `platform`, `nick_name`, `pwd`,`location`,`image`,`sex`,`user_level`,`anchor_exp`,`register_time`,`register_channel`) values ('%s',%d,'%s','%s','%s','%s',%d,1,1,%d,'%s')", account, platform, nickname, pwd, city, headimgurl, sex, time.Now().Unix(), channel_id)
	//godump.Dump(sql)

	s := &User{
		Account:         unionid,
		NickName:        nickname,
		Pwd:             pwd,
		Location:        city,
		Image:           headimgurl,
		Sex:             sex,
		RegisterTime:    time.Now().Unix(),
		RegisterChannel: channel_id,
		Device:          device,
		RegisterFrom:    registerFrom,
		Platform:        platform,
		UserLevel:       1,
		AnchorLevel:     1,
	}
	aff_row, err := orm.Cols("account", "platform", "nick_name", "pwd", "location", "image", "sex", "register_time", "register_channel", "user_level", "anchor_level", "device", "register_from").InsertOne(s)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		return common.ERR_DB_ADD
	}

	/*
		res, err := orm.Exec(sql)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}

		aff_row, err := res.RowsAffected()
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
		if aff_row == 0 {
			return common.ERR_DB_ADD
		}
	*/
	//user, has := GetUserByAccountWithQQPlatform(account)
	user, has := GetUserByAccountAndPlatfrom(unionid, common.PLATFORM_QQ)
	if !has {
		return common.ERR_UNKNOWN
	}

	u := &UserExtra{}
	has, err = orm.Where("uid=?", user.Uid).Get(u)
	if !has {
		res, err := orm.Exec("insert into go_user_extra (uid) values (?)", user.Uid)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
		if aff_row, err := res.RowsAffected(); aff_row == 0 {
			if err != nil {
				common.Log.Errf("orm err is %s", err.Error())
				return common.ERR_UNKNOWN
			}
			return common.ERR_DB_ADD
		}
	} else {
		//godump.Dump("user err)")
	}

	/*
		res, err := orm.Exec("insert into go_user_extra (uid) values (?)", user.Uid)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}

		aff_row, err = res.RowsAffected()
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
		if aff_row == 0 {
			return common.ERR_DB_ADD
		}
	*/
	sql := fmt.Sprintf("insert into go_coupons_month_record (uid,year_months) values (%d,%s)", user.Uid, common.GetCurentYearMonthString())
	res, err := orm.Exec(sql)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	aff_row, err = res.RowsAffected()
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		return common.ERR_DB_ADD
	}
	imgResponse, err := http.Get(headimgurl)
	if err != nil {
		common.Log.Errf("download pic err is %s", err.Error())
		token := common.GenUserToken(user.Uid)
		user.SetToken(token)
		return common.ERR_SUCCESS
	}
	defer imgResponse.Body.Close()
	imgByte, _ := ioutil.ReadAll(imgResponse.Body)
	path := common.Cfg.MustValue("path", "face")
	filename := common.GenNewFileName(user.Uid)
	filepath := common.StaticPath + path + filename
	fh, errimage := os.Create(filepath)
	defer fh.Close()
	if errimage != nil {
		common.Log.Err("qq get image err is ")
		return common.ERR_GET_IMAGE_ERR
	}
	fh.Write(imgByte)
	token := common.GenUserToken(user.Uid)
	user.SetToken(token)

	UploadLocalQiNiu(filepath, filename, user.Uid, user.Token)

	return common.ERR_SUCCESS
}

func CreateAccountBySina(account, pwd, city, headimgurl, nickname string, platform, sex, registerFrom int, times int, channel_id, device string) int {
	_, has := GetUserByNickName(nickname)
	if has {
		if times >= 100 {
			return common.ERR_REPEAT_THIRD_NICKNAME
		}
		num := common.RadnomRange(1, 10000)
		nickname = fmt.Sprintf("%s_%d", nickname, num)
		times++
		return CreateAccountBySina(account, pwd, city, headimgurl, nickname, platform, sex, registerFrom, times, channel_id, device)
	}
	/*
		_, has := GetUserByNickName(nickname)
		if has {
			nickname = fmt.Sprintf("%s_%d", nickname, time.Now().Unix())
			return CreateAccountBySina(account, pwd, city, headimgurl, nickname, platform, sex)
		}
	*/
	//sql := fmt.Sprintf("insert into `user` (`sina_account`, `nick_name`, `pwd`,`location`,`image`,`sex`,`user_level`,`anchor_exp`,`register_time`,`register_channel`,`device`) values ('%s','%s','%s','%s','%s',%d,1,1,%d,'%s','%s')", account, nickname, pwd, city, headimgurl, sex, time.Now().Unix(), channel_id, device)
	sql := fmt.Sprintf("insert into `go_user` (`account`,`platform`, `nick_name`, `pwd`,`location`,`image`,`sex`,`user_level`,`anchor_exp`,`register_time`,`register_channel`,`device`,`register_from`) values ('%s',%d,'%s','%s','%s','%s',%d,1,1,%d,'%s','%s',%d)", account, common.PLATFORM_SINA, nickname, pwd, city, headimgurl, sex, time.Now().Unix(), channel_id, device, registerFrom)

	res, err := orm.Exec(sql)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	aff_row, err := res.RowsAffected()
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		return common.ERR_DB_ADD
	}
	//user, has := GetUserByAccountWithSinaPlatform(account)
	user, has := GetUserByAccountAndPlatfrom(account, common.PLATFORM_SINA)
	if !has {
		return common.ERR_UNKNOWN
	}
	res, err = orm.Exec("insert into go_user_extra (uid) values (?)", user.Uid)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	aff_row, err = res.RowsAffected()
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		return common.ERR_DB_ADD
	}
	sql = fmt.Sprintf("insert into go_coupons_month_record (uid,year_months) values (%d,%s)", user.Uid, common.GetCurentYearMonthString())
	res, err = orm.Exec(sql)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	aff_row, err = res.RowsAffected()
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		return common.ERR_DB_ADD
	}
	imgResponse, err := http.Get(headimgurl)
	if err != nil {
		common.Log.Errf("download pic err is %s", err.Error())
		token := common.GenUserToken(user.Uid)
		user.SetToken(token)
		return common.ERR_SUCCESS
	}
	defer imgResponse.Body.Close()
	imgByte, _ := ioutil.ReadAll(imgResponse.Body)
	path := common.Cfg.MustValue("path", "face")

	filename := common.GenNewFileName(user.Uid)
	filepath := common.StaticPath + path + filename
	fh, errimage := os.Create(filepath)

	defer fh.Close()
	if errimage != nil {
		common.Log.Err("sina get image err is ")
		return common.ERR_GET_IMAGE_ERR
	}
	fh.Write(imgByte)
	token := common.GenUserToken(user.Uid)
	user.SetToken(token)

	UploadLocalQiNiu(filepath, filename, user.Uid, user.Token)
	return common.ERR_SUCCESS
}

func CreateAccountByThirdParty(account, pwd, city, headimgurl, nickname, channel_id, device string, sex, platform, registerFrom, times int) int {
	_, has := GetUserByNickName(nickname)
	if has {
		if times >= 100 {
			return common.ERR_REPEAT_THIRD_NICKNAME
		}
		num := common.RadnomRange(1, 10000)
		nickname = fmt.Sprintf("%s_%d", nickname, num)
		times++
		return CreateAccountByWeiXin(account, "", pwd, city, headimgurl, nickname, platform, sex, registerFrom, times, channel_id, device)
	}

	sql := fmt.Sprintf("insert into `go_user` (`account`, `nick_name`, `pwd`,`location`,`image`,`sex`,`user_level`,`anchor_level`,`register_time`,`register_channel`,`device`,`platform`) values ('%s','%s','%s','%s','%s',%d,1,1,%d,'%s','%s',%d)", account, nickname, pwd, city, headimgurl, sex, time.Now().Unix(), channel_id, device, platform)
	//godump.Dump(sql)
	_, err := orm.Exec(sql)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	user, has := GetUserByAccountAndPlatfrom(account, platform)
	if !has {
		return common.ERR_UNKNOWN
	}
	res, err := orm.Exec("insert into go_user_extra (uid) values (?)", user.Uid)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	aff_row, err := res.RowsAffected()
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		return common.ERR_DB_ADD
	}
	sql = fmt.Sprintf("insert into go_coupons_month_record (uid,year_months) values (%d,%s)", user.Uid, common.GetCurentYearMonthString())
	_, err = orm.Exec(sql)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}

	imgResponse, err := http.Get(headimgurl)
	if err != nil {
		common.Log.Err(err.Error())
		token := common.GenUserToken(user.Uid)
		user.SetToken(token)
		return common.ERR_SUCCESS
	}
	defer imgResponse.Body.Close()
	imgByte, _ := ioutil.ReadAll(imgResponse.Body)
	path := common.Cfg.MustValue("path", "face")
	filename := common.GenNewFileName(user.Uid)
	filepath := common.StaticPath + path + filename
	fh, errimage := os.Create(filepath)
	defer fh.Close()
	if errimage != nil {
		common.Log.Err("ThirdParty get image err is ")
		return common.ERR_GET_IMAGE_ERR
	}
	fh.Write(imgByte)
	token := common.GenUserToken(user.Uid)
	user.SetToken(token)

	UploadLocalQiNiu(filepath, filename, user.Uid, user.Token)
	return common.ERR_SUCCESS
}

func CreateAccountByFacebookOrTwitter(account, pwd, city, headimgurl, nickname, channel_id, device string, sex, platform, times int) int {
	_, has := GetUserByNickName(nickname)
	if has {
		if times >= 100 {
			return common.ERR_REPEAT_THIRD_NICKNAME
		}
		num := common.RadnomRange(1, 10000)
		nickname = fmt.Sprintf("%s_%d", nickname, num)
		times++
		return CreateAccountByFacebookOrTwitter(account, pwd, city, headimgurl, nickname, channel_id, device, sex, platform, times)
	}

	sql := fmt.Sprintf("insert into `go_user` (`account`, `nick_name`, `pwd`,`location`,`image`,`sex`,`user_level`,`anchor_level`,`register_time`,`register_channel`,`device`,`platform`) values ('%s','%s','%s','%s','%s',%d,1,1,%d,'%s','%s',%d)", account, nickname, pwd, city, headimgurl, sex, time.Now().Unix(), channel_id, device, platform)
	//godump.Dump(sql)
	res, err := orm.Exec(sql)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	aff_row, err := res.RowsAffected()
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		return common.ERR_DB_ADD
	}
	user, has := GetUserByAccountAndPlatfrom(account, platform)
	if !has {
		return common.ERR_UNKNOWN
	}
	_, err = orm.Exec("insert into go_user_extra (uid) values (?)", user.Uid)
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	aff_row, err = res.RowsAffected()
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	if aff_row == 0 {
		return common.ERR_DB_ADD
	}
	return common.ERR_SUCCESS
}

func ModifyPwdByTel(tel, newpwd string) int {
	user, has := GetUserByTel(tel)

	if has {
		user.Pwd = newpwd
		_, err := orm.Where("uid=?", user.Uid).Update(user)
		if err != nil {
			common.Log.Errf("orm err is %s", err.Error())
			return common.ERR_UNKNOWN
		}
		return common.ERR_SUCCESS
	}
	return common.ERR_BIND_TEL
}

type AppleChannelService struct {
	Channel      string `xorm:"varchar(30)"`
	Device       string `xorm:"varchar(30)"`
	RegisterTime int64  `xorm:"int(11)`
}

func RegisterNewChannel(channel, device string) int {
	_, err := orm.Exec("insert into apple_channel_service values (?,?,?)", channel, device, time.Now().Unix())
	if err != nil {
		common.Log.Errf("orm err is %s", err.Error())
		return common.ERR_UNKNOWN
	}
	return common.ERR_SUCCESS
}
