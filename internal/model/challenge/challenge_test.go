package challenge

import (
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/parse"
)

func TestParseAndFormatTimestamp(t *testing.T) {
	type dataFormat struct {
		input        string
		wantDuration time.Duration
		wantString   string
	}

	data := []dataFormat{
		{
			input:        "23:2.999",
			wantDuration: time.Duration(23)*time.Minute + time.Duration(2)*time.Second + time.Duration(999)*time.Millisecond,
			wantString:   "23:02.999",
		},
		{
			input:        "12:34.567",
			wantDuration: time.Duration(12)*time.Minute + time.Duration(34)*time.Second + time.Duration(567)*time.Millisecond,
			wantString:   "12:34.567",
		},
		{
			input:        "1:23.45",
			wantDuration: time.Duration(1)*time.Minute + time.Duration(23)*time.Second + time.Duration(450)*time.Millisecond,
			wantString:   " 1:23.450",
		},
		{
			input:        "1:2.3",
			wantDuration: time.Duration(1)*time.Minute + time.Duration(2)*time.Second + time.Duration(300)*time.Millisecond,
			wantString:   " 1:02.300",
		},
	}

	for _, tc := range data {
		gotDuration, err := parse.Timestamp(tc.input)
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}

		if gotDuration != tc.wantDuration {
			t.Errorf("got %s want %s", gotDuration, tc.wantDuration)
		}

		got := formatDuration(gotDuration)

		if got != tc.wantString {
			t.Errorf("got [%s] want [%s]", got, tc.wantString)
		}
	}
}

func TestSortUser(t *testing.T) {
	users := []string{"Bob", "Alice", "Carol"}
	slices.SortFunc(users, func(a, b string) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	})
	want := []string{"Alice", "Bob", "Carol"}

	if !reflect.DeepEqual(users, want) {
		t.Errorf("got %v want %v", users, want)
	}
}
