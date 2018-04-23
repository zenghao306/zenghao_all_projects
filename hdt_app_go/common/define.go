package common

import "time"

/*
const (
	ERR_SUCCESS     = iota // 成功 == 0
	ERR_UNKNOWN            // 未知错误 == 1
	ERR_EXIST              // 账号已经存在 == 2
	ERR_PWD                // 密码错误 == 3
	ERR_DB_ERROR           // 数据库错误 == 4
	ERR_PARAM              //参数错误==5
	ERR_VERIFT_CODE        // 验证码错误 == 6
	ERR_SNS_TIMEOUT        // 验证码超时 == 7

)
*/
const SmsKeepTime = 30000 * time.Second
