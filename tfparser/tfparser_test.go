package tfparser

import (
	"reflect"
	"testing"

	version "github.com/hashicorp/go-version"
)

func TestParseTerraformrc(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    *Credential
		wantErr bool
	}{
		{
			name: "nomarl test",
			path: "./test_config/normal.tf",
			want: &Credential{
				Hostname: "app.terraform.io",
				Token:    "EmN5pXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
			},
			wantErr: false,
		},
		{
			name:    "abnomarl test",
			path:    "./test_config/abnormal.tf",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "not exist",
			path:    "./test_config/not_exist.tf",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTerraformrc(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTerraformrc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTerraformrc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseRemoteBackend(t *testing.T) {
	testRequiredVersion, _ := version.NewConstraint(">= 0.13.0, < 0.13.2")
	tests := []struct {
		name    string
		root    string
		want    *RemoteBackend
		wantErr bool
	}{
		{
			name: "normal test",
			root: "./test_config/normal_backend",
			want: &RemoteBackend{
				Hostname:        "app.terraform.io",
				Organization:    "test-org",
				WorkspaceName:   "test-workspace",
				WorkspacePrefix: "",
				RequiredVersion: testRequiredVersion,
			},
			wantErr: false,
		},
		{
			name:    "abnormal test",
			root:    "./test_config/abnormal_backend",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "not exist",
			root:    "./test_config/not_exist_path",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRemoteBackend(tt.root)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRemoteBackend() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !equalRemoteBackend(got, tt.want) {
				t.Errorf("ParseRemoteBackend() = %v, want %v", got, tt.want)
			}
		})
	}
}

func equalRemoteBackend(a, b *RemoteBackend) bool {
	if a == nil || b == nil {
		return true
	}
	if a.Hostname == b.Hostname &&
		a.Organization == b.Organization &&
		a.WorkspaceName == b.WorkspaceName &&
		a.WorkspacePrefix == b.WorkspacePrefix &&
		a.RequiredVersion.String() == b.RequiredVersion.String() {
		return true
	}
	return false
}
