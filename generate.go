package main

import (
	"fmt"
	"go/ast"
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
	for _, structInfoList := range parserResult.protoMap {
		sortProtoList = append(sortProtoList, structInfoList)
	}
	sort.Slice(sortProtoList, func(i, j int) bool {
		return sortProtoList[i][0].protoName < sortProtoList[j][0].protoName
	})
	for _, structInfoList := range sortProtoList {
		for _, structInfo := range structInfoList {
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
			funcStr = strings.ReplaceAll(funcStr, "{MESSAGE_NAME}", CamelCaseToUpperWords(messageName, "_"))
			// TEST_MESSAGE_XYZ -> TestMessageXyz
			funcStr = strings.ReplaceAll(funcStr, "{CamelMessageName}", UpperWordsToCamelCase(messageName, "_", true))
			// TEST_MESSAGE_XYZ -> Test_Message_Xyz
			funcStr = strings.ReplaceAll(funcStr, "{Camel_Message_Name}", UpperWordsToCamelCase(messageName, "_", false))

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
	writeErr := os.WriteFile(codeTemplate.OutFile, ([]byte)(builder.String()), 0644)
	if writeErr != nil {
		log.Printf("write failed:%v %v", codeTemplate.OutFile, writeErr)
	} else {
		log.Printf("OutFile:%v", codeTemplate.OutFile)
	}
}

func generatePbReader(parserResult *ParserResult) {
	if parserResult.readerTemplates.OutFile == "" {
		return
	}
	builder := strings.Builder{}
	builder.WriteString(strings.Join(parserResult.readerTemplates.Header, "\n"))
	structTemplateStr := `
type %vReader struct {
	v *%v
}
`
	fieldTemplateStr := `
func (r *%vReader) Get%v() %v {
	return r.v.%v
}
`
	for _, structInfo := range parserResult.allProto {
		readerStr := fmt.Sprintf(structTemplateStr,
			structInfo.messageName,
			structInfo.messageName,
			)
		builder.WriteString(readerStr)
		builder.WriteString("\n")
		fieldList := structInfo.structType.Fields
		for _,field := range fieldList.List {
			if len(field.Names) == 0 {
				continue
			}
			fieldName := field.Names[0].Name
			if !ast.IsExported(fieldName) {
				continue
			}
			fieldNameTypeName := getTypeName(field.Type)
			if fieldNameTypeName == "" {
				log.Printf("%v.%v type parse error!", structInfo.messageName, fieldName)
				continue
			}
			fieldStr := fmt.Sprintf(fieldTemplateStr,
				structInfo.messageName,
				fieldName,
				fieldNameTypeName,
				fieldName,
				)
			builder.WriteString(fieldStr)
			builder.WriteString("\n")
		}
	}
	writeErr := os.WriteFile(parserResult.readerTemplates.OutFile, ([]byte)(builder.String()), 0644)
	if writeErr != nil {
		log.Printf("write failed:%v %v", parserResult.readerTemplates.OutFile, writeErr)
	} else {
		log.Printf("OutFile:%v", parserResult.readerTemplates.OutFile)
	}
}

func getTypeName(expr ast.Expr) string {
	switch typ := expr.(type) {
	case *ast.Ident:
		return typ.Name
	case *ast.StarExpr:
		return fmt.Sprintf("*%v", getTypeName(typ.X))
	case *ast.SliceExpr:
		return fmt.Sprintf("[]%v", getTypeName(typ.X))
	case *ast.ArrayType:
		return fmt.Sprintf("[]%v", getTypeName(typ.Elt))
	case *ast.MapType:
		return fmt.Sprintf("map[%v]%v", getTypeName(typ.Key), getTypeName(typ.Value))
	case *ast.SelectorExpr:
		return fmt.Sprintf("%v.%v", getTypeName(typ.X), typ.Sel.Name)
	default:
		log.Printf("unsupport type:%v", typ)
	}
	return ""
}
