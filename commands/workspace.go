package commands

import (
	"strings"

	"github.com/mitchellh/cli"
)

type WorkspaceCommand struct {
	UI cli.Ui
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
	list      List all terraform cloud workspaces
`
