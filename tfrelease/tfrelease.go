package tfrelease

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	version "github.com/hashicorp/go-version"
)

const (
	tfReleaseURL = "https://api.github.com/repos/hashicorp/terraform/releases"
)

// TfRelease represents Terraform release
type TfRelease struct {
	Draft      bool   `json:"draft"`
	PreRelease bool   `json:"prerelease"`
	Tag        string `json:"tag_name"`
	Version    *version.Version
}

// List returns Terraform releases
func List() ([]*TfRelease, error) {
	resp, err := http.Get(tfReleaseURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tfReleases []*TfRelease
	if err = json.Unmarshal(body, &tfReleases); err != nil {
		return nil, err
	}
	for _, v := range tfReleases {
		sv, err := version.NewVersion(v.Tag)
		if err != nil {
			return nil, err
		}
		v.Version = sv
	}
	return tfReleases, nil
}

func Latest() (*TfRelease, error) {
	releases, err := List()
	if err != nil {
		return nil, err
	}

	for _, v := range releases {
		if v.Draft || v.PreRelease {
			continue
		}
		return v, nil
	}

	return nil, fmt.Errorf("Something is wrong to get latest terraform version")
}
