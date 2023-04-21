package main

import (
	"bytes"
	"os"
	"reflect"
	"sync"
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

func TestRenderCompletionStream(t *testing.T) {
	tokens := []string{"this", "is", "a", "test"}
	expected := "thisisatest"
	output := &bytes.Buffer{}

	respCh := make(chan string)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer close(respCh)
		for _, token := range tokens {
			respCh <- token
		}
	}()

	go func() {
		defer wg.Done()
		RenderCompletionStreamResponse(output, respCh)
	}()

	wg.Wait()

	if output.String() != expected {
		t.Errorf("Expected %q, but got %q", expected, output.String())
	}
}
