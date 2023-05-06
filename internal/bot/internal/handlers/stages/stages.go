package stages

import (
	"log"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/builder"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/bwmarrin/discordgo"
)

func Handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content != "!stages" {
		return
	}

	for _, location := range location.List() {
		msg := ""
		for _, stage := range builder.StagesAtLocation(location) {
			msg += stage.FancyString() + "\n"
		}
		_, err := s.ChannelMessageSend(m.ChannelID, msg)
		if err != nil {
			log.Printf("error sending stages message: %v\n", err)
		}
	}
}
