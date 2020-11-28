package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/chroju/tfcloud/tfc"
)

func parseDefaultArgs(args []string) (tfc.TfCloud, error) {
	address := args[len(args)-2]
	token := args[len(args)-1]
	client, err := tfc.NewTfCloud(address, token)
	if err != nil {
		return nil, fmt.Errorf("Terraform Cloud token is not valid")
	}
	return client, nil
}

func askForConfirmation(s string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true, nil
		} else if response == "n" || response == "no" {
			return false, nil
		}
	}
}
