package db

import (
	"encoding/binary"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	. "hdt_app_go/appmeta/log"
	"hdt_app_go/appmeta/model"
	"hdt_app_go/common"
	proto "hdt_app_go/protcol"
	"math"
	"strconv"
	"sync"
	"time"
)

var mutex_mysql sync.Mutex

const (
	APP_NUMBER_PER_PAGE = 10
	SERVER_BASE_PATH    = "http://admin.hudt.io/public/static/upload/img/"
)

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

const Pi = 3.14
const BASE_DEGREEOFDIFFICULTY = 20

//计算难度系数
func (d *Dao) GetHdtDegreeOfDifficulty() float64 {
	p := d.GetHdtPercent()
	degreeOfDifficulty := 1 / math.Pow(1-p, Pi)

	f := fmt.Sprintf("%.5f", BASE_DEGREEOFDIFFICULTY*degreeOfDifficulty)
	degreeOfDifficulty2 := common.ParseFloat(f)

	return degreeOfDifficulty2
}

func (d *Dao) GetMySqlRankingList(endTime int64) (map[string]string, int) {
	ranking := make(map[string]string, 0)

	sql := fmt.Sprintf("SELECT FORMAT(SUM(hdt),5) hdt_sum,tel FROM (SELECT a.appid,a.uid,a.hdt,b.tel FROM go_action_accounts_record a LEFT JOIN go_access_code b ON a.appid = b.appid AND a.uid = b.uid WHERE a.hdt != 0 AND b.tel != '' AND a.accounts_time_end = '%d') c GROUP BY c.tel ORDER BY hdt_sum DESC ", endTime)
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

	sql := fmt.Sprintf("SELECT ifnull(FORMAT(sum(hdt),5),0) hdt_sum FROM (SELECT a.appid,a.uid,a.hdt,b.tel FROM go_action_accounts_record a LEFT JOIN go_access_code b ON a.appid = b.appid AND a.uid = b.uid WHERE a.hdt != 0 AND b.tel != '') c WHERE c.tel = '%s'", tel)
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

func (d *Dao) GetAPPIconNameList(index int32) (err_ int32, lis []*proto.AppListRes_AppNameIcon) {
	mutex_mysql.Lock()
	defer mutex_mysql.Unlock()

	if index < 0 { //index不能小于0
		index = 0
	}
	sql := fmt.Sprintf("SELECT appId,appName,appIcoPath FROM php_app WHERE `status` = '1' limit %d,%d ", index*APP_NUMBER_PER_PAGE, APP_NUMBER_PER_PAGE)
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
		temp.AppIcoPath = SERVER_BASE_PATH + ss["appIcoPath"]
		temp.AppName = ss["appName"]
		appId, _ := strconv.ParseInt(ss["appId"], 10, 64)
		temp.AppId = appId
		icons = append(icons, temp) //追加到icons
	}

	return proto.ERR_OK, icons
}

//获取任务-发布记录
func (d *Dao) GetMinePoolTaskReleaseList() (err_ int32, lis []*proto.MinePoolTaskListRes_MinePoolTask) {
	var hdtBalance_, deliveryNumber_ float64
	mutex_mysql.Lock()
	defer mutex_mysql.Unlock()

	sql := fmt.Sprintf("select b.appId,b.appIcoPath,b.hdtBalance,b.appName,b.`status`,a.deliveryNumber,a.createTime from php_app_delivery a LEFT JOIN php_app b ON a.appId = b.appId where a.createTime = (select max(createTime) from php_app_delivery where appId = a.appId AND `status` = 'success') order by a.appId")
	rowArray, err := d.mysqlCli.Query(sql)
	if err != nil {
		Log.Err(err)
		return proto.ERR_UNKNOWN, nil
	}

	lists := make([]*proto.MinePoolTaskListRes_MinePoolTask, 0) //定义一个slice
	for _, row := range rowArray {
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			ss[colName] = value
		}

		if "1" != ss["status"] {
			continue
		}

		temp := &proto.MinePoolTaskListRes_MinePoolTask{}

		hdtBalance, ok := row["hdtBalance"]
		if ok {
			hdtBalance_ = common.BytesToFloat64(hdtBalance)
		}
		if hdtBalance_ < 0 {
			hdtBalance_ = 0
		}
		f := fmt.Sprintf("%.5f", hdtBalance_)
		hdtBalance2 := common.ParseFloat(f)
		temp.HdtTaskBalance = hdtBalance2

		temp.AppName = ss["appName"]
		appId, _ := strconv.ParseInt(ss["appId"], 10, 64)
		temp.AppId = appId
		temp.AppIcoPath = SERVER_BASE_PATH + ss["appIcoPath"]

		//转换发布额度
		deliveryNumber, ok := rowArray[0]["deliveryNumber"]
		if ok {
			deliveryNumber_ = common.BytesToFloat64(deliveryNumber)
		}
		if deliveryNumber_ < 0 {
			deliveryNumber_ = 0
		}
		f = fmt.Sprintf("%.5f", deliveryNumber_)
		deliveryNumber2 := common.ParseFloat(f)
		temp.Hdt = deliveryNumber2
		temp.Style = 1

		temp.Time, _ = strconv.ParseInt(ss["createTime"], 10, 64)

		lists = append(lists, temp) //追加到icons
	}

	return proto.ERR_OK, lists
}

//获取任务-挖矿记录
func (d *Dao) GetMinePoolTaskDigList() (err_ int32, lis []*proto.MinePoolTaskListRes_MinePoolTask) {
	var hdtBalance_, hdt_ float64
	mutex_mysql.Lock()
	defer mutex_mysql.Unlock()

	stdtime := time.Now()
	yestodayTime := time.Unix(stdtime.Unix()-24*3600, 0)

	sql := fmt.Sprintf("SELECT b.appId,b.appIcoPath,b.hdtBalance,b.appName,b.`status`,SUM(a.hdt) hdt,a.accounts_time_end time FROM go_action_accounts_record a LEFT JOIN php_app b ON a.appid = b.appCode WHERE a.accounts_time_end > %d AND a.accounts_time_end <= %d GROUP BY a.appid ", yestodayTime.Unix(), stdtime.Unix())
	rowArray, err := d.mysqlCli.Query(sql)
	if err != nil {
		Log.Err(err)
		return proto.ERR_UNKNOWN, nil
	}

	lists := make([]*proto.MinePoolTaskListRes_MinePoolTask, 0) //定义一个slice
	for _, row := range rowArray {
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			ss[colName] = value
		}

		if "1" != ss["status"] {
			continue
		}

		temp := &proto.MinePoolTaskListRes_MinePoolTask{}

		hdtBalance, ok := row["hdtBalance"]
		if ok {
			hdtBalance_ = common.BytesToFloat64(hdtBalance)
		}
		if hdtBalance_ < 0 {
			hdtBalance_ = 0
		}
		f := fmt.Sprintf("%.5f", hdtBalance_)
		hdtBalance2 := common.ParseFloat(f)
		temp.HdtTaskBalance = hdtBalance2

		temp.AppName = ss["appName"]
		appId, _ := strconv.ParseInt(ss["appId"], 10, 64)
		temp.AppId = appId
		temp.AppIcoPath = SERVER_BASE_PATH + ss["appIcoPath"]

		//转换发布额度
		hdt, ok := rowArray[0]["hdt"]
		if ok {
			hdt_ = common.BytesToFloat64(hdt)
		}
		if hdt_ < 0 {
			hdt_ = 0
		}
		f = fmt.Sprintf("%.5f", hdt_)
		hdt2 := common.ParseFloat(f)
		temp.Hdt = hdt2
		temp.Style = 0

		temp.Time, _ = strconv.ParseInt(ss["time"], 10, 64)

		lists = append(lists, temp) //追加到icons
	}

	return proto.ERR_OK, lists
}

//获取用户在对应appId上挖到的HDT数量
func (d *Dao) GetUserAppHdt(appId int64, tel string) (int32, float64) {
	var sum_ float64

	sql := fmt.Sprintf("SELECT ifnull(FORMAT(sum(hdt),5),0) hdt_sum FROM (SELECT a.appid,a.uid,a.hdt,b.tel FROM go_action_accounts_record a LEFT JOIN go_access_code b ON a.appid = b.appid AND a.uid = b.uid WHERE a.appid = %d AND a.hdt != 0 AND b.tel != '') c WHERE c.tel = '%s'", appId, tel)
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

//获取APP开发者投放的HDT
func (d *Dao) GetAppHdtTotal(appId int64) (int32, float64) {
	var sum_ float64

	sql := fmt.Sprintf("SELECT ifnull(FORMAT(sum(deliveryNumber),5),0) sum_deliveryNumber FROM (SELECT a.appId,appCode,deliveryNumber,b.`status` FROM php_app a LEFT JOIN php_app_delivery b ON a.appId = b.appId WHERE a.appId = %d AND b.`status` = 'success') c", appId)
	rowArray, err := d.mysqlCli.Query(sql)
	if err != nil {
		Log.Err(err)
		return proto.ERR_UNKNOWN, 0
	}

	sum_ = 0
	if len(rowArray) == 1 {
		sum, ok := rowArray[0]["sum_deliveryNumber"]
		if ok {
			sum_ = common.BytesToFloat64(sum)
		}
		if sum_ < 0 {
			sum_ = 0
		}
	}

	return proto.ERR_OK, sum_
}

func Float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

//获取APP开发者投放的HDT
func (d *Dao) GetAppHdtBalanceTotal() (int32, float64) {
	var sum_ float64

	//sql := fmt.Sprintf("SELECT ifnull(FORMAT(sum(hdtBalance),5),0) sum_hdtBalance FROM php_app")
	sql := fmt.Sprintf("SELECT ifnull(sum(hdtBalance),0) sum_hdtBalance FROM php_app")
	rowArray, err := d.mysqlCli.Query(sql)
	if err != nil {
		Log.Err(err)
		return proto.ERR_UNKNOWN, 0
	}

	sum_ = 0
	if len(rowArray) == 1 {
		sum, ok := rowArray[0]["sum_hdtBalance"]
		if ok {
			sum_ = common.BytesToFloat64(sum)
			f := fmt.Sprintf("%.5f", sum_)
			sum2 := common.ParseFloat(f)

			return proto.ERR_OK, sum2
		}
		if sum_ < 0 {
			sum_ = 0
		}
	}

	return proto.ERR_OK, sum_
}

func (d *Dao) GetAppContent(appId int64) (int32, string, string, string) {
	sql := fmt.Sprintf("SELECT appContent,iosAddress,androidAddress FROM php_app a LEFT JOIN php_app_content b ON a.appId = b.appId WHERE a.appId = %d", appId)
	rowArray, err := d.mysqlCli.Query(sql)
	if err != nil {
		Log.Err(err)
		return proto.ERR_UNKNOWN, "", "", ""
	}

	if len(rowArray) == 1 {
		return proto.ERR_OK, common.BytesToString(rowArray[0]["appContent"]),
			common.BytesToString(rowArray[0]["iosAddress"]),
			common.BytesToString(rowArray[0]["androidAddress"])
	}

	return proto.ERR_OK, "", "", ""
}

func (d *Dao) GeAPPImageList(appId int64) (err_ int32, lis []string) {
	mutex_mysql.Lock()
	defer mutex_mysql.Unlock()

	sql := fmt.Sprintf("SELECT b.appPath FROM php_app a LEFT JOIN php_app_img b ON a.appId = b.appId WHERE a.appId = %d ", appId)
	rowArray, err := d.mysqlCli.Query(sql)
	if err != nil {
		Log.Err(err)
		return proto.ERR_UNKNOWN, nil
	}

	images := make([]string, 0) //定义一个slice
	for _, row := range rowArray {
		ss := make(map[string]string)
		for colName, colValue := range row {
			value := common.BytesToString(colValue)
			ss[colName] = value
		}
		temp := SERVER_BASE_PATH + ss["appPath"]
		images = append(images, temp) //追加到icons
	}

	return proto.ERR_OK, images
}
