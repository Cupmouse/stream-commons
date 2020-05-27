package json

import (
	"fmt"
	"time"
)

var durationBaseTime time.Time

func parseTimestamp(timestamp *string) (*string, error) {
	if timestamp != nil {
		timestampTime, serr := time.Parse(time.RFC3339Nano, *timestamp)
		if serr != nil {
			return nil, serr
		}
		result := fmt.Sprintf("%d", timestampTime.UnixNano())
		return &result, nil
	}
	return nil, nil
}

func parseDuration(duration *string) (*string, error) {
	if duration != nil {
		durationTime, serr := time.Parse(time.RFC3339Nano, *duration)
		if serr != nil {
			return nil, serr
		}
		result := fmt.Sprintf("%d", durationTime.Sub(durationBaseTime).Nanoseconds())
		return &result, nil
	}
	return nil, nil
}
