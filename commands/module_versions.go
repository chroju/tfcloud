package commands

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
)

type ModuleVersionsCommand struct {
	Command
	organization string
	provider     string
	name         string
}

func (c *ModuleVersionsCommand) Run(args []string) int {
	if len(args) != 3 {
		c.UI.Error("Arguments are not valid")
		c.UI.Info(c.Help())
		return 1
	}
	c.organization = args[0]
	c.provider = args[1]
	c.name = args[2]

	result, err := c.Client.ModuleGet(c.organization, c.name, c.provider)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	out := new(bytes.Buffer)
	w := tabwriter.NewWriter(out, 0, 4, 1, ' ', 0)
	fmt.Fprintln(w, "VERSION\tSTATUS\tLINK")
	for _, v := range result.VersionStatuses {
		fmt.Fprintf(w, "%s\t%s\thttps://%s/app/%s/modules/view/%s/%s/%s\n",
			v.Version, v.Status, c.Client.Address(), c.organization, c.name, c.provider, v.Version)
	}
	w.Flush()
	c.UI.Output(out.String())
	return 0
}

func (c *ModuleVersionsCommand) Help() string {
	return strings.TrimSpace(helpModuleVersions)
}

func (c *ModuleVersionsCommand) Synopsis() string {
	return "Show terraform cloud private module all versions"
}

const helpModuleVersions = `
Usage: tfcloud module versions <organization> <provider> <module name>
`
