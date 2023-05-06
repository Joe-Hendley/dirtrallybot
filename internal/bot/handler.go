package bot

import (
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/internal/handlers/cars"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/internal/handlers/ready"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/internal/handlers/stages"
	"github.com/bwmarrin/discordgo"
)

func RegisterHandlers(session *discordgo.Session) {
	session.AddHandler(ready.Handler)
	session.AddHandler(cars.Handler)
	session.AddHandler(stages.Handler)
}
