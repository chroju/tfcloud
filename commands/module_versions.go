package commands

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
)

type ModuleVersionsCommand struct {
	Command
}

func (c *ModuleVersionsCommand) Run(args []string) int {
	organization := args[0]
	provider := args[1]
	name := args[2]

	result, err := c.Client.ModuleGet(organization, name, provider)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	out := new(bytes.Buffer)
	w := tabwriter.NewWriter(out, 0, 4, 1, ' ', 0)
	fmt.Fprintln(w, "VERSION\tSTATUS\tLINK")
	for _, v := range result.VersionStatuses {
		fmt.Fprintf(w, "%s\t%s\thttps://%s/app/%s/modules/view/%s/%s/%s\n",
			v.Version, v.Status, c.Client.Address(), organization, name, provider, v.Version)
	}
	w.Flush()
	c.UI.Output(out.String())
	return 0
}

func (c *ModuleVersionsCommand) Help() string {
	return strings.TrimSpace(helpWorkspaceList)
}

func (c *ModuleVersionsCommand) Synopsis() string {
	return "Show terraform cloud private module all versions"
}

const helpModuleVersions = `
Usage: tfcloud module versions <organization> <provider> <module name>
`
