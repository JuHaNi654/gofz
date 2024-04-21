package ui

import (
	"gofz/internal/assert"
	"gofz/internal/system"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type localModel struct {
	width  int
	height int
	list   list.Model
	active bool

	wd *system.DirectoryCache
}

func newLocalModel() *localModel {
	path, err := os.Getwd()
	assert.Assert("Working directory is not set", err)
	wd := system.InitDirectoryCache(path)

	return &localModel{
		list:   newList(wd.Entries(), wd.GetWd()),
		wd:     wd,
		active: true,
	}
}

func (m *localModel) SetActive(active bool) {
	m.active = active
}

func (m *localModel) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width / 2
		m.height = msg.Height
	case tea.KeyMsg:
		if !m.active {
			return nil
		}

		switch {
		case key.Matches(msg, keys.Enter):
			if 0 == m.list.Index() {
				m.wd.PreviousWd()
				items := loadItems(m.wd.Entries())
				return m.list.SetItems(items)
			}

			selected := m.list.SelectedItem()
			i, _ := selected.(item)

			if i.Entry.IsDir() {
				m.wd.NextWd(i.Entry.Name())
				items := loadItems(m.wd.Entries())
				m.list.Select(0)
				return m.list.SetItems(items)
			}

			return nil
		}
	}

	// Maybe move up before switch
	if !m.active {
		return nil
	}

	m.list, cmd = m.list.Update(msg)
	return cmd
}

func (m *localModel) View() string {
	if m.active {
		return focusBorderStyle.
			Width(m.width - 2).
			Height(m.height - 2).
			Render(m.list.View())
	}
	return borderStyle.
		Width(m.width - 2).
		Height(m.height - 2).
		Render(m.list.View())
}
