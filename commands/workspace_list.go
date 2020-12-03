package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	flag "github.com/spf13/pflag"
)

type WorkspaceListCommand struct {
	Command
	format string
}

func (c *WorkspaceListCommand) Run(args []string) int {
	if len(args) == 0 {
		c.UI.Error("Arguments are not valid.")
		c.UI.Info(c.Help())
		return 1
	}
	organization := args[0]

	f := flag.NewFlagSet("workspace_list", flag.ContinueOnError)
	f.StringVar(&c.format, "output", "table", "output format (table, json)")
	if err := f.Parse(args); err != nil {
		c.UI.Error(fmt.Sprintf("Arguments are not valid: %s", err))
		c.UI.Info(c.Help())
		return 1
	}

	if c.format != "table" && c.format != "json" {
		c.UI.Error("--output must be 'table' or 'json'")
		c.UI.Info(c.Help())
		return 1
	}

	result, err := c.Client.WorkspaceList(organization)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	switch c.format {
	case "table":
		out := new(bytes.Buffer)
		w := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tVERSION\tLINK")
		for _, r := range result {
			fmt.Fprintf(w, "%s\t%s\thttps://%s/app/%s/workspaces/%s\n",
				r.Name, r.TerraformVersion, c.Client.Address(), organization, r.Name)
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
