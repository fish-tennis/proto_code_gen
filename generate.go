package main

import (
	"fmt"
	"go/ast"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"text/template"
)

// 生成相应的辅助代码
func generateCode(parserResult *ParserResult, key string) {
	codeTemplate := parserResult.GetCodeTemplate(key)
	// 排序一下,避免proto文件没改动,生成的代码文件却不一样
	var sortProtoList [][]*ProtoMessageStructInfo
	for _, structInfoList := range parserResult.protoMap {
		sortProtoList = append(sortProtoList, structInfoList)
	}
	sort.Slice(sortProtoList, func(i, j int) bool {
		return sortProtoList[i][0].protoName < sortProtoList[j][0].protoName
	})
	var messageList []*ProtoMessageStructInfo
	for _, structInfoList := range sortProtoList {
		for _, structInfo := range structInfoList {
			if structInfo.keyComment != codeTemplate.KeyComment {
				continue
			}
			messageList = append(messageList, structInfo)
		}
	}
	if len(messageList) == 0 {
		os.Remove(codeTemplate.OutFile)
		return
	}
	tmpl, err := template.ParseFiles(codeTemplate.Template)
	if err != nil {
		log.Printf("parse Template file failed:%v %v", codeTemplate.Template, err)
		return
	}
	err = os.Mkdir(path.Dir(codeTemplate.OutFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Printf("create dir failed:%v %v", path.Dir(codeTemplate.OutFile), err)
		return
	}
	outFile, err := os.OpenFile(codeTemplate.OutFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Printf("open OutFile failed:%v %v", codeTemplate.OutFile, err)
		return
	}
	defer outFile.Close()
	err = tmpl.Execute(outFile, map[string]any{
		"MessageList": messageList,
	})
	if err != nil {
		log.Printf("Execute Template failed:%v %v", codeTemplate.OutFile, err)
		return
	}
	log.Printf("generate code:%v messageCount:%v", codeTemplate.OutFile, len(messageList))
}

func generatePbReader(parserResult *ParserResult) {
	readerConfig := parserResult.readerTemplates
	if readerConfig.OutDir == "" && len(readerConfig.FileFilter) == 0 && len(readerConfig.MessageFilter) == 0 {
		return
	}
	log.Printf("generatePbReader:%v files:%v messages:%v", parserResult.readerTemplates.OutDir,
		readerConfig.FileFilter, readerConfig.MessageFilter)
	os.Mkdir(parserResult.readerTemplates.OutDir, os.ModePerm)
	tmpl, err := template.ParseFiles(readerConfig.Template)
	if err != nil {
		log.Printf("parse Template file failed:%v %v", readerConfig.Template, err)
		return
	}
	generateFileCount := 0
	generateMessageCount := 0
	for protoName, structInfoList := range parserResult.allProto {
		wholeFileMatch := len(readerConfig.FileFilter) > 0 && readerConfig.MatchFile(protoName)
		// proto文件是否使用了anypb.Any,自动import相关的lib
		importAnyPb := hasAnyPbField(structInfoList)
		var messageList []*ProtoMessageStructInfo
		outMessageCount := 0
		for _, structInfo := range structInfoList {
			if !wholeFileMatch {
				if len(readerConfig.MessageFilter) == 0 || !readerConfig.MatchMessage(structInfo.MessageName) {
					continue
				}
			}
			messageList = append(messageList, structInfo)
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
					log.Printf("%v.%v type parse error!", structInfo.MessageName, fieldName)
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
				readerField := &ReaderField{
					FieldName: fieldName,
					FieldType: fieldNameTypeName,
				}
				if isGenericSlice(field.Type) {
					elemTypeName := getElemTypeName(field.Type)
					readerField.IsNormalSlice = true
					readerField.ElemTypeName = elemTypeName
				} else if isStarSlice(field.Type) {
					elemTypeName := getElemTypeName(field.Type)
					if elemTypeName[0] == '*' {
						elemTypeName = elemTypeName[1:]
					}
					readerField.IsPtrSlice = true
					readerField.ElemTypeName = elemTypeName
					// TODO: map[k]v的特殊处理
				} else {
					if isStarField {
						if fieldNameTypeName[0] == '*' {
							readerField.FieldTypeWithoutPtr = fieldNameTypeName[1:]
						}
						readerField.IsStruct = true
					}
				}
				structInfo.Fields = append(structInfo.Fields, readerField)
			}
			outMessageCount++
		}
		outFileName := fmt.Sprintf("%v/%v_reader_gen.go", parserResult.readerTemplates.OutDir, strings.TrimSuffix(protoName, ".pb.go"))
		if outMessageCount > 0 {
			err = os.Mkdir(path.Dir(outFileName), os.ModePerm)
			if err != nil && !os.IsExist(err) {
				log.Printf("create dir failed:%v %v", path.Dir(outFileName), err)
				return
			}
			outFile, err := os.OpenFile(outFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
			if err != nil {
				log.Printf("open OutFile failed:%v %v", outFileName, err)
				return
			}
			defer outFile.Close()
			err = tmpl.Execute(outFile, map[string]any{
				"MessageList": messageList,
				"ImportAnyPb": importAnyPb,
			})
			if err != nil {
				log.Printf("Execute Template failed:%v %v", outFileName, err)
				return
			}
			generateFileCount++
			generateMessageCount += outMessageCount
		} else {
			os.Remove(outFileName)
		}
	}
	log.Printf("generate reader fileCount:%v messageCount:%v", generateFileCount, generateMessageCount)
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
