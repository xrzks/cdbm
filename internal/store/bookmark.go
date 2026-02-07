package store

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type Bookmark struct {
	Name      string
	Directory string
}

var (
	nameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	pathStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	separator = lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render("→")
	nilStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
)

func (b *Bookmark) Pretty() string {
	if b == nil {
		return fmt.Sprint(nilStyle.Render("<nil Bookmark>"))
	}

	name := nameStyle.Render(b.Name)
	directory := pathStyle.Render(b.Directory)
	return fmt.Sprintf("%s %s %s", name, separator, directory)
}
