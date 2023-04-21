package main

import (
	"bytes"
	"flag"
	"os"
	"sync"
	"testing"
)

func TestParseArgs(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected Args
	}{
		{
			name: "no flags",
			args: []string{"arg1", "arg2", "arg3"},
			expected: Args{
				Interactive: false,
				Prompt:      "arg1 arg2 arg3",
			},
		},
		{
			name: "interactive mode with args",
			args: []string{"-i", "arg1", "arg2", "arg3"},
			expected: Args{
				Interactive: true,
				Prompt:      "arg1 arg2 arg3",
			},
		},
		{
			name: "no flags and no args",
			args: []string{},
			expected: Args{
				Interactive: false,
				Prompt:      "",
			},
		},
		{
			name: "interactive mode without args",
			args: []string{"-i"},
			expected: Args{
				Interactive: true,
				Prompt:      "",
			},
		},
		{
			name: "-i is part of prompt when not in the first position",
			args: []string{"arg1", "-i", "arg2"},
			expected: Args{
				Interactive: false,
				Prompt:      "arg1 -i arg2",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Args = append([]string{"main"}, tc.args...)
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			args := parseArgs()
			if args.Interactive != tc.expected.Interactive || args.Prompt != tc.expected.Prompt {
				t.Errorf("Expected %+v, got %+v", tc.expected, args)
			}
		})
	}
}

func TestRenderCompletionStream(t *testing.T) {
	tokens := []string{"this", "is", "a", "test"}
	expected := "thisisatest\n"
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
