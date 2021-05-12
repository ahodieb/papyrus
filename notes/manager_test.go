package notes

import (
	"testing"
	"time"
)

func TestPositionOn_GetsPosition(t *testing.T) {

	m := Manager{Notes: NotesFile{Lines: []string{
		"### Fri 2021/05/14",
		"Entry 4",
		"",
		"### Thu 2021/05/13",
		"Entry 3",
		"",
		"### Wed 2021/05/12",
		"Entry 2",
		"",
		"### Tue 2021/05/11",
		"Entry 1",
		"",
	}}}

	d := time.Date(2021, 5, 12, 0, 0, 0, 0, time.UTC)
	want := 6
	got, _ := m.positionOn(d)
	if got != want {
		t.Errorf("Wanted %d got %d", want, got)
	}
}

func TestPositionBefore_GetsPosition(t *testing.T) {

	m := Manager{Notes: NotesFile{Lines: []string{
		"### Fri 2021/05/14",
		"Entry 4",
		"",
		"### Thu 2021/05/13",
		"Entry 3",
		"",
		"### Wed 2021/05/12",
		"Entry 2",
		"",
		"### Tue 2021/05/11",
		"Entry 1",
		"",
	}}}

	d := time.Date(2021, 5, 13, 0, 0, 0, 0, time.UTC)
	want := 6
	got, _ := m.positionBefore(d)
	if got != want {
		t.Errorf("Wanted %d got %d", want, got)
	}

}

func TestFormatDate(t *testing.T) {
	d := time.Date(2021, 5, 12, 0, 0, 0, 0, time.UTC)
	want := "### Wed 2021/05/12"

	got := FormatDate(d)

	if got != want {
		t.Errorf("Got %s wanted %s", got, want)
	}
}

func TestDatePattern(t *testing.T) {
	if !DATE_PATTERN.MatchString("### Wed 2021/05/12") {
		t.Error("Pattern does not match line")
	}
}
