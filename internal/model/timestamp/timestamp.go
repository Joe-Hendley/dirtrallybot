package timestamp

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func Build(minutes, seconds, milliseconds int) time.Duration {
	return time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second + time.Duration(milliseconds)*time.Millisecond
}

func Parse(s string) (time.Duration, error) {
	minuteString, remainder, ok := strings.Cut(s, ":")
	if !ok {
		return 0, fmt.Errorf("invalid timestamp %s", s)
	}

	var err error
	minutes, err := strconv.Atoi(minuteString)
	if err != nil {
		return 0, fmt.Errorf("invalid timestamp %s: %v", s, err)
	}

	if minutes < 0 || minutes > 100 {
		return 0, fmt.Errorf("invalid timestamp %s", s)
	}

	secondString, msString, _ := strings.Cut(remainder, ".")
	seconds, err := strconv.Atoi(secondString)
	if err != nil {
		return 0, fmt.Errorf("invalid timestamp %s: %v", s, err)
	}

	if seconds < 0 || seconds >= 60 {
		return 0, fmt.Errorf("invalid timestamp %s", s)
	}

	milliseconds := 0
	if msString != "" {
		if len(msString) > 3 {
			msString = msString[:3]
		}
		padded := msString + strings.Repeat("0", 3-len(msString))

		milliseconds, err = strconv.Atoi(padded)
		if err != nil {
			return 0, fmt.Errorf("invalid timestamp %s: %v", s, err)
		}
	}

	if milliseconds < 0 || milliseconds >= 1000 {
		return 0, fmt.Errorf("invalid timestamp %s", s)
	}

	return Build(minutes, seconds, milliseconds), nil
}

// Format renders a duration as mm:ss.mss, the inverse of Parse.
func Format(duration time.Duration) string {
	var (
		minutes      = duration.Truncate(time.Minute)
		seconds      = (duration - minutes).Truncate(time.Second)
		milliseconds = (duration - minutes - seconds).Truncate(time.Millisecond)
	)

	return fmt.Sprintf("%d:%02d.%03d",
		int64(minutes.Minutes()),
		int64(seconds.Seconds()),
		milliseconds.Milliseconds(),
	)
}
