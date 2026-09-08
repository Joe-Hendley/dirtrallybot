package render_test

import (
	"strings"
	"testing"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot/render"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
	"github.com/stretchr/testify/assert"
)

func TestStageConfig(t *testing.T) {
	t.Run("empty config is all random", func(t *testing.T) {
		out := render.StageConfig(challenge.Config{Game: game.DR2})

		assert.Equal(t, "Location: 🎲 Random\nStage: 🎲 Random\nWeather: 🎲 Random", out)
	})

	t.Run("chosen location and weather are named", func(t *testing.T) {
		loc := location.SCO
		wet := weather.WET
		out := render.StageConfig(challenge.Config{Game: game.DR2, Location: &loc, Weather: &wet})

		assert.Contains(t, out, "Location: "+render.Flag(location.SCO)+" Scotland")
		assert.Contains(t, out, "Weather: "+render.WeatherEmoji(weather.WET)+" Wet")
	})

	t.Run("single-weather location hints the likely weather", func(t *testing.T) {
		loc := location.MCO // dry only
		out := render.StageConfig(challenge.Config{Game: game.DR2, Location: &loc})

		assert.Contains(t, out, "probably Dry though")
	})

	t.Run("no unit-separator control characters leak into output", func(t *testing.T) {
		loc := location.WAL
		out := render.StageConfig(challenge.Config{Game: game.DR2, Location: &loc})

		assert.NotContains(t, out, "\x1f")
	})
}

func TestChallenge(t *testing.T) {
	walStages := stage.AtLocation(location.WAL)
	someCar := car.InClass(class.H3, game.DR2)[0]
	c := challenge.NewChallenge(walStages[0], weather.DRY, someCar, nil)

	out := render.Challenge(c)

	assert.Contains(t, out, walStages[0].Name())
	assert.Contains(t, out, someCar.Name())
	assert.True(t, strings.HasSuffix(out, "\n"))
}
