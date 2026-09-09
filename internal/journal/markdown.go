package journal

import (
	"fmt"
	"time"
)

func formatEntry(text string, now time.Time) string {
	timestamp := now.Format("15:04")
	return fmt.Sprintf("\n\n## %s\n\n### %s", timestamp, text)
}

func formatBullet(text string, now time.Time) string {
	timestamp := now.Format("15:04")
	return fmt.Sprintf("\n* %s- %s", timestamp, text)
}

func formatDayHeader(now time.Time) string {
	return fmt.Sprintf("# %s\n\n", now.Format("January 02, 2006 (01/02/06)"))
}
