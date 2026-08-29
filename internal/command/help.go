package command

import (
	"fmt"
	"sort"
	"strings"
)

type HelpCommand struct {
	prefix string
}

func NewHelpCommand(prefix string) *HelpCommand {
	return &HelpCommand{prefix: prefix}
}

func (cmd *HelpCommand) GetName() string {
	return "help"
}

func (cmd *HelpCommand) GetDescription() string {
	return "list all commands"
}

func (command *HelpCommand) Execute(ctx *Context) error {
	cmds := ctx.Router.GetCommands()

	sort.Slice(cmds, func(i, j int) bool {
		return cmds[i].GetName() < cmds[j].GetName()
	})

	var strBuilder strings.Builder
	strBuilder.WriteString("Available Commands:\n")

	for _, cmd := range cmds {
		strBuilder.WriteString(fmt.Sprintf("*`%s%s`* - %s\n", command.prefix, cmd.GetName(), cmd.GetDescription()))
	}

	return ctx.WAService.ReplyText(ctx.Ctx, ctx.Event, strBuilder.String(), true)
}
