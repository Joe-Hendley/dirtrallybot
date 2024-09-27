package challenge

import (
	"log/slog"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/bwmarrin/discordgo"
)

func buildChallengeConfig(session *discordgo.Session, interaction *discordgo.InteractionCreate) (challenge.Config, error) {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Custom DR2 hasn't been implemented yet! Sorry!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		slog.Error("Create Custom DR2 Challenge", "err", err)
		return challenge.Config{}, err
	}

	return challenge.Config{}, nil
}
