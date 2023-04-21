package main

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

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

	respCh := make(chan string)

	var wg sync.WaitGroup
	wg.Add(2)

	client := &OpenAIAdapter{APIKey: config.APIKey}
	go func() {
		defer wg.Done()
		err := client.FetchCompletionStream(CreateRequest(prompt), respCh)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(1)
		}
	}()

	go func() {
		defer wg.Done()
		RenderCompletionStreamResponse(os.Stdout, respCh)
	}()

	wg.Wait()
}

// Print usage information
func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [prompt]", os.Args[0])
}

func CreateRequest(prompt string) openai.CompletionRequest {
	return openai.CompletionRequest{
		Model:       openai.GPT3TextDavinci003,
		MaxTokens:   1024,
		Prompt:      promptPrefix + prompt + promptStartText,
		Stream:      false,
		Temperature: 0.05,
	}
}

func RenderCompletionStreamResponse(w io.Writer, respCh <-chan string) {
	for token := range respCh {
		fmt.Fprint(w, token)
	}
}
