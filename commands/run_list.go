package commands

import (
	"encoding/json"
	"strings"

	"github.com/mitchellh/cli"
)

type RunListCommand struct {
	UI cli.Ui
}

func (c *RunListCommand) Run(args []string) int {
	address := args[0]
	token := args[1]
	workspace := args[2]
	tfc, err := NewTfCloud(address, token)
	if err != nil {
		c.UI.Error("Terraform Cloud token is not valid.")
		return 1
	}

	list, err := tfc.Client.Runs.List(tfc.ctx, workspace, nil)
	if err != nil {
		c.UI.Error("Error")
		return 1
	}
	for _, run := range list {
		c.UI.Output(json.Marshal(run))
	}
	return 0
}

func (c *RunListCommand) Help() string {
	return strings.TrimSpace(helpRunList)
}

func (c *RunListCommand) Synopsis() string {
	return "about terraform runnings"
}

const helpRunList = `
Usage: tfcloud run <subcommand>

SubCommands:
	list    List all current runs
	apply   Apply terraform run needs confirmation
`
