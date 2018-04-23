package client_msg

type UserRegisterReq struct {
	Appid    string
	Uid      string
	Actionid int64
	Value    string
}

type AppUidAccessCodeReq struct {
	Appid string
	Uid   string
}

type AppKeyValidateReq struct {
	Appid string
}

//websocket消息定义
const (
	WebMsgIdLogin = 10000
	WebMsgIdClose = 10001
)

type LoginReq struct {
	UserName string
	Pwd      string
}

type LoginRes struct {
	Errcode int32
	Uid     string
	Name    string
	Token   string
}

type CloseConRes struct {
	Reason string
}

type ActionConfig struct {
	ActionId uint32 `json:"action_id"`
	Action   string `json:"action"`
	Name     string `json:"name"`
}
