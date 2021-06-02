package notes

import (
	"fmt"
	"os"
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
	return m.Editor.Open(m.Notes.Path, m.findLatest(time.Now()))
}

// Open opens notes file in the editor at the specified line index
func (m *Manager) Open(i int) error {
	return m.Editor.Open(m.Notes.Path, i)
}

// AddEntry to notes and return its position
// FIXME still not working properly
// TODO restructure
func (m *Manager) AddEntry(title string, t time.Time) (position int) {
	if _, found := m.findOn(t); !found {
		latest, _ := m.findBefore(t)
		latest = m.Notes.Insert([]string{formatDate(t), ""}, latest)
		position = m.Notes.Insert([]string{formatEntry(title, t)}, latest)
	} else {
		position = m.Notes.Insert([]string{formatEntry(title, t)}, m.findLatest(t)+1)
	}

	if position < len(m.Notes.Lines) && m.Notes.Lines[position] != "" {
		m.Notes.Insert([]string{""}, position)
	}

	if _, err := m.Notes.SaveWithBackup(); err != nil {
		fmt.Fprintf(os.Stderr, "Backup failed: %s", err)
	}

	return
}

// findLatest finds the latest entry in the notes
func (m *Manager) findLatest(t time.Time) int {
	// Looks for the latest entry before the current day entry, and get the position two lines above it
	if position, found := m.findBefore(t); found {
		if position > 2 {
			return position - 2
		}

		return 0
	}

	// If no later entry found look for the current day entry,
	// This could happen if there is only one entry in the file, or change of formats
	// If none were found that means either its an empty file or formats are not recognized
	// and it will default back to the 0th position
	return len(m.Notes.Lines) - 1
}

// findOn finds section created on the specified date
// returns the index and a found bool
func (m *Manager) findOn(t time.Time) (int, bool) {
	return m.Notes.FindContains(formatDate(t))
}

func (m *Manager) findBefore(t time.Time) (int, bool) {
	date := t.Truncate(time.Hour * 24)
	for _, s := range m.Sections() {
		if s.Time.Before(date) {
			return s.Index, true
		}
	}

	return 0, false
}

type Section struct {
	Index   int
	Time    time.Time
	Content []string
}

//FIXME It does not account for other sections in the journal file
// Also i want to think about the terminology (section, entry, ...etc)
func (m *Manager) Sections() []Section {
	var sections []Section

	for i, line := range m.Notes.Lines {
		if DATE_PATTERN.MatchString(line) {
			t, err := parseDate(line)
			if err == nil {
				sections = append(sections, Section{Index: i, Time: t})
			}
		}
	}

	for i, section := range sections {
		if i+1 >= len(sections) {
			sections[i].Content = m.Notes.Lines[section.Index:]
		} else {
			sections[i].Content = m.Notes.Lines[section.Index:sections[i+1].Index]
		}
	}

	sort.Slice(sections, func(i, j int) bool {
		return sections[i].Time.After(sections[j].Time)
	})
	return sections
}

var DATE_PATTERN, _ = regexp.Compile("### [a-zA-Z]{3} (?P<year>[0-9]{4})/(?P<month>[0-9]{2})/(?P<day>[0-9]{2})")

const ROUND_DURATION = 5 * time.Minute

func formatDate(t time.Time) string {
	return t.Format("### Mon 2006/01/02")
}

func parseDate(s string) (time.Time, error) {
	return time.Parse("### Mon 2006/01/02", strings.TrimSpace(s))
}

func formatEntry(title string, start time.Time) string {
	formatedTitle := reverseHashTags(title)
	return fmt.Sprintf("* %s | %s/", formatedTitle, floorTime(start, ROUND_DURATION).Format("15:04"))
}

func floorTime(t time.Time, d time.Duration) time.Time {
	round := t.Round(d)
	if round.After(t) {
		return round.Add(-d)
	}
	return round
}

// reverse hashtags in entry title e.g  something# ->  #something
// creating a new entry starting with #something does not work in bash without escaping it
// as a convince i reverse them something# which works fine on bash
func reverseHashTags(s string) string {
	w := strings.Split(s, " ")

	for i := range w {
		if strings.HasSuffix(w[i], "#") {
			w[i] = "#" + strings.TrimSuffix(w[i], "#")
		}
	}

	return strings.Join(w, " ")
}
