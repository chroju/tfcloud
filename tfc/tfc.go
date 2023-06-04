package tfc

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chroju/tfcloud/tfparser"
	tfe "github.com/hashicorp/go-tfe"
)

var defaultListOptions = &tfe.ListOptions{
	PageNumber: 0,
	PageSize:   100,
}

// Run represents a Terraform workspaces run.
type Run struct {
	ID            *string   `json:"id"`
	Organization  *string   `json:"organization"`
	Workspace     *string   `json:"workspace"`
	Status        *string   `json:"status"`
	IsConfirmable *bool     `json:"is_confirmable"`
	CreatedAt     time.Time `json:"created_at"`
}

// Workspace represents a Terraform Cloud workspace.
type Workspace struct {
	ID               *string   `json:"id"`
	Name             *string   `json:"name"`
	TerraformVersion *string   `json:"terraform_version"`
	ExecutionMode    *string   `json:"execution_mode"`
	AutoApply        *bool     `json:"auto_apply"`
	CurrentRun       *tfe.Run  `json:"current_run,omitempty"`
	VCSRepoName      *string   `json:"vcs_repo"`
	WorkingDirectory *string   `json:"working_directory"`
	ResourceCount    *int      `json:"resource_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// RegistryModule represents a Terraform Cloud registry module.
type RegistryModule struct {
	ID              *string                             `json:"id"`
	Name            *string                             `json:"name"`
	Provider        *string                             `json:"provider"`
	VersionStatuses []tfe.RegistryModuleVersionStatuses `json:"version_statuses"`
	Organization    *string                             `json:"organization"`
	Source          *string                             `json:"source"`
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
	ModuleList(organization string) ([]*RegistryModule, error)
	// ModuleGet returns the specified terraform registry module.
	ModuleGet(organization, name, provider string) (*RegistryModule, error)
}

type tfclient struct {
	address string
	client  *Client
	ctx     context.Context
}

// NewTfCloud creates a new TfCloud interface
func NewTfCloud(address, token string) (TfCloud, error) {
	config, err := NewCredentials("", address, token)
	if err != nil {
		return nil, err
	}

	client, err := NewClient(config)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	return &tfclient{
		address: config.Address,
		client:  client,
		ctx:     ctx,
	}, nil
}

func NewCredentials(filepath, address, token string) (*tfe.Config, error) {
	envs := os.Environ()
	for _, env := range envs {
		if strings.HasPrefix(env, "TF_TOKEN_") {
			kv := strings.Split(env, "=")
			splitted := strings.Split(kv[0], "_")
			return &tfe.Config{
				Address: fmt.Sprintf("https://%s", strings.Join(splitted[2:], ".")),
				Token:   kv[1],
			}, nil
		}
	}
	terraformrcPath := os.Getenv("TF_CLI_CONFIG_FILE")
	if filepath != "" {
		terraformrcPath = filepath
	}
	if terraformrcPath == "" {
		terraformrcPath = os.Getenv("HOME") + "/.terraformrc"
	}

	credential, err := tfparser.ParseTerraformrc(terraformrcPath)
	if err != nil {
		return nil, err
	}
	if address != "" {
		credential.Hostname = address
	}
	if !strings.HasPrefix(credential.Hostname, "https://") {
		credential.Hostname = fmt.Sprintf("https://%s", credential.Hostname)
	}
	if token != "" {
		credential.Token = token
	}

	return &tfe.Config{
		Address: credential.Hostname,
		Token:   credential.Token,
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
				run, err := c.RunGet(*ws.Name, ws.CurrentRun.ID)
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
		ID:            &run.ID,
		Status:        tfe.String(string(run.Status)),
		Workspace:     &workspaceName,
		CreatedAt:     run.CreatedAt,
		IsConfirmable: &run.Actions.IsConfirmable,
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
		Search:      "",
	}

	var workspaces []*tfe.Workspace
	for {
		wslist, err := c.client.Workspaces.List(c.ctx, organization, wlo)
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
			ID:               &v.ID,
			Name:             &v.Name,
			TerraformVersion: &v.TerraformVersion,
			ExecutionMode:    &v.ExecutionMode,
			AutoApply:        &v.AutoApply,
			CurrentRun:       v.CurrentRun,
			VCSRepoName:      &vcsRepoName,
			WorkingDirectory: &v.WorkingDirectory,
			ResourceCount:    &v.ResourceCount,
			CreatedAt:        v.CreatedAt,
			UpdatedAt:        v.UpdatedAt,
		}
	}

	return result, nil
}

func (c *tfclient) WorkspaceGet(organization, workspace string) (*Workspace, error) {
	ws, err := c.client.Workspaces.Read(c.ctx, organization, workspace)
	if err != nil {
		return nil, err
	}
	vcsRepoName := ""
	if ws.VCSRepo != nil {
		vcsRepoName = ws.VCSRepo.Identifier
	}

	return &Workspace{
		ID:               &ws.ID,
		Name:             &ws.Name,
		TerraformVersion: &ws.TerraformVersion,
		ExecutionMode:    &ws.ExecutionMode,
		AutoApply:        &ws.AutoApply,
		CurrentRun:       ws.CurrentRun,
		VCSRepoName:      &vcsRepoName,
		WorkingDirectory: &ws.WorkingDirectory,
		ResourceCount:    &ws.ResourceCount,
		CreatedAt:        ws.CreatedAt,
		UpdatedAt:        ws.UpdatedAt,
	}, nil
}

func (c *tfclient) WorkspaceUpdateVersion(organization, workspace, version string) error {
	wuo := tfe.WorkspaceUpdateOptions{
		TerraformVersion: tfe.String(version),
	}
	_, err := c.client.Workspaces.Update(c.ctx, organization, workspace, wuo)
	return err
}

func (c *tfclient) ModuleList(organization string) ([]*RegistryModule, error) {
	mlo := &tfe.RegistryModuleListOptions{
		ListOptions: *defaultListOptions,
	}

	modulelist, err := c.client.RegistryModules.List(c.ctx, organization, mlo)
	if err != nil {
		return nil, err
	}

	result := make([]*RegistryModule, len(modulelist.Items))
	for i, v := range modulelist.Items {
		result[i] = &RegistryModule{
			ID:   &v.ID,
			Name: &v.Name,
			VersionStatuses: []tfe.RegistryModuleVersionStatuses{
				{
					Version: v.VersionStatuses[0].Version,
				},
			},
			Provider:     &v.Provider,
			Organization: &v.Organization.Name,
			Source:       &v.VCSRepo.Identifier,
		}
	}

	return result, nil
}

func (c *tfclient) ModuleGet(organization, name, provider string) (*RegistryModule, error) {
	moduleID := tfe.RegistryModuleID{
		Organization: organization,
		Name:         name,
		Provider:     provider,
	}
	module, err := c.client.RegistryModules.Read(c.ctx, moduleID)
	if err != nil {
		return nil, err
	}

	return &RegistryModule{
		ID:              &module.ID,
		Name:            &module.Name,
		Provider:        &module.Provider,
		VersionStatuses: module.VersionStatuses,
		Organization:    &module.Organization.Name,
		Source:          &module.VCSRepo.Identifier,
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
