package commands

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

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

	out := new(bytes.Buffer)
	w := tabwriter.NewWriter(out, 0, 4, 1, ' ', 0)
	fmt.Fprintln(w, "VERSION\tSTATUS\tLINK")
	for _, v := range result.VersionStatuses {
		fmt.Fprintf(w, "%s\t%s\thttps://%s/app/%s/modules/view/%s/%s/%s\n", v.Version, v.Status, address, organization, name, provider, v.Version)
	}
	w.Flush()
	c.UI.Output(out.String())
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
