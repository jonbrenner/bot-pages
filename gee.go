package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func main() {

	if len(os.Getenv("DEBUG")) > 0 {
		log.SetLevel(log.DebugLevel)
		if f, err := tea.LogToFile("debug.log", "debug"); err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		} else {
			log.SetOutput(f)
			defer f.Close()
		}
	}
	log.Debug("starting...")
}
