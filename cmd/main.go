package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Joe-Hendley/dirtrallybot/internal/bot"
	"github.com/Joe-Hendley/dirtrallybot/internal/config"
	"github.com/Joe-Hendley/dirtrallybot/internal/store"
	"github.com/bwmarrin/discordgo"
)

func main() {
	cfg := config.New()
	store, err := store.New(cfg)
	if err != nil {
		slog.Error("initialising store:", "err", err)
		os.Exit(1)
	}

	session, err := discordgo.New("Bot " + cfg.Token)

	if err != nil {
		slog.Error("starting session: %w", "err", err)
		os.Exit(1)
	}

	rallyBot, err := bot.New(cfg, store, session)
	if err != nil {
		slog.Error("starting bot:", "err", err)
		os.Exit(1)
	}

	defer rallyBot.Shutdown()

	err = session.Open()
	if err != nil {
		slog.Error("opening connection:", "err", err)
		os.Exit(1)
	}

	slog.Info("Bot is running. Press CTL-C to exit.")
	defer session.Close()

	waitForInterrupt()
	fmt.Println()
	slog.Info("Bot shutting down")
}

func waitForInterrupt() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-done
}
