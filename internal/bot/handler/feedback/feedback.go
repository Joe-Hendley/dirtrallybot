// Package feedback records the 👍 / 👎 votes on a posted challenge, which feed
// the popularity tally that biases future generation.
package feedback

import (
	"context"
	"log/slog"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot/discord"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/port"
	"github.com/bwmarrin/discordgo"
)

// Button custom IDs for the challenge message's feedback buttons.
const (
	GoodID = "feedback-good"
	BadID  = "feedback-bad"
)

// HandleVote applies the clicked 👍 / 👎 to the challenge the button belongs to.
// The challenge is keyed by its message ID, as elsewhere.
func HandleVote(ctx context.Context, store port.Store, session discord.InteractionResponder, interaction *discordgo.InteractionCreate) {
	customID := interaction.MessageComponentData().CustomID

	sentiment, ok := sentimentFor(customID)
	if !ok {
		slog.Error("unknown feedback button", "custom_id", customID)
		return
	}

	challengeID := interaction.Message.ID
	userID := interaction.Member.User.ID

	err := store.RegisterVote(ctx, challengeID, challenge.NewVote(userID, sentiment))
	if err != nil {
		slog.Error("registering vote", "challenge_id", challengeID, "err", err)
		respondEphemeral(session, interaction, "☹️ couldn't record that - the challenge may have expired")
		return
	}

	respondEphemeral(session, interaction, acknowledgement(sentiment))
}

func sentimentFor(customID string) (challenge.Sentiment, bool) {
	switch customID {
	case GoodID:
		return challenge.Up, true
	case BadID:
		return challenge.Down, true
	}
	return 0, false
}

// acknowledgement is deliberately the same whether the click recorded a new vote
// or toggled an existing one off - the user knows which button they pressed.
func acknowledgement(sentiment challenge.Sentiment) string {
	switch sentiment {
	case challenge.Up:
		return "👍 logged - thanks!"
	case challenge.Down:
		return "👎 logged - thanks!"
	}
	return "logged - thanks!"
}

func respondEphemeral(session discord.InteractionResponder, interaction *discordgo.InteractionCreate, content string) {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		slog.Error("sending ephemeral response", "err", err)
	}
}
