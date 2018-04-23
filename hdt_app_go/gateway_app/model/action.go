package model

type UserActionReq struct {
	Appid    string `json:"Appid"`
	Uid      string `json:"Uid"`
	Actionid int64  `json:"Actionid"`
	Value    string `json:"Value,omitempty"`
	Seq      int    `json:"Seq"`
}
