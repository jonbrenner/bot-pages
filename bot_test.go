package main

import (
	"os"
	"reflect"
	"testing"
)

func TestGetUserPrompt(t *testing.T) {
	testCases := []struct {
		input    []string
		expected string
	}{
		{[]string{"arg1", "arg2", "arg3"}, "arg1 arg2 arg3"},
		{[]string{}, ""},
	}

	// Check if args are concatenated into a single string
	for _, tc := range testCases {
		result := getUserPrompt(tc.input)
		if result != tc.expected {
			t.Errorf("Expected '%s', but got '%s' for input: %v", tc.expected, result, tc.input)
		}
	}
}

func TestGetCommandLineArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "arg1", "arg2", "arg3"}
	if !reflect.DeepEqual(os.Args[1:], getCommandLineArgs()) {
		t.Errorf("Expected %v, but got %v", os.Args[1:], getCommandLineArgs())
	}
}
