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
"ranking_hdt_dig":[{"tel":"13631627420","hdt":0.03262931}] //排行榜
}
或者
{"errcode":15}
或者
{"errcode":1}
```

######### IconList

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
```