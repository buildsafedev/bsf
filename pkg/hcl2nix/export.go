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
	Cmd          []string `hcl:"cmd,optional"`
	Entrypoint   []string `hcl:"entrypoint,optional"`
	Publish      *bool    `hcl:"publish"`
	Platform     string   `hcl:"platform"`
	EnvVars      []string `hcl:"envVars,optional"`
	// Credentials or Credential location?
	// todo: we need a field to specify if they want specific directories from current sandoxed directory to be copied over to runtime artifact
}

// Validate validates ExportConfig
func (c *ExportConfig) Validate() *string {
	// todo: maybe we should return hcl.Diagnostic to be consistent
	if !validateArtifactType(c.ArtifactType) {
		return pointerTo(fmt.Sprintf("Invalid artifactType. Valid values are : %s", strings.Join(artifactTypes, ", ")))
	}

	if !validatePlatform(c.Platform) {
		return pointerTo("Invalid platform. Platform cannot contain spaces, commas, or semicolons. Note: multi-platform support will be added in future")
	}

	if len(c.EnvVars) != 0 {
		if !validateEnvVars((c.EnvVars)) {
			return pointerTo("Invalid environment variables, please use 'key=value' format")
		}
	}

	return nil
}

func validatePlatform(platform string) bool {
	if strings.Contains(platform, ",") || strings.Contains(platform, " ") || strings.Contains(platform, ";") {
		return false
	}
	return true
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

func validateEnvVars(envVars []string) bool {
	for _, kv := range envVars {
		keyValuePair := strings.SplitN(kv, "=", 2)
		if len(keyValuePair) != 2 {
			return false
		}
	}
	return true
}
