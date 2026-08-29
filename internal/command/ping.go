package command

import (
	"fmt"
	"time"
)

type PingCommand struct {
}

func (cmd *PingCommand) GetName() string {
	return "ping"
}

func (cmd *PingCommand) GetDescription() string {
	return "replies with pong"
}

func (cmd *PingCommand) Execute(ctx *Context) error {
	latency := time.Since(ctx.Event.Info.Timestamp).Milliseconds()
	if latency < 0 {
		latency = 0
	}

	return ctx.WAService.ReplyText(ctx.Ctx, ctx.Event, fmt.Sprintf("🏓 Pong!\nLatency: `%dms`", latency), true)
}
