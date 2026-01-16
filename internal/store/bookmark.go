package store

import (
	"fmt"
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
	fmt.Printf("%-15s %s\n", "Directory:", b.Directory)
	fmt.Println()
}
