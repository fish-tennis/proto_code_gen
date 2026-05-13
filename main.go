package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	var codeTemplatesConfigFile, protoFiles string
	flag.StringVar(&protoFiles, "proto", "", "pattern to match input *.proto file(s)")
	flag.StringVar(&codeTemplatesConfigFile, "config", "proto_code_gen.yaml", "code template config")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if len(codeTemplatesConfigFile) == 0 {
		log.Fatal("no config file")
	}
	if len(protoFiles) == 0 {
		log.Fatal("no proto file pattern, use -proto")
	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Getwd error")
	}
	log.Printf("cwd:%v", cwd)
	ParseFiles(protoFiles, codeTemplatesConfigFile)
}
