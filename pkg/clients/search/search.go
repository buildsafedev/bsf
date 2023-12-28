package search

import (
	"bytes"
	"context"
	"encoding/json"
	"sort"
	"time"
)

// Package is a struct to define a Package from the Search API
type Package struct {
	Name     string    `json:"name"`
	Revision string    `json:"revision"`
	Version  string    `json:"version"`
	DateTime time.Time `json:"datetime"`
}

// LatestRevision is a struct to define a LatestRevision from the Search API
type LatestRevision struct {
	Revision string    `json:"revision"`
	DateTime time.Time `json:"datetime"`
}

// ListAllPackageNames returns a list of all package names in the Search API
func (c *Client) ListAllPackageNames(ctx context.Context) ([]string, error) {
	url := "/packages"
	resp, err := c.SendGetRequest(ctx, url)
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
func (c *Client) ListPackageVersions(ctx context.Context, packageName string) ([]Package, error) {
	url := "/packages/" + packageName

	resp, err := c.SendGetRequest(ctx, url)
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
func (c *Client) GetPackageVersion(ctx context.Context, packageName, version string) (*Package, error) {
	url := "/packages/" + packageName + "/" + version

	resp, err := c.SendGetRequest(ctx, url)
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
func (c *Client) LatestRevision(ctx context.Context) *LatestRevision {
	url := "/latest-revision"
	resp, err := c.SendGetRequest(ctx, url)
	if err != nil {
		return nil
	}
	latest := LatestRevision{}
	if err := json.NewDecoder(bytes.NewReader(resp)).Decode(&latest); err != nil {
		return nil
	}
	return &latest
}

// SortPackagesWithTimestamp sorts packages with timestamp with latest being the first element
func SortPackagesWithTimestamp(packageVersions []Package) []Package {
	sort.Slice(packageVersions, func(i, j int) bool {
		return packageVersions[i].DateTime.After(packageVersions[j].DateTime)
	})
	return packageVersions
}
