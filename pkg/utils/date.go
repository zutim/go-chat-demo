package utils

import (
	"chat/pkg/constant"
	"time"
)

// GetTime return current local time
func GetTime() time.Time {
	return time.Now().In(GetLocalTimeZone())
}

// GetLocalTimeZone
func GetLocalTimeZone() *time.Location {
	return time.FixedZone("CST", 8*3600) // UTC+8
}

// GetDateStr return current date string,eg: 2019-12-30
func GetDateStr() string {
	return GetTime().Format(constant.DateFormat)
}

// GetTimeStr return current time string,eg:2019-12-30 22:00:00
func GetTimeStr() string {
	return GetTime().Format(constant.TimeFormat)
}

// GetTimeStamp return current timestamp
func GetTimeStamp() int64 {
	return GetTime().Unix()
}
