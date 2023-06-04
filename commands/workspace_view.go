package commands

import (
	"bytes"
	"encoding/json"
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
	rootDir      string
	workspace    string
	organization string
	web          bool
}

func (c *WorkspaceViewCommand) Run(args []string) int {
	var formatOpt string
	currentDir, _ := os.Getwd()
	f := flag.NewFlagSet("workspace_view", flag.ExitOnError)
	f.StringVar(&c.rootDir, "root-path", currentDir, "Terraform config root path")
	f.StringVar(&c.organization, "org", "", "Specify organization name directly, must used with --workspace")
	f.StringVar(&c.workspace, "workspace", "", "Specify workspace name directly, must used with --org")
	f.BoolVarP(&c.web, "web", "w", false, "Show in the web browser")
	f.StringVarP(&formatOpt, "format", "f", "", "Output format. Available formats: json, table")
	if err := f.Parse(args); err != nil {
		c.UI.Error(fmt.Sprintf("Arguments are not valid: %s", err))
		c.UI.Error(err.Error())
		return 1
	}

	// org and workspace are must not be used with --root-path
	if c.rootDir != currentDir && (c.organization != "" || c.workspace != "") {
		c.UI.Error("You can't use --root-path with --organization and -workspace")
		return 1
	}

	// org and workspace are must be combined each other
	if (c.organization == "" && c.workspace != "") || (c.organization != "" && c.workspace == "") {
		c.UI.Error("You must specify both --organization and --workspace")
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

	workspace := c.workspace
	organization := c.organization
	if c.organization == "" && c.workspace == "" {
		rb, err := tfparser.ParseRemoteBackend(c.rootDir)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		workspace = rb.WorkspaceName
		organization = rb.Organization
	}

	ws, err := c.Client.WorkspaceGet(organization, workspace)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	url := fmt.Sprintf("%s/app/%s/workspaces/%s", c.Client.Address(), organization, workspace)

	if c.web {
		openbrowser(url)
	} else if c.Command.Format == FormatJSON {
		out, err := json.MarshalIndent(ws, "", "  ")
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		c.UI.Output(string(out))
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
  You can run this command in the directory where the target terraform file resides.
  You can also specify the target directory or the target workspace directly.

Options:
  --root-path              Terraform config root path. (default: current directory)
  --org					   Specify organization name directly, must used with --workspace
  --workspace              Specify workspace name directly, must used with --org
  --web, -w                View Terraform cloud workspace in a web browser.

  --format, -f             Output format. Available formats: json, table (default: table)
`
