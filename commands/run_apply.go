package commands

import (
	"fmt"
	"strings"

	"github.com/chroju/tfcloud/tfc"
)

type RunApplyCommand struct {
	Command
	runID string
}

func (c *RunApplyCommand) Run(args []string) int {
	client, err := tfc.NewTfCloud("", "")
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	c.Client = client

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
	return "Applies a run that is paused waiting for confirmation after a plan"
}

const helpRunApply = `
Usage: tfcloud run apply <run ID>

  Applies a run that is paused waiting for confirmation after a plan.

Caution:
	This command is not recommended.
	We recommend that you review the results of the "terraform plan"
	in the Terraform Cloud GUI before approving it.
`
