package command

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
)

type Router struct {
	commands map[string]Command
	log      zerolog.Logger
}

func NewRouter(log zerolog.Logger) *Router {
	return &Router{
		commands: make(map[string]Command),
		log:      log,
	}
}

func (r *Router) Register(cmds ...Command) {
	for _, cmd := range cmds {
		r.commands[strings.ToLower(cmd.GetName())] = cmd
		r.log.Debug().Str("command", cmd.GetName()).Msg("registered command")
	}
}

func (r *Router) Dispatch(ctx *Context, cmdName string) error {
	cmd, exists := r.commands[strings.ToLower(cmdName)]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmdName)
	}
	return cmd.Execute(ctx)
}

func (r *Router) GetCommands() []Command {
	list := make([]Command, 0, len(r.commands))

	for _, cmd := range r.commands {
		list = append(list, cmd)
	}

	return list
}
