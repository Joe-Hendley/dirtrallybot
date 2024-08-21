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
)

var (
	r = randomiser.NewSimple()
)

func HandleNewChallenge(store model.Store, s *discordgo.Session, m *discordgo.MessageCreate) {
	slog.Debug("generating new challenge")
	challenge := challenge.New(r)
	slog.Info("new challenge generated", "stage", challenge.Stage.String(), "weather", challenge.Weather.String(), "car", challenge.Car.String())

	msg := &discordgo.MessageSend{
		Content: challenge.FancyString() + "\n",
		Components: []discordgo.MessageComponent{
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
		},
	}

	sent, err := s.ChannelMessageSendComplex(m.ChannelID, msg)
	if err != nil {
		slog.Error("sending challenge message", "id", m.ID, "channel_id", m.ChannelID, "err", err)

		return
	}

	challengeID := sent.ID

	store.Put(challengeID, challenge)
}
