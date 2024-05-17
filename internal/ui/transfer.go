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

	local  *localModel
	remote *remoteModel
}

func newTransferModel() *transferModel {
  return &transferModel{
		local:       newLocalModel(),
		remote:      newRemoteModel(),
	}
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
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.local.Update(msg)
		m.remote.Update(msg)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Back):
      return func() tea.Msg {
        return Menu
      }
		case key.Matches(msg, keys.Next), key.Matches(msg, keys.Prev):
			if m.focus == focusLeft {
				m.focus = focusRight

				m.local.Blur()
				m.remote.Focus()
			} else {
				m.focus = focusLeft

				m.local.Focus()
				m.remote.Blur()
			}

			return nil
		}
	}

	if focusLeft == m.focus {
		cmds = append(cmds, m.local.Update(msg))
	} else if focusRight == m.focus {
		cmds = append(cmds, m.remote.Update(msg))
	}

	return tea.Batch(cmds...)
}

func (m *transferModel) View() string {
	var views []string
	views = append(views, m.local.View())
	views = append(views, m.remote.View())
	return lipgloss.JoinHorizontal(lipgloss.Left, views...)
}
