package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/chroju/tfcloud/tfc"
	flag "github.com/spf13/pflag"
)

type RunListCommand struct {
	Command
	organization string
	localTZ      *time.Location
}

func (c *RunListCommand) Run(args []string) int {
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
	localTZ, _ := time.LoadLocation("Local")
	c.localTZ = localTZ

	runlist, err := c.Client.RunList(c.organization)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	switch c.Command.Format {
	case FormatAlfred:
		alfredItems := make([]AlfredFormatItem, len(runlist))
		for i, v := range runlist {
			localtime := v.CreatedAt.In(c.localTZ)
			alfredItems[i] = AlfredFormatItem{
				Title:        *v.Workspace,
				SubTitle:     fmt.Sprintf("%s / %s", *v.Status, localtime),
				Arg:          fmt.Sprintf("%s/app/%s/workspaces/%s/runs/%s", c.Client.Address(), c.organization, *v.Workspace, *v.ID),
				Match:        strings.ReplaceAll(*v.Workspace, "-", " "),
				AutoComplete: *v.Workspace,
				UID:          *v.ID,
			}
		}
		out, err := AlfredFormatOutput(alfredItems, "No runs found")
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		c.UI.Output(out)
	case FormatJSON:
		out, err := json.MarshalIndent(runlist, "", "  ")
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		c.UI.Output(string(out))
	default:
		out := new(bytes.Buffer)
		w := tabwriter.NewWriter(out, 0, 4, 1, ' ', 0)
		fmt.Fprintln(w, "WORKSPACE\tSTATUS\tNEEDS CONFIRM\tLINK")
		for _, r := range runlist {
			fmt.Fprintf(w, "%s\t%s\t%v\t%s/app/%s/workspaces/%s/runs/%s\n",
				*r.Workspace, *r.Status, *r.IsConfirmable, c.Client.Address(), c.organization, *r.Workspace, *r.ID)
		}
		w.Flush()
		c.UI.Output(out.String())
	}
	return 0
}

func (c *RunListCommand) Help() string {
	return strings.TrimSpace(helpRunList)
}

func (c *RunListCommand) Synopsis() string {
	return "Lists all current terraform runs in the organization."
}

const helpRunList = `
Usage: tfcloud run list [OPTIONS] <organization>

  Lists all current terraform runs in the organization.

Options:
  --format, -f             Output format. Available formats: json, table (default: table)
`
