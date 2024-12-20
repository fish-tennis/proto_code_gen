# proto_code_gen
项目中,我们经常希望能给proto生成message增加一些自定义的设置,如struct tag.

增加自定义的struct tag,可以用[protoc-go-inject-tag](https://github.com/favadi/protoc-go-inject-tag)

但是protoc-go-inject-tag只能给message的字段加struct tag,因为golang并没有给struct结构体提供tag接口.

proto_code_gen提供了一种给message增加类似struct tag的方式

## 应用场景1:生成模板代码
Step1: 在proto文件中,使用自定义tag,示例参考examples/proto/cfg.proto
```proto
// file: examples/proto/cfg.proto
syntax = "proto3";

package pb;
option go_package = "/pb";

// 测试message's struct tag
// @StructTagOfExample
message Example {
  int32 Field1 = 1;
  string Field2 = 2;
}
```

Step2: protoc生成*.pb.go文件,如examples/pb/cfg.pb.go

Step3: 配置代码模板,示例参考examples/message_gen.go.template,examples/reader_gen.go.template

模板使用go自带的text/template

Step4: 运行protoc_code_gen -input=/examples/pb/*.pb.go -config=./code_templates.json

会根据模板生成对应的代码,如examples/gen/example_gen.go,examples/gen/cfg_reader_gen.go
```go
// file: examples/gen/example_gen.go
package gen
import (
  . "github.com/fish-tennis/proto_code_gen/examples/pb"
)
// 测试message's struct tag
func ExampleToString(m *Example) string {
  // just a test code
  return m.String()
}
```

## 应用场景2: 配置数据的只读接口
某些应用场景,会使用protobuf的结构来当作配置数据的格式,proto_code_gen提供了一种生成protobuf只读接口的功能,类似c++中的const.

如examples里的examples/reader_gen.go.template模板对应生成代码examples/gen/cfg_reader_gen.go

如[https://github.com/fish-tennis/gserver/tree/main/gen](https://github.com/fish-tennis/gserver/tree/main/gen)目录下的代码就是使用proto_code_gen生成的

## 使用proto_code_gen
获取
```console
go get github.com/fish-tennis/proto_code_gen
```
运行
```console
protoc_code_gen -input=/dir/*.pb.go -config=./code_templates.json
```
项目需要根据自己的需求,配置对应的代码模板,从而生成不同的代码.

## 原理
protoc_code_gen使用golang的parser库,解析*pb.go文件,读取其中的message结构体上的注释.

## 参考
[protoc-go-inject-tag](https://github.com/favadi/protoc-go-inject-tag)

[text/template](https://pkg.go.dev/text/template)
