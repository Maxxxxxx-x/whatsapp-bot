package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Maxxxxxx-x/whatsapp-bot/internal/command"
	"github.com/Maxxxxxx-x/whatsapp-bot/internal/config"
	"github.com/Maxxxxxx-x/whatsapp-bot/internal/handler"
	"github.com/Maxxxxxx-x/whatsapp-bot/internal/logger"
	"github.com/Maxxxxxx-x/whatsapp-bot/internal/whatsapp"
)

func main() {
	cfg := config.Load()

	log := logger.GetLogger(cfg)
	log.Info().Str("env", cfg.AppEnv).Msg("Starting bot...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	waService, err := whatsapp.NewService(ctx, cfg, log)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initiate whatsapp service")
	}
	defer waService.Close()

	cmdRouter := command.NewRouter(log)
	cmdRouter.Register(
		&command.PingCommand{},
		command.NewHelpCommand(cfg.CommandPrefix),
	)

	msgHandler := handler.NewMessageHandler(waService, cmdRouter, cfg, log)
	waService.RegisterEventHandler(msgHandler.Handle)

	if err := waService.Connect(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to connect to whatsapp")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Info().Msg("Shutting down...")
}
