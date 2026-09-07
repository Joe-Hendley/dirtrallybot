package challenge

import (
	"strings"
	"testing"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
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

// toWireComponents mirrors how discordgo unmarshals message components received
// from Discord: action rows and their children arrive as pointers. The builder
// helpers return value types, so tests must convert before feeding their output
// back into the parsing functions.
func toWireComponents(components []discordgo.MessageComponent) []discordgo.MessageComponent {
	wire := make([]discordgo.MessageComponent, 0, len(components))
	for _, row := range components {
		actionsRow := row.(discordgo.ActionsRow)
		children := make([]discordgo.MessageComponent, 0, len(actionsRow.Components))
		for _, child := range actionsRow.Components {
			switch c := child.(type) {
			case discordgo.SelectMenu:
				menu := c
				children = append(children, &menu)
			case discordgo.Button:
				button := c
				children = append(children, &button)
			default:
				panic("unexpected component type in test helper")
			}
		}
		wire = append(wire, &discordgo.ActionsRow{Components: children})
	}
	return wire
}

func messageComponentInteraction(customID string, values []string, content string, components []discordgo.MessageComponent) *discordgo.InteractionCreate {
	return &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			Type: discordgo.InteractionMessageComponent,
			Data: discordgo.MessageComponentInteractionData{
				CustomID: customID,
				Values:   values,
			},
			Message: &discordgo.Message{
				Content:    content,
				Components: components,
			},
		},
	}
}

func componentCustomID(gameID, component string) string {
	return strings.Join([]string{ChallengeID, gameID, component}, idFieldDelimiter)
}

func TestBuildStageConfigFromInteraction(t *testing.T) {
	// A stored location is recovered from the rendered menus, and the weather
	// the user just picked is applied on top.
	t.Run("recovers stored location and applies chosen weather", func(t *testing.T) {
		loc := location.SCO
		start := challenge.Config{Game: game.DR2, Location: &loc}

		interaction := messageComponentInteraction(
			componentCustomID(DR2ID, weatherID),
			[]string{strings.ToLower(weather.WET.String())},
			"",
			toWireComponents(buildChallengeLocationMessageComponents(start)),
		)

		got, err := buildStageConfigFromInteraction(interaction)

		require.NoError(t, err)
		require.NotNil(t, got.Location)
		assert.Equal(t, location.SCO, *got.Location)
		require.NotNil(t, got.Weather)
		assert.Equal(t, weather.WET, *got.Weather)
	})

	// The submit button carries no Values; stored selections must survive.
	t.Run("submit button preserves stored selections", func(t *testing.T) {
		loc := location.WAL
		wet := weather.WET
		start := challenge.Config{Game: game.DR2, Location: &loc, Weather: &wet}

		interaction := messageComponentInteraction(
			componentCustomID(DR2ID, SubmitLocationAndStageID),
			nil,
			"",
			toWireComponents(buildChallengeLocationMessageComponents(start)),
		)

		got, err := buildStageConfigFromInteraction(interaction)

		require.NoError(t, err)
		require.NotNil(t, got.Location)
		assert.Equal(t, location.WAL, *got.Location)
		require.NotNil(t, got.Weather)
		assert.Equal(t, weather.WET, *got.Weather)
	})

	t.Run("errors on malformed custom ID", func(t *testing.T) {
		interaction := messageComponentInteraction("challenge-dr2", []string{"x"}, "", nil)

		_, err := buildStageConfigFromInteraction(interaction)

		require.Error(t, err)
	})

	t.Run("errors on unknown game", func(t *testing.T) {
		interaction := messageComponentInteraction(
			componentCustomID("rbr", locationID),
			[]string{"x"},
			"",
			nil,
		)

		_, err := buildStageConfigFromInteraction(interaction)

		require.Error(t, err)
	})
}

func TestBuildCarConfigFromInteraction(t *testing.T) {
	// The stage half of the config is not carried in component custom IDs on the
	// car screen; it is recovered by parsing the rendered message content.
	t.Run("recovers stage selections from message content and applies chosen class", func(t *testing.T) {
		loc := location.WAL
		stages := stage.AtLocation(location.WAL)
		chosenStage := stages[1] // "Sweet Lamb" - two words, avoids the distance-parsing heuristic
		start := challenge.Config{Game: game.DR2, Location: &loc, Stage: &chosenStage}

		interaction := messageComponentInteraction(
			componentCustomID(DR2ID, classID),
			[]string{strings.ToLower(class.H3.String())},
			"DR2 Challenge Builder\n"+start.FancyStageString(),
			toWireComponents(buildChallengeCarMessageComponents(start)),
		)

		got, err := buildCarConfigFromInteraction(interaction)

		require.NoError(t, err)
		require.NotNil(t, got.Location)
		assert.Equal(t, location.WAL, *got.Location)
		require.NotNil(t, got.Stage)
		assert.Equal(t, "Sweet Lamb", got.Stage.Name())
		require.NotNil(t, got.Class)
		assert.Equal(t, class.H3, *got.Class)
	})

	// The "Random <distance>" pseudo-stage line is reverse-engineered back into a
	// distance by splitting on spaces. This captures that heuristic.
	t.Run("parses Random distance pseudo-stage from content", func(t *testing.T) {
		loc := location.WAL
		distance := stage.Long
		start := challenge.Config{Game: game.DR2, Location: &loc, Distance: &distance}

		interaction := messageComponentInteraction(
			componentCustomID(DR2ID, drivetrainID),
			[]string{"random"},
			"DR2 Challenge Builder\n"+start.FancyStageString(),
			toWireComponents(buildChallengeCarMessageComponents(start)),
		)

		got, err := buildCarConfigFromInteraction(interaction)

		require.NoError(t, err)
		require.NotNil(t, got.Distance)
		assert.Equal(t, stage.Long, *got.Distance)
	})

	t.Run("errors on malformed custom ID", func(t *testing.T) {
		interaction := messageComponentInteraction("challenge-dr2", []string{"x"}, "", nil)

		_, err := buildCarConfigFromInteraction(interaction)

		require.Error(t, err)
	})
}
