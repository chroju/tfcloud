package tfrelease

import (
	"reflect"
	"testing"

	version "github.com/hashicorp/go-version"
)

var (
	releaseV1225, releaseV13, releaseV13beta1, releaseV13draft, releaseV14 *TfRelease
)

func init() {
	v1225, _ := version.NewVersion("v0.12.25")
	releaseV1225 = &TfRelease{
		Draft:      false,
		PreRelease: false,
		Version:    v1225,
		Tag:        "v0.12.25",
	}
	v13, _ := version.NewVersion("v0.13.0")
	releaseV13 = &TfRelease{
		Draft:      false,
		PreRelease: false,
		Version:    v13,
		Tag:        "v0.13.0",
	}
	v13beta1, _ := version.NewVersion("v0.13.0-beta1")
	releaseV13beta1 = &TfRelease{
		Draft:      false,
		PreRelease: true,
		Version:    v13beta1,
		Tag:        "v0.13.0-beta1",
	}
	v13draft, _ := version.NewVersion("v0.13.0-draft")
	releaseV13draft = &TfRelease{
		Draft:      true,
		PreRelease: false,
		Version:    v13draft,
		Tag:        "v0.13.0-draft",
	}
	v14, _ := version.NewVersion("v0.14.0")
	releaseV14 = &TfRelease{
		Draft:      false,
		PreRelease: false,
		Version:    v14,
		Tag:        "v0.14.0",
	}
}

func TestList(t *testing.T) {
	tests := []struct {
		name        string
		wantVersion *TfRelease
		wantErr     bool
	}{
		{
			name:        "normal release",
			wantVersion: releaseV13,
			wantErr:     false,
		},
		{
			name:        "pre release",
			wantVersion: releaseV13beta1,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := List()
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, v := range got {
				if equalTfRelease(v, tt.wantVersion) {
					return
				}
			}
			t.Errorf("List() = %v, want %v", got, tt.wantVersion)
		})
	}
}

func Test_latest(t *testing.T) {
	tests := []struct {
		name     string
		releases []*TfRelease
		want     *TfRelease
		wantErr  bool
	}{
		{
			name:     "normal",
			releases: []*TfRelease{releaseV14, releaseV13draft, releaseV13beta1, releaseV13},
			want:     releaseV14,
			wantErr:  false,
		},
		{
			name:     "later draft and pre-release exist",
			releases: []*TfRelease{releaseV13draft, releaseV13beta1, releaseV1225},
			want:     releaseV1225,
			wantErr:  false,
		},
		{
			name:     "same version with draft and pre-release exist",
			releases: []*TfRelease{releaseV13draft, releaseV13beta1, releaseV13},
			want:     releaseV13,
			wantErr:  false,
		},
		{
			name:     "no releases",
			releases: []*TfRelease{},
			want:     nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := latest(tt.releases)
			if (err != nil) != tt.wantErr {
				t.Errorf("latest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("latest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func equalTfRelease(a, b *TfRelease) bool {
	if a.Draft == b.Draft &&
		a.PreRelease == b.PreRelease &&
		a.Tag == b.Tag &&
		a.Version.Equal(b.Version) {
		return true
	}
	return false
}
