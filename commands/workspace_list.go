package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/chroju/tfcloud/tfc"
	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
)

type WorkspaceListCommand struct {
	UI cli.Ui
}

func (c *WorkspaceListCommand) Run(args []string) int {
	if len(args) < 3 {
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

	buf := &bytes.Buffer{}
	var format string
	f := flag.NewFlagSet("module_list", flag.ContinueOnError)
	f.SetOutput(buf)
	f.StringVar(&format, "output", "table", "output format (table, json)")
	if err := f.Parse(args); err != nil {
		c.UI.Info(c.Help())
		return 1
	}
	if format != "table" && format != "json" {
		c.UI.Error("--output must be 'table' or 'json'")
		c.UI.Info(c.Help())
		return 1
	}

	result, err := client.WorkspaceList(organization)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	switch format {
	case "table":
		out := new(bytes.Buffer)
		w := tabwriter.NewWriter(out, 0, 4, 1, ' ', 0)
		fmt.Fprintln(w, "NAME\tID\tLINK")
		for _, r := range result {
			fmt.Fprintf(w, "%s\t%s\thttps://%s/app/%s/workspaces/%s\n", r.Name, r.ID, address, organization, r.Name)
		}
		w.Flush()
		c.UI.Output(out.String())
	case "json":
		out, err := json.Marshal(result)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		c.UI.Output(string(out))
	}
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
