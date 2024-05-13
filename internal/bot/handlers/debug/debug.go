package debug

import (
	"log/slog"

	"github.com/Joe-Hendley/dirtrallybot/internal/model/car"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/class"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/location"
	"github.com/Joe-Hendley/dirtrallybot/internal/model/stage"
	"github.com/bwmarrin/discordgo"
)

func HandleCars(s *discordgo.Session, m *discordgo.MessageCreate) {
	slog.Debug("printing cars")

	msg := ""
	for _, class := range class.List() {
		for _, car := range car.InClass(class) {
			buf := car.FancyString() + "\n"

			if len(msg)+len(buf) > 2000 {
				_, err := s.ChannelMessageSend(m.ChannelID, msg)
				if err != nil {
					slog.Error("sending cars message", "err", err)
				}

				msg = buf
			} else {
				msg += buf
			}
		}
	}
	_, err := s.ChannelMessageSend(m.ChannelID, msg)
	if err != nil {
		slog.Error("sending cars message", "err", err)
	}
}

func HandleStages(s *discordgo.Session, m *discordgo.MessageCreate) {
	slog.Debug("sending stages message")

	msg := ""
	for _, location := range location.List() {
		for _, stage := range stage.AtLocation(location) {
			buf := stage.FancyString() + "\n"

			if len(msg)+len(buf) > 2000 {
				_, err := s.ChannelMessageSend(m.ChannelID, msg)
				if err != nil {
					slog.Error("sending stages message", "err", err)
				}

				msg = buf
			} else {
				msg += buf
			}
		}
	}

	_, err := s.ChannelMessageSend(m.ChannelID, msg)
	if err != nil {
		slog.Error("sending stages message", "err", err)
	}
}
