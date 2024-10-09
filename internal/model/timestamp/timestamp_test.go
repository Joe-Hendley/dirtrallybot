package timestamp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	type dataFormat struct {
		input string
		want  time.Duration
	}

	data := []dataFormat{
		{
			input: "23:2.999",
			want:  buildTimestamp(23, 2, 999),
		},
		{
			input: "12:34.567",
			want:  buildTimestamp(12, 34, 567),
		},
		{
			input: "1:23.45",
			want:  buildTimestamp(1, 23, 450),
		},
		{
			input: "1:2.3",
			want:  buildTimestamp(1, 2, 300),
		},
	}

	for _, tc := range data {
		t.Run(tc.input, func(t *testing.T) {
			got, err := Parse(tc.input)
			if assert.Nil(t, err) {
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	type dataFormat struct {
		input time.Duration
		want  string
	}

	data := []dataFormat{
		{
			input: buildTimestamp(23, 2, 999),
			want:  "23:02.999",
		},
		{
			input: buildTimestamp(12, 34, 567),
			want:  "12:34.567",
		},
		{
			input: buildTimestamp(1, 23, 450),
			want:  " 1:23.450",
		},
		{
			input: buildTimestamp(1, 2, 300),
			want:  " 1:02.300",
		},
	}

	for _, tc := range data {
		t.Run(tc.want, func(t *testing.T) {
			got := Format(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}
