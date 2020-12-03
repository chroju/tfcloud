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
	autoApprove   bool
	rootDir       string
	versionString string
	version       *version.Version
}

func (c *WorkspaceUpgradeCommand) Run(args []string) int {
	currentDir, _ := os.Getwd()
	f := flag.NewFlagSet("workspace_upgrade", flag.ExitOnError)
	f.StringVar(&c.rootDir, "root-path", currentDir, "Terraform config root path (default: current directory)")
	f.StringVarP(&c.versionString, "upgrade-version", "u", "latest", "Terraform version to upgrade")
	f.BoolVar(&c.autoApprove, "auto-approve", false, "Automatic approval for upgrade")
	if err := f.Parse(args); err != nil {
		c.UI.Error(fmt.Sprintf("Arguments are not valid: %s", err))
		c.UI.Error(err.Error())
		return 1
	}

	remoteBackend, err := tfparser.ParseRemoteBackend(c.rootDir)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if c.versionString == "latest" {
		latest, err := tfrelease.Latest()
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		c.UI.Info(fmt.Sprintf("Latest terraform version is %s ...", latest.Version.String()))
		c.version = latest.Version
	} else {
		c.version, err = version.NewVersion(c.versionString)
		if err != nil {
			c.UI.Error(fmt.Sprintf("%s is not valid version", c.versionString))
			c.UI.Output(helpMessageUpgrade)
			return 1
		}
	}

	if !remoteBackend.RequiredVersion.Check(c.version) {
		c.UI.Error(fmt.Sprintf("Version %s is not compatible with required version '%s'", c.version.String(), remoteBackend.RequiredVersion.String()))
		return 3
	}

	ws, err := c.Client.WorkspaceGet(remoteBackend.Organization, remoteBackend.WorkspaceName)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if ws.TerraformVersion == c.version.String() {
		c.UI.Warn(fmt.Sprintf("Version %s is already set up.", c.version.String()))
		return 0
	}

	if !c.autoApprove {
		if yn, err := askForConfirmation(fmt.Sprintf("Upgraded: %s -> %s\n ?", ws.TerraformVersion, c.version)); err != nil {
			c.UI.Error(err.Error())
			return 2
		} else if !yn {
			c.UI.Info("Canceled.")
			return 0
		}
	}

	if err = c.Client.WorkspaceUpdateVersion(remoteBackend.Organization, remoteBackend.WorkspaceName, c.version.String()); err != nil {
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
  --upgrade-version, -u    Terraform version to upgrade.
                           It must be in the correct semantic version format like 0.12.1, v0.12.2 .
                           Or you can specify "latest" to automatically upgrade to the latest version.
                           (default: latest)
  --root-path              Terraform config root path. (default: current directory)
  --auto-approve           Skip interactive approval of upgrade.

`
