package cmd

import (
	"fmt"
	"os"

	"github.com/ahodieb/papyrus/editor"
	"github.com/ahodieb/papyrus/notes"
)

const (
	AppName = "papyrus"
)

func Run() {
	cfg := LoadConfig()
	m := notes.Manager{
		Editor: editor.ByName(cfg.Editor),
		Notes:  notes.ReadFromFile(cfg.File),
	}

	// Default action:
	// Open the notes file at the latest position
	err := m.OpenLatest()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
