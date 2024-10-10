package challenge

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot/discord"
	"github.com/Joe-Hendley/dirtrallybot/internal/model"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
	"github.com/Joe-Hendley/dirtrallybot/internal/randomiser"
	"github.com/bwmarrin/discordgo"
)

const (
	idFieldDelimiter = "-"

	CompletedID = "completion-add"
	TimesID     = "completion-display"
	GoodID      = "feedback-good"
	BadID       = "feedback-bad"

	NewChallengeID = "newstage"
	ResponseID     = "response"
	ChallengeID    = "challenge"

	DR2ID = "dr2"
	WRCID = "wrc"

	NewDR2ChallengeID             = NewChallengeID + idFieldDelimiter + DR2ID
	DR2ChallengePrefix            = ChallengeID + idFieldDelimiter + DR2ID
	InitialDR2ChallengeResponseID = DR2ChallengePrefix + idFieldDelimiter + ResponseID

	NewWRCChallengeID             = NewChallengeID + idFieldDelimiter + WRCID
	WRCChallengePrefix            = ChallengeID + idFieldDelimiter + WRCID
	InitialWRCChallengeResponseID = WRCChallengePrefix + idFieldDelimiter + ResponseID
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

func HandleNewDR2Challenge(session discord.InteractionResponder, interaction *discordgo.InteractionCreate) {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			CustomID:   InitialDR2ChallengeResponseID,
			Content:    fmt.Sprintf(baseMessage, game.DR2.String()),
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: buildChallengeLocationMessageComponents(challenge.Config{Game: game.DR2}),
		},
	})

	if err != nil {
		slog.Error("Create Custom DR2 Challenge Initial Message", "err", err)
	}
}

func HandleNewWRCChallenge(session discord.InteractionResponder, interaction *discordgo.InteractionCreate) {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			CustomID:   InitialWRCChallengeResponseID,
			Content:    fmt.Sprintf(baseMessage, game.DR2.String()),
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: buildChallengeLocationMessageComponents(challenge.Config{Game: game.WRC}),
		},
	})

	if err != nil {
		slog.Error("Create Custom DR2 Challenge Initial Message", "err", err)
	}
}

func HandleChallengeBuilderInteraction(store model.Store, session discord.Session, interaction *discordgo.InteractionCreate) {
	split := strings.Split(interaction.MessageComponentData().CustomID, idFieldDelimiter)
	lastField := split[len(split)-1]

	switch lastField {
	case locationID, distanceID, stageID, weatherID:
		updateLocationSelectMessage(session, interaction)
	case SubmitLocationAndStageID, drivetrainID, classID, carID:
		updateCarSelectMessage(session, interaction)
	case SubmitCarID:
		updateSelectMessageAndCreateChallenge(store, session, interaction)
	}
}

func updateLocationSelectMessage(session discord.InteractionResponder, interaction *discordgo.InteractionCreate) {
	config, err := buildStageConfigFromInteraction(interaction)

	if err != nil {
		slog.Error("Create Custom Challenge Location Config", "err", err)
		updateMessageWithError(session, interaction)
		return
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

func updateCarSelectMessage(session discord.InteractionResponder, interaction *discordgo.InteractionCreate) {
	config, err := buildCarConfigFromInteraction(interaction)

	if err != nil {
		slog.Error("Create Custom Challenge Car Config", "err", err)
		updateMessageWithError(session, interaction)
		return
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

func updateSelectMessageAndCreateChallenge(store model.Store, session discord.Session, interaction *discordgo.InteractionCreate) {
	config, err := buildCarConfigFromInteraction(interaction)

	if err != nil {
		slog.Error("Create Custom Challenge Final Config", "err", err)
		updateMessageWithError(session, interaction)
		return
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

func updateMessageWithError(session discord.InteractionResponder, interaction *discordgo.InteractionCreate) {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("☹️ an error occured, please contact an administrator"),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		slog.Error("update interaction message with error", "err", err)
	}
}

func sendChallengeMessage(session discord.ChannelMessageSender, channelID string, challenge challenge.Model) (string, error) {
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
