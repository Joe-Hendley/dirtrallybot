package handler

import (
	"context"
	"strings"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot/discord"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handler/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handler/completion"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handler/feedback"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/port"
	"github.com/bwmarrin/discordgo"
)

func ApplicationCommand(session discord.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch interaction.ApplicationCommandData().Name {
	case challenge.NewDR2ChallengeID, challenge.NewWRCChallengeID:
		challenge.HandleNewChallenge(session, interaction)
	}
}

func InteractionMessageComponent(ctx context.Context, sessions challenge.SessionStore, store port.Store, session discord.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionMessageComponent {
		return
	}

	customID := interaction.MessageComponentData().CustomID

	switch {
	case customID == challenge.DisplayCompletionModalID:
		completion.HandleDisplayEntryModal(session, interaction)
	case customID == challenge.DisplayTimesID:
		completion.HandleDisplayTimes(ctx, store, session, interaction)
	case customID == feedback.GoodID, customID == feedback.BadID:
		feedback.HandleVote(ctx, store, session, interaction)
	case strings.HasPrefix(customID, challenge.ChallengeID):
		challenge.HandleChallengeBuilderInteraction(ctx, sessions, store, session, interaction)
	}
}

func ModalSubmit(ctx context.Context, store port.Store, session discord.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionModalSubmit {
		return
	}

	data := interaction.ModalSubmitData()

	if strings.HasPrefix(data.CustomID, completion.SubmitCompletionPrefix) {
		completion.HandleSubmitModal(ctx, store, session, interaction)
	}
}
