package main

import (
	"encoding/json"
	"hash/crc32"
	"log"
	"os"
)

// 自动生成proto.Message和PacketCommand的映射
func generateCommandMapping(parserResult *ParserResult, outputFile string) error {
	if outputFile == "" {
		return nil
	}
	// 加载之前保存的数据
	mapping := loadCommandMapping(outputFile)
	cmdMapping := make(map[int]string)
	for messageName, cmd := range mapping {
		if _, ok := cmdMapping[cmd]; ok {
			log.Printf("conflict message:%v cmd:%v", messageName, cmd)
			delete(mapping, messageName)
			continue
		}
		cmdMapping[cmd] = messageName
	}
	conflict := make(map[string]int)
	allMessages := make(map[string]int)
	// 遍历所有Message,用hash算法,根据消息名生成消息号
	for _, structInfoList := range parserResult.allProto {
		for _, structInfo := range structInfoList {
			cmd := uint16(crc32.ChecksumIEEE([]byte(structInfo.MessageName)) & 0xFFFF)
			allMessages[structInfo.MessageName] = int(cmd)
		}
	}
	// 删除已经不存在的消息
	for messageName, cmd := range mapping {
		if _, ok := allMessages[messageName]; !ok {
			delete(mapping, messageName)
			delete(cmdMapping, cmd)
		}
	}
	// 记录冲突的消息
	for messageName, cmd := range allMessages {
		oldCmd, ok := mapping[messageName]
		if ok && cmd != oldCmd {
			conflict[messageName] = cmd // hash冲突的message
			continue
		}
		if cmd == 0 {
			conflict[messageName] = cmd // 消息号不能为0
			continue
		}
		mapping[messageName] = cmd
		cmdMapping[cmd] = messageName
	}
	// 冲突的消息,自动查找未使用的消息号,解决冲突
	for messageName, cmd := range conflict {
		log.Printf("conflict message:%v cmd:%v", messageName, cmd)
		hasNewCmd := false
		for i := 1; i < 0xFFFF; i++ {
			if _, ok := cmdMapping[i]; ok {
				continue
			}
			cmdMapping[i] = messageName
			mapping[messageName] = i
			hasNewCmd = true
			log.Printf("conflict message:%v newCmd:%v", messageName, i)
			break
		}
		if !hasNewCmd {
			log.Printf("conflictErr message:%v", messageName)
		}
	}
	fileData, err := json.Marshal(mapping)
	if err != nil {
		log.Printf("generateCommandMapping json.Marshal Err fileName:%v err:%v", outputFile, err)
		return err
	}
	err = os.WriteFile(outputFile, fileData, 0644)
	if err != nil {
		log.Printf("generateCommandMapping os.WriteFile Err fileName:%v err:%v", outputFile, err)
		return err
	}
	return nil
}

func loadCommandMapping(fileName string) map[string]int {
	mapping := make(map[string]int)
	fileData, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return mapping
		}
		log.Printf("loadCommandMapping os.ReadFile Err fileName:%v err:%v", fileName, err)
		return mapping
	}
	err = json.Unmarshal(fileData, &mapping)
	if err != nil {
		log.Printf("loadCommandMapping json.Unmarshal Err fileName:%v err:%v", fileName, err)
		return mapping
	}
	return mapping
}
