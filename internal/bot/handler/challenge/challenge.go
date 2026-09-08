package challenge

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot/discord"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handler/feedback"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/render"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/game"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/popularity"
	"github.com/Joe-Hendley/dirtrallybot/internal/randomiser"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/port"
	"github.com/bwmarrin/discordgo"
)

const (
	idFieldDelimiter = "-"

	DisplayCompletionModalID = "completion-display-modal"
	DisplayTimesID           = "completion-display-times"

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

// Generator builds the randomiser that fills in a generated challenge. The
// biased generator reads the store for the current popularity snapshot; the
// others ignore it.
type Generator func(ctx context.Context, store port.Store, g game.Model) challenge.Randomiser

// newRandomiser is the active generator. SetGenerator overrides it - call once
// at startup, before any interaction is served.
var newRandomiser Generator = BiasedGenerator

// SetGenerator selects the randomiser used for new challenges.
func SetGenerator(g Generator) {
	newRandomiser = g
}

// DeterministicGenerator builds a fixed-seed randomiser.
func DeterministicGenerator(_ context.Context, _ port.Store, g game.Model) challenge.Randomiser {
	return randomiser.NewDeterministic(g)
}

// RandomGenerator builds a uniform randomiser.
func RandomGenerator(_ context.Context, _ port.Store, g game.Model) challenge.Randomiser {
	return randomiser.NewRandom(g)
}

// BiasedGenerator builds a randomiser biased by the current popularity snapshot,
// falling back to uniform if the snapshot cannot be read.
func BiasedGenerator(ctx context.Context, store port.Store, g game.Model) challenge.Randomiser {
	snapshot, err := store.Popularity(ctx)
	if err != nil {
		slog.Warn("loading popularity; generating without bias", "err", err)
		snapshot = popularity.Snapshot{}
	}
	return randomiser.NewBiased(g, snapshot)
}

// HandleNewChallenge opens the challenge builder for whichever game the slash
// command names.
func HandleNewChallenge(session discord.InteractionResponder, interaction *discordgo.InteractionCreate) {
	var g game.Model
	switch interaction.ApplicationCommandData().Name {
	case NewDR2ChallengeID:
		g = game.DR2
	case NewWRCChallengeID:
		g = game.WRC
	default:
		slog.Error("unknown new challenge command", "name", interaction.ApplicationCommandData().Name)
		return
	}

	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			CustomID:   initialChallengeResponseID(g),
			Content:    fmt.Sprintf(baseMessage, g.String()),
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: buildChallengeLocationMessageComponents(challenge.Config{Game: g}),
		},
	})

	if err != nil {
		slog.Error("creating challenge builder message", "game", g.String(), "err", err)
	}
}

func initialChallengeResponseID(g game.Model) string {
	switch g {
	case game.WRC:
		return InitialWRCChallengeResponseID
	case game.DR2:
		return InitialDR2ChallengeResponseID
	}
	return ""
}

// SessionStore holds the in-progress Config for an open challenge builder,
// keyed by the builder message ID.
type SessionStore interface {
	Get(builderID string) (challenge.Config, bool)
	Put(builderID string, config challenge.Config)
	Delete(builderID string)
}

func HandleChallengeBuilderInteraction(ctx context.Context, sessions SessionStore, store port.Store, session discord.Session, interaction *discordgo.InteractionCreate) {
	split := strings.Split(interaction.MessageComponentData().CustomID, idFieldDelimiter)
	lastField := split[len(split)-1]

	switch lastField {
	case locationID, distanceID, stageID, weatherID:
		updateLocationSelectMessage(sessions, session, interaction)
	case SubmitLocationAndStageID, drivetrainID, classID, carID:
		updateCarSelectMessage(sessions, session, interaction)
	case SubmitCarID:
		updateSelectMessageAndCreateChallenge(ctx, sessions, store, session, interaction)
	}
}

func updateLocationSelectMessage(sessions SessionStore, session discord.InteractionResponder, interaction *discordgo.InteractionCreate) {
	config, err := configFromInteraction(sessions, interaction)

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

func updateCarSelectMessage(sessions SessionStore, session discord.InteractionResponder, interaction *discordgo.InteractionCreate) {
	config, err := configFromInteraction(sessions, interaction)

	if err != nil {
		slog.Error("Create Custom Challenge Car Config", "err", err)
		updateMessageWithError(session, interaction)
		return
	}

	err = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    fmt.Sprintf(baseMessage, config.Game.String()) + "\n" + render.StageConfig(config),
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: buildChallengeCarMessageComponents(config),
		},
	})

	if err != nil {
		slog.Error("Update Create Custom Challenge Car Message", "err", err)
	}
}

func updateSelectMessageAndCreateChallenge(ctx context.Context, sessions SessionStore, store port.Store, session discord.Session, interaction *discordgo.InteractionCreate) {
	config, err := configFromInteraction(sessions, interaction)

	if err != nil {
		slog.Error("Create Custom Challenge Final Config", "err", err)
		updateMessageWithError(session, interaction)
		return
	}

	err = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(baseMessage, config.Game.String()) + "\n" + render.StageConfig(config) + "\n" + render.CarConfig(config),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		slog.Error("updating Create Custom Challenge Final Message", "err", err)
	}

	challenge := challenge.NewRandomChallenge(config, newRandomiser(ctx, store, config.Game))
	slog.Info("new challenge generated", "game", config.Game.String(), "stage", challenge.Stage().String(), "weather", challenge.Weather().String(), "car", challenge.Car().String())

	challengeID, err := sendChallengeMessage(session, interaction.ChannelID, challenge)
	if err != nil {
		slog.Error("sending challenge message", "id", interaction.ID, "channel_id", interaction.ChannelID, "err", err)
	}

	err = store.PutChallenge(ctx, challengeID, challenge)
	if err != nil {
		slog.Error("storing challenge", "challenge_id", challengeID, "err", err)
	}

	sessions.Delete(interaction.Message.ID)
}

func updateMessageWithError(session discord.InteractionResponder, interaction *discordgo.InteractionCreate) {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "☹️ an error occured, please contact an administrator",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		slog.Error("update interaction message with error", "err", err)
	}
}

func sendChallengeMessage(session discord.ChannelMessageSender, channelID string, challenge challenge.Model) (string, error) {
	msg := &discordgo.MessageSend{
		Content:    render.Challenge(challenge) + "\n",
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
					CustomID: DisplayCompletionModalID,
				},
				discordgo.Button{
					Emoji:    &discordgo.ComponentEmoji{Name: "📋"},
					Style:    discordgo.SecondaryButton,
					Disabled: false,
					CustomID: DisplayTimesID,
				},
				discordgo.Button{
					Emoji:    &discordgo.ComponentEmoji{Name: "👍"},
					Style:    discordgo.SuccessButton,
					Disabled: false,
					CustomID: feedback.GoodID,
				},
				discordgo.Button{
					Emoji:    &discordgo.ComponentEmoji{Name: "👎"},
					Style:    discordgo.DangerButton,
					Disabled: false,
					CustomID: feedback.BadID,
				},
			},
		},
	}
}
