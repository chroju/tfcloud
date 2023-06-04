package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/chroju/tfcloud/tfc"
	flag "github.com/spf13/pflag"
)

type ModuleListCommand struct {
	organization string
	Command
}

func (c *ModuleListCommand) Run(args []string) int {
	var formatOpt string
	f := flag.NewFlagSet("module_list", flag.ExitOnError)
	f.StringVarP(&formatOpt, "format", "f", "", "Output format. Available formats: json, table")
	if err := f.Parse(args); err != nil {
		c.UI.Error(fmt.Sprintf("Arguments are not valid: %s", err))
		c.UI.Error(err.Error())
		return 1
	}

	if formatOpt != "" {
		c.Command.Format = Format(formatOpt)
	}

	client, err := tfc.NewTfCloud("", "")
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	c.Client = client
	c.organization = f.Arg(0)

	mdlist, err := c.Client.ModuleList(c.organization)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	switch c.Command.Format {
	case FormatAlfred:
		alfredItems := make([]AlfredFormatItem, len(mdlist))
		for i, v := range mdlist {
			alfredItems[i] = AlfredFormatItem{
				Title:        *v.ID,
				SubTitle:     fmt.Sprintf("source: %s", *v.Source),
				Arg:          fmt.Sprintf("%s/app/%s/modules/view/%s/%s/%s", c.Client.Address(), *v.Organization, *v.Name, *v.Provider, v.VersionStatuses[0].Version),
				Match:        *v.Name,
				AutoComplete: *v.Name,
				UID:          *v.ID,
			}
		}
		out, err := AlfredFormatOutput(alfredItems, "No modules found")
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		c.UI.Output(out)
	case FormatJSON:
		out, err := json.MarshalIndent(mdlist, "", "  ")
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		c.UI.Output(string(out))
	default:
		out := new(bytes.Buffer)
		w := tabwriter.NewWriter(out, 0, 4, 1, ' ', 0)
		fmt.Fprintln(w, "NAME\tLATEST\tLINK")
		for _, r := range mdlist {
			latest := r.VersionStatuses[0].Version
			fmt.Fprintf(w, "%s\t%s\t%s/app/%s/modules/view/%s/%s/%s\n",
				*r.Name, latest, c.Client.Address(), *r.Organization, *r.Name, *r.Provider, latest)
		}
		w.Flush()
		c.UI.Output(out.String())
	}
	return 0
}

func (c *ModuleListCommand) Help() string {
	return strings.TrimSpace(helpModuleList)
}

func (c *ModuleListCommand) Synopsis() string {
	return "Lists all terraform cloud private modules"
}

const helpModuleList = `
Usage: tfcloud module list [OPTIONS] <organization>

  Lists all terraform cloud private modules.

Options:
  --format, -f             Output format. Available formats: json, table (default: table)
`
