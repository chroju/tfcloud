package commands

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

// TfCloud represents Terraform Cloud API wrapper
type TfCloud struct {
	*tfe.Client
	ctx context.Context
}

type tfcloudImpl struct {
	*tfe.Client
	ctx context.Context
}

// NewTfCloud creates a new TfCloud interface
func NewTfCloud(address, token string) (*TfCloud, error) {
	config := &tfe.Config{
		Token: token,
	}

	client, err := tfe.NewClient(config)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	return &TfCloud{
		client,
		ctx,
	}, nil
}
