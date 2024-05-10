package newstage

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Joe-Hendley/dirtrallybot/internal/feedback/event"
	"github.com/Joe-Hendley/dirtrallybot/internal/memorystore"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/challenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/randomiser"
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
						Emoji:    &discordgo.ComponentEmoji{Name: "🏁"},
						Style:    discordgo.PrimaryButton,
						Disabled: false,
						CustomID: completedID,
					},
					discordgo.Button{
						Emoji:    &discordgo.ComponentEmoji{Name: "⏱️"},
						Style:    discordgo.SecondaryButton,
						Disabled: true,
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
		log.Printf("error sending challenge msg trigger channel id: %s message id: %s :%v", m.ChannelID, m.ID, err)

		return
	}

	// TODO challenge ID is a bit weird at the mo and defined manually in 3 different places - unify it
	challengeID := m.ChannelID + "_" + sent.ID

	memorystore.DefaultStore.Put(challengeID, challenge)
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
			Flags:    discordgo.MessageFlagsEphemeral,
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
	challengeID := split[1] + "_" + split[2]
	userID := split[3]
	unixTimestampString := split[4]

	unixTimestamp, err := strconv.ParseInt(unixTimestampString, 10, 64)
	if err != nil {
		// if this happens then we have muchos problemos - not sure when it could happen really
		log.Printf("error parsing unix timestamp: %s from event %s", unixTimestampString, eventID)
	}

	rawDuration := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	parsed, err := parseTimestamp(rawDuration)
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

	member := i.Interaction.Member
	displayName := ""
	if member == nil {
		log.Printf("member not found for event %s\n", eventID)
	} else {
		if member.Nick != "" {
			displayName = member.Nick
		} else {
			displayName = member.User.GlobalName
		}
	}

	completionEvent := event.New(eventID, unixTimestamp).AsCompletion(userID, displayName, parsed)
	err = memorystore.DefaultStore.ApplyEvent(challengeID, completionEvent)
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

	go updateTopThree(s, split[1], split[2])

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

func updateTopThree(s *discordgo.Session, channelID, messageID string) {
	challengeID := channelID + "_" + messageID
	challenge, ok := memorystore.DefaultStore.Get(challengeID)
	if !ok {
		log.Printf("could not find challenge id: %s", challengeID)
		return
	}

	edited := discordgo.NewMessageEdit(channelID, messageID).SetContent(challenge.FancyString() + "\n" + challenge.TopThreeFancyString())

	_, err := s.ChannelMessageEditComplex(edited)
	if err != nil {
		log.Printf("error editing challenge id: %s : %v\n", challengeID, err)
		return
	}
}

// TODO - better error messages
func parseTimestamp(s string) (time.Duration, error) {
	minuteString, remainder, ok := strings.Cut(s, ":")
	if !ok {
		return 0, fmt.Errorf("invalid timestamp %s", s)
	}

	var err error
	minutes, err := strconv.Atoi(minuteString)
	if err != nil {
		return 0, fmt.Errorf("invalid timestamp %s: %v", s, err)
	}

	if minutes < 0 || minutes > 100 {
		return 0, fmt.Errorf("invalid timestamp %s", s)
	}

	secondString, msString, _ := strings.Cut(remainder, ".")
	seconds, err := strconv.Atoi(secondString)
	if err != nil {
		return 0, fmt.Errorf("invalid timestamp %s: %v", s, err)
	}

	if seconds < 0 || seconds >= 60 {
		return 0, fmt.Errorf("invalid timestamp %s", s)
	}

	milliseconds := 0
	if msString != "" {
		milliseconds, err = strconv.Atoi(msString)
	}
	if err != nil {
		return 0, fmt.Errorf("invalid timestamp %s: %v", s, err)
	}

	if milliseconds < 0 || milliseconds >= 1000 {
		return 0, fmt.Errorf("invalid timestamp %s", s)
	}

	fmt.Println(minutes, seconds, milliseconds)

	return ((time.Duration(minutes) * time.Minute) +
			(time.Duration(seconds) * time.Second) +
			(time.Duration(milliseconds) * time.Millisecond)),
		nil
}

func displayTimes(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
