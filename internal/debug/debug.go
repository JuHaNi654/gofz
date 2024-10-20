package debug

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func Write(label any, value any) {
	logText := fmt.Sprintf("%v: %+v\n", label, value)
	f, _ := tea.LogToFile("debug.log", "debug")
	defer f.Close()
	f.Write([]byte(logText))
}
