package commands

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/chroju/tfcloud/tfc"
	"github.com/mitchellh/cli"
)

type WorkspaceListCommand struct {
	UI cli.Ui
}

func (c *WorkspaceListCommand) Run(args []string) int {
	if len(args) != 3 {
		c.UI.Error("Arguments is not valid.")
		c.UI.Info(c.Help())
		return 1
	}
	organization := args[0]
	address := args[len(args)-2]
	token := args[len(args)-1]
	client, err := tfc.NewTfCloud(address, token)
	if err != nil {
		c.UI.Error("Terraform Cloud token is not valid.")
		return 1
	}

	result, err := client.WorkspaceList(organization)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	out := new(bytes.Buffer)
	w := tabwriter.NewWriter(out, 0, 4, 1, ' ', 0)
	fmt.Fprintln(w, "NAME\tID\tLINK")
	for _, r := range result {
		fmt.Fprintf(w, "%s\t%s\thttps://%s/app/%s/workspaces/%s\n", r.Name, r.ID, address, organization, r.Name)
	}
	w.Flush()
	c.UI.Output(out.String())
	return 0
}

func (c *WorkspaceListCommand) Help() string {
	return strings.TrimSpace(helpWorkspaceList)
}

func (c *WorkspaceListCommand) Synopsis() string {
	return "List all terraform cloud workspaces"
}

const helpWorkspaceList = `
Usage: tfcloud workspace list <organization>
`
