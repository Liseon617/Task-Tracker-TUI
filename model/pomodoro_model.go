package model

import (
	"kanban/ui"
	"fmt"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/textinput"
)

type PomodoroState int

const (
	PomodoroStopped PomodoroState = iota
	PomodoroRunning
	PomodoroPaused
	PomodoroEditing
)

type Pomodoro struct {
	state         PomodoroState
	workDuration  time.Duration
	breakDuration time.Duration
	timeLeft      time.Duration
	lastTick      time.Time
	workInput     textinput.Model
	breakInput    textinput.Model
	currentInput  *textinput.Model
	animationPos  int
}

func NewPomodoroModel() *Pomodoro {
	p := &Pomodoro{
		workDuration:  25 * time.Minute,
		breakDuration: 5 * time.Minute,
		timeLeft:      25 * time.Minute,
	}

	// Initialize work duration input
	workInput := textinput.New()
	workInput.Placeholder = "25"
	workInput.CharLimit = 3
	workInput.Width = 5
	workInput.Focus()

	// Initialize break duration input
	breakInput := textinput.New()
	breakInput.Placeholder = "5"
	breakInput.CharLimit = 3
	breakInput.Width = 5

	p.workInput = workInput
	p.breakInput = breakInput
	p.currentInput = &p.workInput

	return p
}

func (m Pomodoro) Init() tea.Cmd {
	return nil
}

func (m *Pomodoro) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "p":
			if m.state == PomodoroEditing {
				return m, nil
			}
			return m.togglePause()
		case "s":
			if m.state == PomodoroEditing {
				return m.saveDurations()
			}
			return m.start()
		case "x":
			if m.state == PomodoroEditing {
				return m.cancelEdit()
			}
			return m.stop()
		case "e":
			if m.state != PomodoroEditing {
				return m.editDurations()
			}
		case "ctrl+c", "esc", "q":
			if m.state == PomodoroEditing {
				return m.cancelEdit()
			}
			return Models[MainModel], nil
		case "tab":
			if m.state == PomodoroEditing {
				if m.currentInput == &m.workInput {
					m.currentInput = &m.breakInput
					m.breakInput.Focus()
					m.workInput.Blur()
				} else {
					m.currentInput = &m.workInput
					m.workInput.Focus()
					m.breakInput.Blur()
				}
				return m, nil
			}
		}
	case tickMsg:
		if m.state == PomodoroRunning {
			m.animationPos = (m.animationPos + 1) % 8
			return m.handleTick()
		}
	}

	if m.state == PomodoroEditing {
		var cmd tea.Cmd
		*m.currentInput, cmd = m.currentInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Pomodoro) View() string {
	minutes := int(m.timeLeft.Minutes())
	seconds := int(m.timeLeft.Seconds()) % 60
	stateStr := ""
	animationChar := ""

	switch m.state {
	case PomodoroStopped:
		stateStr = "Stopped"
	case PomodoroRunning:
		stateStr = "Working"
		animationChar = m.getAnimationChar()
	case PomodoroPaused:
		stateStr = "Paused"
	case PomodoroEditing:
		return lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					lipgloss.NewStyle().Bold(true).Render("Edit Pomodoro Timer"),
					"Work Duration (min): "+m.workInput.View(),
					"Break Duration (min): "+m.breakInput.View(),
					"\nControls: s=save, esc=cancel, tab=switch field",
				),
			)
	}

	progressBar := m.renderProgressBar()

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Render("Pomodoro Timer "+animationChar),
				fmt.Sprintf("%02d:%02d", minutes, seconds),
				progressBar,
				stateStr,
				ui.RenderPomodoroHelp(),
			),
		)
}

func (m *Pomodoro) getAnimationChar() string {
	animationChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	return animationChars[m.animationPos]
}

func (m *Pomodoro) renderProgressBar() string {
	totalDuration := m.workDuration
	if m.state == PomodoroRunning && m.timeLeft > m.workDuration {
		totalDuration = m.breakDuration
	}

	progress := 1.0 - float64(m.timeLeft)/float64(totalDuration)
	width := 20
	completed := int(progress * float64(width))

	bar := ""
	for i := 0; i < width; i++ {
		if i < completed {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	return bar
}

type tickMsg time.Time

func (m *Pomodoro) start() (tea.Model, tea.Cmd) {
	m.state = PomodoroRunning
	m.lastTick = time.Now()
	m.timeLeft = m.workDuration
	m.animationPos = 0
	return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *Pomodoro) togglePause() (tea.Model, tea.Cmd) {
	switch m.state {
	case PomodoroRunning:
		m.state = PomodoroPaused
	case PomodoroPaused:
		m.state = PomodoroRunning
		m.lastTick = time.Now()
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}
	return m, nil
}

func (m *Pomodoro) stop() (tea.Model, tea.Cmd) {
	m.state = PomodoroStopped
	m.timeLeft = m.workDuration
	return m, nil
}

func (m *Pomodoro) editDurations() (tea.Model, tea.Cmd) {
	m.state = PomodoroEditing
	m.workInput.SetValue(fmt.Sprintf("%d", int(m.workDuration.Minutes())))
	m.breakInput.SetValue(fmt.Sprintf("%d", int(m.breakDuration.Minutes())))
	m.workInput.Focus()
	m.currentInput = &m.workInput
	return m, nil
}

func (m *Pomodoro) saveDurations() (tea.Model, tea.Cmd) {
	workMinutes, err := strconv.Atoi(m.workInput.Value())
	if err != nil || workMinutes <= 0 {
		workMinutes = 25
	}

	breakMinutes, err := strconv.Atoi(m.breakInput.Value())
	if err != nil || breakMinutes <= 0 {
		breakMinutes = 5
	}

	m.workDuration = time.Duration(workMinutes) * time.Minute
	m.breakDuration = time.Duration(breakMinutes) * time.Minute
	m.timeLeft = m.workDuration
	m.state = PomodoroStopped
	return m, nil
}

func (m *Pomodoro) cancelEdit() (tea.Model, tea.Cmd) {
	m.state = PomodoroStopped
	return m, nil
}

func (m *Pomodoro) handleTick() (tea.Model, tea.Cmd) {
	now := time.Now()
	elapsed := now.Sub(m.lastTick)
	m.lastTick = now

	if m.timeLeft <= elapsed {
		m.timeLeft = 0
		m.state = PomodoroStopped
		// Switch between work and break durations
		if m.timeLeft == 0 && m.workDuration == m.timeLeft {
			m.timeLeft = m.breakDuration
			return m.start()
		}
		return m, nil
	}

	m.timeLeft -= elapsed
	return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}