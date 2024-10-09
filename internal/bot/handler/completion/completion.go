package completion

import (
	"fmt"
	"log"
	"log/slog"
	"slices"
	"strings"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot/discord"
	"github.com/Joe-Hendley/dirtrallybot/internal/model"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/timestamp"
	"github.com/bwmarrin/discordgo"
)

const (
	SubmitCompletionPrefix = "completion-submit"
	CompletionTextInputID  = "completion-input"

	challengeIDIndex = 1
	userIDIndex      = 2
	customIDDelim    = "-"
)

func HandleDisplayEntryModal(session discord.InteractionResponder, interaction *discordgo.InteractionCreate) {
	customIDParts := []string{SubmitCompletionPrefix, "", ""}
	customIDParts[challengeIDIndex] = interaction.Message.ID
	customIDParts[userIDIndex] = interaction.Interaction.Member.User.ID

	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: strings.Join(customIDParts, customIDDelim),
			Title:    "Please input your time",
			Content:  "Please input your time",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    CompletionTextInputID,
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

func HandleSubmitModal(store model.Store, session discord.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionModalSubmit {
		return
	}

	data := interaction.ModalSubmitData()

	if !strings.HasPrefix(data.CustomID, SubmitCompletionPrefix) {
		return
	}

	split := strings.Split(data.CustomID, "-")
	challengeID := split[challengeIDIndex]
	userID := split[userIDIndex]

	rawDuration := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	parsed, err := timestamp.Parse(rawDuration)
	if err != nil {
		respErr := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
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
		slog.Error("submitting timestamp", "challenge-id", challengeID, "err", err)
		respErr := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
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

	go updateTopThree(store, session, interaction.GuildID, split[1], split[2])

	err = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
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

func medal(place int) string {
	switch place {
	case 1:
		return "🥇"
	case 2:
		return "🥈"
	case 3:
		return "🥉"
	default:
		return ""
	}
}

func updateTopThree(store model.Store, session discord.Session, guildID, channelID, messageID string) {
	challengeID := messageID
	challenge, err := store.GetChallenge(challengeID)
	if err != nil {
		slog.Warn("getting challenge", "challengeID", challengeID, "err", err)
		return
	}

	lines := []string{}
	for lineIndex, completion := range challenge.TopThree() {
		lines = append(
			lines,
			fmt.Sprintf(
				"%s **%s**\t%s",
				medal(lineIndex+1),
				timestamp.Format(completion.Duration()),
				discord.GetGuildMemberDisplayName(session, guildID, completion.UserID())),
		)
	}

	topThreeString := strings.Join(lines, "\n")

	edited := discordgo.NewMessageEdit(channelID, messageID).SetContent(challenge.FancyString() + "\n" + topThreeString)

	_, err = session.ChannelMessageEditComplex(edited)
	if err != nil {
		slog.Error("editing challenge id: %s : %v\n", challengeID, err)
		return
	}
}

func HandleDisplayTimes(store model.Store, session discord.Session, interaction *discordgo.InteractionCreate) {
	challengeID := interaction.Message.ID
	challenge, err := store.GetChallenge(challengeID)
	if err != nil {
		slog.Warn("getting challenge", "challengeID", challengeID, "err", err)
		return
	}

	userCompletionMap := challenge.FancyListCompletions()

	type user struct {
		id          string
		displayName string
	}

	users := []user{}
	for userID := range userCompletionMap {
		users = append(users, user{id: userID, displayName: discord.GetGuildMemberDisplayName(session, interaction.GuildID, userID)})
	}

	slices.SortFunc(users, func(a user, b user) int {
		if a.displayName < b.displayName {
			return -1
		}
		if a.displayName > b.displayName {
			return 1
		}
		return 0
	})

	err = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		slog.Error("display times", "err", err)
	}
}
