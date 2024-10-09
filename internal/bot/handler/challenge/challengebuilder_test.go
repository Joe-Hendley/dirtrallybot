package challenge

import (
	"strings"
	"testing"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/stretchr/testify/assert"
)

func TestApplyLocation(t *testing.T) {
	want := location.SCO

	config := challenge.Config{Game: game.DR2}

	config = applyLocation(config, strings.ToLower(want.String()))

	if assert.NotNil(t, config.Location) {
		assert.Equal(t, want, *config.Location)
	}
}
