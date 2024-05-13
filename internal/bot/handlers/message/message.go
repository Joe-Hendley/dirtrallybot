package message

import (
	"log/slog"
	"strings"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handlers/debug"
	"github.com/Joe-Hendley/dirtrallybot/internal/bot/handlers/newchallenge"
	"github.com/bwmarrin/discordgo"
)

func HandleCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	slog.Debug("message", "id", m.ID, "author_id", m.Author.ID, "author_name", m.Author.GlobalName)

	lowercase := strings.ToLower(m.Content)

	switch lowercase {
	case "!cars":
		debug.HandleCars(s, m)

	case "!stages":
		debug.HandleStages(s, m)

	case "!newstage":
		newchallenge.HandleNewChallenge(s, m)

	default:
		slog.Debug("message ignored", "id", m.ID)
	}
}
