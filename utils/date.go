package utils

import (
	"strconv"
	"time"
)

func Date(t time.Time) (int, int, int, int) {
	y, m, d := t.Date()
	yStr := strconv.Itoa(y)
	mStr := strconv.Itoa(int(m))
	dStr := strconv.Itoa(d)
	if m < 10 {
		mStr = "0" + mStr
	}
	if d < 10 {
		dStr = "0" + dStr
	}

	dateStr := yStr + mStr + dStr
	date, _ := strconv.Atoi(dateStr)
	return date, y, int(m), d
}
