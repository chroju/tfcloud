package commands

import (
	"strings"

	"github.com/chroju/tfcloud/tfc"
	"github.com/mitchellh/cli"
)

type ModuleVersionsCommand struct {
	UI cli.Ui
}

func (c *ModuleVersionsCommand) Run(args []string) int {
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

	result, err := client.ModuleVersions(organization, name, provider)
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

func (c *ModuleVersionsCommand) Help() string {
	return strings.TrimSpace(helpWorkspaceList)
}

func (c *ModuleVersionsCommand) Synopsis() string {
	return "Show terraform cloud private module all versions"
}

const helpModuleVersions = `
Usage: tfcloud module versions <organization> <provider> <module name>
`
