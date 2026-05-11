package store

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type Bookmark struct {
	Name      string
	Directory string
}

const (
	ColorCyan  = "6"
	ColorGreen = "10"
	ColorGray  = "245"
	ColorRed   = "1"
)

var (
	nameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorCyan)).Bold(true)
	pathStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorGreen))
	separator = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorGray)).Render("→")
	nilStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorRed))
)

// Pretty returns a formatted string representation of the bookmark.
// If the bookmark is nil, it returns a styled "<nil Bookmark>" message.
// Otherwise, it returns a formatted string with the bookmark name and directory.
func (b *Bookmark) Pretty() string {
	if b == nil {
		return fmt.Sprint(nilStyle.Render("<nil Bookmark>"))
	}

	name := nameStyle.Render(b.Name)
	directory := pathStyle.Render(b.Directory)
	return fmt.Sprintf("%s %s %s", name, separator, directory)
}
