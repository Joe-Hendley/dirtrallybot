package newstage

import (
	"log"

	"github.com/Joe-Hendley/dirtrallybot/internal/memorystore"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/randomiser"
	"github.com/bwmarrin/discordgo"
)

const (
	finishedID = "newstage_finished"

	newTimeID = "newtime"

	timesID = "newstage_time"
	goodID  = "newstage_good"
	badID   = "newstage_bad"
)

var (
	r = randomiser.NewSimple()
)

func Handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content != "!newstage" {
		return
	}

	challenge := challenge.New(r)

	msg := &discordgo.MessageSend{
		Content: challenge.FancyString(),
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Emoji:    discordgo.ComponentEmoji{Name: "🏁"},
						Style:    discordgo.PrimaryButton,
						Disabled: true,
						CustomID: finishedID,
					},
					discordgo.Button{
						Emoji:    discordgo.ComponentEmoji{Name: "⏱️"},
						Style:    discordgo.SecondaryButton,
						Disabled: true,
						CustomID: timesID,
					},
					discordgo.Button{
						Emoji:    discordgo.ComponentEmoji{Name: "👍"},
						Style:    discordgo.SuccessButton,
						Disabled: true,
						CustomID: goodID,
					},
					discordgo.Button{
						Emoji:    discordgo.ComponentEmoji{Name: "👎"},
						Style:    discordgo.DangerButton,
						Disabled: true,
						CustomID: badID,
					},
				},
			},
		},
	}

	resp, err := s.ChannelMessageSendComplex(m.ChannelID, msg)
	if err != nil {
		log.Printf("error sending challenge msg trigger msgid[%s] :%v", m.ID, err)

		return
	}

	memorystore.DefaultStore.Put(resp.ID, challenge)
}

var interactionMap = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	finishedID: finish,
	timesID:    time,
	goodID:     good,
	badID:      bad,
}

func Interaction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	if h, ok := interactionMap[i.MessageComponentData().CustomID]; ok {
		h(s, i)
	}
}

func finish(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "finished",
			Title:    "Please input your time",
			Content:  "Please input your time",
			Flags:    discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    newTimeID,
							Label:       "Minutes:Seconds.Milliseconds",
							Style:       discordgo.TextInputShort,
							Placeholder: "12:34.567",
							Required:    true,
							MaxLength:   9,
							MinLength:   9,
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

func time(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "tick tock",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		log.Println(err)
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
		log.Println(err)
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
		log.Println(err)
	}
}
