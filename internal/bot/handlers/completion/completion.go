package completion

import (
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"strings"

	"github.com/Joe-Hendley/dirtrallybot/internal/model"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/event"
	"github.com/Joe-Hendley/dirtrallybot/internal/parse"
	"github.com/bwmarrin/discordgo"
)

const (
	SubmitCompletionPrefix = "completed"
)

func HandleCompletionRequest(store model.Store, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionModalSubmit {
		return
	}

	data := i.ModalSubmitData()

	if !strings.HasPrefix(data.CustomID, SubmitCompletionPrefix) {
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

	go updateTopThree(store, s, i.GuildID, split[1], split[2])

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

func updateTopThree(store model.Store, s *discordgo.Session, guildID, channelID, messageID string) {
	challengeID := messageID
	challenge, ok := store.Get(challengeID)
	if !ok {
		slog.Warn("could not find challenge", "id", challengeID)
		return
	}

	edited := discordgo.NewMessageEdit(channelID, messageID).SetContent(challenge.FancyString() + "\n" + challenge.TopThreeFancyString(s, guildID))

	_, err := s.ChannelMessageEditComplex(edited)
	if err != nil {
		slog.Error("editing challenge id: %s : %v\n", challengeID, err)
		return
	}
}
