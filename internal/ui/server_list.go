package ui

import (
	"fmt"
	"gofz/internal/assert"
	"gofz/internal/ssh"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var default_config = "/.ssh/config"

type serverItem string

func (i serverItem) FilterValue() string { return "" }

type serverItemDelegate struct{}

func (d serverItemDelegate) Height() int                             { return 1 }
func (d serverItemDelegate) Spacing() int                            { return 0 }
func (d serverItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d serverItemDelegate) Render(
	w io.Writer,
	m list.Model,
	index int,
	listItem list.Item,
) {
	i, ok := listItem.(serverItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render(strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func NewServerList(servers []*ssh.Config) list.Model {
	items := []list.Item{}

	for _, i := range servers {
		items = append(items, serverItem(i.Host))
	}

	l := list.New(items, serverItemDelegate{}, listWidth, listHeight)
	l.Title = "Select server"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle

	return l
}

type serverList struct {
	servers []*ssh.Config
	list    list.Model
}

func newServerList() *serverList {
	var servers []*ssh.Config
	dir, err := os.UserHomeDir()
	assert.Assert("user home dir is not set", err)

	servers = ssh.LoadSSHConfig(dir + default_config)

	return &serverList{
		list:    NewServerList(servers),
		servers: servers,
	}
}

func (m *serverList) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Back):
			return func() tea.Msg {
				return Menu
			}
		case key.Matches(msg, keys.Enter):
			server := m.servers[m.list.Index()]
			return func() tea.Msg {
				return server
			}
		}
	}

	m.list, cmd = m.list.Update(msg)
	return cmd
}

func (m *serverList) View() string {
	return centerContent.
		Width(screen.Width() - borderSpacing).
		Height(screen.Height() - borderSpacing).
		Render(m.list.View())
}
