package main

import (
	"fmt"
	"gofz/internal/debug"
	"gofz/internal/ssh"
	"gofz/internal/ui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func handleIncomingData(recv <-chan any, fn func(msg tea.Msg)) {
	for {
		debug.Write("waiting for recv", "log")
		event := <-recv
		switch msg := event.(type) {
		case []os.FileInfo:
			fn(tea.Msg(msg))
		case bool:
			return
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
