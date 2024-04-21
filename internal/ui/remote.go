package ui

import (
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type remoteModel struct {
	width  int
	height int
	active bool
	list   list.Model
}

func newRemoteModel() *remoteModel {
	return &remoteModel{
		active: false,
		list:   newList(nil, ""),
	}
}

func (m *remoteModel) SetActive(active bool) {
	m.active = active
}

func (m *remoteModel) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case []os.FileInfo:
		items := loadItems(msg)
		return m.list.SetItems(items)
	case tea.WindowSizeMsg:
		m.width = msg.Width / 2
		m.height = msg.Height
	case tea.KeyMsg:
		if !m.active {
			return nil
		}

		switch {
		case key.Matches(msg, keys.Enter):
			return nil
		}
	}

	if !m.active {
		return nil
	}

	m.list, cmd = m.list.Update(msg)
	return cmd
}

func (m *remoteModel) View() string {
	var s string

	s = m.list.View()

	if m.active {
		return focusBorderStyle.
			Width(m.width - 2).
			Height(m.height - 2).
			Render(s)
	}

	return borderStyle.
		Width(m.width - 2).
		Height(m.height - 2).
		Render(s)
}
