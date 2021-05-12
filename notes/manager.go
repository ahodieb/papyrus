package notes

import (
	"time"

	"github.com/ahodieb/papyrus/editor"
)

// TOOD better name ?
type Manager struct {
	Notes  NotesFile
	Editor editor.EditorOpener
}

func (m *Manager) Open() error {
	now := time.Now().UTC()
	position := m.Notes.FindPosition(now)
	return m.Editor.Open(m.Notes.Path, position)
}
