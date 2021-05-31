package date

import "time"

const (
	dateFormat = "2006-01-02T15:04:05Z"
)

// GetNow ... returns the current UTC time
func GetNow() time.Time {
	// get current datetime in UTC
	return time.Now().UTC()
}

// GetNowString ... returns the current UTC time in the selected string format
func GetNowString() string {
	// format according to setting
	return GetNow().Format(dateFormat)
}
