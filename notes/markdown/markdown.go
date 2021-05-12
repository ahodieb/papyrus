package markdown

import "time"

func FormatDate(t time.Time) string {
	return t.Format("### Mon 2006/01/02")
}
