package notes

import (
	"strings"
	"testing"
	"time"
)

var ReferenceNote = `
## Todo

* [ ] 🖋 Write a blob post about papyrus
* [ ] 📕 Read 10 pages

## Daily Routine

* Have fun 🤪
* Make history 📜

## Today

If i can only do 1 thing today it would be: Watch new episode of the madisonian

---

### Fri 2021/05/14

* #meeting Where are we are we going to eat for lunch | 11:00/11:15
* #code Write a program to randomly select where to eat | 11:15/11:45
  * Call API to pull lists of food places
  * Random list generator
* #lunch | 11:45/12:30


### Thu 2021/05/13

* #meeting Council of Elrond | 10:00/11:00
  * https://lotr.fandom.com/wiki/Council_of_Elrond
  * Introductions
  * Who will accompany the ring
  * Conclusions: 
    * Frodo volunteered, he's the ring bearer
`
var TestNote = NotesFile{Lines: strings.Split(ReferenceNote, "\n")}

type FindResult struct {
	index int
	found bool
}
type FindTest struct {
	time time.Time
	want FindResult
}

func TestFindOn(t *testing.T) {
	tests := []FindTest{
		{time.Date(2021, 5, 14, 0, 0, 0, 0, time.UTC), FindResult{17, true}},
		{time.Date(2021, 5, 13, 0, 0, 0, 0, time.UTC), FindResult{26, true}},
		{time.Date(2021, 5, 14, 10, 24, 0, 0, time.UTC), FindResult{17, true}},
		{time.Date(2021, 5, 15, 0, 0, 0, 0, time.UTC), FindResult{0, false}},
	}
	m := Manager{Notes: TestNote}
	for _, tt := range tests {
		if index, found := m.findOn(tt.time); index != tt.want.index || found != tt.want.found {
			t.Errorf("n.findOn(%v) = (%d, %t) want (%d, %t)", tt.time, index, found, tt.want.index, tt.want.found)
		}
	}
}

func TestFindBefore(t *testing.T) {
	tests := []FindTest{
		{time.Date(2021, 5, 14, 0, 0, 0, 0, time.UTC), FindResult{26, true}},
		{time.Date(2021, 5, 14, 10, 24, 0, 0, time.UTC), FindResult{26, true}},
		{time.Date(2021, 5, 13, 0, 0, 0, 0, time.UTC), FindResult{0, false}},
		{time.Date(2021, 5, 20, 0, 0, 0, 0, time.UTC), FindResult{17, true}},
		{time.Date(2021, 5, 1, 0, 0, 0, 0, time.UTC), FindResult{0, false}},
	}
	m := Manager{Notes: TestNote}
	for _, tt := range tests {
		if index, found := m.findBefore(tt.time); index != tt.want.index || found != tt.want.found {
			t.Errorf("n.findBefore(%v) = (%d, %t) want (%d, %t)", tt.time, index, found, tt.want.index, tt.want.found)
		}
	}
}

func TestFindLatest(t *testing.T) {
	tests := []struct {
		t time.Time
		i int
	}{
		{time.Date(2021, 5, 14, 0, 0, 0, 0, time.UTC), 24},
		{time.Date(2021, 5, 14, 10, 24, 0, 0, time.UTC), 24},
		{time.Date(2021, 5, 13, 0, 0, 0, 0, time.UTC), 34},
		{time.Date(2021, 5, 20, 0, 0, 0, 0, time.UTC), 15},
	}
	m := Manager{Notes: TestNote}
	for _, tt := range tests {
		if index := m.findLatest(tt.t); index != tt.i {
			t.Errorf("n.findLatest(%v) = (%d) want (%d), line: %q", tt.t, index, tt.i, m.Notes.Lines[index])
		}
	}
}

func TestFormatDate(t *testing.T) {
	d := time.Date(2021, 5, 12, 0, 0, 0, 0, time.UTC)
	want := "### Wed 2021/05/12"

	if got := formatDate(d); got != want {
		t.Errorf("formatDate(%v) = (%q) want (%q)", d, got, want)
	}
}

func TestDatePattern(t *testing.T) {
	if !DATE_PATTERN.MatchString("### Wed 2021/05/12") {
		t.Error("Pattern does not match line")
	}
}

// func TestFormatEntry(t *testing.T) {
// 	start := time.Date(2021, 5, 12, 10, 13, 0, 0, time.UTC)
// 	end := time.Date(2021, 5, 12, 19, 55, 0, 0, time.UTC)

// 	got := formatEntry("title", start)
// 	want := "title | 10:13/"
// 	if got != want {
// 		t.Errorf("Got %s wanted %s", got, want)
// 	}

// 	got = formatCompleteEntry("title", start, end)
// 	want = "title | 10:13/19:55"
// 	if got != want {
// 		t.Errorf("Got %s wanted %s", got, want)
// 	}
// }
