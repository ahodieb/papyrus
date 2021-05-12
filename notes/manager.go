package notes

import (
	"fmt"
	"sort"
	"time"

	"github.com/ahodieb/papyrus/editor"
	"github.com/ahodieb/papyrus/notes/markdown"
)

type Manager struct {
	Notes  NotesFile
	Editor editor.EditorOpener
}

func (m *Manager) Open() error {
	now := time.Now().UTC()
	position := m.FindPosition(now)

	fmt.Println(position)
	return m.Editor.Open(m.Notes.Path, position)
}

// Find the latest entry in the Journal file
func (m *Manager) FindPosition(t time.Time) int {

	// Looks for the latest entry before the current day entry, and get the position two lines above it
	position, found := m.positionBefore(t)
	if found {
		if position > 1 {
			return position - 1
		}

		return 0
	}

	// If no later entry found look for the current day entry,
	// This could happen if there is only one entry in the file, or change of formats
	// If none were found that means either its an empty file or formats are not recognized
	// and it will default back to the 0th position
	position, _ = m.positionOn(t)
	return position
}

func (m *Manager) positionOn(t time.Time) (int, bool) {
	return m.Notes.FindContaining(markdown.FormatDate(t))
}

func (m *Manager) positionBefore(t time.Time) (int, bool) {

	for _, e := range m.TimeEntries() {
		if e.Time.Before(t) {
			return e.Position, true
		}
	}

	return 0, false
}

type TimeEntryPosition struct {
	Position int
	Time     time.Time
}

func (m *Manager) TimeEntries() []TimeEntryPosition {
	var positions []TimeEntryPosition

	for i, line := range m.Notes.Lines {
		if markdown.DATE_PATTERN.MatchString(line) {
			t, err := markdown.ParseDate(line)
			if err == nil {
				positions = append(positions, TimeEntryPosition{Position: i, Time: t})
			}
		}
	}

	sort.Slice(positions, func(i, j int) bool {
		return positions[i].Time.After(positions[j].Time)
	})

	return positions
}
