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
	runID := args[0]
	address := args[1]
	token := args[2]
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
	return "about terraform runnings"
}

const helpRunApply = `
Usage: tfcloud run <subcommand>

SubCommands:
	list    List all current runs
	apply   Apply terraform run needs confirmation
`
