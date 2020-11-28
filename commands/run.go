package commands

import (
	"strings"
)

type RunCommand struct {
	Command
}

func (c *RunCommand) Run(args []string) int {
	c.UI.Output(strings.TrimSpace(helpRun))
	return 2
}

func (c *RunCommand) Help() string {
	return strings.TrimSpace(helpRun)
}

func (c *RunCommand) Synopsis() string {
	return "about terraform runnings"
}

const helpRun = `
Usage: tfcloud run <subcommand>

SubCommands:
	list    List all current terraform runs
	apply   Apply terraform run which needs a confirmation
`
