package main

import (
	"fmt"
	"os"

	"github.com/chroju/tfcloud/commands"
	"github.com/chroju/tfcloud/tfparser"
	"github.com/mitchellh/cli"
)

const (
	app     = "tfcloud"
	version = "0.0.1"
)

func main() {
	c := cli.NewCLI(app, version)
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}
	defaultUI := &cli.ColoredUi{Ui: ui, WarnColor: cli.UiColorYellow, ErrorColor: cli.UiColorRed}

	var terraformrcPath string
	if terraformrcPath = os.Getenv("TF_CLI_CONFIG_FILE"); terraformrcPath == "" {
		terraformrcPath = os.Getenv("HOME") + "/.terraformrc"
	}
	credential, err := tfparser.ParseTerraformrc(terraformrcPath)
	if err != nil {
		ui.Error(fmt.Sprintf("Error: %s", err))
	}

	// c.Args[0] is Terraform cloud endpoint, [1] is API token, [2:] are command line arguments.
	commonArgs := []string{credential.Name, credential.Token}
	c.Args = append(os.Args[1:], commonArgs...)

	c.Commands = map[string]cli.CommandFactory{
		"run": func() (cli.Command, error) {
			return &commands.RunCommand{UI: defaultUI}, nil
		},
		"run list": func() (cli.Command, error) {
			return &commands.RunListCommand{UI: defaultUI}, nil
		},
		"run apply": func() (cli.Command, error) {
			return &commands.RunApplyCommand{UI: defaultUI}, nil
		},
		"workspace": func() (cli.Command, error) {
			return &commands.WorkspaceCommand{UI: defaultUI}, nil
		},
		"workspace list": func() (cli.Command, error) {
			return &commands.WorkspaceListCommand{UI: defaultUI}, nil
		},
		"workspace upgrade": func() (cli.Command, error) {
			return &commands.WorkspaceUpgradeCommand{UI: defaultUI}, nil
		},
		"module": func() (cli.Command, error) {
			return &commands.ModuleCommand{UI: defaultUI}, nil
		},
		"module list": func() (cli.Command, error) {
			return &commands.ModuleListCommand{UI: defaultUI}, nil
		},
		"module versions": func() (cli.Command, error) {
			return &commands.ModuleVersionsCommand{UI: defaultUI}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		ui.Error(fmt.Sprintf("Error: %s", err))
	}

	os.Exit(exitStatus)
}
