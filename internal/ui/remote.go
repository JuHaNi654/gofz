package ui

import (
	"gofz/internal/ssh"
	"gofz/internal/system"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type remoteModel struct {
	width      int
	height     int
	focus      bool
	list       list.Model
  reloadList bool
}

func newRemoteModel() *remoteModel {
	return &remoteModel{
		list: newList(nil, ""),
	}
}

func (m *remoteModel) Focus() {
	m.focus = true
}

func (m *remoteModel) Blur() {
	m.focus = false
}

func (m *remoteModel) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
  case ssh.RecvEvent:
		switch msg.Event {
	  case ssh.Put: {
      m.reloadList = true
      return func() tea.Msg {
        return SendEvent{
          Event:   ssh.List,
          Payload: remoteDirectory.GetWd(),
        }
      }
    }	
    case ssh.List:
      if m.reloadList {
        m.reloadList = false 
        entries, _ := msg.Payload.([]os.FileInfo)
        items := compareItems(
          entries, 
          m.list.Items(),
        )

        return m.list.SetItems(items) 
      }
			items, _ := msg.Payload.([]os.FileInfo)
			return m.list.SetItems(loadItems(items))
		case ssh.Wd:
			path, _ := msg.Payload.(string)
			remoteDirectory = system.InitDirectoryCache(path)
			m.list.Title = path
			return nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width / 2
		m.height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Transfer):
			selected := m.list.SelectedItem()
			i, _ := selected.(item)

			return func() tea.Msg {
				return SendEvent{
					Event:   ssh.Get,
					Payload: i.Entry,
				}
			}
		case key.Matches(msg, keys.Enter):
			if 0 == m.list.Index() {
				remoteDirectory.PreviousWd()
				m.list.Title = remoteDirectory.GetWd()
				return func() tea.Msg {
					return SendEvent{
						Event:   ssh.List,
						Payload: remoteDirectory.GetWd(),
					}
				}
			}

			selected := m.list.SelectedItem()
			i, _ := selected.(item)

			if i.Entry.IsDir() {
				remoteDirectory.NextWd(i.Entry.Name())
				m.list.Title = remoteDirectory.GetWd()
				m.list.Select(0)
				return func() tea.Msg {
					return SendEvent{
						Event:   ssh.List,
						Payload: remoteDirectory.GetWd(),
					}
				}
			}

			return nil
		}
	}

	m.list, cmd = m.list.Update(msg)
	return cmd
}

func (m *remoteModel) View() string {
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
