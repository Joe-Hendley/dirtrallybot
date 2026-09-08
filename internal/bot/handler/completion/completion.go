package completion

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot/discord"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/render"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/timestamp"
	"github.com/Joe-Hendley/dirtrallybot/internal/store/port"
	"github.com/bwmarrin/discordgo"
)

const (
	SubmitCompletionPrefix = "completion-submit"
	CompletionTextInputID  = "completion-submit-input"
	InvalidSubmissionID    = "completion-submit-input-invalid"
	ValidSubmissionID      = "completion-submit-input-valid"

	DisplayTimesResponseID = "completion-display-response"

	challengeIDIndex = 1
	userIDIndex      = 2
	customIDDelim    = "_"

	genericErrorMessage = "☹️ something went wrong, please contact an administrator"
)

// respondEphemeral sends a private message back to whoever triggered the
// interaction, logging if even that fails.
func respondEphemeral(session discord.InteractionResponder, interaction *discordgo.InteractionCreate, content string) {
	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		slog.Error("sending ephemeral response", "err", err)
	}
}

func HandleDisplayEntryModal(session discord.InteractionResponder, interaction *discordgo.InteractionCreate) {
	customIDParts := []string{SubmitCompletionPrefix, "", ""}
	customIDParts[challengeIDIndex] = interaction.Message.ID
	customIDParts[userIDIndex] = interaction.Member.User.ID

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

func HandleSubmitModal(ctx context.Context, store port.Store, session discord.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionModalSubmit {
		return
	}

	data := interaction.ModalSubmitData()

	if !strings.HasPrefix(data.CustomID, SubmitCompletionPrefix) {
		return
	}

	split := strings.Split(data.CustomID, customIDDelim)
	challengeID := split[challengeIDIndex]
	userID := split[userIDIndex]

	rawDuration := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	parsed, err := timestamp.Parse(rawDuration)
	if err != nil {
		respErr := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content:  "Timestamp needs to be submitted in format: mm:ss.mss",
				Flags:    discordgo.MessageFlagsEphemeral,
				CustomID: InvalidSubmissionID,
			},
		})

		if respErr != nil {
			slog.Error("responding to invalid timestamp", "err", respErr)
		}
		return
	}

	displayName := discord.GetGuildMemberDisplayName(session, interaction.GuildID, userID)
	completion := challenge.NewCompletionAt(userID, displayName, parsed, time.Now().UTC())
	err = store.RegisterCompletion(ctx, challengeID, completion)

	if err != nil {
		slog.Error("registering completion", "challenge_id", challengeID, "err", err)
		respondEphemeral(session, interaction, genericErrorMessage)
		return
	}

	err = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:  fmt.Sprintf("Submitted time %s", rawDuration),
			Flags:    discordgo.MessageFlagsEphemeral,
			CustomID: ValidSubmissionID,
		},
	})

	if err != nil {
		slog.Error("responding to valid timestamp", "err", err)
	}

	updateTopThree(ctx, store, session, interaction.GuildID, interaction.ChannelID, challengeID)
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

func updateTopThree(ctx context.Context, store port.Store, session discord.Session, guildID, channelID, messageID string) {
	challengeID := messageID
	challenge, err := store.GetChallenge(ctx, challengeID)
	if err != nil {
		slog.Warn("getting challenge", "challenge_id", challengeID, "err", err)
		return
	}

	lines := []string{}
	for lineIndex, completion := range challenge.TopThree() {
		name := completion.DisplayName()
		if name == "" {
			name = discord.GetGuildMemberDisplayName(session, guildID, completion.UserID())
		}

		lines = append(
			lines,
			fmt.Sprintf(
				"%s **%s**\t%s",
				medal(lineIndex+1),
				timestamp.Format(completion.Duration()),
				name),
		)
	}

	topThreeString := strings.Join(lines, "\n")

	edited := discordgo.NewMessageEdit(channelID, messageID).SetContent(render.Challenge(challenge) + "\n" + topThreeString)

	_, err = session.ChannelMessageEditComplex(edited)
	if err != nil {
		slog.Error("editing challenge message", "challenge_id", challengeID, "err", err)
		return
	}
}

func HandleDisplayTimes(ctx context.Context, store port.Store, session discord.Session, interaction *discordgo.InteractionCreate) {
	challengeID := interaction.Message.ID
	challenge, err := store.GetChallenge(ctx, challengeID)
	if err != nil {
		slog.Warn("getting challenge", "challenge_id", challengeID, "err", err)
		respondEphemeral(session, interaction, "☹️ couldn't find this challenge - it may have expired")
		return
	}

	userCompletionMap := challenge.UserCompletions()

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

	buf := strings.Builder{}

	for _, user := range users {
		buf.Write([]byte("**" + user.displayName + "**\n"))
		for _, completion := range userCompletionMap[user.id] {
			buf.Write([]byte("\t" + timestamp.Format(completion) + "\n"))
		}
	}

	const discordMaxMessageLength = 2000
	if buf.Len() > discordMaxMessageLength {
		buf.Reset()
		buf.Write([]byte("☹️ unable to display all completions, please contact an administrator"))
	}

	err = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:  buf.String(),
			Flags:    discordgo.MessageFlagsEphemeral,
			CustomID: DisplayTimesResponseID,
		},
	})

	if err != nil {
		slog.Error("display times", "err", err)
	}
}
