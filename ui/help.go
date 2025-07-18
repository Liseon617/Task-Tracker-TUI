package ui

import "github.com/charmbracelet/lipgloss"

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

func RenderKanbanHelp() string {
	return helpStyle.Render(`
	Controls: 
←/a: Prev board  →/d: Next board   j: Move ←   l: Move →
n: New task	 delete/backspace: Delete task   
c+c Clear tasks	  t: Pomodoro timer	  esc/q/ctrl+c: Quit
`)
}

func RenderPomodoroHelp() string {
	return helpStyle.Render(`
Controls: 
s: Start   p: Unpause/Pause   x: Stop
e: Edit   esc/q: Back   ctrl+c: Quit`,
	)
}
