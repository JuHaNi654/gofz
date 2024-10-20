package ui

import (
	"gofz/internal/system"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

type Header struct {
	content string
}

func (h *Header) SetContent(input string) {
	h.content = input
}

func (h *Header) Render() string {
	return headerStyle.
		Width(screen.width - borderSpacing).
		Render(h.content)
}

func (h *Header) Height() int {
	return lipgloss.Height(h.Render())
}

type Footer struct {
	content help.Model
}

func (f *Footer) SetWidth(w int) {
	f.content.Width = w
}

func (f *Footer) Render() string {
	return footerStyle.
		Width(screen.width - borderSpacing).
		Render(f.content.View(keys))
}

func (f *Footer) Height() int {
	return lipgloss.Height(f.Render())
}

type Screen struct {
	height int
	width  int
}

func (s Screen) Height() int {
	return s.height - (footer.Height() + header.Height())
}

func (s Screen) Width() int {
	return s.width
}

var (
	localDirectory  *system.DirectoryCache
	remoteDirectory *system.DirectoryCache
	screen          = Screen{}
	footer          = Footer{content: help.New()}
	header          = Header{}
)
