package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	Divisor      = 4
	columnStyle  = lipgloss.NewStyle().Padding(1, 2)
	focusedStyle = lipgloss.NewStyle().Padding(0, 1).
			Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62"))
)

func SetDimensions(width, height int) {
	columnStyle.Width(width / Divisor)
	focusedStyle.Width(width / Divisor)
	columnStyle.Height(height - Divisor)
	focusedStyle.Height(height - Divisor)
}

func RenderColumn(l list.Model, focused bool) string {
	if focused {
		return focusedStyle.Render(l.View())
	}
	return columnStyle.Render(l.View())
}
