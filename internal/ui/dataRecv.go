package ui

import (
	"gofz/internal/ssh"

	tea "github.com/charmbracelet/bubbletea"
)

func HandleIncomingData(recv <-chan ssh.RecvEvent, fn func(msg tea.Msg)) {
	for {
		event := <-recv
		switch event.Event {
		case ssh.Quit:
		  fn(tea.Msg(Connected(false)))	
      return
		default:
			fn(tea.Msg(event))
		}
	}
}
