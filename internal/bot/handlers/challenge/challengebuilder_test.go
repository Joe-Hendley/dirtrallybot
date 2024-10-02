package challenge

import (
	"strings"
	"testing"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
)

func TestApplyLocation(t *testing.T) {
	wantLocation := location.SCO

	config := challenge.Config{}

	config = applyLocation(config, strings.ToLower(wantLocation.String()))

	if config.Location == nil {
		t.Errorf("nil location")
	} else if *config.Location != wantLocation {
		t.Errorf("got %s want %s", config.Location.String(), wantLocation.String())
	}
}
