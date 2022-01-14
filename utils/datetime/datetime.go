package datetime

import (
	"math"
	"time"
)

func FormatSeconds(ts int64) string {
	return FromSeconds(ts).Format("2006-01-02 15:04:05")
}

func ParseISOStr(iso string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05-0700", iso)
}

func FormatMillis(ts int64) string {
	return time.UnixMilli(ts).Format("2006-01-02 15:04:05")
}

func Duration(before, after int64) time.Time {
	return FromSeconds(int64(math.Abs(float64(after - before))))
}

func FromSeconds(seconds int64) time.Time {
	return time.Unix(seconds, 0)
}
