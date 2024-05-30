package util

import (
	"time"
)

func unmarshalTime(layout string, data []byte) (time.Time, error) {
	return time.Parse(layout, string(data))
}
