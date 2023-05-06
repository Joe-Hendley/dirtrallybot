package cars

import (
	"log"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/builder"
	"github.com/bwmarrin/discordgo"
)

func Handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content != "!cars" {
		return
	}

	cars := builder.Cars()
	msg := ""
	for car := range cars {
		msg += car.String() + "\n"
	}

	_, err := s.ChannelMessageSend(m.ChannelID, msg)
	if err != nil {
		log.Printf("error sending cars message: %v\n", err)
	}
}
