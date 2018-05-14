package wrap

import (
	"github.com/martini-contrib/render"
)

type Call struct {
	Request interface{}
	Reply   interface{}
	Render  render.Render
	Done    chan *Call //用于结果返回时,消息通知,使用者必须依靠这个来获取真正的结果。
	Uid     int
	GiftID  int
	Num     int
	RevId   int
	Token   string
}

/*
var (
	msg_proc chan *Call
)
*/

func (call *Call) DoneV2() {
	select {
	case call.Done <- call:
		// ok
	default:
		// 阻塞情况处理,这里忽略
	}
}

/*
func GO(req int,reply *int,done chan *Call)*Call{
	if done==nil{
		done=make(chan *Call,10)
	}else{
		if cap(done)==0{
			fmt.Println("chan容量为0,无法返回结果,退出此次计算!")
			return nil
		}
	}
	call:=&Call{
		Request:req,
		Reply:reply,
		Done:done,
	}
	//调用一个可能比较耗时的计算，注意用"go"
	go caculate(call)
	return call
}

func caculate(call *Call){
	//假定运算一次需要耗时1秒

	call.done()
}
*/
/*
func init() {
	msg_proc = make(chan *Call, 10)
	go func() {
		for {
			select {

			case msg, ok := <-msg_proc:
				if ok {
					SendGift(msg)
				}
			}
		}
	}()
}

*/
