package store

import (
	"fmt"
	"path/filepath"
)

type Bookmark struct {
	Name      string
	Directory string
}

func (b *Bookmark) Pretty() {
	if b == nil {
		fmt.Println("<nil Bookmark>")
		return
	}
	fmt.Printf("%-15s %s\n", "Name:", b.Name)
	fmt.Printf("%-15s %s\n", "Directory:", filepath.Base(b.Directory))
}
