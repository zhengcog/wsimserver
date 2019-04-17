# webim

基于WebSocket协议,使用golang官方websocket包

- [详细文档](https://godoc.org/golang.org/x/net/websocket)

## Installation

```
//启动websocket服务
cd yourserverpath
go mod tidy
go build .
[root@i-9yfnjavm ~]# ./wsimserver
```

## Usage
连接websocket服务器地址

```
ws://127.0.0.1:9090/server

```

登录
```
token=Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEwMDAwMDAwMCwiZXhwIjoxNTE4ODU1MDI2LCJpc3MiOiJ6aGVuZ2NvZ0BnbWFpbC5jb20iLCJuYmYiOjE1MDMzMDMwMjZ9.JIOOrmV1Izy0Gj5IedXA7Gy7ejHaKi3t7Dj4aTSPWpo
```

发送消息:

```
文本消息：
{
    "target_type" : "user", // user 单聊 目前只有单聊
    "target": 100000034,// 消息接收者
    "msg" : {
        "type" : "txt",//消息类型 txt文本消息 pic图片 video视频 audio 语音 （目前只有文本消息类型）
        "content" : "hello from rest" //消息内容
    },
    "ext":{ //扩展属性，自己定义。可以没有这个字段，但是如果有，值不能是"ext:null"这种形式，否则出错
        "sendername":"发送者昵称",
        "senderavatar":"发送者头像",
        "receivername":"接受者昵称",
        "receiveravatar":"接受者头像",
        "attr1...":"v1..."   // 消息的扩展内容，可以增加字段，扩展消息主要解析部分，必须是基本类型数据。
    },
    "from" : "100000000" //表示消息发送者。无from请求失败
}
```

接收消息:

```
文本消息：
{
    "target_type": "user",//user单聊 目前只有单聊
    "target": "100000034",//消息接收者
    "msg": {
        "type": "txt", //消息类型 txt文本消息 pic图片 video视频 audio 语音 （目前只有文本消息类型）
        "content": "hello word" //消息内容
    },
    "ext": {
        "sendername":"发送者昵称",
        "senderavatar":"发送者头像",
        "receivername":"接受者昵称",
        "receiveravatar":"接受者头像",
        "attr1":"v1" //扩展属性，自己定义。可以没有这个字段，但是如果有，值不能是"ext:null"这种形式，否则出错
    },
    "from": "100000000",//消息发送者
    "timestamp": 1503388323 //消息发送时间戳，距离1970.1的秒数
}
```
