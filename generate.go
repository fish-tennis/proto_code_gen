package main

import (
	"log"
	"os"
	"sort"
	"strings"
)

// 生成相应的辅助代码
func generateCode(parserResult *ParserResult, key string) {
	codeTemplate := parserResult.GetCodeTemplate(key)
	builder := strings.Builder{}
	builder.WriteString(strings.Join(codeTemplate.Header, "\n"))
	// 排序一下,避免proto文件没改动,生成的代码文件却不一样
	var sortProtoList [][]*ProtoMessageStructInfo
	for _,structInfoList := range parserResult.protoMap {
		sortProtoList = append(sortProtoList, structInfoList)
	}
	sort.Slice(sortProtoList, func(i, j int) bool {
		return sortProtoList[i][0].protoName < sortProtoList[j][0].protoName
	})
	for _,structInfoList := range sortProtoList {
		for _,structInfo := range structInfoList {
			if structInfo.keyComment != codeTemplate.KeyComment {
				continue
			}
			// test.pb.go -> test
			protoFileName := strings.TrimSuffix(structInfo.protoName, ".pb.go")
			protoName := protoFileName
			// 首字母大写
			// test -> Test
			ProtoName := strings.ToUpper(protoFileName[:1]) + protoFileName[1:]
			messageName := structInfo.messageName
			funcStr := strings.Join(codeTemplate.FuncTemplate, "\n")
			// 替换掉代码模板中的关键字
			funcStr = strings.ReplaceAll(funcStr, "{MessageName}", messageName)
			// TestMessageXyz -> TESTMESSAGEXYZ
			funcStr = strings.ReplaceAll(funcStr, "{MESSAGENAME}", strings.ToUpper(messageName))
			// TestMessageXyz -> TEST_MESSAGE_XYZ
			funcStr = strings.ReplaceAll(funcStr, "{MESSAGE_NAME}", CamelCaseToUpperWords(messageName,"_"))
			// TEST_MESSAGE_XYZ -> TestMessageXyz
			funcStr = strings.ReplaceAll(funcStr, "{CamelMessageName}", UpperWordsToCamelCase(messageName,"_", true))
			// TEST_MESSAGE_XYZ -> Test_Message_Xyz
			funcStr = strings.ReplaceAll(funcStr, "{Camel_Message_Name}", UpperWordsToCamelCase(messageName,"_", false))

			funcStr = strings.ReplaceAll(funcStr, "{protoName}", protoName)
			// test -> Test
			funcStr = strings.ReplaceAll(funcStr, "{ProtoName}", ProtoName)

			funcStr = strings.ReplaceAll(funcStr, "{PackageName}", structInfo.pbPackageName)
			funcStr = strings.ReplaceAll(funcStr, "{Value}", structInfo.keyCommentValue)
			funcStr = strings.ReplaceAll(funcStr, "{Comment}", structInfo.normalComment)
			builder.WriteString(funcStr)
			builder.WriteString("\n")
		}
	}
	builder.WriteString(strings.Join(codeTemplate.Tail, "\n"))
	os.WriteFile(codeTemplate.OutFile, ([]byte)(builder.String()), 0644)
	log.Printf("OutFile:%v", codeTemplate.OutFile)
}