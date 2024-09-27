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
	slog.Info("new challenge generated", "stage", challenge.Stage.String(), "weather", challenge.Weather.String(), "car", challenge.Car.String())

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

	challengeID, err := sendChallengeMessage(*challenge, session, invocation)
	if err != nil {
		slog.Error("sending challenge message", "id", invocation.id, "channel_id", invocation.channelID, "err", err)
	}

	store.Put(challengeID, challenge)
}

func HandleCreateDR2ChallengeCustom(store model.Store, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	challengeConfig, err := buildChallengeConfig(session, interaction)

	slog.Debug("generating new challenge")
	challenge := challenge.New(challengeConfig, r)
	slog.Info("new challenge generated", "stage", challenge.Stage.String(), "weather", challenge.Weather.String(), "car", challenge.Car.String())

	if err != nil {
		slog.Error("Create Custom DR2 Challenge", "err", err)
	}

	challengeID, err := sendChallengeMessage(*challenge, session, NewInvocationFromInteractionCreate(*interaction))
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

// Not sure how to do this as I want to key everything to do with a challenge from the RESPONSE id rather than the INTERACTION ID
// Doesn't seem too annoying to change later on, but one thing at a time
// func sendChallengeInteraction(challenge challenge.Model, session *discordgo.Session, invocation invocation) (string, error) {
// 	respData := &discordgo.InteractionResponseData{
// 		Content:    challenge.FancyString() + "\n",
// 		Components: getChallengeButtons(),
// 	}

// 	sent, err := session.InteractionRespond(invocation.interaction, &discordgo.InteractionResponse{
// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
// 		Data: respData,
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	return sent.ID, nil
// }

func sendChallengeMessage(challenge challenge.Model, session *discordgo.Session, invocation invocation) (string, error) {
	msg := &discordgo.MessageSend{
		Content:    challenge.FancyString() + "\n",
		Components: getChallengeButtons(),
	}

	sent, err := session.ChannelMessageSendComplex(invocation.channelID, msg)
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
