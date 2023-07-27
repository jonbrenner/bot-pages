package main

import (
	"embed"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

//go:embed prompt_prefix*.txt
var promptPrefixFS embed.FS

const promptStartText = "\n[EXAMPLES]\n"

var err error

func main() {

	args := parseArgs()
	if len(args.Prompt) == 0 {
		usage()
		os.Exit(1)
	}

	var promptPrefix []byte
	switch {
	case args.Mode == "command":
		promptPrefix, err = promptPrefixFS.ReadFile("prompt_prefix_command.txt")
		errHandler(err, "Error reading prompt_prefix_command")
	case args.Mode == "config":
		promptPrefix, err = promptPrefixFS.ReadFile("prompt_prefix_config.txt")
		errHandler(err, "Error reading prompt_prefix_config")
	}

	homeDir, err := os.UserHomeDir()
	errHandler(err, "Error retrieving home directory")

	// TODO: distinguish between config for bot program and config as a command flag
	configPath := filepath.Join(homeDir, configFilename)
	config, err := loadConfig(configPath)
	errHandler(err, "Error loading config")

	respCh := make(chan string)

	var wg sync.WaitGroup
	wg.Add(2)

	client := &OpenAIAdapter{APIKey: config.APIKey}
	go func() {
		defer wg.Done()
		err := client.FetchCompletionStream(CreateRequest(string(promptPrefix), args.Prompt), respCh)
		errHandler(err)
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
	Help        bool
	Interactive bool
	Mode        string
	Prompt      string
}

func parseArgs() *Args {
	help := flag.Bool("help", false, "Display usage options")
	interactive := flag.Bool("i", false, "Enter interactive mode (optional)")
	mode := flag.String("mode", "command", "Query config file details")
	flag.Parse()
	prompt := strings.Join(flag.Args(), " ")
	return &Args{Help: *help, Interactive: *interactive, Mode: *mode, Prompt: prompt}
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func RenderCompletionStreamResponse(w io.Writer, respCh <-chan string) {
	for token := range respCh {
		fmt.Fprint(w, token)
	}
	fmt.Fprintln(w, "")
}

func errHandler(err error, args ...string) {
	errString := "Error"
	if len(args) > 0 {
		errString = args[0]
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v", errString, err)
		os.Exit(1)
	}
}
