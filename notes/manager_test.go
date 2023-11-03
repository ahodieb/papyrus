package notes

import (
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
var TestNote = FromText(ReferenceNote)

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
		if index, found := m.FindOn(tt.time); index != tt.want.index || found != tt.want.found {
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
	if !DatePattern.MatchString("### Wed 2021/05/12") {
		t.Error("Pattern does not match line")
	}
}

var AddEntryTests = []struct {
	text    string
	title   string
	time    time.Time
	newText string
	index   int
}{
	{
		title: "First entry in an empty file",
		time:  time.Date(2021, 5, 14, 0, 0, 0, 0, time.UTC),
		index: 3,
		text:  "",
		newText: `### Fri 2021/05/14

* First entry in an empty file | 00:00/`,
	},
	{
		title: "First entry in a new day",
		time:  time.Date(2021, 5, 17, 11, 15, 0, 0, time.UTC),
		index: 10,
		text: `
## Today

If i can only do 1 thing today it would be: Watch new episode of the madisonian

---

### Fri 2021/05/14

* #code Write a program to randomly select where to eat | 11:15/11:45
  * Call API to pull lists of food places
`,
		newText: `
## Today

If i can only do 1 thing today it would be: Watch new episode of the madisonian

---

### Mon 2021/05/17

* First entry in a new day | 11:15/

### Fri 2021/05/14

* #code Write a program to randomly select where to eat | 11:15/11:45
  * Call API to pull lists of food places
`,
	},
	{
		title: "Second entry in same day",
		time:  time.Date(2021, 5, 14, 12, 30, 0, 0, time.UTC),
		index: 11,
		text: `
## Today

If i can only do 1 thing today it would be: Watch new episode of the madisonian

---

### Fri 2021/05/14

* First entry | 11:15/11:45`,
		newText: `
## Today

If i can only do 1 thing today it would be: Watch new episode of the madisonian

---

### Fri 2021/05/14

* First entry | 11:15/11:45
* Second entry in same day | 12:30/`,
	},
}

func TestAddEntry(t *testing.T) {

	split := cmpopts.AcyclicTransformer("Lines", func(s string) []string {
		return strings.Split(s, "\n")
	})

	for _, tt := range AddEntryTests {
		m := Manager{Notes: FromText(tt.text)}
		index := m.AddEntry(tt.title, tt.time)
		diff := cmp.Diff(tt.newText, m.Notes.String(), split)

		if index != tt.index {
			t.Errorf("m.AddEntry(%q, %v) = (%d) want (%d)\n%s", tt.title, tt.time, index, tt.index, m.Notes.StringWithLineNumbers())
		}

		if diff != "" {
			t.Errorf("m.AddEntry(%q, %v)\n notes content  mismatch (-want +got):\n%s\n--- Got: \n%s", tt.title, tt.time, diff, m.Notes.StringWithLineNumbers())

		}
	}
}

func TestFloorTime(t *testing.T) {
	tests := []struct {
		t    time.Time
		want time.Time
	}{
		{time.Date(2021, 5, 14, 0, 0, 0, 0, time.UTC), time.Date(2021, 5, 14, 0, 0, 0, 0, time.UTC)},
		{time.Date(2021, 5, 14, 10, 0, 0, 0, time.UTC), time.Date(2021, 5, 14, 10, 0, 0, 0, time.UTC)},
		{time.Date(2021, 5, 14, 10, 15, 0, 0, time.UTC), time.Date(2021, 5, 14, 10, 15, 0, 0, time.UTC)},
		{time.Date(2021, 5, 14, 10, 12, 0, 0, time.UTC), time.Date(2021, 5, 14, 10, 10, 0, 0, time.UTC)},
		{time.Date(2021, 5, 14, 10, 13, 0, 0, time.UTC), time.Date(2021, 5, 14, 10, 10, 0, 0, time.UTC)},
		{time.Date(2021, 5, 14, 10, 17, 0, 0, time.UTC), time.Date(2021, 5, 14, 10, 15, 0, 0, time.UTC)},
		{time.Date(2021, 5, 14, 10, 19, 0, 0, time.UTC), time.Date(2021, 5, 14, 10, 15, 0, 0, time.UTC)},
	}

	for _, tt := range tests {
		if got := floorTime(tt.t, time.Minute*5); got != tt.want {
			t.Errorf("floorTime(%v) = (%v) want (%v)", tt.t, got, tt.want)
		}
	}
}

func TestFormatEntry(t *testing.T) {
	tests := []struct {
		s    string
		t    time.Time
		want string
	}{
		{"", time.Date(2021, 5, 14, 10, 13, 13, 0, time.UTC), "*  | 10:10/"},
		{"Council of Elrond", time.Date(2021, 5, 14, 10, 13, 13, 0, time.UTC), "* Council of Elrond | 10:10/"},
		{"Council of Elrond #meeting", time.Date(2021, 5, 14, 10, 13, 13, 0, time.UTC), "* Council of Elrond #meeting | 10:10/"},
		{"Council of Elrond meeting#", time.Date(2021, 5, 14, 10, 13, 13, 0, time.UTC), "* Council of Elrond #meeting | 10:10/"},
	}

	for _, tt := range tests {
		if got := formatEntry(tt.s, tt.t); got != tt.want {
			t.Errorf("formatEntry(%q, %v) = (%q) want (%q)", tt.s, tt.t, got, tt.want)
		}
	}
}
