package timeago

import (
	"github.com/xeonx/timeago"
	"time"
)

var gTimeAgoConfig = timeago.Config{
	PastSuffix:   " ago",
	FuturePrefix: "in ",
	Periods: []timeago.FormatPeriod{
		{time.Second, "a sec", "%d sec"},
		{time.Minute, "a min", "%d min"},
		{time.Hour, "an hour", "%d hrs"},
		{timeago.Day, "one day", "%d days"},
		{timeago.Month, "one mon", "%d mons"},
		{timeago.Year, "one year", "%d yrs"},
	},
	Zero:          "moments",
	Max:           99 * timeago.Year,
	DefaultLayout: "2006-01-02",
}

func TimeAgo(t time.Time) string {
	return gTimeAgoConfig.Format(t)
}
