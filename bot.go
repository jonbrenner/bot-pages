package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sashabaranov/go-openai"
)

//go:embed prompt_prefix.txt
var promptPrefix string

const promptStartText = "\n[EXAMPLES]\n"

func main() {

	prompt := getUserPrompt(getCommandLineArgs())
	if len(prompt) == 0 {
		usage()
		os.Exit(1)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving home directory: %v", err)
		os.Exit(1)
	}

	configPath := filepath.Join(homeDir, configFilename)
	config, err := loadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v", err)
		os.Exit(1)
	}

	fetchCompletion(config, prompt)
}

// Print usage information
func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [prompt]", os.Args[0])
}

func fetchCompletion(config Config, prompt string) {
	c := openai.NewClient(config.APIKey)
	ctx := context.Background()

	req := openai.CompletionRequest{
		Model:       openai.GPT3TextDavinci003,
		MaxTokens:   1024,
		Prompt:      promptPrefix + prompt + promptStartText,
		Stream:      false,
		Temperature: 0.05,
	}

	stream, err := c.CreateCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("CompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	fmt.Println("")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("")
			return
		}

		if err != nil {
			fmt.Printf("Stream error: %v\n", err)
			return
		}

		fmt.Printf(response.Choices[0].Text)
	}
}
