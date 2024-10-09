package handler

import (
	"strings"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot/discord"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handler/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handler/completion"
	"github.com/Joe-Hendley/dirtrallybot/internal/model"
	"github.com/bwmarrin/discordgo"
)

func ApplicationCommand(session discord.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch interaction.ApplicationCommandData().Name {
	case challenge.NewDR2ChallengeID:
		challenge.HandleNewDR2Challenge(session, interaction)
	case challenge.NewWRCChallengeID:
		challenge.HandleNewWRCChallenge(session, interaction)
	}
}

func InteractionMessageComponent(store model.Store, session discord.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionMessageComponent {
		return
	}

	customID := interaction.MessageComponentData().CustomID

	switch {
	case customID == challenge.CompletedID:
		completion.HandleDisplayEntryModal(session, interaction)
	case customID == challenge.TimesID:
		completion.HandleDisplayTimes(store, session, interaction)
	case strings.HasPrefix(customID, challenge.DR2ChallengePrefix) || strings.HasPrefix(customID, challenge.WRCChallengePrefix):
		challenge.HandleChallengeBuilderInteraction(store, session, interaction)
	}
}

func ModalSubmit(store model.Store, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionModalSubmit {
		return
	}

	data := interaction.ModalSubmitData()

	if strings.HasPrefix(data.CustomID, completion.SubmitCompletionPrefix) {
		completion.HandleSubmitModal(store, session, interaction)
	}
}
