# proto_code_gen
使用protobuf的项目,proto_code_gen可以直接解析原始的*.proto文件,基于bufbuild/protocompile库,做一些预处理代码生成工作,无需依赖protoc生成的*.pb.go文件。

proto_code_gen目前提供2个功能:
 - 生成proto的Message的消息号映射(应用场景:proto当作网络消息协议)
 - 根据自定义模板生成代码

## 命令行参数

| 参数 | 必填 | 说明 |
|------|------|------|
| `-p` | 是 | proto文件目录(支持相对路径),程序扫描该目录下所有*.proto文件 |
| `-m` | 否 | 输出的CommandMapping文件名 |
| `-t` | 否 | 代码生成的模板文件路径 |
| `-c` | 否 | 代码生成的输出文件路径 |
| `-e` | 否 | 解析时排除的proto文件名(逗号分隔,可不带.proto后缀) |
| `-oe` | 否 | 输出时排除的proto文件名(逗号分隔,可不带.proto后缀) |

`-m` 和 `-t/-c` 至少有一个有效。`-t` 和 `-c` 必须同时提供才启用代码生成。

## 应用场景1: 自动生成网络消息协议号
proto_code_gen可以根据proto的message名自动生成网络消息协议号,并解决消息号冲突的问题。

```console
proto_code_gen -p ./proto -m ./gen/message_command_mapping.json
```

生成的映射文件格式:
```json
{"Child":31533,"Example":41475,"Example2":48661,"ExampleWithoutTag":43057}
```

## 应用场景2: 生成模板代码
项目中,我们经常希望能给proto的message增加一些自定义的设置,如struct tag。

Step1: 在proto文件中,使用自定义tag,示例参考examples/proto/cfg.proto
```proto
// file: examples/proto/cfg.proto
syntax = "proto3";

package pb;
option go_path = "/pb";

// 测试message's struct tag
// @StructTagOfExample
message Example {
  int32 Field1 = 1;
  string Field2 = 2;
}
```

Step2: 配置代码模板,示例参考examples/tmpl/csharp_cs.template

模板使用go自带的text/template

Step3: 运行proto_code_gen
```console
proto_code_gen -p ./examples/proto -t ./examples/tmpl/csharp_cs.template -c ./examples/gen/csharp.cs
```

同时生成消息号映射和模板代码:
```console
proto_code_gen -p ./proto -m ./gen/message_command_mapping.json -t ./template.tmpl -c ./gen/output.go
```

## 排除proto文件
`-e`和`-oe`的区别:
- `-e`: 解析proto文件时排除,这些文件完全不参与后续处理
- `-oe`: 输出时排除,这些文件的message参与了消息号计算,但在保存CommandMapping和生成代码时不包含它们

```console
proto_code_gen -p ./proto -m ./gen/mapping.json -e "test,internal" -oe "server_only"
```

## 使用proto_code_gen
获取
```console
go get github.com/fish-tennis/proto_code_gen
```
运行
```console
proto_code_gen -p ./proto -m ./gen/message_command_mapping.json
```

## 参考
[text/template](https://pkg.go.dev/text/template)  
[bufbuild/protocompile](https://github.com/bufbuild/protocompile)
