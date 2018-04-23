package model

type ActionSendRecord struct {
	Id         int64
	Appid      string
	Uid        string
	ActionId   int64
	Power      float64
	Value      int64
	CreateTime int64
}

type AccessCode struct {
	Appid    string
	Uid      string
	Tel      string
	BindCode string
	Seq      uint64
}

type Member struct {
	Id               int64
	Account          string
	Password         string
	NickName         string  //昵称
	Avatar           string  //头像地址
	Sign             string  //签名
	Email            string  //Email地址
	LoginIp          string  //登陆IP地址
	RegIp            string  //注册IP
	RegisterTime     int64   //注册时间
	RegisterFrom     int8    //注册来源（欢迎补充）: 0, web; 1, Android; 2, iOS
	LoginTime        int64   //最后一次登陆时间
	LoginNum         int64   //用户登陆次数
	Status           int8    //用户状态
	HdtBalance       float64 //HDT余额
	RealnameAuth     int8    //实名认证: 1,已认证; 0,未认证
	TelegramCode     string  //电报群验证码
	BindTelegramName string  //绑定电报帐号名
}
