package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var customBorder = lipgloss.Border{
	Top:         "\u2550",
	Bottom:      "\u2550",
	Left:        "\u2551",
	Right:       "\u2551",
	TopLeft:     "\u2554",
	TopRight:    "\u2557",
	BottomLeft:  "\u255a",
	BottomRight: "\u255d",
}

var (
	headerStyle = lipgloss.
			NewStyle().
			Border(customBorder)
	borderStyle = lipgloss.NewStyle().
			Border(customBorder)
	focusBorderStyle = lipgloss.NewStyle().
				Border(customBorder).
				BorderForeground(lipgloss.Color("63"))

	centerContentFocus = lipgloss.NewStyle().
				Border(customBorder).
				BorderForeground(lipgloss.Color("63")).
				AlignHorizontal(lipgloss.Center).
				AlignVertical(lipgloss.Center)
	centerContent = lipgloss.NewStyle().
			Border(customBorder).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center)
	centerText = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center)
)

// List styles
const listHeight = 14
const listWidth = 200

var (
	titleStyle = lipgloss.NewStyle().Bold(true)
	itemStyle  = lipgloss.NewStyle().
			PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(lipgloss.Color("255"))
	paginationStyle = list.DefaultStyles().PaginationStyle.
			PaddingLeft(4)
)
