package challenge

import (
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	type dataFormat struct {
		input time.Duration
		want  string
	}

	data := []dataFormat{
		{
			input: time.Duration(23)*time.Minute + time.Duration(2)*time.Second + time.Duration(999)*time.Millisecond,
			want:  "23:02.999",
		},
	}

	for _, tc := range data {
		got := formatDuration(tc.input)

		if got != tc.want {
			t.Errorf("got %s want %s", got, tc.want)
		}
	}
}
