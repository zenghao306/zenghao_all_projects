# API


| reason | 错误码说明 |
| ----- | ----- |
| 0 | 正常返回 |
| 1 | 服务器错误 |
| 2 | 参数错误 |
| 3 | 不存在的Action |
| 4 | JSON解析错误 |
| 9 | 短信交易过期 |
| 10 | 短信验证失败 |
| 11 | 验证码错误 |
| 12 | 账号已存在 |
| 13 | 账号不存在 |
| 14 | 密码错误 |
| 15 | TOKEN过期 |
######### 发送短信验证

**url**

{{ host }}/send/sns

**method**

POST

**参数**

| 参数名 | 类型 | 说明 |
| ----- | ----- | ------ |
| tel | string | 接收短信验证码的手机号 |


**返回值**

```
{"code":"9869",
"errcode":0 //错误码
}
```

######### 用户注册

**url**

{{ host }}/register

**method**

POST

| 参数名 | 类型 | 说明 |
| ----- | ----- | ------ |
| tel | string | 注册手机号 |
| pwd | string | 注册密码 |
| code | string | 注册状态码 |
| register_from | int | 注册来源（1-安卓，2-IOS） |

**示例**

```
{"tel":"13631627420","pwd":"123456","code":"9279","register_from":1}

**返回值**

```
{"errcode":0}
或者
{"errcode":12}
或者
{"errcode":9}
或者
{"errcode":1}
```

######### 用户登陆

**url**

{{ host }}/login

**method**

POST

| 参数名 | 类型 | 说明 |
| ----- | ----- | ------ |
| tel | string | 注册手机号 |
| pwd | string | 注册密码 |

**示例**

```
{"tel":"13631627420","pwd":"123456"}

**返回值**

```
{"errcode":0,
"token":"79e27bec7ba1371ef046c9630fb38ff8",
"user_info":{
    "uid":10001,                                             //用户ID
    "nickName":"希望明天",                                   //昵称
    "avatar":"http://face.shangtv.cn/100001149611045044914", //头像地址
    "sign":"Sky",                                            //签名
    "email":"zenghao306@163.com",                            //email地址
    "regTime":1523239229,                                    //注册时间戳
    "hdtBalance":1.50000001                                  //互动币余额
}
}
或者
{"errcode":13}
或者
{"errcode":14}
或者
{"errcode":1}
```

######### 更改密码

**url**

{{ host }}/modify/pwd

**method**

POST

| 参数名 | 类型 | 说明 |
| ----- | ----- | ------ |
| tel | string | 注册手机号 |
| pwd | string | 注册密码 |
| code | string | 短信验证码 |

**示例**

```
{"tel":"13631627420","pwd":"123456","code":"1234"}

**返回值**

```
{"errcode":0,
"token":"79e27bec7ba1371ef046c9630fb38ff8",
"user_info":{
    "uid":10001,                                             //用户ID
    "nickName":"希望明天",                                   //昵称
    "avatar":"http://face.shangtv.cn/100001149611045044914", //头像地址
    "sign":"Sky",                                            //签名
    "email":"zenghao306@163.com",                            //email地址
    "regTime":1523239229,                                    //注册时间戳
    "hdtBalance":1.50000001                                  //互动币余额
}
}
或者
{"errcode":13}
或者
{"errcode":14}
或者
{"errcode":1}
```

######### 产品首页排行榜以及用户HDT账户信息

**url**

{{ host }}/ranking/info

**method**

POST

| 参数名 | 类型 | 说明 |
| ----- | ----- | ------ |
| tel | string | 注册手机号 |
| token | string | 登陆的token |

**示例**

```
{"tel":"13631627420","token":"f5de14f22eb27bea7425183ad80d0926"}

**返回值**

```
{"difficulty":3.064729539719619,                           //当前难度系数
"errcode":0,                                               //错误码
"hdt_mining_last":0.03262931,                              //上结算周期挖到的HDT
"hdt_mining_total":10.24925397,                            //已挖到的HDT总计
"mining_index":1,                                          //挖矿排名
}
或者
{"errcode":15}
或者
{"errcode":1}
```

######### 产品首页排行榜信息

**url**

{{ host }}/ranking/hdt/dig

**method**

POST

| 参数名 | 类型 | 说明 |
| ----- | ----- | ------ |
| tel | string | 注册手机号 |
| token | string | 登陆的token |

**示例**

```
{"tel":"13631627420","token":"f5de14f22eb27bea7425183ad80d0926"}

**返回值**

```
{
"errcode":0,                                               //错误码
"ranking_hdt_dig":[{"tel":"13631627420","hdt":0.03262931}] //排行榜
}
或者
{"errcode":15}
或者
{"errcode":1}
```

######### APP列表

**url**

{{ host }}/app/list

**method**

POST

| 参数名 | 类型 | 说明 |
| ----- | ----- | ------ |
| tel | string | 注册手机号 |
| token | string | 登陆的token |
| index | int | 分页index（从0开始计数） |

**示例**

```
{"tel":"13631627420","token":"417097f543f34458461453751f0db2c1","index":1}

**返回值**

```
{"errcode":0,"list":[
{"AppName":"测试应用",       //APP名
"AppIcoPath":".../8.png"},]} //APP图片路径
或者
{"errcode":15}
或者
{"errcode":1}
```

######### APP详情

**url**

{{ host }}/app/detail/info

**method**

POST

| 参数名 | 类型 | 说明 |
| ----- | ----- | ------ |
| tel | string | 注册手机号 |
| token | string | 登陆的token |
| app_id | int64 | 对应APP的ID |

**示例**

```
{"tel":"13631627420","token":"50ad9046474d679a7eb79fbd616e9509","app_id":15}

**返回值**

```
{"app_android_address":"0",
"app_content":"友电-全新视频社交领航者！实时的一对一视频语音社交约聊社区。这里有公司白领、美女学生、商务人士、模特空姐，各行各业的人士聚集在此，当你想：认识更多朋友、相遇知己、差旅疲惫、孤独寂寞、有烦恼却没人倾述的时候，打开友电和你心动的那个Ta通上一个电话吧，你会发现收获生活中的快乐其实也很简单！让相遇变简单，你我相遇，尽在友电!\u003cbr/\u003e友电的特点：\u003cbr/\u003e【真实社交】友电一直认为社交应该是真实的，做陌生人的真实社交，让相遇变得简单！\u003cbr/\u003e【同城交友】他乡遇故知，天涯共此时。找寻附近同城投缘的知己，在线来电。\u003cbr/\u003e【多样畅聊】支持文字、语音电话、视频聊天，想怎么聊就怎么聊！\u003cbr/\u003e【私密聊天】一对一视频语音聊天，私密一对一，拒绝尴尬！\u003cbr/\u003e【拒绝照骗】人工审核照片、视频、身份、行业等多种资料，安全有保障。\u003cbr/\u003e",
"app_hdt_total":0,
"app_imgs":["http://admin.hudt.io/public/static/upload/img/","http://admin.hudt.io/public/static/upload/img/","http://admin.hudt.io/public/static/upload/img/","http://admin.hudt.io/public/static/upload/img/","http://admin.hudt.io/public/static/upload/img/","http://admin.hudt.io/public/static/upload/img/","http://admin.hudt.io/public/static/upload/img/","http://admin.hudt.io/public/static/upload/img/","http://admin.hudt.io/public/static/upload/img/","http://admin.hudt.io/public/static/upload/img/","http://admin.hudt.io/public/static/upload/img/","http://admin.hudt.io/public/static/upload/img/","http://admin.hudt.io/public/static/upload/img/","http://admin.hudt.io/public/static/upload/img/","http://admin.hudt.io/public/static/upload/img/"],"app_ios_address":"   0","errcode":0,"user_app_hdt":0}
或者
{"errcode":15}
或者
{"errcode":1}
```

### 矿池首页信息

**说明**：返回矿池首页信息。

**url**
I
{{ host }}/mine/pool/info

**method**

POST

上传格式-JSON：
{"tel":"13631627420","token":"0b80150fbd2789251e6012117c346644"}

**参数**

| 参数名 | 类型 | 说明 |
| ----- | ----- | ------ |
| tel   | string | 用户手机号 |
| token| string | 用户token |

**返回值**

```
{
"errcode":0,
"app_hdt_balance_total":10693.0582,     //任务池余额
"degree_difficulty":3.06473,            //难度系数
"hdt_supply_limit":1000000000000000000, //HDT矿池总量
"hdt_total_supply":300000000000000000   //已产出总量
}
```

### 矿池任务列表

**说明**：返回任务列表信息。

**url**
I
{{ host }}/mine/pool/tast/list

**method**

POST

上传格式-JSON：
{"tel":"13631627420","token":"0b80150fbd2789251e6012117c346644"}

**参数**

| 参数名 | 类型 | 说明 |
| ----- | ----- | ------ |
| tel   | string | 用户手机号 |
| token| string | 用户token |

**返回参数说明**

| 参数名 | 类型 | 说明 |
| ----- | ----- | ------ |
| AppName   | string | AppName |
| AppIcoPath| string | App图片路径 |
| AppId| int | AppId |
| Time| int64 | 对应时间戳 |
| Style| int | 1为发布，2为挖出 |
| Hdt| double | 对应HDT数量|
| HdtTaskBalance| double | APP对应矿池余额|

**返回值**

```
{
	"errcode": 0,
	"list": [{
		"AppName": "友电",
		"AppIcoPath": "http://admin.hudt.io/public/static/upload/img/20180409/6ed0c9f540ea941056aa66aa2fc84ca2.png",
		"AppId": 15,
		"Time": 1523498358,
		"Style": 1,
		"Hdt": 500,
		"HdtTaskBalance": 9979.5252
	},{
		"AppName": "友电",
		"AppIcoPath": "http://admin.hudt.io/public/static/upload/img/20180409/6ed0c9f540ea941056aa66aa2fc84ca2.png",
		"AppId": 15,
		"Time": 1526014800,
		"Style": 2,
		"Hdt": 0.03263,
		"HdtTaskBalance": 9979.5252
	}]
}
```