package markdown

import (
	"testing"
	"time"
)

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
