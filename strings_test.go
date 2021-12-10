package main

import "testing"

func TestStrings(t *testing.T) {
	println(ToCamelCase("test"))
	println(ToCamelCase("TEST"))
	println(UpperWordsToCamelCase("TEST_MESSAGE_XXX","_", true))
	println(UpperWordsToCamelCase("TEST_MESSAGE_XXX","_", false))
	println(SplitWords("TestMessageXxx"))
	println(CamelCaseToUpperWords("TestMessageXxx", "_"))
}
