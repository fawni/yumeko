package common

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/logrusorgru/aurora/v3"
)

func Fatal(err interface{}) {
	fmt.Println(aurora.Red(err))
	os.Exit(1)
}

func FmtDuration(d time.Duration) string {
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	if h == 0 {
		return fmt.Sprintf("%02d:%02d", m, s)
	}
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func FmtTime(t time.Time) string {
	ap := "AM"
	hour := t.Hour()
	if hour > 11 {
		hour -= 12
		ap = "PM"
	}
	if hour == 0 {
		hour = 12
	}
	return fmt.Sprintf("%02d:%02d:%02d %s, %02d %s %d", hour, t.Minute(), t.Second(), ap, t.Day(), t.Month(), t.Year())
}

func SI(input float64, object string) string {
	value, prefix := humanize.ComputeSI(input)
	return humanize.Ftoa(value) + strings.ToUpper(prefix) + " " + object
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
