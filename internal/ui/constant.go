package ui

import (
	"gofz/internal/ssh"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	focusLeft  = 0
	focusRight = 1
)

type Passphrase string
type ActiveView int
type Connected bool

const (
	Menu ActiveView = iota
	ServerList
	Transfer
)

type ViewModel interface {
	Update(msg tea.Msg) tea.Cmd
	View() string
}

type SendEvent struct {
	Event   ssh.EventType
	Payload any
}
