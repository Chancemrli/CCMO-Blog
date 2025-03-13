package util

import (
	"math"
	"time"
)

func Hot(ups, downs int, date time.Time) float64 {
	s := float64(ups - downs)
	order := math.Log(math.Max(math.Abs(s), 1))
	var sign float64
	if s > 0 {
		sign = 1
	} else if s == 0 {
		sign = 0
	} else {
		sign = -1
	}
	seconds := float64(date.Unix() - 1577808000)
	return math.Round(sign*order + seconds/43200)
}
