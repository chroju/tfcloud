package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	flag "github.com/spf13/pflag"
)

type RunListCommand struct {
	Command
	organization string
	format       string
}

func (c *RunListCommand) Run(args []string) int {
	if len(args) < 3 {
		c.UI.Error("Arguments are not valid.")
		c.UI.Info(c.Help())
		return 1
	}

	f := flag.NewFlagSet("run_list", flag.ContinueOnError)
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

	c.organization = args[0]
	result, err := c.Client.RunList(c.organization)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	switch c.format {
	case "table":
		out := new(bytes.Buffer)
		w := tabwriter.NewWriter(out, 0, 4, 1, ' ', 0)
		fmt.Fprintln(w, "WORKSPACE\tSTATUS\tNEEDS CONFIRM\tLINK")
		for _, r := range result {
			fmt.Fprintf(w, "%s\t%s\t%v\thttps://%s/app/%s/workspaces/%s/runs/%s\n",
				r.Workspace, r.Status, r.IsConfirmable, c.Client.Address(), c.organization, r.Workspace, r.ID)
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

func (c *RunListCommand) Help() string {
	return strings.TrimSpace(helpRunList)
}

func (c *RunListCommand) Synopsis() string {
	return "List all current terraform runs"
}

const helpRunList = `
Usage: tfcloud run list <organization>
`
