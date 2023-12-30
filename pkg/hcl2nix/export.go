package hcl2nix

import (
	"fmt"
	"strings"
)

var (
	// we'll support nydus, estargz, etc in future
	artifactTypes = []string{"oci"}
)

// ExportConfig to export Nix package outputs to an artifact
type ExportConfig struct {
	Environment  string   `hcl:"environment,label"`
	ArtifactType string   `hcl:"artifactType"`
	Name         string   `hcl:"name"`
	Platform     []string `hcl:"platform"`
	Publish      *bool    `hcl:"publish"`
	// Credentials or Credential location?
}

// Validate validates ExportConfig
func (c *ExportConfig) Validate() *string {
	// todo: maybe we should return hcl.Diagnostic to be consistent
	if !validateArtifactType(c.ArtifactType) {
		return pointerTo(fmt.Sprintf("Invalid artifactType. Valid values are : %s", strings.Join(artifactTypes, ", ")))
	}

	return nil
}

func validateArtifactType(artifactType string) bool {
	for _, at := range artifactTypes {
		if at == artifactType {
			return true
		}
	}
	return false
}

func pointerTo[T any](value T) *T {
	return &value
}
