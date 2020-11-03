package tfc

import (
	"context"
	"time"

	tfe "github.com/hashicorp/go-tfe"
)

var ListOptions = &tfe.ListOptions{
	PageNumber: 0,
	PageSize:   100,
}

type tfclient struct {
	client         *tfe.Client
	registryClient *RegistryClient
	ctx            context.Context
}

type Run struct {
	ID            string
	Organization  string
	Workspace     string
	Status        string
	IsConfirmable bool
	CreatedAt     time.Time
}

type Workspace struct {
	ID               string
	Name             string
	TerraformVersion string
	CurrentRun       *tfe.Run
}

type RegistryModule struct {
	ID              string
	Name            string
	Provider        string
	VersionStatuses []tfe.RegistryModuleVersionStatuses
	Organization    string
	VCSRepo         string
}

type RegistryModuleVersionStatuses struct {
	Version string
}

// Client represents Terraform Cloud API client
type TfCloud interface {
	RunList(organization string) ([]*Run, error)
	RunGet(workspaceID, WorkspaceName string) (*Run, error)
	RunApply(RunID string) error
	WorkspaceList(organization string) ([]*Workspace, error)
	ModuleList() ([]*RegistryModule, error)
	ModuleVersions(organization, name, provider string) (*RegistryModule, error)
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

	registryClient, err := NewRegistryClient(config)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	return &tfclient{
		client,
		registryClient,
		ctx,
	}, nil
}

func (c *tfclient) RunList(organization string) ([]*Run, error) {
	type result struct {
		Error    error
		Response *Run
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
			run, err := c.RunGet(ws.ID, ws.Name)
			resultChan <- result{Error: err, Response: run}
		}(ws)
	}

	var rtn []*Run
	for range wslist.Items {
		run := <-resultChan
		if run.Error != nil {
			return nil, err
		}
		if run.Response != nil {
			rtn = append(rtn, run.Response)
		}
	}

	return rtn, nil
}

func (c *tfclient) RunGet(workspaceID, WorkspaceName string) (*Run, error) {
	rlo := &tfe.RunListOptions{
		ListOptions: *ListOptions,
	}

	runlist, err := c.client.Runs.List(c.ctx, workspaceID, *rlo)
	if err != nil {
		return nil, err
	}

	for _, run := range runlist.Items {
		if checkRunCompleted(run) {
			continue
		}
		return &Run{
			ID:            run.ID,
			Status:        string(run.Status),
			Workspace:     WorkspaceName,
			CreatedAt:     run.CreatedAt,
			IsConfirmable: run.Actions.IsConfirmable,
		}, nil
	}

	return nil, nil
}

func (c *tfclient) RunApply(runID string) error {
	rao := &tfe.RunApplyOptions{
		Comment: tfe.String("Apply via tfcloud"),
	}

	if err := c.client.Runs.Apply(c.ctx, runID, *rao); err != nil {
		return err
	}

	return nil
}

func (c *tfclient) WorkspaceList(organization string) ([]*Workspace, error) {
	wlo := &tfe.WorkspaceListOptions{
		ListOptions: *ListOptions,
		Search:      nil,
	}

	wslist, err := c.client.Workspaces.List(c.ctx, organization, *wlo)
	if err != nil {
		return nil, err
	}

	result := make([]*Workspace, len(wslist.Items))
	for i, v := range wslist.Items {
		result[i] = &Workspace{
			ID:               v.ID,
			Name:             v.Name,
			TerraformVersion: v.TerraformVersion,
		}
	}

	return result, nil
}

func (c *tfclient) ModuleList() ([]*RegistryModule, error) {
	mlo := &RegistryModuleListOptions{
		Limit: 100,
	}

	modulelist, err := c.registryClient.RegistryModules.List(c.ctx, *mlo)
	if err != nil {
		return nil, err
	}

	result := make([]*RegistryModule, len(modulelist.Items))
	for i, v := range modulelist.Items {
		result[i] = &RegistryModule{
			ID:   v.ID,
			Name: v.Name,
			VersionStatuses: []tfe.RegistryModuleVersionStatuses{
				{
					Version: v.VersionStatuses[0].Version,
				},
			},
			Provider:     v.Provider,
			Organization: v.Organization.Name,
		}
	}

	return result, nil
}

func (c *tfclient) ModuleVersions(organization, name, provider string) (*RegistryModule, error) {
	module, err := c.client.RegistryModules.Read(c.ctx, organization, name, provider)
	if err != nil {
		return nil, err
	}

	return &RegistryModule{
		ID:              module.ID,
		Name:            module.Name,
		Provider:        module.Provider,
		VersionStatuses: module.VersionStatuses,
		Organization:    module.Organization.Name,
		VCSRepo:         module.VCSRepo.Identifier,
	}, nil
}

func checkRunCompleted(run *tfe.Run) bool {
	if run.Status == tfe.RunApplied ||
		run.Status == tfe.RunCanceled ||
		run.Status == tfe.RunErrored ||
		run.Status == tfe.RunDiscarded ||
		run.Status == tfe.RunPolicySoftFailed ||
		run.Status == tfe.RunPlannedAndFinished ||
		run.Status == "" {
		return true
	}
	return false
}
