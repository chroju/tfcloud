package main

import (
	"fmt"
	"os"

	"github.com/chroju/tfcloud/commands"
	"github.com/mitchellh/cli"
)

const (
	app         = "tfcloud"
	version     = "0.0.1"
	tfcEndpoint = "app.terraform.io"
)

func main() {
	c := cli.NewCLI(app, version)
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	var terraformrcPath string
	if terraformrcPath = os.Getenv("TF_CLI_CONFIG_FILE"); terraformrcPath == "" {
		terraformrcPath = os.Getenv("HOME") + "/.terraformrc"
	}
	token, err := commands.ParseTerraformrc(terraformrcPath)
	if err != nil {
		ui.Error(fmt.Sprintf("Error: %s", err))
	}

	commonArgs := []string{tfcEndpoint, token}
	c.Args = append(os.Args[1:], commonArgs...)

	c.Commands = map[string]cli.CommandFactory{
		"run": func() (cli.Command, error) {
			return &commands.RunCommand{UI: &cli.ColoredUi{Ui: ui, WarnColor: cli.UiColorYellow, ErrorColor: cli.UiColorRed}}, nil
		},
		"run list": func() (cli.Command, error) {
			return &commands.RunListCommand{UI: &cli.ColoredUi{Ui: ui, WarnColor: cli.UiColorYellow, ErrorColor: cli.UiColorRed}}, nil
		},
		"run apply": func() (cli.Command, error) {
			return &commands.RunApplyCommand{UI: &cli.ColoredUi{Ui: ui, WarnColor: cli.UiColorYellow, ErrorColor: cli.UiColorRed}}, nil
		},
		"workspace list": func() (cli.Command, error) {
			return &commands.WorkspaceListCommand{UI: &cli.ColoredUi{Ui: ui, WarnColor: cli.UiColorYellow, ErrorColor: cli.UiColorRed}}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		ui.Error(fmt.Sprintf("Error: %s", err))
	}

	os.Exit(exitStatus)
}
