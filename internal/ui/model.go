package ui

import (
	"fmt"
	"gofz/internal/assert"
	"gofz/internal/ssh"
	"gofz/internal/system"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	localDirectory  *system.DirectoryCache
	remoteDirectory *system.DirectoryCache
)

func loadView(view ActiveView) ViewModel {
	switch view {
	case Transfer:
		return newTransferModel()
	case ServerList:
		return newServerList()
	default:
		return newMenuModel()
	}
}

type model struct {
	width     int
	height    int
	view      ActiveView
	viewModel ViewModel
	config    *ssh.Config
	msg       string
	connected Connected
	client    *ssh.SftpClient
	help      help.Model
	err       error
}

func (m *model) connect() tea.Cmd {
	connected, err := ssh.Connect(m.client, m.config)

	if err != nil {
		return func() tea.Msg {
			return err
		}
	}

	m.connected = Connected(connected)
	m.client.Getwd()
	m.client.List(".")
	return func() tea.Msg {
		return m.connected
	}
}

func (m model) Init() tea.Cmd {
	path, err := os.Getwd()
	assert.Assert("Working directory is not set", err)
	localDirectory = system.InitDirectoryCache(path)

	return nil
}

func NewModel(client *ssh.SftpClient) model {
	return model{
		view:      Menu,
		viewModel: loadView(Menu),
		client:    client,
		help:      help.New(),
	}
}

func (m model) UpdateViewPort() {
	header := headerStyle.
		Width(m.width - 2).
		Render(m.msg)
	headerHeight := lipgloss.Height(header)

	footer := footerStyle.
		Width(m.width - 2).
		Render(m.help.View(keys))
	footerHeight := lipgloss.Height((footer))

	m.viewModel.Update(tea.WindowSizeMsg{
		Width:  m.width,
		Height: m.height - (headerHeight + footerHeight),
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SendEvent:
		switch msg.Event {
		case ssh.Get:
			entry, _ := msg.Payload.(os.FileInfo)

			if entry.IsDir() {
				return m, func() tea.Msg {
					return fmt.Errorf("cannot get directory from remote host")
				}
			}

			m.client.Get(
				remoteDirectory.GetEntryPath(entry.Name()),
				localDirectory.GetEntryPath(entry.Name()),
			)

			return m, nil
		case ssh.Put:
			entry, _ := msg.Payload.(os.FileInfo)
			if entry.IsDir() {
				return m, func() tea.Msg {
					return fmt.Errorf("cannot put directory to remote host")
				}
			}

			m.client.Put(
				localDirectory.GetEntryPath(entry.Name()),
				remoteDirectory.GetEntryPath(entry.Name()),
			)
			return m, nil
		case ssh.List:
			path, _ := msg.Payload.(string)
			m.client.List(path)
		}
		return m, nil
	case ssh.RecvEvent:
		if msg.Event == ssh.Get {
			val, _ := msg.Payload.(string)
			m.msg = val
			m.UpdateViewPort()
			return m, m.viewModel.Update(tea.Msg(ReloadLocal))
		} else if msg.Event == ssh.Put {
			val, _ := msg.Payload.(string)
			m.msg = val
			m.UpdateViewPort()
		}
	case error:
		m.err = msg
		m.msg = msg.Error()
		m.UpdateViewPort()
		return m, nil
	case *ssh.Config:
		m.config = msg
		return m, func() tea.Msg {
			return Transfer
		}
	case Passphrase:
		m.client.SetPassphrase(string(msg))
		return m, m.connect()
	case ActiveView:
		m.view = msg
		m.viewModel = loadView(msg)
		m.msg = ""
		m.UpdateViewPort()
		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Back):
			if m.view == Transfer && m.connected {
				m.connected = false
				m.client.Quit()
			}
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			m.UpdateViewPort()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		m.UpdateViewPort()
		return m, nil
	}

	var cmds []tea.Cmd
	cmds = append(cmds, m.viewModel.Update(msg))
	return m, tea.Batch(cmds...)
}

func (m *model) header() string {
	return headerStyle.Width(m.width - 2).Render(m.msg)
}

func (m *model) footer() string {
	return footerStyle.Width(m.width - 2).Render(m.help.View(keys))
}

func (m model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.header(),
		m.viewModel.View(),
		m.footer(),
	)
}
