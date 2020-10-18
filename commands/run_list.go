package commands

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/chroju/tfcloud/tfc"
	"github.com/mitchellh/cli"
)

type RunListCommand struct {
	UI cli.Ui
}

func (c *RunListCommand) Run(args []string) int {
	organization := args[0]
	address := args[1]
	token := args[2]
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

	out := new(bytes.Buffer)
	w := tabwriter.NewWriter(out, 0, 4, 1, ' ', 0)
	fmt.Fprintln(w, "WORKSPACE\tSTATUS\tNEEDS CONFIRM\tLINK")
	for _, r := range result {
		fmt.Fprintf(w, "%s\t%s\t%v\thttps://%s/app/%s/workspaces/%s/runs/%s\n", r.Workspace, r.Status, r.IsConfirmable, address, organization, r.Workspace, r.ID)
	}
	w.Flush()
	c.UI.Output(out.String())
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
