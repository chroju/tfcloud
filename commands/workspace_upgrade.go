package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/chroju/tfcloud/tfparser"
	"github.com/chroju/tfcloud/tfrelease"
	version "github.com/hashicorp/go-version"
	flag "github.com/spf13/pflag"
)

type WorkspaceUpgradeCommand struct {
	Command
}

func (c *WorkspaceUpgradeCommand) Run(args []string) int {
	var approve bool
	var root, versionString string
	var updateVer *version.Version

	currentDir, _ := os.Getwd()
	f := flag.NewFlagSet("workspace_upgrade", flag.ExitOnError)
	f.StringVar(&root, "root-path", currentDir, "Terraform config root path (default: current directory)")
	f.StringVar(&versionString, "tfversion", "latest", "Terraform version to upgrade")
	f.BoolVar(&approve, "auto-approve", false, "Automatic approval for upgrade")
	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	remoteBackend, err := tfparser.ParseRemoteBackend(root)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if versionString == "latest" {
		latest, err := tfrelease.Latest()
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		c.UI.Info(fmt.Sprintf("Latest terraform version is %s ...", latest.Version.String()))
		updateVer = latest.Version
	} else {
		updateVer, err = version.NewVersion(versionString)
		if err != nil {
			c.UI.Error(fmt.Sprintf("%s is not valid version", versionString))
			c.UI.Output(helpMessageUpgrade)
			return 1
		}
	}

	if !remoteBackend.RequiredVersion.Check(updateVer) {
		c.UI.Error(fmt.Sprintf("Version %s is not compatible with required version '%s'", updateVer.String(), remoteBackend.RequiredVersion.String()))
		return 3
	}

	ws, err := c.Client.WorkspaceGet(remoteBackend.Organization, remoteBackend.WorkspaceName)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if ws.TerraformVersion == updateVer.String() {
		c.UI.Warn(fmt.Sprintf("Version %s is already set up.", updateVer.String()))
		return 0
	}

	if !approve {
		if yn, err := askForConfirmation(fmt.Sprintf("Upgraded: %s -> %s\n ?", ws.TerraformVersion, updateVer)); err != nil {
			c.UI.Error(err.Error())
			return 2
		} else if !yn {
			c.UI.Info("Canceled.")
			return 0
		}
	}

	if err = c.Client.WorkspaceUpdateVersion(remoteBackend.Organization, remoteBackend.WorkspaceName, updateVer.String()); err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	c.UI.Info("Successfully upgraded.")
	return 0
}

func (c *WorkspaceUpgradeCommand) Help() string {
	return strings.TrimSpace(helpMessageUpgrade)
}

func (c *WorkspaceUpgradeCommand) Synopsis() string {
	return "Upgrade Terraform cloud workspace terraform version"
}

const helpMessageUpgrade = `
Usage: tfcloud workspace upgrade [OPTION]

Notes:
  This command works by reading the remote config in the current directory.

Options:
  --root-path       Terraform config root path. (default: current directory)
  --tfversion       Terraform version to upgrade.
				    It must be in the correct semantic version format like 0.12.1, v0.12.2 .
				    Or you can specify "latest" to automatically upgrade to the latest version.
				    (default: latest)
  --auto-approve    Skip interactive approval of upgrade.

`
