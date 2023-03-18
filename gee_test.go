package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadConfig(t *testing.T) {
	jsonData := `{"api-key": "test-api-key"}`

	config, err := readConfig(strings.NewReader(jsonData))
	if err != nil {
		t.Errorf("readConfig returned an error: %v", err)
	}

	expectedAPIKey := "test-api-key"
	if config.APIKey != expectedAPIKey {
		t.Errorf("Expected APIKey to be %s, but got %s", expectedAPIKey, config.APIKey)
	}
}

func TestCreateConfigFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testConfig := Config{APIKey: "test-api-key"}
	tempFilePath := filepath.Join(tempDir, ".gee")
	err = createConfigFile(tempFilePath, testConfig)
	if err != nil {
		t.Errorf("createConfigFile returned an error: %v", err)
	}

	// Check if the file exists and has the correct content
	data, err := os.ReadFile(tempFilePath)
	if err != nil {
		t.Errorf("Failed to read temporary file: %v", err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		t.Errorf("Failed to unmarshal config: %v", err)
	}

	if config.APIKey != testConfig.APIKey {
		t.Errorf("Expected APIKey to be %s, but got %s", testConfig.APIKey, config.APIKey)
	}

	// Check if the file has the correct permissions
	fileInfo, err := os.Stat(tempFilePath)
	if err != nil {
		t.Errorf("Failed to get file info: %v", err)
	}

	expectedMode := os.FileMode(0600)
	if fileInfo.Mode().Perm() != expectedMode {
		t.Errorf("Expected file mode to be %v, but got %v", expectedMode, fileInfo.Mode().Perm())
	}
}

func TestReadCommandLineArgs(t *testing.T) {
	testCases := []struct {
		input    []string
		expected string
	}{
		{[]string{"arg1", "arg2", "arg3"}, "arg1 arg2 arg3"},
		{[]string{}, ""},
	}

	// Check if args are concatenated into a single string
	for _, tc := range testCases {
		result := readCommandLineArgs(tc.input)
		if result != tc.expected {
			t.Errorf("Expected '%s', but got '%s' for input: %v", tc.expected, result, tc.input)
		}
	}
}
