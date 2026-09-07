package challenge

import (
	"strings"
	"testing"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
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
