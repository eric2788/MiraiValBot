package datetime

import "time"

func Format(ts int64) string {
	return time.UnixMilli(ts * 1000).Format("2006-01-02 15:04:05")
}
