package ready

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func Handler(s *discordgo.Session, r *discordgo.Ready) {
	log.Println("Bot is ready")
}
