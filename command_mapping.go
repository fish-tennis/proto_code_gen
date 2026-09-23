package main

import (
	"encoding/json"
	"hash/crc32"
	"log"
	"os"
	"sort"
)

// resolveCommandMapping 生成消息名->命令字映射
//
// 命令字分配规则(保证幂等:同一输入永远得到同一输出,消除历史值漂移):
//  1. 默认命令字 = crc32(消息名)的低16位(确定性哈希,同名永远同值)
//  2. crc值空闲 → 采用crc值;
//     历史上因碰撞被分配人工值、且碰撞已消除的消息,在此自愈回归crc值
//  3. crc值被其他消息占用 → 已有人工值则保持不变(不再重新分配,消除每轮漂移);
//     新消息则从小到大找空闲值分配
//  4. 采用新值前先释放自身旧值占坑,避免旧值"自我碰撞"把新值越推越大
//  5. 所有遍历按消息名排序进行,与map随机遍历序无关,结果确定
func resolveCommandMapping(allMessageNames []string, existingMapping map[string]int) map[string]int {
	// 第1步: 旧映射反转cmd->messageName;cmd重复时按消息名序保留先到者,剔除后到者
	cmdMapping := make(map[int]string)
	sortedOld := make([]string, 0, len(existingMapping))
	for messageName := range existingMapping {
		sortedOld = append(sortedOld, messageName)
	}
	sort.Strings(sortedOld)
	for _, messageName := range sortedOld {
		cmd := existingMapping[messageName]
		if holder, ok := cmdMapping[cmd]; ok {
			log.Printf("duplicate cmd:%v message:%v holder:%v", cmd, messageName, holder)
			delete(existingMapping, messageName)
			continue
		}
		cmdMapping[cmd] = messageName
	}
	// 第2步: 计算所有消息的crc16
	conflict := make(map[string]int)
	allMessages := make(map[string]int)
	for _, messageName := range allMessageNames {
		cmd := uint16(crc32.ChecksumIEEE([]byte(messageName)) & 0xFFFF)
		allMessages[messageName] = int(cmd)
	}
	// 第3步: 删除proto里已不存在的旧消息,释放其cmd占用
	for messageName, cmd := range existingMapping {
		if _, ok := allMessages[messageName]; !ok {
			delete(existingMapping, messageName)
			delete(cmdMapping, cmd)
		}
	}
	// 第4步: 按消息名排序分配,结果与map遍历序无关
	sortedNames := make([]string, 0, len(allMessages))
	for name := range allMessages {
		sortedNames = append(sortedNames, name)
	}
	sort.Strings(sortedNames)
	for _, messageName := range sortedNames {
		cmd := allMessages[messageName]
		if cmd == 0 {
			conflict[messageName] = cmd
			continue
		}
		if holder, ok := cmdMapping[cmd]; !ok || holder == messageName {
			// crc值空闲(或已被自己持有):采用crc值
			// 旧值是历史碰撞分配的人工值时,先释放旧值占坑再自愈回归crc值
			if oldCmd, ok := existingMapping[messageName]; ok && oldCmd != cmd {
				if cmdMapping[oldCmd] == messageName {
					delete(cmdMapping, oldCmd)
				}
			}
			existingMapping[messageName] = cmd
			cmdMapping[cmd] = messageName
			continue
		}
		// crc值被其他消息占用
		if _, ok := existingMapping[messageName]; ok {
			// 已有人工值:保持不变(第1步已占坑),不再每轮重新分配
			continue
		}
		// 新消息且crc被占:待冲突分配
		conflict[messageName] = cmd
	}
	// 第5步: 冲突消息按消息名排序后从小到大分配空闲cmd(结果确定)
	conflictNames := make([]string, 0, len(conflict))
	for name := range conflict {
		conflictNames = append(conflictNames, name)
	}
	sort.Strings(conflictNames)
	for _, messageName := range conflictNames {
		allocated := false
		for i := 1; i < 0xFFFF; i++ {
			if _, ok := cmdMapping[i]; ok {
				continue
			}
			cmdMapping[i] = messageName
			existingMapping[messageName] = i
			allocated = true
			log.Printf("conflict message:%v crcCmd:%v newCmd:%v", messageName, conflict[messageName], i)
			break
		}
		if !allocated {
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
