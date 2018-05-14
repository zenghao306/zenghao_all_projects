package common

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"

	"fmt"
	//"github.com/liudng/godump"
	"io"
	"math/rand"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
	"unsafe"
	//"syscall"
)

var Log *Logger
var Logpath string

//var ErrFile *File

func SetLog() {
	var out io.Writer
	if Cfg.MustBool("", "log_file", false) {
		nowtime := time.Now()
		filename := nowtime.Format("20060102")
		finalname := Logpath + filename

		f, _ := os.OpenFile(finalname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
		out = io.MultiWriter(f)
		//syscall.Dup2(int(f.Fd()), 2)

	} else {
		out = os.Stdout
	}

	Log = New(out, "[go-blog]", Lshortfile|Ldate|Lmicroseconds)
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func ParseInt(value string) int {
	if value == "" {
		return 0
	}
	val, _ := strconv.Atoi(value)
	return val
}

func IntString(value int) string {
	return strconv.Itoa(value)
}

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Atoa(str string) string {
	var result string
	for i := 0; i < len(str); i++ {
		c := rune(str[i])
		if 'A' <= c && c <= 'Z' && i > 0 {
			result = result + "_" + strings.ToLower(string(str[i]))
		} else {
			result = result + string(str[i])
		}
	}
	return result
}

func GetRemoteIp(r *http.Request) (ip string) {
	ip = r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = r.RemoteAddr
	}
	ip = strings.Split(ip, ":")[0]
	if len(ip) < 7 || ip == "127.0.0.1" {
		ip = "localhost"
	}
	return
}

/* Test Helpers */
func Expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func GenToken() string {
	nano := time.Now().UnixNano()
	rand.Seed(nano)
	rndNum := rand.Int63()
	uuid := Md5(Md5(strconv.FormatInt(nano, 10)) + Md5(strconv.FormatInt(rndNum, 10)))
	return uuid
}

func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

func BytesToInt(b []byte) int {
	res := BytesToString(b)
	value, _ := strconv.Atoi(res)
	return value
}

func BytesToFloat32(b []byte) float32 {
	res := BytesToString(b)
	i, err := strconv.ParseFloat(res, 32)
	if err != nil {
		Log.Errf("parse err is %s", err.Error())
		return 0.0
	}
	return float32(i)
}

func StringToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{sh.Data, sh.Len, 0}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func RadnomRange(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min) + min
}

func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}

func RandomRangeArr(start int, end int, count int) []int {
	//范围检查
	if end < start || (end-start) < count {
		return nil
	}

	//存放结果的slice
	nums := make([]int, 0)
	//随机数生成器，加入时间戳保证每次生成的随机数不一样
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		//生成随机数
		num := r.Intn((end - start)) + start

		//查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}

		if !exist {
			nums = append(nums, num)
		}
	}

	return nums
}

func RandnomRange64(min, max int64) int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int63n(max-min) + min
}

func RandomRangeArr64(min, max int64, count int) []int64 {
	s := make([]int64, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < count; i++ {

		rvaule := r.Int63n(max-min) + min
		for j := 0; j < len(s); j++ {
			if s[j] == rvaule {
				continue
			}
		}
		s = append(s, rvaule)
	}
	return s
}

func GetTelToekn() int {
	return RadnomRange(VERIFT_CODE_MIN, VERIFY_CODE_MAX)
}

func GenUserToken(uid int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Intn(1000000) + time.Now().Nanosecond()
	srand := fmt.Sprintf("%d", n)
	suid := fmt.Sprintf("%d", uid)
	rr := suid + srand[2:]
	return Md5(rr)
}

func GenWeiXinRandom() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Intn(1000000) + time.Now().Nanosecond()
	srand := fmt.Sprintf("%d", n)
	return Md5(srand)
}

func Upload(req *http.Request, path string, uid int) (string, error) {
	file, handler, err2 := req.FormFile("file")
	if err2 != nil {
		return "", err2
	}
	defer file.Close()
	reg := regexp.MustCompile(`\..+`)

	suffix := reg.FindAllString(handler.Filename, -1)
	if len(suffix) == 0 {
		return "", nil
	}
	if suffix[0] != ".png" && suffix[0] != "jpeg" {
		return "", nil
	}

	filename := fmt.Sprintf("%d_%d_%s", uid, time.Now().Unix(), handler.Filename)

	f, err := os.Create(StaticPath + path + filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	io.Copy(f, file)
	return filename, err
}

func UploadBinding(file *multipart.FileHeader, path string, uid int) (string, bool) {
	if file == nil {
		return "", false
	}

	filename := fmt.Sprintf("%d_face", uid)
	f, err := os.Create(StaticPath + path + filename)
	if err != nil {
		return "", false
	}
	defer f.Close()

	pic, err := file.Open()
	defer pic.Close()
	io.Copy(f, pic)
	return filename, true
}

func GetRoomKey(rid int) string {
	return fmt.Sprintf("%s%d", RoomKey, rid)
}

func GetListKey(rid int) string {
	return fmt.Sprintf("%s%d", RoomKey, rid)
}

func IsDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
	panic("not reached")
}

func GetOnlyId(uid int) string {
	randnum := RadnomRange(1, 99999)
	//key := fmt.Sprintf("%d%d", uid, randnum)
	key := fmt.Sprintf("%d%d%d", uid, time.Now().Unix(), randnum)
	return key
}
func GenNewFileName(uid int) string {
	randnum := RadnomRange(1, 99999)
	//key := fmt.Sprintf("%d%d", uid, randnum)
	key := fmt.Sprintf("%d%d%d", uid, time.Now().Unix(), randnum)
	return key
}

func GetFormartTime() string {
	timestamp := time.Now().Unix()
	return time.Unix(timestamp, 0).Format("20060102")
	//return time.Unix(timestamp, 0).Format("2006-01-02")
}

func GetFormartTime2() string {
	timestamp := time.Now().Unix()
	return time.Unix(timestamp, 0).Format("2006-01-02")
	//return time.Unix(timestamp, 0).Format("2006-01-02")
}

const (
	base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

var coder = base64.NewEncoding(base64Table)

func Base64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}

func Base64Decode(src []byte) ([]byte, error) {
	return coder.DecodeString(string(src))
}

func UnicodeEncode(rs string) string {
	json := ""
	html := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			json += string(r)
			html += string(r)
		} else {
			json += "\\u" + strconv.FormatInt(int64(rint), 16) // json
			html += "&#" + strconv.Itoa(int(r)) + ";"          // 网页
		}
	}
	return json
}

func PrintPanicStack(extras ...interface{}) {
	if x := recover(); x != nil {
		i := 0
		funcName, file, line, ok := runtime.Caller(i)

		for ok {
			Log.Errf("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			i++
			funcName, file, line, ok = runtime.Caller(i)
		}
	}
}

func GetConverting(t int) int {
	switch t {
	case TRADE_TYPE_APPLE:
		return APPLE_TRADE_CONVERT
	default:
		return 1
	}
}

func GetCurentYearMonthString() string {
	timeNow := time.Now()
	month := timeNow.Month()
	yearMonth := strconv.Itoa(timeNow.Year())
	if month < 10 {
		yearMonth += "0"
	}

	yearMonth += strconv.Itoa(int(month))
	return yearMonth
}

func GetLastYearMonth() (int, int) {
	timeNow := time.Now()
	year := timeNow.Year()
	month := timeNow.Month()
	if month == 1 {
		return year - 1, 12
	} else {
		return year, int(month - 1)
	}
}

func GetLastYearMonthString() string {
	lastYear, lastMonth := GetLastYearMonth()
	yearMonth := strconv.Itoa(lastYear)
	if lastMonth < 10 {
		yearMonth += "0"
	}
	yearMonth += strconv.Itoa(lastMonth)
	return yearMonth
}

/*
* 函数名
*   GetCurentWeekFirstUnixTime
*
* 说明
*       获取本周一的凌晨零点时间戳
*
* 参数说明
*
* RETURNS
*   UNIX时间戳
 */
func GetCurentWeekFirstUnixTime() int64 {
	stdtime := time.Now()

	t := stdtime
	t2 := stdtime
	days := stdtime.Weekday()
	switch days {
	case 0:
		t = stdtime.AddDate(0, 0, -6)
	case 1:
		t = stdtime.AddDate(0, 0, 0)
	case 2:
		t = stdtime.AddDate(0, 0, -1)
	case 3:
		t = stdtime.AddDate(0, 0, -2)
	case 4:
		t = stdtime.AddDate(0, 0, -3)
	case 5:
		t = stdtime.AddDate(0, 0, -4)
	case 6:
		t = stdtime.AddDate(0, 0, 5)
	}
	t2 = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)

	return t2.Unix()
}

// 按照特定格式（"2006-01-02"）字符串返还本周一的日期
func GetCurentWeekFirstDate() string {
	stdtime := time.Now()

	t := stdtime
	t2 := stdtime
	days := stdtime.Weekday()
	switch days {
	case 0:
		t = stdtime.AddDate(0, 0, -6)
	case 1:
		t = stdtime.AddDate(0, 0, 0)
	case 2:
		t = stdtime.AddDate(0, 0, -1)
	case 3:
		t = stdtime.AddDate(0, 0, -2)
	case 4:
		t = stdtime.AddDate(0, 0, -3)
	case 5:
		t = stdtime.AddDate(0, 0, -4)
	case 6:
		t = stdtime.AddDate(0, 0, -5)
	}
	t2 = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)

	return t2.Format("2006-01-02")
}

type MapSorter []Item

type Item struct {
	Uid      int `json:"uid"`
	WinCoins int `json:"win_coins"`
}

func NewMapSorter(m map[int]int) MapSorter {
	ms := make(MapSorter, 0, len(m))

	for k, v := range m {
		ms = append(ms, Item{k, v})
	}

	return ms
}

func (ms MapSorter) Len() int {
	return len(ms)
}

func (ms MapSorter) Less(i, j int) bool {
	return ms[i].WinCoins > ms[j].WinCoins // 按值逆排序
}

func (ms MapSorter) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

func GetLocalIpStr() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return ""
	}

	localIP := ""
	i := 0
	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {

				fmt.Println(ipnet.IP.String())

				if i == 2 {
					localIP = ipnet.IP.String()
				}
			}

		}
		i++
	}

	return localIP
}

func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}
