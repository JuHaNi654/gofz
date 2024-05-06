package ui

import (
	"gofz/internal/ssh"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type transferModel struct {
	width       int
	height      int
	focus       int
	passphrase  string
	inputActive bool
	input       *passphrase

	local  *localModel
	remote *remoteModel
}

func newTransferModel() *transferModel {
	return &transferModel{
		local:       newLocalModel(),
		remote:      newRemoteModel(),
		input:       newInput(true, "Passphrase ..."),
		inputActive: true,
	}
}

func (m *transferModel) ClearPassphrase() {
	m.passphrase = ""
}

func (m *transferModel) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
  case ViewEvent:
    if (msg == ReloadLocal) {
      return m.local.Update(msg)
    } else {
      return m.remote.Update(msg) 
    }
  case ssh.RecvEvent:
		return m.remote.Update(msg)
	case Connected:
		if msg {
			m.inputActive = false
		} else {
			m.inputActive = true
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.input.Update(msg)
		m.local.Update(msg)
		m.remote.Update(msg)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Back):
			if m.focus == focusRight && m.inputActive {
				m.remote.Update(msg)

				if len(m.passphrase) > 0 {
					m.passphrase = m.passphrase[:len(m.passphrase)-1]
				}

			} else {
				return func() tea.Msg {
					return Menu
				}
			}
		case key.Matches(msg, keys.Next), key.Matches(msg, keys.Prev):
			m.ClearPassphrase()
			if m.focus == focusLeft {
				m.focus = focusRight

				m.local.Blur()
				m.remote.Focus()
				m.input.Focus()

			} else {
				m.focus = focusLeft

				m.local.Focus()
				m.remote.Blur()
				m.input.Blur()
			}

			return nil
		case key.Matches(msg, keys.Enter):
			if m.focus == focusRight && m.inputActive {
				return func() tea.Msg {
					return Passphrase(m.passphrase)
				}
			}
		default:
			if m.inputActive && msg.Type == tea.KeyRunes {
				m.passphrase = m.passphrase + msg.String()
			}
		}
	}

	if focusLeft == m.focus {
		cmds = append(cmds, m.local.Update(msg))
	} else if focusRight == m.focus && m.inputActive {
		cmds = append(cmds, m.input.Update(msg))
	} else {
		cmds = append(cmds, m.remote.Update(msg))
	}

	return tea.Batch(cmds...)
}

func (m *transferModel) View() string {
	var views []string
	views = append(views, m.local.View())

	if m.inputActive {
		views = append(views, m.input.View())
	} else {
		views = append(views, m.remote.View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, views...)
}
