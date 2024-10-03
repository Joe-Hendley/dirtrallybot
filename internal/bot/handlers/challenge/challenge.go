package challenge

import (
	"log/slog"

	"github.com/Joe-Hendley/dirtrallybot/internal/model"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/randomiser"
	"github.com/bwmarrin/discordgo"
)

const (
	CompletedID = "newstage_completed"
	TimesID     = "newstage_time"
	GoodID      = "newstage_good"
	BadID       = "newstage_bad"

	NewstageDR2DefaultID = "newstage"
	NewstageWRCDefaultID = "newstagewrc"
	NewstageDR2CustomID  = "newstage-custom"
	NewstageWRCCustomID  = "newstagewrc-custom"
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
	r = randomiser.NewSimple()
)

func HandleCreateDR2ChallengeDefault(store model.Store, session *discordgo.Session, invocation invocation) {
	slog.Debug("generating new challenge")
	challenge := challenge.New(challenge.Config{}, r)
	slog.Info("new challenge generated", "stage", challenge.Stage.LongString(), "weather", challenge.Weather.String(), "car", challenge.Car.LongString())

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

	challengeID, err := sendChallengeMessage(session, invocation.channelID, *challenge)
	if err != nil {
		slog.Error("sending challenge message", "id", invocation.id, "channel_id", invocation.channelID, "err", err)
	}

	store.Put(challengeID, challenge)
}

func HandleCreateDR2ChallengeCustom(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    buildChallengeBuilderContent(challenge.Config{}),
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: buildChallengeLocationMessageComponents(challenge.Config{}),
		},
	})

	if err != nil {
		slog.Error("Create Custom DR2 Challenge Initial Message", "err", err)
	}
}

func HandleCustomDR2ChallengeInteraction(store model.Store, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	switch interaction.MessageComponentData().CustomID {
	case LocationSelectID, StageSelectID, WeatherSelectID:
		updateDR2LocationSelectMessage(session, interaction)
	case SubmitLocationAndStageID, DrivetrainSelectID, ClassSelectID, CarSelectID:
		updateDR2CarSelectMessage(session, interaction)
	case SubmitCarID:
		updateDR2SelectMessageAndCreateChallenge(store, session, interaction)
	}
}

func updateDR2LocationSelectMessage(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	config, err := buildStageConfigFromInteraction(interaction)

	if err != nil {
		slog.Error("Create Custom DR2 Challenge Config", "err", err)
	}

	err = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    buildChallengeBuilderContent(config),
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: buildChallengeLocationMessageComponents(config),
		},
	})

	if err != nil {
		slog.Error("Update Create Custom DR2 Location Challenge Message", "err", err)
	}
}

func updateDR2CarSelectMessage(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	config, err := buildCarConfigFromInteraction(interaction)

	if err != nil {
		slog.Error("Create Custom DR2 Challenge Config", "err", err)
	}

	err = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    buildChallengeBuilderContent(config) + "\n" + config.FancyStageString(),
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: buildChallengeCarMessageComponents(config),
		},
	})

	if err != nil {
		slog.Error("Update Create Custom DR2 Car Challenge Message", "err", err)
	}
}

func updateDR2SelectMessageAndCreateChallenge(store model.Store, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	config, err := buildCarConfigFromInteraction(interaction)

	if err != nil {
		slog.Error("Create Custom DR2 Challenge Config", "err", err)
	}

	err = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: buildChallengeBuilderContent(config) + "\n" + config.FancyStageString() + "\n" + config.FancyCarString(),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		slog.Error("Update Create Custom DR2 Final Challenge Message", "err", err)
	}

	challenge := challenge.New(config, r)
	slog.Info("new challenge generated", "stage", challenge.Stage.LongString(), "weather", challenge.Weather.String(), "car", challenge.Car.LongString())

	challengeID, err := sendChallengeMessage(session, interaction.ChannelID, *challenge)
	if err != nil {
		slog.Error("sending challenge message", "id", interaction.ID, "channel_id", interaction.ChannelID, "err", err)
	}

	store.Put(challengeID, challenge)
}

func HandleCreateWRCChallengeDefault(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "WRC hasn't been implemented yet! Sorry!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		slog.Error("Create Default WRC Challenge", "err", err)
	}
}

func HandleCreateWRCChallengeCustom(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "WRC hasn't been implemented yet! Sorry!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		slog.Error("Create Custom WRC Challenge", "err", err)
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
