package tfc

import (
	"context"
	"encoding/json"

	tfe "github.com/hashicorp/go-tfe"
)

var ListOptions = &tfe.ListOptions{
	PageNumber: 0,
	PageSize:   100,
}

type tfclient struct {
	client *tfe.Client
	ctx    context.Context
}

// Client represents Terraform Cloud API client
type TfCloud interface {
	RunList(organization string) ([]byte, error)
	RunGet(workspaceID string) ([]byte, error)
}

// NewTfCloud creates a new TfCloud interface
func NewTfCloud(address, token string) (TfCloud, error) {
	config := &tfe.Config{
		Token: token,
	}

	client, err := tfe.NewClient(config)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	return &tfclient{
		client,
		ctx,
	}, nil
}

func (c *tfclient) RunList(organization string) ([]byte, error) {
	type result struct {
		Error    error
		Response []byte
	}
	wlo := &tfe.WorkspaceListOptions{
		ListOptions: *ListOptions,
		Search:      nil,
	}

	wslist, err := c.client.Workspaces.List(c.ctx, organization, *wlo)
	if err != nil {
		return nil, err
	}

	resultChan := make(chan result)
	for _, ws := range wslist.Items {
		go func(ws *tfe.Workspace) {
			runs, err := c.RunGet(ws.ID)
			resultChan <- result{Error: err, Response: runs}
		}(ws)
	}

	var rtn []byte
	for range wslist.Items {
		run := <-resultChan
		if run.Error != nil {
			return nil, err
		}
		rtn = append(rtn, run.Response...)
	}

	return rtn, nil
}

func (c *tfclient) RunGet(workspaceID string) ([]byte, error) {
	rlo := &tfe.RunListOptions{
		ListOptions: *ListOptions,
	}

	runlist, err := c.client.Runs.List(c.ctx, workspaceID, *rlo)
	if err != nil {
		return nil, err
	}

	var result []byte
	for _, run := range runlist.Items {
		if !checkRunInAction(run) {
			continue
		}
		runJSON, err := json.Marshal(run)
		if err != nil {
			return nil, err
		}
		result = append(result, runJSON...)
	}

	return result, nil
}

func checkRunInAction(run *tfe.Run) bool {
	if run.Status == tfe.RunApplied ||
		run.Status == tfe.RunCanceled ||
		run.Status == tfe.RunErrored ||
		run.Status == tfe.RunDiscarded ||
		run.Status == tfe.RunPlannedAndFinished ||
		run.Status == "" {
		return false
	}
	return true
}
