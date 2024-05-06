package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type passphrase struct {
	width     int
	height    int
	hidden    bool
	focus     bool
	textInput textinput.Model
}

func newInput(hidden bool, placeholder string) *passphrase {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 255
	ti.Width = 20

	return &passphrase{
		textInput: ti,
		hidden:    hidden,
	}
}

func (m *passphrase) Focus() {
	m.textInput.Focus()
	m.focus = true
}

func (m *passphrase) Blur() {
	m.textInput.Blur()
	m.textInput.Reset()
	m.focus = false
}

func (m *passphrase) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width / 2
	case tea.KeyMsg:
		if m.hidden && msg.Type == tea.KeyRunes {
			star := tea.Key{
				Type:  tea.KeyRunes,
				Runes: []rune{rune('*')},
				Alt:   false,
			}

			m.textInput, cmd = m.textInput.Update(tea.KeyMsg(star))
		} else {
			m.textInput, cmd = m.textInput.Update(msg)
		}
	}

	return cmd
}

func (m *passphrase) View() string {
	s := fmt.Sprintf(
		"Passphrase required:\n%s\n",
		m.textInput.View(),
	)

	if m.focus {
		return centerContentFocus.
			Width(m.width - 2).
			Height(m.height - 2).
			Render(s)
	}

	return centerContent.
		Width(m.width - 2).
		Height(m.height - 2).
		Render(s)
}
