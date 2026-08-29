package handler

import (
	"context"
	"regexp"
	"strings"

	"github.com/Maxxxxxx-x/whatsapp-bot/internal/command"
	"github.com/Maxxxxxx-x/whatsapp-bot/internal/config"
	"github.com/Maxxxxxx-x/whatsapp-bot/internal/whatsapp"
	"github.com/rs/zerolog"
	"go.mau.fi/whatsmeow/types/events"
)

type MessageHandler struct {
	waService *whatsapp.Service
	router    *command.Router
	cfg       *config.Config
	log       zerolog.Logger
}

func NewMessageHandler(waService *whatsapp.Service, router *command.Router, cfg *config.Config, log zerolog.Logger) *MessageHandler {
	return &MessageHandler{
		waService: waService,
		router:    router,
		cfg:       cfg,
		log:       log,
	}
}

func (handler *MessageHandler) Handle(evt any) {
	msgEvt, ok := evt.(*events.Message)
	if !ok || msgEvt.Info.IsNewsletterStatus || msgEvt.Info.Timestamp.Before(handler.waService.StartTime) {
		return
	}

	text := whatsapp.ExtractText(msgEvt.Message)
	if text == "" {
		return
	}

	if !strings.HasPrefix(text, handler.cfg.CommandPrefix) {
		senderNum, err := handler.waService.GetPhoneNumber(context.Background(), msgEvt.Info.Sender)
		if err != nil {
			return
		}

		if senderNum != "000000000" && senderNum != "000000000" {
			return
		}

		hiRegex := regexp.MustCompile(`^(?i)h+i+`)
		if hiRegex.MatchString(text) {
			handler.waService.ReplyText(context.Background(), msgEvt, "haiiiii", true)
		}
		return
	}

	senderNum, isAuthorized := handler.isAuthorized(msgEvt)
	if !isAuthorized {
		handler.log.Warn().
			Str("sender", msgEvt.Info.Sender.User).
			Str("text", text).
			Msg("attempted to execute commands (unauthorized)")
		return
	}

	rawCommand := strings.TrimPrefix(text, handler.cfg.CommandPrefix)
	args := strings.Fields(rawCommand)
	if len(args) == 0 {
		return
	}
	cmdName := args[0]
	cmdArgs := args[1:]
	rawArgs := strings.TrimSpace(strings.TrimPrefix(rawCommand, cmdName))

	cmdCtx := &command.Context{
		Ctx:       context.Background(),
		WAService: handler.waService,
		Router:    handler.router,
		Event:     msgEvt,
		Args:      cmdArgs,
		RawArgs:   rawArgs,
		Log:       handler.log,
	}

	if !msgEvt.Info.IsGroup {
		number, err := handler.waService.GetPhoneNumber(context.Background(), msgEvt.Info.Chat)
		if err != nil {
			return
		}
		handler.log.Info().Str("sender", senderNum).Str("command", cmdName).Str("args", rawArgs).Str("chat", number).Msg("test")
	}

	if err := handler.router.Dispatch(cmdCtx, cmdName); err != nil {
		handler.waService.ReplyText(context.Background(), msgEvt, "No such command", true)
		handler.log.Debug().Err(err).Str("command", cmdName).Msg("skipped or failed command execution")
	}
	handler.log.Info().Str("command", cmdName).Str("caller", senderNum).Msg("executed command successfully!")
}

func (handler *MessageHandler) isAuthorized(evt *events.Message) (string, bool) {
	senderNum, err := handler.waService.GetPhoneNumber(context.Background(), evt.Info.Sender)
	if err != nil {
		return senderNum, false
	}

	if evt.Info.IsFromMe || senderNum == "000000000" || senderNum == "000000000" || senderNum == "000000000" {
		return senderNum, true
	}

	return senderNum, false
}
