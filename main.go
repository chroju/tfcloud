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
	c.Args = os.Args[1:]
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	c.Commands = map[string]cli.CommandFactory{
		"run": func() (cli.Command, error) {
			return &commands.RunCommand{UI: &cli.ColoredUi{Ui: ui, WarnColor: cli.UiColorYellow, ErrorColor: cli.UiColorRed}}, nil
		},
		"run list": func() (cli.Command, error) {
			return &commands.RunListCommand{UI: &cli.ColoredUi{Ui: ui, WarnColor: cli.UiColorYellow, ErrorColor: cli.UiColorRed}}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		ui.Error(fmt.Sprintf("Error: %s", err))
	}

	os.Exit(exitStatus)
}
