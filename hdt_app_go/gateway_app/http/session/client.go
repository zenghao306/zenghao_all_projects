package session

import (
	"encoding/json"
	"github.com/olahol/melody"
	cli "hdt_app_go/gateway_app/client_msg"
	. "hdt_app_go/gateway_app/log"
	proto "hdt_app_go/protcol"
	"sync"
)

//客户端管理器
var (
	user_manager map[string]*UserConn
	utex         sync.Mutex
)

func AddUser(uid string, session *melody.Session) bool {
	u := &UserConn{
		Uid:  uid,
		Sess: session,
	}
	utex.Lock()
	defer utex.Unlock()
	r, ok := user_manager[uid]
	if ok {
		if ret := r.CloseCon(); ret != proto.ERR_OK {
			return false
		}
	}
	user_manager[uid] = u
	return true
}

//单个客户端连接
type UserConn struct {
	Uid  string
	Sess *melody.Session
}

func (s *UserConn) CloseCon() int {
	m := &cli.CloseConRes{}
	b, err := json.Marshal(m)
	if err != nil {
		Log.Err(err.Error())
		return proto.ERR_JOSN
	}
	err = s.Sess.CloseWithMsg(b)
	if err != nil {
		Log.Err(err.Error())
		return proto.ERR_UNKNOWN
	}
	return proto.ERR_OK
}
func ParseMsg() {

}
