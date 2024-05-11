package hcl2nix

import (
	"strings"
)

// OCIArtifact to export Nix package outputs to an artifact
type OCIArtifact struct {
	Environment string `hcl:"environment,label"`
	// Name of the image . Ex: ttl.sh/myproject/app:1h
	Name string `hcl:"name"`
	// Cmd defines the default arguments to the entrypoint of the container.
	Cmd []string `hcl:"cmd,optional"`
	// Entrypoint defines a list of arguments to use as the command to execute when the container starts.
	Entrypoint []string `hcl:"entrypoint,optional"`
	// Env is a list of environment variables to be used in a container.
	EnvVars []string `hcl:"envVars,optional"`
	// ExposedPorts a set of ports to expose from a container running this image. Ex: ["80/tcp", "443/tcp"]
	ExposedPorts []string `hcl:"exposedPorts,optional"`
	// Names of configs to import
	ImportConfigs []string `hcl:"importConfigs,optional"`
	// DevDeps defines if development dependencies should be present in the image. By default, it is false.
	DevDeps bool `hcl:"devDeps,optional"`

	Platform string `hcl:"platform"`
}

// Validate validates ExportConfig
func (c *OCIArtifact) Validate() *string {
	// todo: maybe we should return hcl.Diagnostic to be consistent
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
