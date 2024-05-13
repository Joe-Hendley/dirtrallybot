package newchallenge

import (
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/event"
	"github.com/Joe-Hendley/dirtrallybot/internal/parse"
	"github.com/Joe-Hendley/dirtrallybot/internal/randomiser"
	"github.com/Joe-Hendley/dirtrallybot/internal/store"
	"github.com/bwmarrin/discordgo"
)

const (
	completedID           = "newstage_completed"
	completionEventPrefix = "completed"

	timesID = "newstage_time"
	goodID  = "newstage_good"
	badID   = "newstage_bad"
)

var (
	r = randomiser.NewSimple()
)

func HandleNewChallenge(s *discordgo.Session, m *discordgo.MessageCreate) {
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
						CustomID: completedID,
					},
					discordgo.Button{
						Emoji:    &discordgo.ComponentEmoji{Name: "📋"},
						Style:    discordgo.SecondaryButton,
						Disabled: false,
						CustomID: timesID,
					},
					discordgo.Button{
						Emoji:    &discordgo.ComponentEmoji{Name: "👍"},
						Style:    discordgo.SuccessButton,
						Disabled: true,
						CustomID: goodID,
					},
					discordgo.Button{
						Emoji:    &discordgo.ComponentEmoji{Name: "👎"},
						Style:    discordgo.DangerButton,
						Disabled: true,
						CustomID: badID,
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

var interactionMap = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	completedID: completion,
	timesID:     displayTimes,
	goodID:      good,
	badID:       bad,
}

func Interaction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	if h, ok := interactionMap[i.MessageComponentData().CustomID]; ok {
		h(s, i)
	}
}

func completion(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: strings.Join([]string{completionEventPrefix, i.ChannelID, i.Message.ID, i.Interaction.Member.User.ID, strconv.FormatInt(time.Now().Unix(), 10)}, "_"),
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
		log.Println(err)
	}
}

func HandleCompletionRequest(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionModalSubmit {
		return
	}

	data := i.ModalSubmitData()

	if !strings.HasPrefix(data.CustomID, completionEventPrefix) {
		return
	}

	split := strings.Split(data.CustomID, "_")
	eventID := split[0]
	challengeID := split[2]
	userID := split[3]
	unixTimestampString := split[4]

	unixTimestamp, err := strconv.ParseInt(unixTimestampString, 10, 64)
	if err != nil {
		// if this happens then we have muchos problemos - not sure when it could happen really
		log.Printf("error parsing unix timestamp: %s from event %s", unixTimestampString, eventID)
	}

	rawDuration := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	parsed, err := parse.Timestamp(rawDuration)
	if err != nil {
		respErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Timestamp needs to be submitted in format: mm:ss.mss",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})

		if respErr != nil {
			log.Printf("error sending response to bad timestamp: %v\n", respErr)
		}
		return
	}

	completionEvent := event.New(eventID, unixTimestamp).AsCompletion(userID, parsed)
	err = store.ApplyEvent(challengeID, completionEvent)
	if err != nil {
		log.Printf("error submitting timestamp for eventID: %s challengeID: %s: %v", eventID, challengeID, err)
		respErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "error submitting timestamp - contact an administrator",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})

		if respErr != nil {
			log.Printf("error sending response: %v\n", respErr)
		}
		return
	}

	go updateTopThree(s, i.GuildID, split[1], split[2])

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Submitted time %s", rawDuration),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		log.Printf("error sending response to valid timestamp: %v\n", err)
	}
}

func updateTopThree(s *discordgo.Session, guildID, channelID, messageID string) {
	challengeID := messageID
	challenge, ok := store.Get(challengeID)
	if !ok {
		slog.Warn("could not find challenge", "id", challengeID)
		return
	}

	edited := discordgo.NewMessageEdit(channelID, messageID).SetContent(challenge.FancyString() + "\n" + challenge.TopThreeFancyString(s, guildID))

	_, err := s.ChannelMessageEditComplex(edited)
	if err != nil {
		log.Printf("error editing challenge id: %s : %v\n", challengeID, err)
		return
	}
}

func displayTimes(s *discordgo.Session, i *discordgo.InteractionCreate) {
	challengeID := i.Message.ID
	challenge, ok := store.Get(challengeID)
	if !ok {
		slog.Warn("could not find challenge", "id", challengeID)
		return
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: challenge.FancyListCompletions(s, i.GuildID),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		slog.Error("display times", "err", err)
	}
}

func good(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

func bad(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
