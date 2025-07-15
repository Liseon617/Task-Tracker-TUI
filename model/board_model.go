package model

import (
	"kanban/ui"
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Board struct {
	focused Status
	lists   []list.Model
	loaded  bool
	err     error
}

type errMsg struct {
	error
}

func NewBoard() *Board {
	return &Board{}
}

func (m Board) Init() tea.Cmd { return nil }

func (m *Board) initLists(width, height int) {
	defaultList := list.New(nil, list.NewDefaultDelegate(), width/ui.Divisor, height/2)
	defaultList.SetShowHelp(false)

	m.lists = []list.Model{defaultList, defaultList, defaultList}
	m.lists[Todo].Title = "To Do"
	m.lists[InProgress].Title = "In Progress"
	m.lists[Done].Title = "Done"

	tasks, err := LoadTasks()
	if err != nil {
		m.err = err
		log.Printf("Error loading tasks: %v", err)
		return
	}

	todoItems := []list.Item{}
	inProgressItems := []list.Item{}
	doneItems := []list.Item{}

	for _, task := range tasks {
		switch task.status {
		case Todo:
			todoItems = append(todoItems, task)
		case InProgress:
			inProgressItems = append(inProgressItems, task)
		case Done:
			doneItems = append(doneItems, task)
		}
	}
	// m.lists[Todo].SetItems([]list.Item{
	// 	NewTask(Todo, "Task 1", "Description for Task 1"),
	// 	NewTask(Todo, "Task 2", "Description for Task 2"),
	// })
	// m.lists[InProgress].SetItems([]list.Item{
	// 	NewTask(InProgress, "Task 3", "Description for Task 3"),
	// })
	// m.lists[Done].SetItems([]list.Item{
	// 	NewTask(Done, "Task 4", "Description for Task 4"),
	// })
	m.lists[Todo].SetItems(todoItems)
	m.lists[InProgress].SetItems(inProgressItems)
	m.lists[Done].SetItems(doneItems)
}

func (m Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			ui.SetDimensions(msg.Width, msg.Height)
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		case "left", "a":
			m.focused = (m.focused + 2) % 3
		case "right", "d":
			m.focused = (m.focused + 1) % 3
		case "j":
			m.moveToPrev()
		case "l":
			m.moveToNext()
		case "n":
			Models[MainModel] = m
			Models[FormModel] = NewFormModel(m.focused)
			return Models[FormModel].Update(nil)
		case "delete", "backspace":
			m.deleteTask()
		case "c":
			return m, m.ClearBoard()
		case "t":
			Models[MainModel] = m
    		return Models[PomodoroModel], nil
		}
	case Task:
		task := msg
		if err := SaveTask(task); err != nil {
			m.err = err
			log.Printf("Error saving task: %v", err)
			return m, nil
		}
		return m, m.lists[task.status].InsertItem(len(m.lists[task.status].Items()), task)
	case errMsg:
		m.err = msg.error
		log.Printf("Error: %v", msg.error)
		return m, nil
	}
	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m Board) View() string {
	if !m.loaded {
		return "Loading...\n"
	}

	view := lipgloss.JoinHorizontal(
		lipgloss.Left,
		ui.RenderColumn(m.lists[Todo], m.focused == Todo),
		ui.RenderColumn(m.lists[InProgress], m.focused == InProgress),
		ui.RenderColumn(m.lists[Done], m.focused == Done),
	)

	if m.err != nil {
		errorView := lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Padding(1).
			Render("Error: " + m.err.Error())
		return lipgloss.JoinVertical(lipgloss.Left, view, errorView, ui.RenderKanbanHelp())
	}
	return lipgloss.JoinVertical(lipgloss.Left, view, ui.RenderKanbanHelp())
}

func (m *Board) moveToNext() {
	if len(m.lists[m.focused].Items()) == 0 || m.focused == Done {
		return
	}

	task := m.lists[m.focused].SelectedItem().(Task)
	// First delete the old task from database
	if err := DeleteTask(task); err != nil {
        m.err = err
        log.Printf("Error deleting task: %v", err)
        return
    }
	
	// remove from current list
	m.lists[task.status].RemoveItem(m.lists[m.focused].Index())

	// update status and save new version
	task.Next()
	if err := SaveTask(task); err != nil {
		m.err = err
		log.Printf("Error saving task: %v", err)
		return
	}

	// add to new list
	m.lists[task.status].InsertItem(len(m.lists[task.status].Items()), task)
}

func (m *Board) moveToPrev() {
	if len(m.lists[m.focused].Items()) == 0 || m.focused == Todo {
		return
	}
	task := m.lists[m.focused].SelectedItem().(Task)
	if err := DeleteTask(task); err != nil {
        m.err = err
        log.Printf("Error deleting task: %v", err)
        return
    }

	m.lists[task.status].RemoveItem(m.lists[m.focused].Index())
	
	task.Prev()
	if err := SaveTask(task); err != nil {
		m.err = err
		log.Printf("Error saving task: %v", err)
		return
	}
	m.lists[task.status].InsertItem(len(m.lists[task.status].Items()), task)
}

func (m *Board) deleteTask() {
	if item := m.lists[m.focused].SelectedItem(); item != nil {
		task := item.(Task)
		if err := DeleteTask(task); err != nil {
			m.err = err
			log.Printf("Error deleting task: %v", err)
			return
		}
		m.lists[m.focused].RemoveItem(m.lists[m.focused].Index())
	}
}

func (m *Board) ClearBoard() tea.Cmd {
	return func() tea.Msg {
		if err := ClearAllTasks(); err != nil {
			return errMsg{err}
		}
		// Clear UI lists
		for i := range m.lists {
			m.lists[i].SetItems([]list.Item{})
		}
		return nil
	}
}
