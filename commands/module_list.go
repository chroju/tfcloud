package commands

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/chroju/tfcloud/tfc"
	"github.com/mitchellh/cli"
)

type ModuleListCommand struct {
	UI cli.Ui
}

func (c *ModuleListCommand) Run(args []string) int {
	address := args[len(args)-2]
	token := args[len(args)-1]
	client, err := tfc.NewTfCloud(address, token)
	if err != nil {
		c.UI.Error("Terraform Cloud token is not valid.")
		return 1
	}

	result, err := client.ModuleList()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	out := new(bytes.Buffer)
	w := tabwriter.NewWriter(out, 0, 4, 1, ' ', 0)
	fmt.Fprintln(w, "NAME\tLATEST\tLINK")
	for _, r := range result {
		latest := r.VersionStatuses[0].Version
		fmt.Fprintf(w, "%s\t%s\thttps://%s/app/%s/modules/view/%s/%s/%s\n", r.Name, latest, address, r.Organization, r.Name, r.Provider, latest)
	}
	w.Flush()
	c.UI.Output(out.String())
	return 0
}

func (c *ModuleListCommand) Help() string {
	return strings.TrimSpace(helpWorkspaceList)
}

func (c *ModuleListCommand) Synopsis() string {
	return "List all terraform cloud private modules"
}

const helpModuleList = `
Usage: tfcloud module list
`
