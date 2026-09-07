package challenge

import (
	"strings"
	"testing"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/drivetrain"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/buildersession"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyLocation(t *testing.T) {
	want := location.SCO

	config := challenge.Config{Game: game.DR2}

	config = applyLocation(config, strings.ToLower(want.String()))

	if assert.NotNil(t, config.Location) {
		assert.Equal(t, want, *config.Location)
	}
}

func componentCustomID(gameID, component string) string {
	return strings.Join([]string{ChallengeID, gameID, component}, idFieldDelimiter)
}

// The handler's game ID fragments must stay in step with the game package, which
// owns the slug <-> Model mapping.
func TestGameIDFragmentsMatchGamePackage(t *testing.T) {
	assert.Equal(t, DR2ID, game.DR2.ID())
	assert.Equal(t, WRCID, game.WRC.ID())
}

func componentInteraction(customID string, values []string, builderID string) *discordgo.InteractionCreate {
	return &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionMessageComponent,
			Data: discordgo.MessageComponentInteractionData{
				CustomID: customID,
				Values:   values,
			},
			Message: &discordgo.Message{ID: builderID},
		},
	}
}

func optionByValue(menu discordgo.SelectMenu, value string) (discordgo.SelectMenuOption, bool) {
	for _, option := range menu.Options {
		if option.Value == value {
			return option, true
		}
	}
	return discordgo.SelectMenuOption{}, false
}

func defaultOption(menu discordgo.SelectMenu) (discordgo.SelectMenuOption, bool) {
	for _, option := range menu.Options {
		if option.Default {
			return option, true
		}
	}
	return discordgo.SelectMenuOption{}, false
}

func TestBuildMenus(t *testing.T) {
	t.Run("location menu offers random plus every location and defaults to random", func(t *testing.T) {
		menu := buildLocationsMenu(challenge.Config{Game: game.DR2})

		assert.Equal(t, componentCustomID(DR2ID, locationID), menu.CustomID)
		assert.Len(t, menu.Options, 1+len(location.List(game.DR2)))
		def, ok := defaultOption(menu)
		require.True(t, ok)
		assert.Equal(t, RandomID, def.Value)
		assert.False(t, menu.Disabled)
	})

	t.Run("location menu defaults to the selected location", func(t *testing.T) {
		loc := location.SCO
		menu := buildLocationsMenu(challenge.Config{Game: game.DR2, Location: &loc})

		selected, ok := optionByValue(menu, strings.ToLower(location.SCO.String()))
		require.True(t, ok)
		assert.True(t, selected.Default)
		random, ok := optionByValue(menu, RandomID)
		require.True(t, ok)
		assert.False(t, random.Default)
	})

	t.Run("distance menu includes the 16 sector option only for WRC", func(t *testing.T) {
		reallyLong := strings.ToLower(stage.ReallyLong.String())

		_, inDR2 := optionByValue(buildDistanceMenu(challenge.Config{Game: game.DR2}), reallyLong)
		_, inWRC := optionByValue(buildDistanceMenu(challenge.Config{Game: game.WRC}), reallyLong)

		assert.False(t, inDR2)
		assert.True(t, inWRC)
	})

	t.Run("stage menu is disabled until a location is chosen", func(t *testing.T) {
		menu := buildStageMenu(challenge.Config{Game: game.DR2})

		assert.Len(t, menu.Options, 1)
		assert.True(t, menu.Disabled)
	})

	t.Run("weather menu collapses to a single disabled option for a one-weather location", func(t *testing.T) {
		loc := location.MCO // Monte Carlo, dry only
		menu := buildWeatherMenu(challenge.Config{Game: game.DR2, Location: &loc})

		require.Len(t, menu.Options, 1)
		assert.True(t, menu.Disabled)
		assert.Equal(t, strings.ToLower(weather.DRY.String()), menu.Options[0].Value)
		assert.True(t, menu.Options[0].Default)
	})

	t.Run("class menu is limited to the chosen drivetrain", func(t *testing.T) {
		dt := drivetrain.RWD
		menu := buildClassMenu(challenge.Config{Game: game.DR2, Drivetrain: &dt})

		assert.Len(t, menu.Options, 1+len(class.WithDrivetrain(drivetrain.RWD, game.DR2)))
	})
}

func TestConfigFromInteraction(t *testing.T) {
	t.Run("first interaction seeds the config from the game in the custom ID", func(t *testing.T) {
		sessions := buildersession.New()

		got, err := configFromInteraction(sessions, componentInteraction(
			componentCustomID(DR2ID, locationID),
			[]string{strings.ToLower(location.SCO.String())},
			"builder-1",
		))

		require.NoError(t, err)
		assert.Equal(t, game.DR2, got.Game)
		require.NotNil(t, got.Location)
		assert.Equal(t, location.SCO, *got.Location)
	})

	t.Run("selections accumulate across interactions via the session", func(t *testing.T) {
		sessions := buildersession.New()

		_, err := configFromInteraction(sessions, componentInteraction(
			componentCustomID(DR2ID, locationID),
			[]string{strings.ToLower(location.WAL.String())},
			"b",
		))
		require.NoError(t, err)

		got, err := configFromInteraction(sessions, componentInteraction(
			componentCustomID(DR2ID, weatherID),
			[]string{strings.ToLower(weather.WET.String())},
			"b",
		))
		require.NoError(t, err)

		require.NotNil(t, got.Location)
		assert.Equal(t, location.WAL, *got.Location)
		require.NotNil(t, got.Weather)
		assert.Equal(t, weather.WET, *got.Weather)
	})

	t.Run("changing location clears a previously selected stage", func(t *testing.T) {
		sessions := buildersession.New()
		walStages := stage.AtLocation(location.WAL)

		_, err := configFromInteraction(sessions, componentInteraction(
			componentCustomID(DR2ID, locationID),
			[]string{strings.ToLower(location.WAL.String())},
			"b",
		))
		require.NoError(t, err)

		withStage, err := configFromInteraction(sessions, componentInteraction(
			componentCustomID(DR2ID, stageID),
			[]string{strings.ToLower(walStages[0].Name())},
			"b",
		))
		require.NoError(t, err)
		require.NotNil(t, withStage.Stage)

		relocated, err := configFromInteraction(sessions, componentInteraction(
			componentCustomID(DR2ID, locationID),
			[]string{strings.ToLower(location.SCO.String())},
			"b",
		))
		require.NoError(t, err)
		assert.Nil(t, relocated.Stage)
	})

	t.Run("submit button preserves the accumulated config", func(t *testing.T) {
		sessions := buildersession.New()

		_, err := configFromInteraction(sessions, componentInteraction(
			componentCustomID(DR2ID, locationID),
			[]string{strings.ToLower(location.SCO.String())},
			"b",
		))
		require.NoError(t, err)

		got, err := configFromInteraction(sessions, componentInteraction(
			componentCustomID(DR2ID, SubmitLocationAndStageID),
			nil,
			"b",
		))
		require.NoError(t, err)
		require.NotNil(t, got.Location)
		assert.Equal(t, location.SCO, *got.Location)
	})

	t.Run("errors on malformed custom ID", func(t *testing.T) {
		sessions := buildersession.New()

		_, err := configFromInteraction(sessions, componentInteraction("challenge-dr2", []string{"x"}, "b"))

		require.Error(t, err)
	})

	t.Run("errors on unknown game", func(t *testing.T) {
		sessions := buildersession.New()

		_, err := configFromInteraction(sessions, componentInteraction(
			componentCustomID("rbr", locationID),
			[]string{"x"},
			"b",
		))

		require.Error(t, err)
	})

	t.Run("errors when a non-submit component carries no value", func(t *testing.T) {
		sessions := buildersession.New()

		_, err := configFromInteraction(sessions, componentInteraction(
			componentCustomID(DR2ID, locationID),
			nil,
			"b",
		))

		require.Error(t, err)
	})
}
