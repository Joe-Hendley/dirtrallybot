package ready

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func Handler(s *discordgo.Session, r *discordgo.Ready) {
	slog.Info("Bot is ready")
}
