package main

import (
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"text/template"
)

func generateCodes(parserResult *ParserResult, templateFile string, outFile string, outputExcludeFiles []string) {
	oeSet := make(map[string]bool)
	for _, f := range outputExcludeFiles {
		oeSet[f] = true
	}
	var sortProtoList []string
	for protoName := range parserResult.allProto {
		if oeSet[protoName] {
			continue
		}
		sortProtoList = append(sortProtoList, UpperWordsToCamelCase(strings.TrimSuffix(protoName, ".proto"), "_", true))
	}
	sort.Slice(sortProtoList, func(i, j int) bool {
		return sortProtoList[i] < sortProtoList[j]
	})
	var messageList []*ProtoMessageStructInfo
	for protoName, structInfoList := range parserResult.allProto {
		if oeSet[protoName] {
			continue
		}
		for _, structInfo := range structInfoList {
			messageList = append(messageList, structInfo)
		}
	}
	if len(messageList) == 0 {
		os.Remove(outFile)
		return
	}
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		log.Printf("parse Template file failed:%v %v", templateFile, err)
		return
	}
	err = os.MkdirAll(path.Dir(outFile), os.ModePerm)
	if err != nil {
		log.Printf("create dir failed:%v %v", path.Dir(outFile), err)
		return
	}
	outF, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Printf("open OutFile failed:%v %v", outFile, err)
		return
	}
	defer outF.Close()
	err = tmpl.Execute(outF, map[string]any{
		"MessageList": messageList,
		"ProtoList":   sortProtoList,
	})
	if err != nil {
		log.Printf("Execute Template failed:%v %v", outFile, err)
		return
	}
	log.Printf("generate code:%v protoFileCount:%v", outFile, len(sortProtoList))
}
