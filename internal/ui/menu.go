package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var choices = []string{
	"Connect from ssh config",
}

type menuModel struct {
	width  int
	height int
	cursor int
}

func newMenuModel() *menuModel {
	return &menuModel{}
}

func (m *menuModel) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Up):
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		case key.Matches(msg, keys.Down):
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}
		case key.Matches(msg, keys.Enter):
			return func() tea.Msg {
				return ServerList
			}
		}
	}

	return cmd
}

func (m *menuModel) View() string {
	s := strings.Builder{}
	s.WriteString("Connect to the server:\n\n")

	for i := 0; i < len(choices); i++ {
		if m.cursor == i {
			s.WriteString("🡺 ")
		} else {
			s.WriteString(" ")
		}

		s.WriteString(choices[i])
		s.WriteString("\n")
	}

	return centerContent.
		Width(m.width - 2).
		Height(m.height - 2).
		Render(s.String())
}
