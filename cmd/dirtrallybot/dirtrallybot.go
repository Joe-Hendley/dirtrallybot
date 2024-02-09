package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot"
	"github.com/bwmarrin/discordgo"
)

func main() {
	botConfig := bot.NewConfig()
	session, err := discordgo.New("Bot " + botConfig.Token)

	if err != nil {
		log.Fatalf("error creating session: %v", err)
	}

	bot.RegisterHandlers(session, true)
	bot.CreateCommands(botConfig, session)

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentGuildMessageReactions

	err = session.Open()
	if err != nil {
		log.Fatalf("error opening connection: %v", err)
	}

	log.Println("Bot is running. Press CTL-C to exit.")
	defer session.Close()

	waitForInterrupt()

	log.Println("\nBot shutting down")

	bot.CleanupCommands(botConfig, session)
}

func waitForInterrupt() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-done
}
