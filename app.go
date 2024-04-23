package main

import (
	"fmt"
	"gofz/internal/ssh"
	"gofz/internal/ui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func handleIncomingData(recv <-chan ssh.RecvEvent, fn func(msg tea.Msg)) {
	for {
		event := <-recv
		switch event.Event {
		case ssh.Quit:
			return
		default:
			fn(tea.Msg(event))
		}
	}
}

func main() {
	client := ssh.NewSftpClient()
	p := tea.NewProgram(ui.NewModel(client), tea.WithAltScreen())

	go handleIncomingData(client.Recv, p.Send)

	if err := p.Start(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
}
