package main

import (
	"context"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/bufbuild/protocompile"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ProtoMessageStructInfo struct {
	protoName   string
	MessageName string
	Comment     string
	PackageName string
}

type ParserResult struct {
	allProto map[string][]*ProtoMessageStructInfo
}

func ParseFiles(protoDir string, excludeFiles []string) *ParserResult {
	files, err := filepath.Glob(filepath.Join(protoDir, "*.proto"))
	if err != nil {
		log.Fatal(err)
	}
	if len(files) == 0 {
		log.Fatal("no proto files found in directory: " + protoDir)
	}

	excludeSet := make(map[string]bool)
	for _, e := range excludeFiles {
		excludeSet[e] = true
	}

	var filteredFiles []string
	for _, f := range files {
		if excludeSet[filepath.Base(f)] {
			log.Printf("excluding proto file: %s", f)
			continue
		}
		filteredFiles = append(filteredFiles, f)
	}
	if len(filteredFiles) == 0 {
		log.Fatal("no proto files to parse after exclusion")
	}

	importPath := findCommonParent(filteredFiles)

	compiler := protocompile.Compiler{
		Resolver: &protocompile.SourceResolver{
			ImportPaths: []string{importPath},
		},
		SourceInfoMode: protocompile.SourceInfoStandard,
	}

	var relFiles []string
	for _, f := range filteredFiles {
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

	return parserResult
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
