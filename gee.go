package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

const (
	configFilename = ".gee"
)

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

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error retrieving home directory.", "error", err)
	}

	configPath := filepath.Join(homeDir, configFilename)
	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatal("Error loading config.", "error", err)
	}

	log.Debug(config)
	log.Debug("", "args", readCommandLineArgs(getCommandLineArgs()))

}

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
	fmt.Printf("Config file written to %s\n", file)

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
