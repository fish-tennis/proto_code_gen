package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/bufbuild/protocompile"
	"google.golang.org/protobuf/reflect/protoreflect"
	"gopkg.in/yaml.v3"
)

type ProtoMessageStructInfo struct {
	protoName   string
	MessageName string
	Comment     string
	PackageName string
}

type CodeTemplate struct {
	Template string `yaml:"Template"`
	OutDir   string `yaml:"OutDir"`
}

type CommandMapping struct {
	OutFile string `yaml:"OutFile"`
}

type Configs struct {
	ProtoCodes     *CodeTemplate   `yaml:"ProtoCodes"`
	CommandMapping *CommandMapping `yaml:"CommandMapping"`
}

type ParserResult struct {
	allProto map[string][]*ProtoMessageStructInfo
}

func ParseFiles(protoFilePattern string, codeTemplatesConfig string) {
	configs := loadConfig(codeTemplatesConfig)

	files, err := filepath.Glob(protoFilePattern)
	if err != nil {
		log.Fatal(err)
	}
	if len(files) == 0 {
		log.Fatal("no proto files matched")
	}

	importPath := findCommonParent(files)

	compiler := protocompile.Compiler{
		Resolver: &protocompile.SourceResolver{
			ImportPaths: []string{importPath},
		},
		SourceInfoMode: protocompile.SourceInfoStandard,
	}

	var relFiles []string
	for _, f := range files {
		rel, err := filepath.Rel(importPath, f)
		if err != nil {
			log.Fatal(err)
		}
		relFiles = append(relFiles, filepath.ToSlash(rel))
	}

	fds, err := compiler.Compile(context.Background(), relFiles...)
	if err != nil {
		log.Fatalf("compile proto files failed: %v", err)
	}
	log.Printf("proto file count:%v", len(fds))

	parserResult := &ParserResult{
		allProto: map[string][]*ProtoMessageStructInfo{},
	}

	for _, fd := range fds {
		protoName := path.Base(fd.Path())
		msgs := fd.Messages()
		for i := 0; i < msgs.Len(); i++ {
			md := msgs.Get(i)
			if md.IsMapEntry() {
				continue
			}
			structInfo := extractStructInfo(protoName, md)
			allList := parserResult.allProto[protoName]
			parserResult.allProto[protoName] = append(allList, structInfo)
		}
	}

	if configs.CommandMapping != nil {
		err := generateCommandMapping(parserResult, configs.CommandMapping.OutFile)
		if err != nil {
			log.Fatal(fmt.Sprintf("generateCommandMappingErr:%v", err))
			return
		}
	}
	if configs.ProtoCodes != nil {
		generateCodes(parserResult, configs.ProtoCodes)
	}
}

func findCommonParent(files []string) string {
	if len(files) == 0 {
		return "."
	}
	dir := filepath.Dir(files[0])
	for _, f := range files[1:] {
		d := filepath.Dir(f)
		for !strings.HasPrefix(filepath.ToSlash(d)+"/", filepath.ToSlash(dir)+"/") && dir != "." && dir != "" {
			dir = filepath.Dir(dir)
		}
	}
	if dir == "" {
		return "."
	}
	return dir
}

func extractStructInfo(protoName string, md protoreflect.MessageDescriptor) *ProtoMessageStructInfo {
	messageName := string(md.Name())

	comment := extractMessageComment(md)

	structInfo := &ProtoMessageStructInfo{
		protoName:   protoName,
		MessageName: messageName,
		Comment:     comment,
		PackageName: string(md.ParentFile().Package()),
	}

	return structInfo
}

func extractMessageComment(md protoreflect.MessageDescriptor) string {
	fd := md.ParentFile()
	loc := fd.SourceLocations().ByDescriptor(md)
	if loc.LeadingComments == "" {
		return ""
	}
	return strings.TrimSpace(loc.LeadingComments)
}

func loadConfig(config string) *Configs {
	fileData, err := os.ReadFile(config)
	if err != nil {
		panic("read config file err")
	}
	configs := &Configs{
		ProtoCodes:     &CodeTemplate{},
		CommandMapping: &CommandMapping{},
	}
	err = yaml.Unmarshal(fileData, configs)
	if err != nil {
		panic(err)
	}
	//if configs.Reader != nil {
	//	autoCheckDir(&configs.Reader.OutDir)
	//}
	if configs.ProtoCodes != nil {
		autoCheckDir(&configs.ProtoCodes.OutDir)
	}
	return configs
}

func autoCheckDir(dir *string) {
	if *dir == "" {
		return
	}
	if !strings.HasSuffix(*dir, "/") && !strings.HasSuffix(*dir, "\\") {
		*dir = *dir + "/"
	}
}
