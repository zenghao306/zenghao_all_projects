package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	. "hdt_app_go/appmeta/log"
	"hdt_app_go/appmeta/model"
	"hdt_app_go/common"
	proto "hdt_app_go/protcol"
	"sync"
	"time"
)

var mutex_mysql sync.Mutex

const APP_NUMBER_PER_PAGE = 10

func NewMysql(server, username, password, dbName, dbPort string) *xorm.Engine {
	//common.Log.Info("db initializing...")
	var err error
	/*
		server := Cfg.MustValue("db", "server")
		username := Cfg.MustValue("db", "username")
		password := Cfg.MustValue("db", "password")
		dbName := Cfg.MustValue("db", "db_name")
		dbPort := Cfg.MustValue("db", "db_port")
	*/
	//fmt.Println(server, username, password, dbName)
	//common.Log.Infof(server, username, password, dbName, dbPort)
	orm, err := xorm.NewEngine("mysql", username+":"+password+"@tcp("+server+":"+dbPort+")/"+dbName+"?charset=utf8mb4&loc=Local")
	//common.PanicIf(err)
	err = orm.Ping()
	if err != nil {
		Log.Fatalf(err.Error())
	}
	common.PanicIf(err)
	orm.SetMaxIdleConns(3000)
	orm.SetMaxOpenConns(5000)
	orm.TZLocation = time.Local
	orm.ShowSQL(false)

	tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, "go_")
	orm.SetTableMapper(tbMapper)

	err = orm.Sync(new(model.ActionSendRecord))
	if err != nil {
		Log.Fatal("sync error  table is go_action_send_record")
	}

	err = orm.Sync(new(model.AccessCode))
	if err != nil {
		Log.Fatal("sync error  table is go_access_code")
	}

	err = orm.Sync(new(model.Member))
	if err != nil {
		Log.Fatal("sync error  table is go_member")
	}

	/*
		_,err=orm.Insert(&model.ActionSendRecord{
			Appid:"1",
			Uid:"2",
			Action:"3",
			CreateTime:time.Now().Unix(),
		})
		if err != nil {
			godump.Dump(err.Error())
			Log.Fatal("sync error  table is RoomList")
		}
	*/
	//ExecExtraSql()

	return orm
}

func (d *Dao) GetUserActionLimitByDb(uid uint32, action string) (count int, err error) {
	return
}

func (d *Dao) CreateAccountByTel(tel, pwd, regIp string, registerFrom int) int32 {
	m := &model.Member{}
	has, err := d.mysqlCli.Where("account=?", tel).Get(m)
	if err != nil {
		Log.Err(err.Error())
		return proto.ERR_UNKNOWN
	} else if has {
		return proto.ERR_EXIST_ACCOUNT
	}

	pwdMd5 := common.TPMd5(pwd) //生成TPMD5暗文

	_, err2 := d.mysqlCli.Insert(&model.Member{
		Account:      tel,
		Password:     pwdMd5,
		RegisterTime: time.Now().Unix(),
		RegIp:        regIp,
		RegisterFrom: int8(registerFrom),
		LoginNum:     0,
		HdtBalance:   0,
		RealnameAuth: 0,
	})

	if err2 != nil {
		Log.Err(err.Error())
		return proto.ERR_UNKNOWN
	}

	return proto.ERR_OK
}

func (d *Dao) GetUserAccountInfo(tel, pwd string) (int32, *proto.UserInfo) {
	pwdMd5 := common.TPMd5(pwd)

	a := &proto.UserInfo{}
	m := &model.Member{}
	has, err := d.mysqlCli.Where("account=?", tel).Get(m)
	if err != nil {
		Log.Err(err.Error())
		return proto.ERR_UNKNOWN, nil
	} else if !has {
		return proto.ERR_NOT_EXIST_ACCOUNT, nil
	} else if pwdMd5 != m.Password {
		return proto.ERR_PASSWORD, nil
	}

	m.LoginTime = time.Now().Unix()

	a.Uid = m.Id
	a.NickName = m.NickName
	a.Avatar = m.Avatar
	a.Sign = m.Sign
	a.Email = m.Email
	a.RegTime = m.RegisterTime
	a.HdtBalance = m.HdtBalance
	a.RealnameAuth = int32(m.RealnameAuth)

	return proto.ERR_OK, a
}

func (d *Dao) ModifyPwdByTel(tel, pwd string) int32 {
	m := &model.Member{}
	m.LoginTime = time.Now().Unix()

	m.Password = common.TPMd5(pwd)
	_, err := d.mysqlCli.Where("account=?", tel).Cols("password").Update(m)
	if err != nil {
		Log.Err(err.Error())
		return proto.ERR_UNKNOWN
	}
	return proto.ERR_OK
}

func (d *Dao) UpdateUserLoginTime(tel string) {
	m := &model.Member{}
	m.LoginTime = time.Now().Unix()
	d.mysqlCli.Where("account=?", tel).Cols("login_time").Update(m)
}

func (d *Dao) GetMinedInfo() (limit float64, total float64, mtime int) {
	var supplyLimit, totalSupply_ float64
	var minedTime int
	supplyLimit = 1000000000000000000
	totalSupply_ = 300000000000000000
	sql := fmt.Sprintf("SELECT supply_limit,total_supply,mined_time FROM go_mined_info")

	rowArray, err := d.mysqlCli.Query(sql)
	if err != nil {
		Log.Err(err.Error())
		return
	}
	if len(rowArray) == 1 {
		supply_limit, ok := rowArray[0]["supply_limit"]
		if ok {
			supplyLimit = common.BytesToFloat64(supply_limit)
		}
		totalSupply, ok := rowArray[0]["total_supply"]
		if ok {
			totalSupply_ = common.BytesToFloat64(totalSupply)
		}
		mined_time, ok := rowArray[0]["mined_time"]
		if ok {
			minedTime = common.BytesToInt(mined_time)
		}
	}

	return supplyLimit, totalSupply_, minedTime
}

func (d *Dao) GetHdtPercent() float64 {
	var supplyLimit, totalSupply, sum_, percent float64
	var minedTime int
	supplyLimit = 1000000000000000000
	totalSupply = 300000000000000000
	supplyLimit, totalSupply, minedTime = d.GetMinedInfo()

	sql := fmt.Sprintf("select ifnull(sum(hdt),0) AS sum from go_action_accounts_record where accounts_time_end >= %d", minedTime)
	rowArray, err := d.mysqlCli.Query(sql)
	if err != nil {
		Log.Err(err.Error())
	}
	if len(rowArray) == 1 {
		sum, ok := rowArray[0]["sum"]
		if ok {
			sum_ = common.BytesToFloat64(sum)
		}
		if sum_ < 0 {
			sum_ = 0
		}
	}

	percent = (totalSupply + sum_) / supplyLimit
	if percent < 0.3 || percent >= 1 {
		percent = 0.3
	}
	return percent
}

//p := db.GetHdtPercent()
//w := (1 / math.Pow(1-p, Pi))

/*
//SELECT SUM(hdt),tel FROM (SELECT a.appid,a.uid,a.hdt,b.tel FROM go_action_accounts_record a LEFT JOIN go_access_code b ON a.appid = b.appid AND a.uid = b.uid WHERE a.hdt != 0 AND b.tel != "") c GROUP BY c.tel
//以下是计算挖矿排名
// end_time挖矿的计算时间
func (d *Dao) GetMySqlRankingList(endTime int64) ([]map[string]string, int) {
	ranking := make([]map[string]string, 0)

	sql := fmt.Sprintf("SELECT SUM(hdt) hdt_sum,tel FROM (SELECT a.appid,a.uid,a.hdt,b.tel FROM go_action_accounts_record a LEFT JOIN go_access_code b ON a.appid = b.appid AND a.uid = b.uid WHERE a.hdt != 0 AND b.tel != '' AND a.accounts_time_end = '%d') c GROUP BY c.tel ORDER BY hdt_sum DESC ",endTime)
	rowArray, err := d.mysqlCli.Query(sql)
	if err != nil {
		Log.Err(err)
		return nil, proto.ERR_UNKNOWN
	}

	//for _, row := range rowArray {
	//	ss := make(map[string]string)
	//	ss2 := make(map[string]string)
	//	for colName, colValue := range row {
	//		value := common.BytesToString(colValue)
	//		ss[colName] = value
	//	}
	//	hdtSum := ss["hdt_sum"]
	//	tel := ss["tel"]
	//	ss2[tel] = hdtSum
	//	ranking = append(ranking, ss2)
	//}
	for _, row := range rowArray {
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			ss[colName] = value
		}
		ranking = append(ranking, ss)
	}
	return ranking, proto.ERR_OK
}
*/
func (d *Dao) GetMySqlRankingList(endTime int64) (map[string]string, int) {
	ranking := make(map[string]string, 0)

	sql := fmt.Sprintf("SELECT SUM(hdt) hdt_sum,tel FROM (SELECT a.appid,a.uid,a.hdt,b.tel FROM go_action_accounts_record a LEFT JOIN go_access_code b ON a.appid = b.appid AND a.uid = b.uid WHERE a.hdt != 0 AND b.tel != '' AND a.accounts_time_end = '%d') c GROUP BY c.tel ORDER BY hdt_sum DESC ", endTime)
	rowArray, err := d.mysqlCli.Query(sql)
	if err != nil {
		Log.Err(err)
		return nil, proto.ERR_UNKNOWN
	}

	for _, row := range rowArray {
		ss := make(map[string]string)
		//ss2 := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			ss[colName] = value
		}
		hdtSum := ss["hdt_sum"]
		tel := ss["tel"]
		//ss2[tel] = hdtSum
		//ranking = append(ranking, ss2)
		ranking[tel] = hdtSum
	}

	return ranking, proto.ERR_OK
}

//根据手机号查询用户挖到的HDT数量总额
func (d *Dao) GetUserHdtMiningTotalByTel(tel string) (int, float64) {
	var sum_ float64

	sql := fmt.Sprintf("SELECT ifnull(FORMAT(sum(hdt),8),0) hdt_sum FROM (SELECT a.appid,a.uid,a.hdt,b.tel FROM go_action_accounts_record a LEFT JOIN go_access_code b ON a.appid = b.appid AND a.uid = b.uid WHERE a.hdt != 0 AND b.tel != '') c WHERE c.tel = '%s'", tel)
	rowArray, err := d.mysqlCli.Query(sql)
	if err != nil {
		Log.Err(err)
		return proto.ERR_UNKNOWN, 0
	}

	sum_ = 0
	if len(rowArray) == 1 {
		sum, ok := rowArray[0]["hdt_sum"]
		if ok {
			sum_ = common.BytesToFloat64(sum)
		}
		if sum_ < 0 {
			sum_ = 0
		}
	}

	return proto.ERR_OK, sum_
}

func (d *Dao) GeAPPIconNameList(index int32) (err_ int32, lis []*proto.AppListRes_AppNameIcon) {
	mutex_mysql.Lock()
	defer mutex_mysql.Unlock()

	if index < 0 { //index不能小于0
		index = 0
	}
	sql := fmt.Sprintf("SELECT appName,appIcoPath FROM php_app limit %d,%d ", index*APP_NUMBER_PER_PAGE, APP_NUMBER_PER_PAGE)
	rowArray, err := d.mysqlCli.Query(sql)
	if err != nil {
		Log.Err(err)
		return proto.ERR_UNKNOWN, nil
	}

	icons := make([]*proto.AppListRes_AppNameIcon, 0) //定义一个slice
	for _, row := range rowArray {
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			ss[colName] = value
		}
		temp := &proto.AppListRes_AppNameIcon{}
		temp.AppIcoPath = "http://admin.hudt.io/public/static/upload/img/" + ss["appIcoPath"]
		temp.AppName = ss["appName"]
		icons = append(icons, temp) //追加到icons
	}

	return proto.ERR_OK, icons
}
