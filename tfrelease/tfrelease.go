package tfrelease

import (
	"context"
	"fmt"

	"github.com/google/go-github/v43/github"
	version "github.com/hashicorp/go-version"
)

const (
	repoOwner = "hashicorp"
	repoName  = "terraform"
)

// TfRelease represents a Terraform release.
type TfRelease struct {
	Draft      bool   `json:"draft"`
	PreRelease bool   `json:"prerelease"`
	Tag        string `json:"tag_name"`
	Version    *version.Version
}

// List returns all Terraform releases.
func List() ([]*TfRelease, error) {
	var tfReleases []*TfRelease
	client := github.NewClient(nil)

	ctx := context.Background()
	lo := &github.ListOptions{
		PerPage: 100,
	}
	for {
		releases, resp, err := client.Repositories.ListReleases(ctx, repoOwner, repoName, lo)
		if err != nil {
			return nil, err
		}

		for _, v := range releases {
			sv, err := version.NewVersion(*v.TagName)
			if err != nil {
				return nil, err
			}
			tfReleases = append(tfReleases, &TfRelease{
				Draft:      *v.Draft,
				PreRelease: *v.Prerelease,
				Tag:        *v.TagName,
				Version:    sv,
			})
		}
		if resp.NextPage == 0 {
			break
		}
		lo.Page = resp.NextPage
	}

	return tfReleases, nil
}

// Latest returns the latest terraform release (not draft or pre release version).
func Latest() (*TfRelease, error) {
	releases, err := List()
	if err != nil {
		return nil, err
	}
	return latest(releases)
}

func latest(releases []*TfRelease) (*TfRelease, error) {
	for _, v := range releases {
		if v.Draft || v.PreRelease {
			continue
		}
		return v, nil
	}

	return nil, fmt.Errorf("Something is wrong to get latest terraform version")
}
