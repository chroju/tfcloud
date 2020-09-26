package commands

import (
	"strings"

	"github.com/chroju/tfcloud/tfc"
	"github.com/mitchellh/cli"
)

type RunListCommand struct {
	UI cli.Ui
}

func (c *RunListCommand) Run(args []string) int {
	address := args[0]
	token := args[1]
	organization := args[2]
	client, err := tfc.NewTfCloud(address, token)
	if err != nil {
		c.UI.Error("Terraform Cloud token is not valid.")
		return 1
	}

	result, err := client.RunList(organization)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	c.UI.Output(string(result))
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
