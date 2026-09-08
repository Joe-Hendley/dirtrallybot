package feedback_test

import (
	"context"
	"testing"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handler/feedback"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/weather"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/memorystore"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// responderFunc adapts a function to discord.InteractionResponder.
type responderFunc func(*discordgo.Interaction, *discordgo.InteractionResponse, ...discordgo.RequestOption) error

func (f responderFunc) InteractionRespond(i *discordgo.Interaction, r *discordgo.InteractionResponse, o ...discordgo.RequestOption) error {
	return f(i, r, o...)
}

func voteInteraction(customID, challengeID, userID string) *discordgo.InteractionCreate {
	return &discordgo.InteractionCreate{Interaction: &discordgo.Interaction{
		Type:    discordgo.InteractionMessageComponent,
		Data:    discordgo.MessageComponentInteractionData{CustomID: customID},
		Member:  &discordgo.Member{User: &discordgo.User{ID: userID}},
		Message: &discordgo.Message{ID: challengeID},
	}}
}

func captureResponse(into **discordgo.InteractionResponse) responderFunc {
	return func(_ *discordgo.Interaction, r *discordgo.InteractionResponse, _ ...discordgo.RequestOption) error {
		*into = r
		return nil
	}
}

func TestHandleVoteRecordsAgainstTheChallenge(t *testing.T) {
	ctx := context.Background()
	store := memorystore.New()
	require.NoError(t, store.PutChallenge(ctx, "c1", challenge.NewChallenge(
		stage.New("Sweet Lamb", location.WAL, stage.Long), weather.DRY, car.Model{}, nil, nil)))

	var response *discordgo.InteractionResponse
	feedback.HandleVote(ctx, store, captureResponse(&response), voteInteraction(feedback.BadID, "c1", "alice"))

	require.NotNil(t, response)
	assert.Equal(t, discordgo.MessageFlagsEphemeral, response.Data.Flags)

	stored, err := store.GetChallenge(ctx, "c1")
	require.NoError(t, err)
	assert.Equal(t, -1, stored.Score())
}

func TestHandleVoteOnMissingChallengeApologises(t *testing.T) {
	var response *discordgo.InteractionResponse
	feedback.HandleVote(context.Background(), memorystore.New(), captureResponse(&response),
		voteInteraction(feedback.GoodID, "gone", "alice"))

	require.NotNil(t, response)
	assert.Contains(t, response.Data.Content, "☹️")
}
