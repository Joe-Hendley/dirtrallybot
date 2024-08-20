package bot

import (
	"log"

	"github.com/Joe-Hendley/dirtrallybot/internal/config"
	"github.com/bwmarrin/discordgo"
)

var (
	commands []discordgo.ApplicationCommand
	cmdIDs   map[string]string
)

func registerCommands() {
}

func CreateCommands(config config.Config, session *discordgo.Session) {
	registerCommands()

	cmdIDs = make(map[string]string, len(commands))

	for _, cmd := range commands {
		rcmd, err := session.ApplicationCommandCreate(config.App, "", &cmd)
		if err != nil {
			log.Fatalf("Cannot create slash command %q: %v", cmd.Name, err)
		}

		cmdIDs[rcmd.ID] = rcmd.Name
	}
}

func CleanupCommands(config config.Config, session *discordgo.Session) {
	for id, name := range cmdIDs {
		err := session.ApplicationCommandDelete(config.App, "", id)
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", name, err)
		}
	}
}
