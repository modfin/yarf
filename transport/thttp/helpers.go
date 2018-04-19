package thttp

import "time"

func stringOr(s string, or string) string {
	if s == "" {
		return or
	}
	return s
}

func durationOr(t time.Duration, or time.Duration) time.Duration {
	if t == 0 {
		return or
	}
	return t
}
