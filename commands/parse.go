package commands

import (
	"fmt"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclparse"
)

type terraformrc struct {
	Credentials []credential `hcl:"credentials,block"`
}

type credential struct {
	Name  string `hcl:"name,label"`
	Token string `hcl:"token"`
}

func ParseTerraformrc(path string) (string, error) {
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		return "", fmt.Errorf("Parse %s failed", path)
	}

	var tfrc terraformrc
	diags = gohcl.DecodeBody(f.Body, nil, &tfrc)
	if diags.HasErrors() {
		return "", fmt.Errorf("Decode %s failed", path)
	}

	return tfrc.Credentials[0].Token, nil
}
