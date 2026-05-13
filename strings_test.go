package main

import (
	"path"
	"regexp"
	"testing"
)

func TestStrings(t *testing.T) {
	println(ToCamelCase("test"))
	println(ToCamelCase("TEST"))
	println(UpperWordsToCamelCase("TEST_MESSAGE_XXX", "_", true))
	println(UpperWordsToCamelCase("TEST_MESSAGE_XXX", "_", false))
	println(SplitWords("TestMessageXxx"))
	println(CamelCaseToUpperWords("TestMessageXxx", "_"))
}

func TestRegexp(t *testing.T) {
	ok, err := regexp.MatchString(".*.pb.go", "cfg.pb.go")
	t.Logf("ok:%v, err:%v", ok, err)

	ok, err = regexp.MatchString(".*.pb.go", "cfg.pb")
	t.Logf("ok:%v, err:%v", ok, err)

	ok, err = regexp.MatchString(".*.pb.go", "pb.go")
	t.Logf("ok:%v, err:%v", ok, err)
}

func TestProtoName(t *testing.T) {
	baseName := path.Base("path/to/file.proto")
	t.Logf("baseName:%v", baseName)
}
