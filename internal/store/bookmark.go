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

func (b *Bookmark) Pretty() {
	if b == nil {
		fmt.Println(nilStyle.Render("<nil Bookmark>"))
		return
	}

	name := nameStyle.Render(b.Name)
	directory := pathStyle.Render(b.Directory)
	fmt.Printf("%s %s %s\n", name, separator, directory)
}
