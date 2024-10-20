package ui

import (
	"fmt"
	"gofz/internal/assert"
	"gofz/internal/ssh"
	"gofz/internal/system"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func loadView(view ActiveView) ViewModel {
	switch view {
	case Transfer:
		return newTransferModel()
	case ServerList:
		return newServerList()
	case PassphraseInput:
		return newPassphraseInput()
	default:
		return newMenuModel()
	}
}

type model struct {
	msg       string
	ready     bool
	err       error
	view      ActiveView
	viewModel ViewModel
	config    *ssh.Config
	connected Connected
	channel   *ssh.SftpChannel
}

func (m *model) connect() tea.Cmd {
	m.channel.OpenSender()
	connected, err := ssh.Connect(m.channel, m.config)

	if err != nil {
		return func() tea.Msg {
			return err
		}
	}

	m.connected = Connected(connected)
	m.channel.Getwd()
	m.channel.List(".")

	return nil
}

func (m model) Init() tea.Cmd {
	path, err := os.Getwd()
	assert.Assert("Working directory is not set", err)
	localDirectory = system.InitDirectoryCache(path)

	return nil
}

func NewModel(channel *ssh.SftpChannel) model {
	return model{
		view:      Menu,
		channel:   channel,
		viewModel: loadView(Menu),
		ready:     true,
	}
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

			m.channel.Get(
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

			m.channel.Put(
				localDirectory.GetEntryPath(entry.Name()),
				remoteDirectory.GetEntryPath(entry.Name()),
			)
			return m, nil
		case ssh.List:
			path, _ := msg.Payload.(string)
			m.channel.List(path)
		}
		return m, nil
	case ssh.RecvEvent:
		switch msg.Event {
		case ssh.Get:
			val, _ := msg.Payload.(string)
			header.SetContent(val)
			return m, m.viewModel.Update(tea.Msg(ReloadLocal))
		case ssh.Put:
			val, _ := msg.Payload.(string)
			header.SetContent(val)
		case ssh.Connected:
			// TODO update ready state
		}
	case *ssh.Config:
		m.config = msg
		if msg.Protected {
			return m, func() tea.Msg {
				return PassphraseInput
			}
		}

		return m, func() tea.Msg {
			return Transfer
		}
	case Passphrase:
		m.config.Passphrase = string(msg)
		return m, func() tea.Msg {
			return Transfer
		}
	case ActiveView:
		m.view = msg
		m.viewModel = loadView(msg)

		if msg == Transfer {
			return m, m.connect()
		}

		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Back):
			if m.view == Transfer && m.connected {
				m.connected = false
				m.channel.Quit()
			}
		case key.Matches(msg, keys.Quit):
			if m.view == Transfer && m.connected {
				m.connected = false
				m.channel.Quit()
			}
			return m, tea.Quit
		case key.Matches(msg, keys.Help):
			footer.content.ShowAll = !footer.content.ShowAll
		}
	case tea.WindowSizeMsg:
		screen.width = msg.Width
		screen.height = msg.Height
		footer.SetWidth(msg.Width)
		return m, nil
	case error:
		m.err = msg
		header.SetContent(msg.Error())

		switch msg.(type) {
		case ssh.PassphraseError:
		case ssh.SSHClientError:
			return m, func() tea.Msg {
				return ServerList
			}
		}

		return m, nil
	}

	var cmds []tea.Cmd
	cmds = append(cmds, m.viewModel.Update(msg))
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		header.Render(),
		m.viewModel.View(),
		footer.Render(),
	)
}
