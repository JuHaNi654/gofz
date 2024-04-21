package ui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type item struct {
	Label string
	Entry os.FileInfo
}

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(
	w io.Writer,
	m list.Model,
	index int,
	listItem list.Item,
) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.Label)
	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render(">", strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func loadItems(entries []os.FileInfo) []list.Item {
	items := []list.Item{}

	items = append(items, item{Label: ".."})

	for _, entry := range entries {
		items = append(items, item{Label: entry.Name(), Entry: entry})
	}

	return items
}

func newList(entries []os.FileInfo, path string) list.Model {
	items := loadItems(entries)

	l := list.New(items, itemDelegate{}, listWidth, listHeight)
	l.Title = path
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle

	return l
}
