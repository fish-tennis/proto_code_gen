package main

import (
	"flag"
	"log"
)

func main() {
	var codeTemplatesConfigFile,inputFiles string
	flag.StringVar(&inputFiles, "input", "", "pattern to match input *.pb.go file(s)")
	flag.StringVar(&codeTemplatesConfigFile, "config", "", "code template config")
	flag.Parse()

	if len(inputFiles) == 0 {
		log.Fatal("no input file")
	}
	ParseFiles(inputFiles, codeTemplatesConfigFile)
}
