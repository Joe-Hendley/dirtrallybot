package debug

import (
	"log"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/bwmarrin/discordgo"
)

func Handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	switch m.Content {
	case "!cars":
		handleCars(s, m)

	case "!stages":
		handleStages(s, m)
	}
}

func handleCars(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := ""
	for _, class := range class.List() {
		for _, car := range car.InClass(class) {
			buf := car.FancyString() + "\n"

			if len(msg)+len(buf) > 2000 {
				_, err := s.ChannelMessageSend(m.ChannelID, msg)
				if err != nil {
					log.Printf("error sending cars message: %v\n", err)
				}

				msg = buf
			} else {
				msg += buf
			}
		}
	}
	_, err := s.ChannelMessageSend(m.ChannelID, msg)
	if err != nil {
		log.Printf("error sending cars message: %v\n", err)
	}
}

func handleStages(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := ""

	for _, location := range location.List() {
		for _, stage := range stage.AtLocation(location) {
			buf := stage.FancyString() + "\n"

			if len(msg)+len(buf) > 2000 {
				_, err := s.ChannelMessageSend(m.ChannelID, msg)
				if err != nil {
					log.Printf("error sending stages message: %v\n", err)
				}

				msg = buf
			} else {
				msg += buf
			}
		}
	}

	_, err := s.ChannelMessageSend(m.ChannelID, msg)
	if err != nil {
		log.Printf("error sending stages message: %v\n", err)
	}
}
