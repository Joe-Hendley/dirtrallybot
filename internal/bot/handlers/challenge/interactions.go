package challenge

import (
	"log"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/model"
	"github.com/bwmarrin/discordgo"
)

func HandleInteractionMessageComponent(store model.Store, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionMessageComponent {
		return
	}

	switch interaction.MessageComponentData().CustomID {
	case completedID:
		completion(session, interaction)
	case timesID:
		displayTimes(store, session, interaction)
	case goodID:
		good(session, interaction)
	case badID:
		bad(session, interaction)
	}
}

func HandleModalSubmit(store model.Store, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionModalSubmit {
		return
	}

	data := interaction.ModalSubmitData()

	if strings.HasPrefix(data.CustomID, completionEventPrefix) {
		HandleCompletionRequest(store, session, interaction)
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
		log.Printf("error editing challenge id: %s : %v\n", challengeID, err)
		return
	}
}

func displayTimes(store model.Store, session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	challengeID := interaction.Message.ID
	challenge, ok := store.Get(challengeID)
	if !ok {
		slog.Warn("could not find challenge", "id", challengeID)
		return
	}

	err := session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: challenge.FancyListCompletions(session, interaction.GuildID),
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
