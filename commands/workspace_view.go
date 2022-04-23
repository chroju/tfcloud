package commands

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/chroju/tfcloud/tfc"
	"github.com/chroju/tfcloud/tfparser"
	flag "github.com/spf13/pflag"
)

const timeFormat = "2006-01-02 15:04:05 (MST)"

type WorkspaceViewCommand struct {
	Command
	rootDir string
	web     bool
}

func (c *WorkspaceViewCommand) Run(args []string) int {
	currentDir, _ := os.Getwd()
	f := flag.NewFlagSet("workspace_view", flag.ExitOnError)
	f.StringVar(&c.rootDir, "root-path", currentDir, "Terraform config root path (default: current directory)")
	f.BoolVarP(&c.web, "web", "w", false, "Show in the web browser")
	if err := f.Parse(args); err != nil {
		c.UI.Error(fmt.Sprintf("Arguments are not valid: %s", err))
		c.UI.Error(err.Error())
		return 1
	}

	client, err := tfc.NewTfCloud("", "")
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	c.Client = client

	rb, err := tfparser.ParseRemoteBackend(c.rootDir)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	ws, err := c.Client.WorkspaceGet(rb.Organization, rb.WorkspaceName)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	url := fmt.Sprintf("%s/app/%s/workspaces/%s", c.Client.Address(), rb.Organization, rb.WorkspaceName)

	if c.web {
		openbrowser(url)
	} else {
		out := new(bytes.Buffer)
		w := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
		fmt.Fprintf(w, "%s\t%s\n", "Name:", *ws.Name)
		fmt.Fprintf(w, "%s\t%s\n", "Terraform Version:", *ws.TerraformVersion)
		fmt.Fprintf(w, "%s\t%s\n", "VCS Repo:", *ws.VCSRepoName)
		fmt.Fprintf(w, "%s\t%s\n", "Working Directory:", *ws.WorkingDirectory)
		fmt.Fprintf(w, "%s\t%s\n", "Execution mode:", *ws.ExecutionMode)
		fmt.Fprintf(w, "%s\t%v\n", "Auto Apply:", *ws.AutoApply)
		fmt.Fprintf(w, "%s\t%v\n", "Resource count:", *ws.ResourceCount)
		fmt.Fprintf(w, "%s\t%s\n", "Created at:", ws.CreatedAt.Format(timeFormat))
		fmt.Fprintf(w, "%s\t%s\n", "Updated at:", ws.UpdatedAt.Format(timeFormat))
		fmt.Fprintf(w, "%s\t%s\n", "URL:", url)
		w.Flush()
		c.UI.Output(out.String())
	}

	return 0
}

func (c *WorkspaceViewCommand) Help() string {
	return strings.TrimSpace(helpMessageWorkspaceView)
}

func (c *WorkspaceViewCommand) Synopsis() string {
	return "Views Terraform cloud workspace details"
}

const helpMessageWorkspaceView = `
Usage: tfcloud workspace view [OPTION]

  Views Terraform cloud workspace details.

Notes:
  This command works by reading the remote config in the current directory.
  You must run this command in the directory where the target terraform file resides.
  Or you can specify the target directory with the --root-path option.

Options:
  --root-path              Terraform config root path. (default: current directory)
  --web, -w                View Terraform cloud workspace in a web browser.

`
