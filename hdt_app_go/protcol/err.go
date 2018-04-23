package external

//protoc --proto_path=$GOPATH/src:. --micro_out=. --go_out=. /data/src/hdt/*.proto
//protoc --proto_path=$GOPATH/src:. --micro_out=. --go_out=. ./*.proto

/*
windows
protoc -I$GOPATH/src  --go_out=plugins=micro:rpc  $GOPATH/src/HDTapSDK_Go/protcol/*.proto
protoc -I$GOPATH/src  --go_out=$GOPATH/src/HDTapSDK_Go/protcol/   --plugins=micro.exe:$GOPATH/bin    $GOPATH/src/HDTapSDK_Go/protcol/*.proto

protoc -I rpc rpc/*.proto --plugin=micro=$GOPATH/bin/micro.exe --go_out=rpc

protoc --go_out=plugins=micro.exe:$GOPATH/src $GOPATH/src/HDTapSDK_Go/protcol/*.proto
protoc --plugin=protoc-gen-NAME=path/to/mybinary.exe --go_out=$GOPATH/bin
protoc -I$GOPATH/src    --plugin=micro=$GOPATH/bin --go_out=$GOPATH/src/HDTapSDK_Go/protcol/   $GOPATH/src/HDTapSDK_Go/protcol/*.proto

protoc -I$GOPATH/src -I$GOPATH/bin  --go_out=plugins=micro.exe:E:/mygo/src/HDTapSDK_Go/protcol $GOPATH/src/HDTapSDK_Go/protcol/*.proto

protoc -I$GOPATH/src --go_out=plugins=micro.exe:$GOPATH/src $GOPATH/src/HDTapSDK_Go/protcol/*.proto


protoc -I$GOPATH/src/github.com/micro/protobuf/protoc-gen-go/micro -I rpc  rpc/*.proto  --go_out=plugins=micro:rpc
linux
protoc -I$GOPATH/src --go_out=plugins=micro:$GOPATH/src  $GOPATH/src/hdt/*.proto
protoc  --go_out=plugins=micro:.  $GOPATH/src/hdt/*.proto
*/
const (
	ERR_OK                = 0
	ERR_UNKNOWN           = 1
	ERR_PARAM             = 2
	ERR_EXIST_ACTION      = 3
	ERR_JOSN              = 4
	ERR_LIMIT             = 5
	ERR_EXIST_APPID       = 6
	ERR_REGISTER          = 7
	ERR_PARAM_SEQ         = 8
	ERR_SNS_TIMEOUT       = 9
	ERR_SNS_CORRECT       = 10
	ERR_VERIFT_CODE       = 11 //验证错误码错误
	ERR_EXIST_ACCOUNT     = 12
	ERR_NOT_EXIST_ACCOUNT = 13
	ERR_PASSWORD          = 14
	ERR_EXPIRATION        = 15 //TOKEN已过期
)
