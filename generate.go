package main

import (
	"fmt"
	"go/ast"
	"log"
	"os"
	"path"
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
	os.Mkdir(path.Dir(codeTemplate.OutFile), os.ModePerm)
	writeErr := os.WriteFile(codeTemplate.OutFile, ([]byte)(builder.String()), 0644)
	if writeErr != nil {
		log.Printf("write failed:%v %v", codeTemplate.OutFile, writeErr)
	} else {
		log.Printf("OutFile:%v", codeTemplate.OutFile)
	}
}

func generatePbReader(parserResult *ParserResult) {
	readerConfig := parserResult.readerTemplates
	if readerConfig.OutDir == "" && len(readerConfig.FileFilter) == 0 && len(readerConfig.MessageFilter) == 0 {
		return
	}
	log.Printf("generatePbReader:%v files:%v messages:%v", parserResult.readerTemplates.OutDir,
		readerConfig.FileFilter, readerConfig.MessageFilter)
	os.Mkdir(parserResult.readerTemplates.OutDir, os.ModePerm)

	// xxxReader模板
	structTemplateStr := `
type {MessageName}Reader struct {
	v *{MessageName}
}

func New{MessageName}Reader(src *{MessageName}) *{MessageName}Reader {
	return &{MessageName}Reader{v:src}
}

func (r *{MessageName}Reader) Raw() *{MessageName} {
	return r.v
}
`

	// 普通字段模板
	fieldTemplateStr := `
func (r *{MessageName}Reader) Get{FieldName}() {FieldType} {
	return r.v.Get{FieldName}()
}
`

	// proto2的星号字段模板
	starFieldTemplateStr := `
func (r *{MessageName}Reader) Get{FieldName}() {FieldType}Reader {
	return New{FieldTypeNoStar}Reader(r.v.Get{FieldName}())
}
`

	// 基础类型数组字段模板
	genericSliceFieldTemplateStr := `
func (r *{MessageName}Reader) Get{FieldName}() []{FieldType} {
	src := r.v.Get{FieldName}()
	if src == nil {
		return nil
	}
	copySlice := make([]{FieldType},len(src))
	copy(copySlice, src)
	return copySlice
}
`

	// 星号数组字段模板
	starSliceFieldTemplateStr := `
func (r *{MessageName}Reader) Get{FieldName}() []{FieldType}Reader {
	src := r.v.Get{FieldName}()
	if src == nil {
		return nil
	}
	sliceReader := make([]{FieldType}Reader,len(src))
	for i,v := range src {
		sliceReader[i] = New{FieldTypeNoStar}Reader(v)
	}
	return sliceReader
}
`

	for protoName, structInfoList := range parserResult.allProto {
		wholeFileMatch := len(readerConfig.FileFilter) > 0 && readerConfig.MatchFile(protoName)
		builder := strings.Builder{}
		builder.WriteString(strings.Join(parserResult.readerTemplates.Header, "\n"))
		// proto文件是否使用了anypb.Any,自动import相关的lib
		importAnyPb := hasAnyPbField(structInfoList)
		if len(parserResult.readerTemplates.Import) > 0 || importAnyPb {
			builder.WriteString("import (\n")
			for _, v := range parserResult.readerTemplates.Import {
				builder.WriteString("\t")
				builder.WriteString(v)
				builder.WriteString("\n")
			}
			if importAnyPb {
				builder.WriteString("\t")
				builder.WriteString("\"google.golang.org/protobuf/types/known/anypb\"")
				builder.WriteString("\n")
			}
			builder.WriteString(")\n")
		}
		outMessageCount := 0
		for _, structInfo := range structInfoList {
			if !wholeFileMatch {
				if len(readerConfig.MessageFilter) == 0 || !readerConfig.MatchMessage(structInfo.messageName) {
					continue
				}
			}
			readerStr := strings.ReplaceAll(structTemplateStr, "{MessageName}", structInfo.messageName)
			builder.WriteString(readerStr)
			builder.WriteString("\n")
			fieldList := structInfo.structType.Fields
			for _, field := range fieldList.List {
				if len(field.Names) == 0 {
					continue
				}
				fieldName := field.Names[0].Name
				if !ast.IsExported(fieldName) {
					continue
				}
				if parserResult.readerTemplates.ProtoV2 {
					if fieldName == "XXX_unrecognized" {
						continue
					}
				}
				fieldNameTypeName := getTypeName(field.Type)
				if fieldNameTypeName == "" {
					log.Printf("%v.%v type parse error!", structInfo.messageName, fieldName)
					continue
				}
				isStarField := isStar(field.Type)
				if parserResult.readerTemplates.ProtoV2 {
					// proto2的基础类型的引用类型转换成值类型
					if isGenericTypeName(fieldNameTypeName, true) {
						isStarField = false
						fieldNameTypeName = fieldNameTypeName[1:]
					}
				}
				fieldStr := ""
				if isGenericSlice(field.Type) {
					elemTypeName := getElemTypeName(field.Type)
					fieldStr = strings.ReplaceAll(genericSliceFieldTemplateStr, "{MessageName}", structInfo.messageName)
					fieldStr = strings.ReplaceAll(fieldStr, "{FieldName}", fieldName)
					fieldStr = strings.ReplaceAll(fieldStr, "{FieldType}", elemTypeName)
				} else if isStarSlice(field.Type) {
					elemTypeName := getElemTypeName(field.Type)
					fieldStr = strings.ReplaceAll(starSliceFieldTemplateStr, "{MessageName}", structInfo.messageName)
					fieldStr = strings.ReplaceAll(fieldStr, "{FieldName}", fieldName)
					fieldStr = strings.ReplaceAll(fieldStr, "{FieldType}", elemTypeName)
					if elemTypeName[0] == '*' {
						fieldStr = strings.ReplaceAll(fieldStr, "{FieldTypeNoStar}", elemTypeName[1:])
					} else {
						fieldStr = strings.ReplaceAll(fieldStr, "{FieldTypeNoStar}", elemTypeName)
					}
					// TODO: map[k]v的特殊处理
				} else {
					if isStarField {
						fieldStr = strings.ReplaceAll(starFieldTemplateStr, "{MessageName}", structInfo.messageName)
						fieldStr = strings.ReplaceAll(fieldStr, "{FieldName}", fieldName)
						fieldStr = strings.ReplaceAll(fieldStr, "{FieldType}", fieldNameTypeName)
						if fieldNameTypeName[0] == '*' {
							fieldStr = strings.ReplaceAll(fieldStr, "{FieldTypeNoStar}", fieldNameTypeName[1:])
						} else {
							fieldStr = strings.ReplaceAll(fieldStr, "{FieldTypeNoStar}", fieldNameTypeName)
						}
					} else {
						fieldStr = strings.ReplaceAll(fieldTemplateStr, "{MessageName}", structInfo.messageName)
						fieldStr = strings.ReplaceAll(fieldStr, "{FieldName}", fieldName)
						fieldStr = strings.ReplaceAll(fieldStr, "{FieldType}", fieldNameTypeName)
					}
				}
				builder.WriteString(fieldStr)
			}
			outMessageCount++
		}
		outFileName := fmt.Sprintf("%v/%v_reader_gen.go", parserResult.readerTemplates.OutDir, strings.TrimSuffix(protoName, ".pb.go"))
		if outMessageCount > 0 {
			writeErr := os.WriteFile(outFileName, ([]byte)(builder.String()), 0644)
			if writeErr != nil {
				log.Printf("write failed:%v %v", outFileName, writeErr)
			} else {
				log.Printf("OutFile:%v", outFileName)
			}
		} else {
			os.Remove(outFileName)
		}
	}
}

// 是否是基础类型(bool int uint float string)
func isGenericTypeName(typeName string, isStar bool) bool {
	genericTypes := []string{"bool", "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "float32", "float64", "string"}
	for _, v := range genericTypes {
		if isStar {
			if typeName == fmt.Sprintf("*%v", v) {
				return true
			}
		} else {
			if typeName == v {
				return true
			}
		}
	}
	return false
}

// 类型string
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

// 是否是[]*xxx格式的数组
func isStarSlice(expr ast.Expr) bool {
	switch typ := expr.(type) {
	case *ast.SliceExpr:
		return isStar(typ.X)
	case *ast.ArrayType:
		return isStar(typ.Elt)
	}
	return false
}

// 是否是星号字段
func isStar(expr ast.Expr) bool {
	switch expr.(type) {
	case *ast.StarExpr:
		return true
	}
	return false
}

// 是否是基础类型的数组([]int []string)
func isGenericSlice(expr ast.Expr) bool {
	switch typ := expr.(type) {
	case *ast.SliceExpr:
		return isGenericTypeName(getTypeName(typ.X), false)
	case *ast.ArrayType:
		return isGenericTypeName(getTypeName(typ.Elt), false)
	}
	return false
}

// slice的elem类型
func getElemTypeName(expr ast.Expr) string {
	switch typ := expr.(type) {
	case *ast.SliceExpr:
		return getTypeName(typ.X)
	case *ast.ArrayType:
		return getTypeName(typ.Elt)
	}
	return ""
}

// 是否有anypb.Any类型的字段
func hasAnyPbField(structInfoList []*ProtoMessageStructInfo) bool {
	for _, structInfo := range structInfoList {
		fieldList := structInfo.structType.Fields
		for _, field := range fieldList.List {
			if len(field.Names) == 0 {
				continue
			}
			fieldName := field.Names[0].Name
			if !ast.IsExported(fieldName) {
				continue
			}
			fieldNameTypeName := getTypeName(field.Type)
			if strings.HasSuffix(fieldNameTypeName, "anypb.Any") {
				return true
			}
		}
	}
	return false
}
