package commands

import (
	"strings"
)

type ModuleCommand struct {
	Command
}

func (c *ModuleCommand) Run(args []string) int {
	c.UI.Output(strings.TrimSpace(helpModule))
	return 2
}

func (c *ModuleCommand) Help() string {
	return strings.TrimSpace(helpModule)
}

func (c *ModuleCommand) Synopsis() string {
	return "about terraform private modules"
}

const helpModule = `
Usage: tfcloud module <subcommand>

SubCommands:
	list        List all terraform private modules in your account
	versions    List all terraform private module versions
`
