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
  Entry  os.FileInfo
	Label  string
  notify bool	
}

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int  { return 1 }
func (d itemDelegate) Spacing() int { return 0 }
func (d itemDelegate) Update(_ tea.Msg, list *list.Model) tea.Cmd {
	if list.Index() != 0 {
		i := list.SelectedItem().(item)
		if i.notify {
			i.notify = false
			list.SetItem(list.Index(), i)
		}
	}

	return nil
}
func (d itemDelegate) Render(
	w io.Writer,
	m list.Model,
	index int,
	listItem list.Item,
) {
	var str string

	i, ok := listItem.(item)
	if !ok {
		return
	}

	if i.notify {
		str = fmt.Sprintf("%s*", i.Label)
	} else {
		str = fmt.Sprintf("%s", i.Label)
	}

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

	// TODO: think another way insert this, so that there is no entry with nil value
	items = append(items, item{Label: ".."})
	for _, entry := range entries {
		items = append(items, item{Label: entry.Name(), Entry: entry})
	}

	return items
}

func compareItems(entries []os.FileInfo, listItems []list.Item) []list.Item {
	items := []list.Item{}

	items = append(items, item{Label: ".."})
	for _, entry := range entries {
		notify := fileChanged(entry, &listItems)
		items = append(items, item{Label: entry.Name(), Entry: entry, notify: notify})
	}

	return items
}

func fileChanged(entry os.FileInfo, listItems *[]list.Item) bool {
	newEntry := true

	for _, listItem := range *listItems {
		i, _ := listItem.(item)
		if i.Entry == nil || i.Entry.Name() != entry.Name() {
			continue
		}

		newEntry = false
		diff := i.Entry.ModTime().Compare(entry.ModTime())
		if diff != 0 || i.notify {
			return true
		}
	}

	if newEntry {
		return true
	}

	return false
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
