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
	focus  bool

	wd *system.DirectoryCache
}

func newLocalModel() *localModel {
	path, err := os.Getwd()
	assert.Assert("Working directory is not set", err)
	wd := system.InitDirectoryCache(path)

	return &localModel{
		focus: true,
		list:  newList(wd.Entries(), wd.GetWd()),
		wd:    wd,
	}
}

func (m *localModel) Focus() {
	m.focus = true
}

func (m *localModel) Blur() {
	m.focus = false
}

func (m *localModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width / 2
		m.height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Enter):
			if 0 == m.list.Index() {
				m.wd.PreviousWd()
				m.list.Title = m.wd.GetWd()
				items := loadItems(m.wd.Entries())
				return m.list.SetItems(items)
			}

			selected := m.list.SelectedItem()
			i, _ := selected.(item)

			if i.Entry.IsDir() {
				m.wd.NextWd(i.Entry.Name())
				items := loadItems(m.wd.Entries())
				m.list.Title = m.wd.GetWd()
				m.list.Select(0)
				return m.list.SetItems(items)
			}

			return nil
		}
	}

	m.list, _ = m.list.Update(msg)
	return nil
}

func (m *localModel) View() string {
	if m.focus {
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
