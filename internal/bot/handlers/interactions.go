package handlers

import (
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handlers/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handlers/completion"
	"github.com/Joe-Hendley/dirtrallybot/internal/model"
	"github.com/bwmarrin/discordgo"
)

func ApplicationCommand(store model.Store, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch interaction.ApplicationCommandData().Name {
	case challenge.NewstageDR2DefaultID:
		challenge.HandleCreateDR2ChallengeDefault(store, session, challenge.NewInvocationFromInteractionCreate(*interaction))
	case challenge.NewstageDR2CustomID:
		challenge.HandleCreateDR2ChallengeCustom(session, interaction)
	case challenge.NewstageWRCDefaultID:
		challenge.HandleCreateWRCChallengeDefault(session, interaction)
	case challenge.NewstageWRCCustomID:
		challenge.HandleCreateWRCChallengeCustom(session, interaction)
	}
}

func InteractionMessageComponent(store model.Store, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionMessageComponent {
		return
	}

	customID := interaction.MessageComponentData().CustomID

	switch {
	case customID == challenge.CompletedID:
		handleCompletion(session, interaction)
	case customID == challenge.TimesID:
		handleDisplayTimes(store, session, interaction)
	case customID == challenge.GoodID:
		handleGood(session, interaction)
	case customID == challenge.BadID:
		handleBad(session, interaction)
	case strings.HasPrefix(customID, challenge.DR2ChallengePrefix):
		challenge.HandleCustomDR2ChallengeInteraction(store, session, interaction)
	}
}

func ModalSubmit(store model.Store, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionModalSubmit {
		return
	}

	data := interaction.ModalSubmitData()

	if strings.HasPrefix(data.CustomID, completion.SubmitCompletionPrefix) {
		completion.HandleCompletionRequest(store, session, interaction)
	}
}

func handleCompletion(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: strings.Join([]string{completion.SubmitCompletionPrefix, i.ChannelID, i.Message.ID, i.Interaction.Member.User.ID, strconv.FormatInt(time.Now().Unix(), 10)}, "_"),
			Title:    "Please input your time",
			Content:  "Please input your time",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "rawDuration",
							Label:       "Minutes:Seconds.Milliseconds",
							Style:       discordgo.TextInputShort,
							Placeholder: "12:34.567",
							Required:    true,
							MaxLength:   9, // 12:34.567
							MinLength:   5, // 3:2.1
						},
					},
				},
			},
		},
	})

	if err != nil {
		slog.Error("responding to completion interaction", "err", err)
	}
}

func handleDisplayTimes(store model.Store, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	challengeID := interaction.Message.ID
	challenge, ok := store.Get(challengeID)
	if !ok {
		slog.Warn("could not find challenge", "id", challengeID)
		return
	}

	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: challenge.FancyListCompletions(session, interaction.GuildID),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		slog.Error("display times", "err", err)
	}
}

func handleGood(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Yay!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		slog.Error("good", "err", err)
	}
}

func handleBad(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "oh no",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		slog.Error("bad", "err", err)
	}
}
