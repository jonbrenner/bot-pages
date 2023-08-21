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

/*
- The "go:embed" directive must immediately precede a line containing the declaration of a single variable
- Therefore, this program embeds all files in the current folder that match prompt_prefix*.txt into a file system variable "promptPrefixFS"
- Then, depending on the flags passed to the command, the actual "promptPrefix" is defined by reading from "promptPrefixFS"
- This method makes use of the embed directive while allowing the final "promptPrefix" to be conditionally defined
*/

//go:embed prompt_prefix*.txt
var promptPrefixFS embed.FS

const promptStartText = "\n[EXAMPLES]\n"

var err error

func main() {

	args := parseArgs()

	var promptPrefix []byte
	switch args.Config {
	case false:
		promptPrefix, err = promptPrefixFS.ReadFile("prompt_prefix_command.txt")
		errHandler(err, "Error reading prompt_prefix_command")
	case true:
		promptPrefix, err = promptPrefixFS.ReadFile("prompt_prefix_config.txt")
		errHandler(err, "Error reading prompt_prefix_config")
	}

	homeDir, err := os.UserHomeDir()
	errHandler(err, "Error retrieving home directory")

	botConfigPath := filepath.Join(homeDir, configFilename)
	botConfig, err := loadConfig(botConfigPath)
	errHandler(err, "Error loading config")

	respCh := make(chan string)

	var wg sync.WaitGroup
	wg.Add(2)

	client := &OpenAIAdapter{APIKey: botConfig.APIKey}
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
	Config      bool
	Prompt      string
}

func parseArgs() *Args {

	help := flag.Bool("help", false, "Display usage options")
	interactive := flag.Bool("i", false, "Enter interactive mode (coming soon!)")
	config := flag.Bool("config", false, "Optional flag to operate as config bot")
	flag.Parse()
	prompt := strings.Join(flag.Args(), " ")

	if *help || len(os.Args) < 2 || len(prompt) == 0 {
		Usage()
	}

	return &Args{Help: *help, Interactive: *interactive, Config: *config, Prompt: prompt}
}

func Usage() {
	fmt.Println("Simple command-line utility for looking up command usage or config file examples.")
	fmt.Println("Flags:")
	flag.PrintDefaults()
	fmt.Println("Example usage:")
	fmt.Println("  bot nc")
	fmt.Println("  bot -config cron")
	fmt.Println("  bot -help")
	os.Exit(1)
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
