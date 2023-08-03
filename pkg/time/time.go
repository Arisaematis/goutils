package time

import "time"

var timeTemplates = []string{
	"2006-01-02T15:04:05Z",
	// "2006-01-02 15:04:05", //常规类型
	// "2006/01/02 15:04:05",
	"2006-01-02",
	// "2006/01/02",
	// "15:04:05",
}

func StringToGoTime(tm string) time.Time {
	for i := range timeTemplates {
		t, err := time.ParseInLocation(timeTemplates[i], tm, time.UTC)
		if nil == err && !t.IsZero() {
			return t
		}
	}
	return time.Time{}
}
