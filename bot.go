package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

//go:embed prompt_prefix.txt
var promptPrefix string

const promptStartText = "\n[EXAMPLES]\n"

func main() {

	args := parseArgs()
	if len(args.Prompt) == 0 {
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
		err := client.FetchCompletionStream(CreateRequest(args.Prompt), respCh)
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

type Args struct {
	Interactive bool
	Prompt      string
}

func parseArgs() *Args {
	interactive := flag.Bool("i", false, "Enter interactive mode (optional)")
	flag.Parse()
	prompt := strings.Join(flag.Args(), " ")
	return &Args{Interactive: *interactive, Prompt: prompt}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [prompt]", os.Args[0])
}

func RenderCompletionStreamResponse(w io.Writer, respCh <-chan string) {
	for token := range respCh {
		fmt.Fprint(w, token)
	}
	fmt.Fprintln(w, "")
}
