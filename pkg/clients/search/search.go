package search

import (
	"bytes"
	"encoding/json"
	"time"
)

// Package is a struct to define a Package from the Search API
type Package struct {
	Name     string
	Revision string
	Version  string
	DateTime time.Time
}

// LatestRevision is a struct to define a LatestRevision from the Search API
type LatestRevision struct {
	Revision string
	DateTime time.Time
}

// ListAllPackageNames returns a list of all package names in the Search API
func (c *Client) ListAllPackageNames() ([]string, error) {
	url := "/packages"
	resp, err := c.SendGetRequest(url)
	if err != nil {
		return nil, err
	}

	packages := make([]string, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&packages); err != nil {
		return nil, err
	}

	return packages, nil
}

// ListPackageVersions returns a list of all packages in the Search API
func (c *Client) ListPackageVersions(packageName string) ([]Package, error) {
	url := "/packages/" + packageName

	resp, err := c.SendGetRequest(url)
	if err != nil {
		return nil, err
	}

	packages := make([]Package, 0)
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&packages); err != nil {
		return nil, err
	}

	return packages, nil
}

// GetPackageVersion returns a specific package version from the Search API
func (c *Client) GetPackageVersion(packageName, version string) (*Package, error) {
	url := "/packages/" + packageName + "/" + version

	resp, err := c.SendGetRequest(url)
	if err != nil {
		return nil, err
	}

	pkg := Package{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&pkg); err != nil {
		return nil, err
	}

	return &pkg, nil
}

// LatestRevision returns the latest revision of the Search API that indicates when the index was last updated
func (c *Client) LatestRevision() *LatestRevision {
	url := "/latest-revision"
	resp, err := c.SendGetRequest(url)
	if err != nil {
		return nil
	}
	latest := LatestRevision{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&latest); err != nil {
		return nil
	}
	return &latest
}
