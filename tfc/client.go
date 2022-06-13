package tfc

import (
	"fmt"

	"net/url"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// Client is the Terraform Enterprise and Terraform Cloud API client.
type Client struct {
	*tfe.Client
}

// NewClient creates a new Terraform Enterprise API client.
func NewClient(cfg *tfe.Config) (*Client, error) {
	config := tfe.DefaultConfig()

	// Layer in the provided config for any non-blank values.
	if cfg != nil {
		if cfg.Address != "" {
			config.Address = cfg.Address
		}
		if cfg.BasePath != "" {
			config.BasePath = cfg.BasePath
		}
		if cfg.Token != "" {
			config.Token = cfg.Token
		}
		for k, v := range cfg.Headers {
			config.Headers[k] = v
		}
		if cfg.HTTPClient != nil {
			config.HTTPClient = cfg.HTTPClient
		}
		if cfg.RetryLogHook != nil {
			config.RetryLogHook = cfg.RetryLogHook
		}
	}

	// Parse the address to make sure its a valid URL.
	baseURL, err := url.Parse(config.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %v", err)
	}

	baseURL.Path = config.BasePath
	if !strings.HasSuffix(baseURL.Path, "/") {
		baseURL.Path += "/"
	}

	tfeClient, err := tfe.NewClient(config)
	if err != nil {
		return nil, err
	}

	// Create the client.
	client := &Client{
		Client: tfeClient,
	}

	return client, nil
}
