package commands

import (
	"strings"
)

type WorkspaceCommand struct {
	Command
}

func (c *WorkspaceCommand) Run(args []string) int {
	c.UI.Output(strings.TrimSpace(helpWorkspace))
	return 2
}

func (c *WorkspaceCommand) Help() string {
	return strings.TrimSpace(helpWorkspace)
}

func (c *WorkspaceCommand) Synopsis() string {
	return "about terraform cloud workspaces"
}

const helpWorkspace = `
Usage: tfcloud workspace <subcommand>

SubCommands:
	list      Lists all terraform cloud workspaces
	upgrade   Upgrades the terraform version of the current workspace
`
