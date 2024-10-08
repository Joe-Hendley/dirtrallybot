package challenge

import (
	"fmt"
	"log/slog"

	"github.com/Joe-Hendley/dirtrallybot/internal/model"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
	"github.com/Joe-Hendley/dirtrallybot/internal/randomiser"
	"github.com/bwmarrin/discordgo"
)

const (
	CompletedID = "newstage_completed"
	TimesID     = "newstage_time"
	GoodID      = "newstage_good"
	BadID       = "newstage_bad"

	NewstageDR2CustomID = "newstage-custom"
	NewstageWRCCustomID = "newstagewrc-custom"
)

func NewInvocationFromMessageCreate(m discordgo.MessageCreate) invocation {
	return invocation{
		id:        m.ID,
		channelID: m.ChannelID,
	}
}

func NewInvocationFromInteractionCreate(i discordgo.InteractionCreate) invocation {
	return invocation{
		id:          i.ID,
		channelID:   i.ChannelID,
		interaction: i.Interaction,
	}
}

type invocation struct {
	id          string
	channelID   string
	interaction *discordgo.Interaction
}

var (
	DR2Randomiser = randomiser.NewSimple(game.DR2)
)

func HandleCreateDR2ChallengeDefault(store model.Store, session *discordgo.Session, invocation invocation) {
	slog.Debug("generating new challenge")
	challenge := challenge.NewRandomChallenge(challenge.Config{}, DR2Randomiser)
	slog.Info("new challenge generated", "stage", challenge.Stage().String(), "weather", challenge.Weather().String(), "car", challenge.Car().String())

	if invocation.interaction != nil {
		// HAVE to respond to an interaction
		err := session.InteractionRespond(invocation.interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Have fun! 🟠🔴🔴🔴🔴🔴",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})

		if err != nil {
			slog.Error("Responding to Default DR2 Challenge", "err", err)
		}
	}

	challengeID, err := sendChallengeMessage(session, invocation.channelID, challenge)
	if err != nil {
		slog.Error("sending challenge message", "id", invocation.id, "channel_id", invocation.channelID, "err", err)
	}

	err = store.PutChallenge(challengeID, challenge)
	if err != nil {
		slog.Error("storing default dr2 challenge", "err", err)
	}
}

func HandleCreateDR2ChallengeCustom(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    fmt.Sprintf(baseMessage, game.DR2.String()),
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: buildChallengeLocationMessageComponents(challenge.Config{Game: game.DR2}),
		},
	})

	if err != nil {
		slog.Error("Create Custom DR2 Challenge Initial Message", "err", err)
	}
}

func HandleCreateWRCChallengeCustom(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    fmt.Sprintf(baseMessage, game.DR2.String()),
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: buildChallengeLocationMessageComponents(challenge.Config{Game: game.WRC}),
		},
	})

	if err != nil {
		slog.Error("Create Custom DR2 Challenge Initial Message", "err", err)
	}
}

func HandleCustomChallengeInteraction(store model.Store, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	stripped := interaction.MessageComponentData().CustomID[4:]
	switch stripped {
	case LocationSelectID, StageSelectID, WeatherSelectID:
		updateLocationSelectMessage(session, interaction)
	case SubmitLocationAndStageID, DrivetrainSelectID, ClassSelectID, CarSelectID:
		updateCarSelectMessage(session, interaction)
	case SubmitCarID:
		updateSelectMessageAndCreateChallenge(store, session, interaction)
	}
}

func updateLocationSelectMessage(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	config, err := buildStageConfigFromInteraction(interaction)

	if err != nil {
		slog.Error("Create Custom Challenge Location Config", "err", err)
	}

	err = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    fmt.Sprintf(baseMessage, config.Game.String()),
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: buildChallengeLocationMessageComponents(config),
		},
	})

	if err != nil {
		slog.Error("Update Create Custom Challenge Location Message", "err", err)
	}
}

func updateCarSelectMessage(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	config, err := buildCarConfigFromInteraction(interaction)

	if err != nil {
		slog.Error("Create Custom Challenge Car Config", "err", err)
	}

	err = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    fmt.Sprintf(baseMessage, config.Game.String()) + "\n" + config.FancyStageString(),
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: buildChallengeCarMessageComponents(config),
		},
	})

	if err != nil {
		slog.Error("Update Create Custom Challenge Car Message", "err", err)
	}
}

func updateSelectMessageAndCreateChallenge(store model.Store, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	config, err := buildCarConfigFromInteraction(interaction)

	if err != nil {
		slog.Error("Create Custom Challenge Final Config", "err", err)
	}

	err = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(baseMessage, config.Game.String()) + "\n" + config.FancyStageString() + "\n" + config.FancyCarString(),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		slog.Error("updating Create Custom Challenge Final Message", "err", err)
	}

	challenge := challenge.NewRandomChallenge(config, DR2Randomiser)
	slog.Info("new challenge generated", "stage", challenge.Stage().String(), "weather", challenge.Weather().String(), "car", challenge.Car().String())

	challengeID, err := sendChallengeMessage(session, interaction.ChannelID, challenge)
	if err != nil {
		slog.Error("sending challenge message", "id", interaction.ID, "channel_id", interaction.ChannelID, "err", err)
	}

	err = store.PutChallenge(challengeID, challenge)
	if err != nil {
		slog.Error("storing custom dr2 challenge", "err", err)
	}
}

func sendChallengeMessage(session *discordgo.Session, channelID string, challenge challenge.Model) (string, error) {
	msg := &discordgo.MessageSend{
		Content:    challenge.FancyString() + "\n",
		Components: getChallengeButtons(),
	}

	sent, err := session.ChannelMessageSendComplex(channelID, msg)
	if err != nil {
		return "", err
	}

	return sent.ID, nil
}

func getChallengeButtons() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Emoji:    &discordgo.ComponentEmoji{Name: "⏱️"},
					Style:    discordgo.PrimaryButton,
					Disabled: false,
					CustomID: CompletedID,
				},
				discordgo.Button{
					Emoji:    &discordgo.ComponentEmoji{Name: "📋"},
					Style:    discordgo.SecondaryButton,
					Disabled: false,
					CustomID: TimesID,
				},
				discordgo.Button{
					Emoji:    &discordgo.ComponentEmoji{Name: "👍"},
					Style:    discordgo.SuccessButton,
					Disabled: true,
					CustomID: GoodID,
				},
				discordgo.Button{
					Emoji:    &discordgo.ComponentEmoji{Name: "👎"},
					Style:    discordgo.DangerButton,
					Disabled: true,
					CustomID: BadID,
				},
			},
		},
	}
}
