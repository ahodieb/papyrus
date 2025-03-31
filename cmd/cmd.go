package cmd

import (
	"os"
	"strings"
	"time"

	"github.com/ahodieb/papyrus/editor"
	"github.com/ahodieb/papyrus/notes"
)

const (
	AppName = "papyrus"
)

func Run() error {
	cfg := LoadConfig()
	m := notes.Manager{
		Editor: editor.ByName(cfg.Editor),
		Notes:  notes.ReadFromFile(cfg.File),
	}

	args := os.Args[1:]
	if len(args) == 0 {
		// Default action:
		// Open the notes file at the latest position
		return m.OpenOrCreateLatest()
	}

	switch args[0] {
	default:
		p := m.AddEntry(strings.Join(args, " "), time.Now())
		return m.Open(p)
	}
}
