package commands

import (
	"github.com/chroju/tfcloud/tfc"
	"github.com/mitchellh/cli"
)

type Command struct {
	Client tfc.TfCloud
	UI     cli.Ui
	Format Format
}
