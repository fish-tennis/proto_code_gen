package main

import (
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"text/template"
)

func generateCodes(parserResult *ParserResult, codeTemplate *CodeTemplate) {
	var sortProtoList []string
	for protoName := range parserResult.allProto {
		sortProtoList = append(sortProtoList, UpperWordsToCamelCase(strings.TrimSuffix(protoName, ".proto"), "_", true))
	}
	sort.Slice(sortProtoList, func(i, j int) bool {
		return sortProtoList[i] < sortProtoList[j]
	})
	var messageList []*ProtoMessageStructInfo
	for _, structInfoList := range parserResult.allProto {
		for _, structInfo := range structInfoList {
			messageList = append(messageList, structInfo)
		}
	}
	outFileName := codeTemplate.OutDir + strings.TrimSuffix(path.Base(codeTemplate.Template), ".template")
	if len(messageList) == 0 {
		os.Remove(outFileName)
		return
	}
	tmpl, err := template.ParseFiles(codeTemplate.Template)
	if err != nil {
		log.Printf("parse Template file failed:%v %v", codeTemplate.Template, err)
		return
	}
	err = os.Mkdir(path.Dir(outFileName), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Printf("create dir failed:%v %v", path.Dir(outFileName), err)
		return
	}
	outFile, err := os.OpenFile(outFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Printf("open OutFile failed:%v %v", outFile, err)
		return
	}
	defer outFile.Close()
	err = tmpl.Execute(outFile, map[string]any{
		"MessageList": messageList,
		"ProtoList":   sortProtoList,
	})
	if err != nil {
		log.Printf("Execute Template failed:%v %v", outFileName, err)
		return
	}
	log.Printf("generate code:%v protoFileCount:%v", outFileName, len(sortProtoList))
}
