package ui

import (
	"gofz/internal/debug"
	"gofz/internal/ssh"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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

	connected Connected
	client    *ssh.SftpClient

	err error
}

func (m *model) connect() tea.Cmd {
	connected, err := ssh.Connect(m.client, m.config)

	if err != nil {
		return func() tea.Msg {
			return err
		}
	}

	m.connected = Connected(connected)
	m.client.List(".")
	return func() tea.Msg {
		return m.connected
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func NewModel(client *ssh.SftpClient) model {
	return model{
		view:      Menu,
		viewModel: loadView(Menu),
		client:    client,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		debug.Write(msg, "Error")
		m.err = msg
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
		m.viewModel.Update(tea.WindowSizeMsg{
			Width:  m.width,
			Height: m.height,
		})
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Back):
			if m.view == Transfer && m.connected {
				m.connected = false
				m.client.Quit()
			}
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	var cmds []tea.Cmd
	cmds = append(cmds, m.viewModel.Update(msg))
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return m.viewModel.View()
}
