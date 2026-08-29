package command

import (
	"context"

	"github.com/Maxxxxxx-x/whatsapp-bot/internal/whatsapp"
	"github.com/rs/zerolog"
	"go.mau.fi/whatsmeow/types/events"
)

type Context struct {
	Ctx       context.Context
	WAService *whatsapp.Service
	Router    *Router
	Event     *events.Message
	Args      []string
	RawArgs   string
	Log       zerolog.Logger
}

type Command interface {
	GetName() string
	GetDescription() string
	Execute(ctx *Context) error
}
