package model

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Form struct {
	focused     Status
	title       textinput.Model
	description textarea.Model
}

func NewFormModel(focused Status) *Form {
	form := &Form{focused: focused}
	form.title = textinput.New()
	form.title.Focus()
	form.description = textarea.New()
	return form
}

func (m Form) Init() tea.Cmd { return nil }

func (m Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if m.title.Focused() {
				m.title.Blur()
				m.description.Focus()
				return m, textarea.Blink
			}
			Models[FormModel] = m
			return Models[MainModel], m.createTask()
		}
	}
	var cmd tea.Cmd
	if m.title.Focused() {
		m.title, cmd = m.title.Update(msg)
	} else {
		m.description, cmd = m.description.Update(msg)
	}
	return m, cmd
}

func (m Form) View() string {
	return m.title.View() + "\n" + m.description.View()
}

func (m Form) createTask() tea.Cmd {
	return func() tea.Msg {
		return NewTask(m.focused, m.title.Value(), m.description.Value())
	}
}
