package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	var codeTemplatesConfigFile, inputFiles string
	flag.StringVar(&inputFiles, "input", "", "pattern to match input *.pb.go file(s)")
	flag.StringVar(&codeTemplatesConfigFile, "config", "proto_code_gen.yaml", "code template config")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if len(inputFiles) == 0 {
		log.Fatal("no input file")
	}
	if len(codeTemplatesConfigFile) == 0 {
		log.Fatal("no config file")
	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Getwd error")
	}
	log.Printf("cwd:%v", cwd)
	ParseFiles(inputFiles, codeTemplatesConfigFile)
}
