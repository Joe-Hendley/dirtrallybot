package main

import (
	"context"
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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("loading config", "err", err)
		os.Exit(1)
	}

	slog.Info("starting with config", "store", cfg.Store, "randomiser", cfg.Randomiser)

	store, err := store.New(cfg)
	if err != nil {
		slog.Error("initialising store", "err", err)
		os.Exit(1)
	}

	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		slog.Error("creating session", "err", err)
		os.Exit(1)
	}

	rallyBot, err := bot.New(ctx, cfg, store, session)
	if err != nil {
		slog.Error("starting bot", "err", err)
		os.Exit(1)
	}

	defer func() {
		if err := rallyBot.Shutdown(); err != nil {
			slog.Error("shutting down", "err", err)
		}
	}()

	go func() {
		if err := rallyBot.ServeWeb(); err != nil {
			slog.Error("challenge viewer stopped", "err", err)
		}
	}()

	if err := session.Open(); err != nil {
		slog.Error("opening connection", "err", err)
		os.Exit(1)
	}
	defer session.Close()

	slog.Info("Bot is running. Press CTRL-C to exit.")

	<-ctx.Done()
	fmt.Println()
	slog.Info("Bot shutting down")
}
