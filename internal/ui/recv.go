package ui

import (
	"gofz/internal/ssh"

	tea "github.com/charmbracelet/bubbletea"
)

func StartRecv(recv <-chan ssh.RecvEvent, fn func(msg tea.Msg)) {
	for {
		event, open := <-recv

    if open {
			fn(tea.Msg(event))
    } else {
		  fn(tea.Msg(Connected(false)))	
      return
    }
	}
}
