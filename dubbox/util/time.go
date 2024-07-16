package util

import (
	"errors"
	"strconv"
	"time"
)

func ParseDubboxDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, errors.New("time string can not be empty")
	}

	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return time.Duration(i) * time.Millisecond, nil
	}
	return time.ParseDuration(s)
}
