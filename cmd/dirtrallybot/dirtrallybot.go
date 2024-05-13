package main

import (
	"log/slog"
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
		slog.Error("error creating session", "err", err)
		os.Exit(1)
	}

	bot.RegisterHandlers(session)
	bot.CreateCommands(botConfig, session)

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentGuildMessageReactions

	err = session.Open()
	if err != nil {
		slog.Error("error opening connection", "err", err)
		os.Exit(1)
	}

	slog.Info("Bot is running. Press CTL-C to exit.")
	defer session.Close()

	waitForInterrupt()

	slog.Info("Bot shutting down")

	bot.CleanupCommands(botConfig, session)
}

func waitForInterrupt() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-done
}
