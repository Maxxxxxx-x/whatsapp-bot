package command

type WhitelistCommand struct {
}

func (cmd *WhitelistCommand) GetName() string {
	return "whitelist"
}

func (cmd *WhitelistCommand) GetDescription() string {
	return "Whitelists a user in the dm / mentioned to use certian command categories"
}
