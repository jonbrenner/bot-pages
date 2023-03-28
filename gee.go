package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/sashabaranov/go-openai"
)

const (
	configFilename = ".gee"
)

//go:embed prompt_prefix.txt
var promptPrefix string

const promptStartText = "\n[EXAMPLES]\n"

type Config struct {
	APIKey string `json:"api-key"`
}

func main() {

	if len(os.Getenv("DEBUG")) > 0 {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
		if f, err := tea.LogToFile("debug.log", "debug"); err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		} else {
			log.SetOutput(f)
			defer f.Close()
		}
		log.Debug("Starting...")
	}

	prompt := readCommandLineArgs(getCommandLineArgs())
	if len(prompt) == 0 {
		usage()
		os.Exit(1)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error retrieving home directory.", "error", err)
	}

	configPath := filepath.Join(homeDir, configFilename)
	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatal("Error loading config.", "error", err)
	}

	log.Debug("", "args", readCommandLineArgs(getCommandLineArgs()))

	p := tea.NewProgram(initialModel(prompt, &config))
	if _, err := p.Run(); err != nil {
		log.Fatal("Error running program.", "error", err)
	}
}

// Print usage information
func usage() {
	fmt.Printf("Usage: %s [prompt]", os.Args[0])
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
	} else {
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
			log.Warn("Warning: File permissions should be 600. Please update the file permissions.")
		}
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
	log.Info("Config file written to %s\n", file)

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

func readCommandLineArgs(args []string) string {
	return strings.Join(args, " ")
}

func getCommandLineArgs() []string {
	return os.Args[1:]
}

type Completion string

type Model struct {
	config     *Config
	Prompt     string
	Completion Completion
}

func initialModel(prompt string, config *Config) Model {
	return Model{Prompt: prompt, config: config}
}

func (m Model) Init() tea.Cmd {
	return m.fetchCompletion
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case Completion:
		m.Completion = msg
	}
	return m, nil
}

func (m Model) View() string {
	return fmt.Sprintf("%s\n", m.Completion)
}

func (m Model) fetchCompletion() tea.Msg {
	c := openai.NewClient(m.config.APIKey)
	ctx := context.Background()
	log.Debug("Creating completion request.", "prompt", promptPrefix+m.Prompt+promptStartText)
	req := openai.CompletionRequest{
		Model:       openai.GPT3TextDavinci003,
		MaxTokens:   1024,
		Prompt:      promptPrefix + m.Prompt + promptStartText,
		Stream:      false,
		Temperature: 0.05,
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		log.Fatal("Error creating completion.", "error", err)
	}

	return Completion(resp.Choices[0].Text)
}
