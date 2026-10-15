//go:build !production

package github

import "time"

func SetClockForTest(now func() time.Time, sleep func(time.Duration)) {
	timeNow = now
	timeSleep = sleep
}

func ResetClockForTest() {
	timeNow = time.Now
	timeSleep = time.Sleep
}

func SetTimeoutForTest(t time.Duration) {
	pollTimeout = t
}
