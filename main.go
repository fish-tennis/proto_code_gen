package main

import (
	"flag"
	"log"
)

func main() {
	var codeTemplatesConfigFile,inputFiles string
	flag.StringVar(&inputFiles, "input", "", "pattern to match input *.pb.go file(s)")
	flag.StringVar(&codeTemplatesConfigFile, "config", "code_templates.json", "code template config")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if len(inputFiles) == 0 {
		log.Fatal("no input file")
	}
	if len(codeTemplatesConfigFile) == 0 {
		log.Fatal("no config file")
	}
	ParseFiles(inputFiles, codeTemplatesConfigFile)
}
