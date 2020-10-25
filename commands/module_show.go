package commands

import (
	"strings"

	"github.com/chroju/tfcloud/tfc"
	"github.com/mitchellh/cli"
)

type ModuleShowCommand struct {
	UI cli.Ui
}

func (c *ModuleShowCommand) Run(args []string) int {
	organization := args[0]
	provider := args[1]
	name := args[2]

	address := args[len(args)-2]
	token := args[len(args)-1]
	client, err := tfc.NewTfCloud(address, token)
	if err != nil {
		c.UI.Error("Terraform Cloud token is not valid.")
		return 1
	}

	result, err := client.ModuleShow(organization, name, provider)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	var out string
	for _, v := range result.VersionStatuses {
		out = out + v.Version + "\n"
	}
	c.UI.Output(out)
	return 0
}

func (c *ModuleShowCommand) Help() string {
	return strings.TrimSpace(helpWorkspaceList)
}

func (c *ModuleShowCommand) Synopsis() string {
	return "Show terraform cloud private module details"
}

const helpModuleShow = `
Usage: tfcloud module show <organization> <provider> <module name>
`
