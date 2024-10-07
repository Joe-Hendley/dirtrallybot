package completion

import (
	"fmt"
	"log"
	"log/slog"
	"strings"

	"github.com/Joe-Hendley/dirtrallybot/internal/model"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
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

	completion := challenge.NewCompletion(userID, parsed)
	err = store.RegisterCompletion(challengeID, completion)

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
	challenge, err := store.GetChallenge(challengeID)
	if err != nil {
		slog.Warn("getting challenge", "challengeID", challengeID, "err", err)
		return
	}

	edited := discordgo.NewMessageEdit(channelID, messageID).SetContent(challenge.FancyString() + "\n" + challenge.TopThreeFancyString(s, guildID))

	_, err = s.ChannelMessageEditComplex(edited)
	if err != nil {
		slog.Error("editing challenge id: %s : %v\n", challengeID, err)
		return
	}
}
