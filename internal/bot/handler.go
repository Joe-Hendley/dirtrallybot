package bot

import (
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handlers/message"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handlers/newchallenge"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handlers/ready"
	"github.com/bwmarrin/discordgo"
)

func RegisterHandlers(s *discordgo.Session) {
	s.AddHandler(ready.Handler)

	s.AddHandler(message.HandleCreate)
	s.AddHandler(newchallenge.Interaction)
	s.AddHandler(newchallenge.HandleCompletionRequest)
}
