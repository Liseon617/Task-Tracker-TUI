package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"kanban/model"
)

func main() {
	if err := model.InitModels(); err != nil {
		fmt.Printf("Failed to initialize: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(model.Models[model.MainModel])
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
