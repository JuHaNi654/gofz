package main

import (
	"fmt"
	"gofz/internal/ssh"
	"gofz/internal/ui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	client := ssh.NewSftpClient()
	p := tea.NewProgram(ui.NewModel(client), tea.WithAltScreen())

	go ui.HandleIncomingData(client.Recv, p.Send)

	if err := p.Start(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
}
