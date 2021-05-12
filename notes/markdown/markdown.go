package markdown

import (
	"regexp"
	"strings"
	"time"
)

var DATE_PATTERN, _ = regexp.Compile("### [a-zA-Z]{3} (?P<year>[0-9]{4})/(?P<month>[0-9]{2})/(?P<day>[0-9]{2})")

func FormatDate(t time.Time) string {
	return t.Format("### Mon 2006/01/02")
}

func ParseDate(s string) (time.Time, error) {
	return time.Parse("### Mon 2006/01/02", strings.TrimSpace(s))
}
