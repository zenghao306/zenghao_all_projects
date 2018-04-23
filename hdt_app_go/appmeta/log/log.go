package log

import (
	"hdt_app_go/common/log"
	"hdt_app_go/common/timer"
	"io"
	"os"
	"runtime"
	"time"
)

var Log *EasyLog.Logger

func InitLog(Logpath string) {
	var out io.Writer

	nowtime := time.Now()
	filename := nowtime.Format("20060102")
	finalname := Logpath + filename

	f, err := os.OpenFile(finalname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {

	}
	out = io.MultiWriter(f)

	Log = EasyLog.New(out, "[go-blog]", EasyLog.Lshortfile|EasyLog.Ldate|EasyLog.Lmicroseconds)

	go TimerTask(Logpath)
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

func TimerTask(Logpath string) {
	defer PrintPanicStack()
	for {
		d := timer.NewDispatcher(1)
		stdtime := time.Now()
		tomorrow := stdtime.AddDate(0, 0, 1)
		t := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, time.Local)
		du := time.Duration(t.Unix()-stdtime.Unix()) * time.Second

		d.AfterFunc(du, func() {
			nowtime := time.Now()
			filename := nowtime.Format("20060102")
			finalname := Logpath + filename
			f, _ := os.OpenFile(finalname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
			Log.SetNewOutPutFile(f)

		})
		(<-d.ChanTimer).Cb()
	}
}
