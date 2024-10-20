package ui

import (
	"gofz/internal/ssh"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type localModel struct {
	focus bool
	list  list.Model
}

func newLocalModel() *localModel {
	return &localModel{
		focus: true,
		list:  newList(localDirectory.Entries(), localDirectory.GetWd()),
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
	case ViewEvent:
		if msg == ReloadLocal {
			items := compareItems(
				localDirectory.Entries(),
				m.list.Items(),
			)

			m.list.Select(0)
			return m.list.SetItems(items)
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Transfer):
			selected := m.list.SelectedItem()
			i, _ := selected.(item)

			return func() tea.Msg {
				return SendEvent{
					Event:   ssh.Put,
					Payload: i.Entry,
				}
			}

		case key.Matches(msg, keys.Enter):
			if 0 == m.list.Index() {
				localDirectory.PreviousWd()
				m.list.Title = localDirectory.GetWd()
				items := loadItems(localDirectory.Entries())
				return m.list.SetItems(items)
			}

			selected := m.list.SelectedItem()
			i, _ := selected.(item)

			if i.Entry.IsDir() {
				localDirectory.NextWd(i.Entry.Name())
				items := loadItems(localDirectory.Entries())
				m.list.Title = localDirectory.GetWd()
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
	w := (screen.Width() / 2) - borderSpacing
	h := screen.Height() - borderSpacing

	if m.focus {
		return focusBorderStyle.
			Width(w).
			Height(h).
			Render(m.list.View())
	}

	return borderStyle.
		Width(w).
		Height(h).
		Render(m.list.View())
}
