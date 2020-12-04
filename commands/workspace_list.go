package commands

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
)

type WorkspaceListCommand struct {
	Command
	organization string
}

func (c *WorkspaceListCommand) Run(args []string) int {
	if len(args) != 1 {
		c.UI.Error("Arguments are not valid.")
		c.UI.Info(c.Help())
		return 1
	}
	c.organization = args[0]

	wslist, err := c.Client.WorkspaceList(c.organization)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	switch c.Command.Format {
	case "alfred":
		alfredItems := make([]AlfredFormatItem, len(wslist))
		for i, v := range wslist {
			alfredItems[i] = AlfredFormatItem{
				Title:        v.Name,
				SubTitle:     fmt.Sprintf("vcs repo: %s", v.VCSRepoName),
				Arg:          fmt.Sprintf("https://%s/app/%s/workspaces/%s", c.Client.Address(), c.organization, v.Name),
				Match:        strings.ReplaceAll(v.Name, "-", " "),
				AutoComplete: v.Name,
				UID:          v.ID,
			}
		}
		out, err := AlfredFormatOutput(alfredItems, "No workspaces found")
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		c.UI.Output(out)
	default:
		out := new(bytes.Buffer)
		w := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tVERSION\tLINK")
		for _, v := range wslist {
			fmt.Fprintf(w, "%s\t%s\thttps://%s/app/%s/workspaces/%s\n",
				v.Name, v.TerraformVersion, c.Client.Address(), c.organization, v.Name)
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
Usage: tfcloud workspace list <organization>

  Lists all terraform cloud workspaces
`
