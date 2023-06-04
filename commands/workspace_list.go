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

type WorkspaceListCommand struct {
	Command
	organization string
}

func (c *WorkspaceListCommand) Run(args []string) int {
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

	c.organization = f.Arg(0)

	client, err := tfc.NewTfCloud("", "")
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	c.Client = client

	wslist, err := c.Client.WorkspaceList(c.organization)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	switch c.Command.Format {
	case FormatAlfred:
		alfredItems := make([]AlfredFormatItem, len(wslist))
		for i, v := range wslist {
			alfredItems[i] = AlfredFormatItem{
				Title:        *v.Name,
				SubTitle:     fmt.Sprintf("vcs repo: %s", *v.VCSRepoName),
				Arg:          fmt.Sprintf("%s/app/%s/workspaces/%s", c.Client.Address(), c.organization, *v.Name),
				Match:        strings.ReplaceAll(*v.Name, "-", " "),
				AutoComplete: *v.Name,
				UID:          *v.ID,
			}
		}
		out, err := AlfredFormatOutput(alfredItems, "No workspaces found")
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		c.UI.Output(out)
	case FormatJSON:
		out, err := json.MarshalIndent(wslist, "", "  ")
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		c.UI.Output(string(out))
	default:
		out := new(bytes.Buffer)
		w := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tVERSION\tLINK")
		for _, v := range wslist {
			fmt.Fprintf(w, "%s\t%s\t%s/app/%s/workspaces/%s\n",
				*v.Name, *v.TerraformVersion, c.Client.Address(), c.organization, *v.Name)
		}
		w.Flush()
		c.UI.Output(out.String())
	}
	return 0
}

func (c *WorkspaceListCommand) Help() string {
	return strings.TrimSpace(helpWorkspaceList)
}

func (c *WorkspaceListCommand) Synopsis() string {
	return "Lists all terraform cloud workspaces"
}

const helpWorkspaceList = `
Usage: tfcloud workspace list [OPTIONS] <organization>

  Lists all terraform cloud workspaces

Options:
  --format, -f             Output format. Available formats: json, table (default: table)
`
