package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type passphrase struct {
	width     int
	height    int
  value     string
	textInput textinput.Model
}

func newPassphraseInput() *passphrase {
	ti := textinput.New()
	ti.Placeholder = "Passphrase ..."
	ti.CharLimit = 255
	ti.Width = 20
  ti.Focus()

	return &passphrase{
		textInput: ti,
	}
}

func (m *passphrase) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
    switch {
    case key.Matches(msg, keys.Enter):
      return func() tea.Msg {
        return Passphrase(m.value)
      }
    case key.Matches(msg, keys.Back):
      if len(m.value) > 0 {
        m.value = m.value[:len(m.value)-1] 
      }
      
      m.textInput, cmd = m.textInput.Update(msg)
    default: 
      if msg.Type == tea.KeyRunes {
        star := tea.Key{
          Type:  tea.KeyRunes,
          Runes: []rune{rune('*')},
          Alt:   false,
        }

        m.value += msg.String()
        m.textInput, cmd = m.textInput.Update(tea.KeyMsg(star))
      } else {
        m.textInput, cmd = m.textInput.Update(msg)
      }
    }
	}

	return cmd
}

func (m *passphrase) View() string {
	s := fmt.Sprintf(
		"Passphrase required:\n%s\n",
		m.textInput.View(),
	)

	return centerContent.
		Width(m.width - 2).
		Height(m.height - 2).
		Render(s)
}
