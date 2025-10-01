package main

import (
	"log"
	"os"

	"dotfiles-manager/tui"
)

func main() {
	// Run the TUI
	err := tui.RunTUI()
	if err != nil {
		log.Fatal(err)
	}
}