# proto_code_gen
项目中,我们经常希望能给proto生成message增加一些自定义的设置,如struct tag.

增加自定义的struct tag,可以用[protoc-go-inject-tag](https://github.com/favadi/protoc-go-inject-tag)

但是protoc-go-inject-tag只能给message的字段加struct tag,因为golang并没有给struct结构体提供tag接口.

proto_code_gen提供了一种给message增加类似struct tag的方式

## 游戏项目常用的proto风格
```proto
// file: test.proto
syntax = "proto3";

package pb;
option go_package = "/pb";

enum CmdTest {
    Cmd_None = 0;
    Cmd_Req = 1; // 客户端的请求消息
    Cmd_Res = 2; // 服务器回复消息
}

// 客户端的请求消息
message Req {
  string content = 1;
}

// 服务器回复消息
message Res {
  string result = 1;
}
```

## 服务器常见代码风格
```go
// 回调函数
func OnReq(conn Connection, m Message) {
	// 手动转换message类型
	req := m.(*pb.Req)
	// ...
	res := new(pb.Res)
	// 手动填写消息号
	conn.Send(pb.CmdTest_Cmd_Res, res)
}
// 手动注册消息回调,手动填写消息号
register(pb.CmdTest_Cmd_Req, OnReq)
```

## 使用proto_code_gen可以生成的辅助代码
```proto
// file: test.proto
syntax = "proto3";

package pb;
option go_package = "/pb";

enum CmdTest {
    Cmd_None = 0;
    Cmd_Req = 1; // 客户端的请求消息
	Cmd_Res = 2; // 服务器回复消息
}

// 客户端的请求消息
// @Handler用来自动注册回调函数
// @Handler
message Req {
  string content = 1;
}

// 服务器回复消息
// @Player用来自动生成发送消息的接口
// @Player
message Res {
  string result = 1;
}
```
```go
// file: auto_register_gen.go
// 工具生成的自动注册函数,无需为每个消息单独注册
func auto_register() {
	// 工具生成的自动注册函数,应用层无需手动填写消息号
    register(pb.CmdTest_Cmd_Req, OnReq)
    // ...
}
```
```go
// file: send_gen.go
// 工具生成的发送pb.Res的函数,应用层调用更友好
// 无需手动填写消息号,且增强了消息类型检查
func sendRes(conn Connection, res *pb.Res) {
    conn.Send(pb.CmdTest_Cmd_Res, res)
}
```

## 使用proto_code_gen
获取
```console
go get github.com/fish-tennis/proto_code_gen
```
运行
```console
// windows系统: -input=\\dir\\*.pb.go
protoc_code_gen -input=/dir/*.pb.go -config=./code_templates.json
```
项目需要根据自己的需求,修改code_templates.json里面的内容,从而生成不同的代码.

proto_code_gen项目里自带的code_templates.json是针对[gserver](https://github.com/fish-tennis/gserver)和[gnet](https://github.com/fish-tennis/gnet)用的

## 原理
protoc_code_gen使用golang的parser库,解析*pb.go文件,读取其中的message结构体上的注释.

如果code_templates.json配置了对应的{ReplaceKey}关键字,则按照代码模板进行相关的文本替换,并自动生成代码文件.

目前支持的{ReplaceKey}:
- {MessageName}: 消息结构体的名字
- {MESSAGENAME}: TestMessageXyz -> TESTMESSAGEXYZ
- {MESSAGE_NAME}: TestMessageXyz -> TEST_MESSAGE_XYZ
- {CamelMessageName}: TEST_MESSAGE_XYZ -> TestMessageXyz
- {Camel_Message_Name}: TEST_MESSAGE_XYZ -> Test_Message_Xyz
- {protoName}: proto文件名
- {ProtoName}: 首字母大写的proto文件名
- {PackageName}: *.pb.go文件的package名
- {Value}: 注释里可以加一个Value值,如@Player:TheValue
- {Comment}: message的注释(排除了@key)

## 参考
[protoc-go-inject-tag](https://github.com/favadi/protoc-go-inject-tag)