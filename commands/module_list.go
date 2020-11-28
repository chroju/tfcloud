package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	flag "github.com/spf13/pflag"
)

type ModuleListCommand struct {
	Command
}

func (c *ModuleListCommand) Run(args []string) int {
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

	result, err := c.Client.ModuleList()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	switch format {
	case "table":
		out := new(bytes.Buffer)
		w := tabwriter.NewWriter(out, 0, 4, 1, ' ', 0)
		fmt.Fprintln(w, "NAME\tLATEST\tLINK")
		for _, r := range result {
			latest := r.VersionStatuses[0].Version
			fmt.Fprintf(w, "%s\t%s\thttps://%s/app/%s/modules/view/%s/%s/%s\n",
				r.Name, latest, c.Client.Address(), r.Organization, r.Name, r.Provider, latest)
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

func (c *ModuleListCommand) Help() string {
	return strings.TrimSpace(helpWorkspaceList)
}

func (c *ModuleListCommand) Synopsis() string {
	return "List all terraform cloud private modules"
}

const helpModuleList = `
Usage: tfcloud module list
`
