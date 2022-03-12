package commands

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/chroju/tfcloud/tfc"
)

type ModuleListCommand struct {
	Command
}

func (c *ModuleListCommand) Run(args []string) int {
	if len(args) != 0 {
		c.UI.Error("Arguments are not valid.")
		c.UI.Info(c.Help())
		return 1
	}

	client, err := tfc.NewTfCloud("", "")
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	c.Client = client

	mdlist, err := c.Client.ModuleList()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	switch c.Command.Format {
	case "alfred":
		alfredItems := make([]AlfredFormatItem, len(mdlist))
		for i, v := range mdlist {
			alfredItems[i] = AlfredFormatItem{
				Title:        v.ID,
				SubTitle:     fmt.Sprintf("source: %s", v.Source),
				Arg:          fmt.Sprintf("%s/app/%s/modules/view/%s/%s/%s", c.Client.Address(), v.Organization, v.Name, v.Provider, v.VersionStatuses[0].Version),
				Match:        v.Name,
				AutoComplete: v.Name,
				UID:          v.ID,
			}
		}
		out, err := AlfredFormatOutput(alfredItems, "No modules found")
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		c.UI.Output(out)
	default:
		out := new(bytes.Buffer)
		w := tabwriter.NewWriter(out, 0, 4, 1, ' ', 0)
		fmt.Fprintln(w, "NAME\tLATEST\tLINK")
		for _, r := range mdlist {
			latest := r.VersionStatuses[0].Version
			fmt.Fprintf(w, "%s\t%s\thttps://%s/app/%s/modules/view/%s/%s/%s\n",
				r.Name, latest, c.Client.Address(), r.Organization, r.Name, r.Provider, latest)
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
Usage: tfcloud module list

  Lists all terraform cloud private modules.
`
