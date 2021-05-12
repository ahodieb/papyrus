package notes

import (
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/ahodieb/papyrus/editor"
)

type Manager struct {
	Notes  NotesFile
	Editor editor.EditorOpener
}

// Open notes file in the editor at the end of the latest time entry
func (m *Manager) OpenLatest() error {
	return m.Editor.Open(m.Notes.Path, m.FindPosition(time.Now()))
}

// Open notes file in the editor at the specified position
func (m *Manager) Open(position int) error {
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
	return m.Notes.FindContaining(FormatDate(t))
}

func (m *Manager) positionBefore(t time.Time) (int, bool) {

	for _, e := range m.Entries() {
		if e.Time.Before(t) {
			return e.Position, true
		}
	}

	return 0, false
}

type Entry struct {
	Position int
	Time     time.Time
	Content  []string
}

func (m *Manager) Entries() []Entry {
	var entries []Entry

	for i, line := range m.Notes.Lines {
		if DATE_PATTERN.MatchString(line) {
			t, err := ParseDate(line)
			if err == nil {
				entries = append(entries, Entry{Position: i, Time: t})
			}
		}
	}

	for i, entry := range entries {
		if i+1 >= len(entries) {
			entries[i].Content = m.Notes.Lines[entry.Position:]
		} else {
			entries[i].Content = m.Notes.Lines[entry.Position:entries[i+1].Position]
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Time.After(entries[j].Time)
	})

	return entries
}

var DATE_PATTERN, _ = regexp.Compile("### [a-zA-Z]{3} (?P<year>[0-9]{4})/(?P<month>[0-9]{2})/(?P<day>[0-9]{2})")

func FormatDate(t time.Time) string {
	return t.Format("### Mon 2006/01/02")
}

func ParseDate(s string) (time.Time, error) {
	return time.Parse("### Mon 2006/01/02", strings.TrimSpace(s))
}
