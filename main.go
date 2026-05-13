package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

func normalizeProtoNames(names []string) []string {
	result := make([]string, 0, len(names))
	for _, n := range names {
		n = strings.TrimSpace(n)
		if n == "" {
			continue
		}
		if !strings.HasSuffix(n, ".proto") {
			n = n + ".proto"
		}
		result = append(result, n)
	}
	return result
}

func main() {
	var protoDir, mappingFile, templateFile, codeOutFile string
	var excludeStr, outputExcludeStr string

	flag.StringVar(&protoDir, "p", "", "proto file directory (required)")
	flag.StringVar(&mappingFile, "m", "", "output CommandMapping file name")
	flag.StringVar(&templateFile, "t", "", "code generation template file")
	flag.StringVar(&codeOutFile, "c", "", "code generation output file name")
	flag.StringVar(&excludeStr, "e", "", "exclude proto files when parsing (comma separated)")
	flag.StringVar(&outputExcludeStr, "oe", "", "exclude proto files when generating output (comma separated)")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if protoDir == "" {
		log.Fatal("-p is required: proto file directory")
	}

	hasMapping := mappingFile != ""
	hasCodeGen := templateFile != "" && codeOutFile != ""
	if !hasMapping && !hasCodeGen {
		log.Fatal("at least one of -m or -t/-c must be provided")
	}

	var excludeFiles []string
	if excludeStr != "" {
		excludeFiles = normalizeProtoNames(strings.Split(excludeStr, ","))
	}

	var outputExcludeFiles []string
	if outputExcludeStr != "" {
		outputExcludeFiles = normalizeProtoNames(strings.Split(outputExcludeStr, ","))
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Getwd error")
	}
	log.Printf("cwd:%v", cwd)

	parserResult := ParseFiles(protoDir, excludeFiles)

	if hasMapping {
		err := generateCommandMapping(parserResult, mappingFile, excludeFiles, outputExcludeFiles)
		if err != nil {
			log.Fatalf("generateCommandMappingErr:%v", err)
		}
	}
	if hasCodeGen {
		generateCodes(parserResult, templateFile, codeOutFile, outputExcludeFiles)
	}
}
