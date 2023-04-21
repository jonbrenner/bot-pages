package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const (
	configFilename = ".bot-pages"
)

type Config struct {
	APIKey string `json:"api-key"`
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

func validateConfig(config Config) error {
	if len(config.APIKey) == 0 {
		return fmt.Errorf("API Key is not set")
	}

	return nil
}
