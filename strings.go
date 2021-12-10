package main

import "strings"

// 转换成首字母大写的单词
// test -> Test
// TEST -> Test
func ToCamelCase(singleWord string) (camelCaseString string) {
	camelCaseString = strings.ToUpper(singleWord[:1]) + strings.ToLower(singleWord[1:])
	return
}

// TEST_MESSAGE -> Test_Message or TestMessage
func UpperWordsToCamelCase(upperStr,sep string, removeSep bool ) string {
	words := strings.Split(upperStr, sep)
	for i := 0; i < len(words); i++ {
		words[i] = ToCamelCase(words[i])
	}
	if removeSep {
		return strings.Join(words, "")
	} else {
		return strings.Join(words, sep)
	}
}

// TestMessage -> TEST_MESSAGE
func CamelCaseToUpperWords(str,sep string) string {
	words := SplitWords(str)
	for i := 0; i < len(words); i++ {
		words[i] = strings.ToUpper(words[i])
	}
	return strings.Join(words, sep)
}

// 根据大小写分解单词
// TestMessage -> {Test,Message}
func SplitWords(str string) (words []string) {
	beginIdx := 0
	for i := 1; i < len(str); i++ {
		c := string(str[i])
		if strings.ToUpper(c) == c {
			word := str[beginIdx:i]
			words = append(words, word)
			beginIdx = i
		}
	}
	word := str[beginIdx:]
	if word != "" {
		words = append(words, word)
	}
	return
}