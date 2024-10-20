package main

import (
	"fmt"
	"gofz/internal/ssh"
	"gofz/internal/ui"
	"os"


	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	channel := ssh.NewSftpChannel()
	p := tea.NewProgram(ui.NewModel(channel), tea.WithAltScreen())

	go ui.StartRecv(channel.Recv, p.Send)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error while running program:", err)
		os.Exit(1)
	}
}
