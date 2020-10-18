package commands

import (
	"strings"

	"github.com/chroju/tfcloud/tfc"
	"github.com/mitchellh/cli"
)

type RunApplyCommand struct {
	UI cli.Ui
}

func (c *RunApplyCommand) Run(args []string) int {
	if len(args) != 3 {
		c.UI.Error("Arguments is not valid.")
		c.UI.Info(c.Help())
		return 1
	}
	runID := args[0]
	address := args[len(args)-2]
	token := args[len(args)-1]
	client, err := tfc.NewTfCloud(address, token)
	if err != nil {
		c.UI.Error("Terraform Cloud token is not valid.")
		return 1
	}

	if err := client.RunApply(runID); err != nil {
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
