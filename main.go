package main

import (
	"fmt"
	"os"

	"github.com/chroju/tfcloud/commands"
	"github.com/mitchellh/cli"
)

const (
	app     = "tfcloud"
	version = "0.3.0"
)

var (
	// UI is a cli.Ui
	UI cli.Ui
	// Format is a format of output
	Format commands.Format
)

func init() {
	UI = &cli.ColoredUi{
		Ui: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
		WarnColor:  cli.UiColorYellow,
		ErrorColor: cli.UiColorRed,
	}
}

func main() {
	format := commands.Format(os.Getenv("TFCLOUD_FORMAT"))

	command := commands.Command{
		UI:     UI,
		Format: format,
	}

	c := cli.NewCLI(app, version)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"run": func() (cli.Command, error) {
			return &commands.RunCommand{Command: command}, nil
		},
		"run list": func() (cli.Command, error) {
			return &commands.RunListCommand{Command: command}, nil
		},
		"run apply": func() (cli.Command, error) {
			return &commands.RunApplyCommand{Command: command}, nil
		},
		"workspace": func() (cli.Command, error) {
			return &commands.WorkspaceCommand{Command: command}, nil
		},
		"workspace list": func() (cli.Command, error) {
			return &commands.WorkspaceListCommand{Command: command}, nil
		},
		"workspace upgrade": func() (cli.Command, error) {
			return &commands.WorkspaceUpgradeCommand{Command: command}, nil
		},
		"workspace view": func() (cli.Command, error) {
			return &commands.WorkspaceViewCommand{Command: command}, nil
		},
		"module": func() (cli.Command, error) {
			return &commands.ModuleCommand{Command: command}, nil
		},
		"module list": func() (cli.Command, error) {
			return &commands.ModuleListCommand{Command: command}, nil
		},
		"module versions": func() (cli.Command, error) {
			return &commands.ModuleVersionsCommand{Command: command}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		UI.Error(fmt.Sprintf("Error: %s", err))
	}

	os.Exit(exitStatus)
}
