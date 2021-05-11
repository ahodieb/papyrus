package journal

import (
	"time"

	"github.com/ahodieb/papyrus/editor"
)

// TOOD better name ?
type Manager struct {
	Journal Journal
	Editor  editor.EditorOpener
}

func (m *Manager) Open() error {
	now := time.Now().UTC()
	position := m.Journal.FindPosition(now)
	return m.Editor.Open(m.Journal.Path, position)
}
