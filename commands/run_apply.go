package commands

import (
	"strings"
)

type RunApplyCommand struct {
	Command
}

func (c *RunApplyCommand) Run(args []string) int {
	if len(args) != 3 {
		c.UI.Error("Arguments is not valid.")
		c.UI.Info(c.Help())
		return 1
	}

	runID := args[0]

	if err := c.Client.RunApply(runID); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	return 0
}

func (c *RunApplyCommand) Help() string {
	return strings.TrimSpace(helpRunApply)
}

func (c *RunApplyCommand) Synopsis() string {
	return "Apply terraform run which needs a confirmation"
}

const helpRunApply = `
Usage: tfcloud run apply <run ID>
`
