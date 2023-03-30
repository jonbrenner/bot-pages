package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const (
	configFilename = ".botpages"
)

//go:embed prompt_prefix.txt
var promptPrefix string

const promptStartText = "\n[EXAMPLES]\n"

type Config struct {
	APIKey string `json:"api-key"`
}

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

// loadConfig loads the config file from the user's home directory and returns
// the config. If the config file does not exist, it will be created. If the
// file permissions are not 600, a warning will be printed.
func loadConfig(file string) (Config, error) {
	var config Config

	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		fmt.Print("Enter your API Key: ")
		fmt.Scan(&config.APIKey)

		err := createConfigFile(file, config)
		if err != nil {
			return config, fmt.Errorf("error creating config file: %w", err)
		}
	} else if err != nil {
		return config, fmt.Errorf("error checking config file: %w", err)
	}

	config, err = readConfigFromFile(file)
	if err != nil {
		return config, fmt.Errorf("error reading config file: %w", err)
	}

	fileInfo, err := os.Stat(file)
	if err != nil {
		return config, fmt.Errorf("error getting file info: %w", err)
	}

	fileMode := fileInfo.Mode().Perm()
	if fileMode != 0600 {
		fmt.Fprintf(os.Stderr, "Warning: File permissions should be 600. Please update the file permissions.")
	}

	if err := validateConfig(config); err != nil {
		return config, fmt.Errorf("error validating config: %w", err)
	}

	return config, nil
}

func createConfigFile(file string, config Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(file, data, 0600)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Config file written to %s\n", file)

	return nil
}

func readConfigFromFile(filePath string) (Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Config{}, err
	}

	config, err := readConfig(file)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func readConfig(reader io.Reader) (Config, error) {
	var config Config

	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func getUserPrompt(args []string) string {
	return strings.Join(args, " ")
}

func getCommandLineArgs() []string {
	return os.Args[1:]
}

func validateConfig(config Config) error {
	if len(config.APIKey) == 0 {
		return fmt.Errorf("API Key is not set")
	}

	return nil
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

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("Stream finished")
			return
		}

		if err != nil {
			fmt.Printf("Stream error: %v\n", err)
			return
		}

		fmt.Printf(response.Choices[0].Text)
	}
}
