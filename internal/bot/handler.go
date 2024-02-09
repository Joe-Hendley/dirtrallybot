package bot

import (
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handlers/debug"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handlers/newstage"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handlers/ready"
	"github.com/bwmarrin/discordgo"
)

func RegisterHandlers(s *discordgo.Session, useDebug bool) {
	s.AddHandler(ready.Handler)

	if useDebug {
		s.AddHandler(debug.Handler)
	}

	s.AddHandler(newstage.Handler)
	s.AddHandler(newstage.Interaction)
}
