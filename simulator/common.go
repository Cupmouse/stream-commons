package simulator

import "time"

func unixMillisec(millisec int64) time.Time {
	timestamp := time.Duration(millisec) * time.Millisecond
	return time.Unix(int64(timestamp/time.Second), int64(timestamp%time.Second))
}
