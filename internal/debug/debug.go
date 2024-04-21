package debug

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func Write(value any, msg string) {
	logText := fmt.Sprintf("%v: %+v\n", msg, value)
	f, _ := tea.LogToFile("debug.log", "debug")
	defer f.Close()
	f.Write([]byte(logText))
}
