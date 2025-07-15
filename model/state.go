package model

import (

	tea "github.com/charmbracelet/bubbletea"
)
type ModelType int

const (
	MainModel ModelType = iota
	FormModel
	PomodoroModel
)

var Models = make([]tea.Model, 3)

func InitModels() error {
	if err := InitDB(); err != nil {
		return err
	}
	Models[MainModel] = NewBoard()
	Models[FormModel] = NewFormModel(Todo)
	Models[PomodoroModel] = NewPomodoroModel()
	return nil
}
