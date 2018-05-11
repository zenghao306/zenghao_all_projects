package session

import (
	//"encoding/json"
	//"github.com/olahol/melody"
	"hdt_app_go/gateway_app/http/melody"
	"github.com/tidwall/gjson"
	cli "hdt_app_go/gateway_app/client_msg"
	. "hdt_app_go/gateway_app/log"
	//. "hdt_app_go/gateway_app/rpc"
	//proto "hdt_app_go/protcol"
	"encoding/json"
	"sync/atomic"
)

var (
	MSocket  *melody.Melody
	UUid     int32
	sessions map[int32]*melody.Session
)

func InitWebScoket() {
	MSocket = melody.New()
	MSocket.HandleConnect(WebSocketConnect)
	MSocket.HandleMessage(WebSocketHandMessage)
	UUid = 0
}

func WebSocketConnect(s *melody.Session) {
	atomic.AddInt32(&UUid, 1)
	Log.Debugf("connect a new conn id is %d", UUid)
	s.Set("cid", UUid)
}

func WebSocketHandMessage(s *melody.Session, msg []byte) {
	m := gjson.GetBytes(msg, "msg_id")
	if m.Exists() == false {
		Log.Err("conn send valid messsage")
		return
	}

	switch m.Int() {
	case cli.WebMsgIdLogin:
		//b := gjson.GetBytes(msg, "body")
		//login := b.Value().(cli.LoginReq)
		//r := &cli.LoginRes{}
		//r.Errcode, r.Uid, r.Name, r.Token = RpcClient.Register.Login(login.UserName, login.Pwd)
		//if r.Errcode == proto.ERR_OK {
		//	AddUser(r.Uid, s)
		//}
		//res, _ := json.Marshal(r)
		//s.WriteBinary(res)
	default:
		return
	}

}

func SendMsgToUser(uid int32, msg interface{}) bool {
	user := GetUserSessByUid(uid)
	if user != nil {
		if b, err := json.Marshal(msg); err == nil {
			MSocket.SendToSelf(b, sessions[uid])
			return true
		}
	}
	return false
}

func GetUserSessByUid(uid int32) *melody.Session {
	//mutex_session_guard.RLock()
	//defer mutex_session_guard.RUnlock()
	v, ok := sessions[uid]
	if !ok {
		return nil
	}
	return v
}