package main

import (
	"encoding/json"
	"hash/crc32"
	"log"
	"os"
)

func resolveCommandMapping(allMessageNames []string, existingMapping map[string]int) map[string]int {
	cmdMapping := make(map[int]string)
	for messageName, cmd := range existingMapping {
		if _, ok := cmdMapping[cmd]; ok {
			log.Printf("conflict message:%v cmd:%v", messageName, cmd)
			delete(existingMapping, messageName)
			continue
		}
		cmdMapping[cmd] = messageName
	}
	conflict := make(map[string]int)
	allMessages := make(map[string]int)
	for _, messageName := range allMessageNames {
		cmd := uint16(crc32.ChecksumIEEE([]byte(messageName)) & 0xFFFF)
		allMessages[messageName] = int(cmd)
	}
	for messageName, cmd := range existingMapping {
		if _, ok := allMessages[messageName]; !ok {
			delete(existingMapping, messageName)
			delete(cmdMapping, cmd)
		}
	}
	for messageName, cmd := range allMessages {
		oldCmd, ok := existingMapping[messageName]
		if ok && cmd != oldCmd {
			conflict[messageName] = cmd
			continue
		}
		if cmd == 0 {
			conflict[messageName] = cmd
			continue
		}
		existingMapping[messageName] = cmd
		cmdMapping[cmd] = messageName
	}
	for messageName, cmd := range conflict {
		log.Printf("conflict message:%v cmd:%v", messageName, cmd)
		hasNewCmd := false
		for i := 1; i < 0xFFFF; i++ {
			if _, ok := cmdMapping[i]; ok {
				continue
			}
			cmdMapping[i] = messageName
			existingMapping[messageName] = i
			hasNewCmd = true
			log.Printf("conflict message:%v newCmd:%v", messageName, i)
			break
		}
		if !hasNewCmd {
			log.Printf("conflictErr message:%v", messageName)
		}
	}
	return existingMapping
}

func saveCommandMapping(mapping map[string]int, outputFile string) error {
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

func generateCommandMapping(parserResult *ParserResult, outputFile string, excludeFiles []string, outputExcludeFiles []string) error {
	if outputFile == "" {
		return nil
	}
	mapping := loadCommandMapping(outputFile)
	excludeSet := make(map[string]bool)
	for _, f := range excludeFiles {
		excludeSet[f] = true
	}
	var allMessageNames []string
	for protoName, structInfoList := range parserResult.allProto {
		if excludeSet[protoName] {
			continue
		}
		for _, structInfo := range structInfoList {
			allMessageNames = append(allMessageNames, structInfo.MessageName)
		}
	}
	mapping = resolveCommandMapping(allMessageNames, mapping)
	oeSet := make(map[string]bool)
	for _, f := range outputExcludeFiles {
		oeSet[f] = true
	}
	for protoName, structInfoList := range parserResult.allProto {
		if !oeSet[protoName] {
			continue
		}
		for _, structInfo := range structInfoList {
			delete(mapping, structInfo.MessageName)
		}
	}
	return saveCommandMapping(mapping, outputFile)
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
