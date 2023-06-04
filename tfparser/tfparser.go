package tfparser

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	version "github.com/hashicorp/go-version"
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/hcl/v2/json"
)

type terraformrc struct {
	Credentials []*Credential `hcl:"credentials,block"`
}

// Credential represents a Terraform Cloud credential.
type Credential struct {
	Hostname string `hcl:"name,label"`
	Token    string `hcl:"token"`
}

// RemoteBackend represents a Terraform remote backend config.
type RemoteBackend struct {
	Token           string
	Hostname        string
	Organization    string
	WorkspaceName   string
	WorkspacePrefix string
	RequiredVersion version.Constraints
}

// ParseTerraformrc parses terraformrc file and returns a Terraform Cloud credential.
func ParseTerraformrc(path string) (*Credential, error) {
	tfrcJSONPath := os.Getenv("HOME") + "/.terraform.d/credentials.tfrc.json"
	var b *hcl.File
	var diags hcl.Diagnostics

	if _, err := os.Stat(tfrcJSONPath); err == nil {
		b, diags = json.ParseFile(tfrcJSONPath)
		if diags.HasErrors() {
			return nil, fmt.Errorf("%s", diags.Error())
		}
	} else if _, err = os.Stat(path); err == nil {
		parser := hclparse.NewParser()
		b, diags = parser.ParseHCLFile(path)
		if diags.HasErrors() {
			return nil, fmt.Errorf("%s", diags.Error())
		}
	} else {
		return nil, fmt.Errorf("terraform credential file not found")
	}

	return parseTerraformrcConfig(b)
}

func parseTerraformrcConfig(configFile *hcl.File) (*Credential, error) {
	var tfrc terraformrc
	diags := gohcl.DecodeBody(configFile.Body, nil, &tfrc)
	if diags.HasErrors() {
		return nil, diags
	}

	return tfrc.Credentials[0], nil
}

// ParseRemoteBackend parses remote backend config in the specified directory and returns values.
func ParseRemoteBackend(root string) (*RemoteBackend, error) {
	var config *RemoteBackend
	paths, err := filepath.Glob(fmt.Sprintf("%s/*.tf", root))
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		src, err := ioutil.ReadFile(path)
		if err != nil {
			continue
		}

		file, diags := hclwrite.ParseConfig(src, path, hcl.InitialPos)
		if diags.HasErrors() {
			continue
		}

		for _, block := range file.Body().Blocks() {
			if block.Type() == "terraform" {
				for _, subBlock := range block.Body().Blocks() {
					if subBlock.Type() == "backend" && subBlock.Labels()[0] == "remote" {
						subBlockBody := subBlock.Body()

						requiredVersion, err := parseRequiredVersion(block.Body().GetAttribute("required_version"))
						if err != nil {
							continue
						}

						cfg := &RemoteBackend{
							Organization:    parseAttribute(subBlockBody.GetAttribute("organization")),
							Hostname:        parseAttribute(subBlockBody.GetAttribute("hostname")),
							WorkspaceName:   parseAttribute(subBlockBody.Blocks()[0].Body().GetAttribute("name")),
							WorkspacePrefix: parseAttribute(subBlockBody.Blocks()[0].Body().GetAttribute("prefix")),
							RequiredVersion: requiredVersion,
						}

						if cfg != nil && config != nil {
							return nil, fmt.Errorf("Remote backend config is duplicated")
						}

						config = cfg
					}
				}
			}
		}
	}

	if config == nil {
		return nil, fmt.Errorf("Remote backend config is not found")
	}

	return config, nil
}

func parseAttribute(a *hclwrite.Attribute) string {
	if a == nil {
		return ""
	}
	return string(a.Expr().BuildTokens(nil)[1].Bytes)
}

func parseRequiredVersion(a *hclwrite.Attribute) (version.Constraints, error) {
	if a == nil {
		return nil, nil
	}
	return version.NewConstraint(parseAttribute(a))
}
