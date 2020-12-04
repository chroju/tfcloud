package tfc

import (
	"context"
	"time"

	tfe "github.com/hashicorp/go-tfe"
)

var defaultListOptions = &tfe.ListOptions{
	PageNumber: 0,
	PageSize:   100,
}

// Run represents a Terraform workspaces run.
type Run struct {
	ID            string
	Organization  string
	Workspace     string
	Status        string
	IsConfirmable bool
	CreatedAt     time.Time
}

// Workspace represents a Terraform Cloud workspace.
type Workspace struct {
	ID               string
	Name             string
	TerraformVersion string
	CurrentRun       *tfe.Run
	VCSRepoName      string
}

// RegistryModule represents a Terraform Cloud registry module.
type RegistryModule struct {
	ID              string
	Name            string
	Provider        string
	VersionStatuses []tfe.RegistryModuleVersionStatuses
	Organization    string
	VCSRepo         string
}

// TfCloud represents Terraform Cloud API client.
type TfCloud interface {
	// Address returns a Terraform Cloud / Enterprise API endpoint addres.
	Address() string
	// RunList returns all the terraform workspace current runs.
	RunList(organization string) ([]*Run, error)
	// RunGet returns the specified terraform workspace run.
	RunGet(workspaceName, runID string) (*Run, error)
	// RunApply applys the specified terraform workspace run.
	RunApply(RunID string) error
	// WorkspaceList returns all the terraform workspaces in an organization.
	WorkspaceList(organization string) ([]*Workspace, error)
	// WorkspaceGet returns the specified terraform workspace.
	WorkspaceGet(organization, workspace string) (*Workspace, error)
	// WorkspaceUpdateVersion updates the terraform version config in the specified workspace.
	WorkspaceUpdateVersion(organization, workspace, version string) error
	// ModuleList returns all the terraform registry modules.
	ModuleList() ([]*RegistryModule, error)
	// ModuleGet returns the specified terraform registry module.
	ModuleGet(organization, name, provider string) (*RegistryModule, error)
}

type tfclient struct {
	address        string
	client         *tfe.Client
	registryClient *RegistryClient
	ctx            context.Context
}

// NewTfCloud creates a new TfCloud interface
func NewTfCloud(address, token string) (TfCloud, error) {
	config := &tfe.Config{
		Address: address,
		Token:   token,
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
		address,
		client,
		registryClient,
		ctx,
	}, nil
}

func (c *tfclient) Address() string {
	return c.address
}

func (c *tfclient) RunList(organization string) ([]*Run, error) {
	type result struct {
		Error    error
		Response *Run
	}

	workspaces, err := c.WorkspaceList(organization)
	if err != nil {
		return nil, err
	}

	resultChan := make(chan result)
	for _, ws := range workspaces {
		go func(ws *Workspace) {
			if ws.CurrentRun == nil {
				resultChan <- result{Error: nil, Response: nil}
			} else {
				run, err := c.RunGet(ws.Name, ws.CurrentRun.ID)
				resultChan <- result{Error: err, Response: run}
			}
		}(ws)
	}

	var rtn []*Run
	for range workspaces {
		run := <-resultChan
		if run.Error != nil {
			return nil, run.Error
		}
		if run.Response != nil {
			rtn = append(rtn, run.Response)
		}
	}

	return rtn, nil
}

func (c *tfclient) RunGet(workspaceName, runID string) (*Run, error) {
	run, err := c.client.Runs.Read(c.ctx, runID)
	if err != nil {
		return nil, err
	}

	if checkRunCompleted(run) {
		return nil, nil
	}

	return &Run{
		ID:            run.ID,
		Status:        string(run.Status),
		Workspace:     workspaceName,
		CreatedAt:     run.CreatedAt,
		IsConfirmable: run.Actions.IsConfirmable,
	}, nil
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
		ListOptions: *defaultListOptions,
		Search:      nil,
	}

	var workspaces []*tfe.Workspace
	for {
		wslist, err := c.client.Workspaces.List(c.ctx, organization, *wlo)
		if err != nil {
			return nil, err
		}
		workspaces = append(workspaces, wslist.Items...)
		if wslist.CurrentPage == wslist.TotalPages {
			break
		}
		wlo.PageNumber = wslist.NextPage
	}

	result := make([]*Workspace, len(workspaces))
	for i, v := range workspaces {
		vcsRepoName := ""
		if v.VCSRepo != nil {
			vcsRepoName = v.VCSRepo.Identifier
		}
		result[i] = &Workspace{
			ID:               v.ID,
			Name:             v.Name,
			TerraformVersion: v.TerraformVersion,
			CurrentRun:       v.CurrentRun,
			VCSRepoName:      vcsRepoName,
		}
	}

	return result, nil
}

func (c *tfclient) WorkspaceGet(organization, workspace string) (*Workspace, error) {
	ws, err := c.client.Workspaces.Read(c.ctx, organization, workspace)
	if err != nil {
		return nil, err
	}

	return &Workspace{
		ID:               ws.ID,
		Name:             ws.Name,
		TerraformVersion: ws.TerraformVersion,
		CurrentRun:       ws.CurrentRun,
	}, nil
}

func (c *tfclient) WorkspaceUpdateVersion(organization, workspace, version string) error {
	wuo := tfe.WorkspaceUpdateOptions{
		TerraformVersion: tfe.String(version),
	}
	_, err := c.client.Workspaces.Update(c.ctx, organization, workspace, wuo)
	return err
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

func (c *tfclient) ModuleGet(organization, name, provider string) (*RegistryModule, error) {
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
