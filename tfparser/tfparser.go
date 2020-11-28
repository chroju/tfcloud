package tfparser

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	version "github.com/hashicorp/go-version"
	hcl "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type terraformrc struct {
	Credentials []*Credential `hcl:"credentials,block"`
}

// Credential represents a Terraform Cloud credential.
type Credential struct {
	Name  string `hcl:"name,label"`
	Token string `hcl:"token"`
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
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		return nil, fmt.Errorf("%s", diags.Error())
	}

	var tfrc terraformrc
	diags = gohcl.DecodeBody(f.Body, nil, &tfrc)
	if diags.HasErrors() {
		return nil, fmt.Errorf("%s", diags.Error())
	}

	return tfrc.Credentials[0], nil
}

// ParseRemoteBackend parses remote backend config in the specified directory and returns values.
func ParseRemoteBackend(root string) (*RemoteBackend, error) {
	var config *RemoteBackend
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() || !strings.HasSuffix(info.Name(), ".tf") {
				return nil
			}

			src, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			file, diags := hclwrite.ParseConfig(src, path, hcl.InitialPos)
			if diags.HasErrors() {
				return diags
			}

			for _, block := range file.Body().Blocks() {
				if block.Type() == "terraform" {
					for _, subBlock := range block.Body().Blocks() {
						if subBlock.Type() == "backend" && subBlock.Labels()[0] == "remote" {
							subBlockBody := subBlock.Body()
							ver, err := version.NewConstraint(parseAttribute(block.Body().GetAttribute("required_version")))
							if err != nil {
								return err
							}
							config = &RemoteBackend{
								Organization:    parseAttribute(subBlockBody.GetAttribute("organization")),
								Hostname:        parseAttribute(subBlockBody.GetAttribute("hostname")),
								WorkspaceName:   parseAttribute(subBlockBody.Blocks()[0].Body().GetAttribute("name")),
								WorkspacePrefix: parseAttribute(subBlockBody.Blocks()[0].Body().GetAttribute("prefix")),
								RequiredVersion: ver,
							}
							return nil
						}
					}
				}
			}
			return nil
		})

	if err != nil {
		return nil, err
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
