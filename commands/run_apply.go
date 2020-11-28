package commands

import (
	"fmt"
	"strings"
)

type RunApplyCommand struct {
	Command
	runID string
}

func (c *RunApplyCommand) Run(args []string) int {
	if len(args) != 3 {
		c.UI.Error("Arguments are not valid.")
		c.UI.Info(c.Help())
		return 1
	}

	c.runID = args[0]

	if err := c.Client.RunApply(c.runID); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Info(fmt.Sprintf("ID %s has been applied successfully.", c.runID))

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
