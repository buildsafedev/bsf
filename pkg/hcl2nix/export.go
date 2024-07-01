package hcl2nix

import (
	"strconv"
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
	// Base defines if the OCI artifact that needs to be built is for a Base image. By default, it is false.
	Base bool `hcl:"base,optional"`
	// BaseDeps defines the type of dependencies which needs to included in the base image. Can be either "runtime" or "dev".
	BaseDeps string `hcl:"baseDeps,optional"`
}

// Validate validates ExportConfig
func (c *OCIArtifact) Validate(conf *Config) *string {
	if len(c.EnvVars) != 0 {
		if !validateEnvVars((c.EnvVars)) {
			return pointerTo("Invalid environment variables, please use 'key=value' format")
		}
	}

	if len(c.ImportConfigs) != 0 {
		if !validateImportConfigs(c.ImportConfigs, conf) {
			return pointerTo("Invalid import configs, please specify a valid config name")
		}
	}

	if len(c.ExposedPorts) != 0 {
		if !validateExposedPorts(c.ExposedPorts) {
			return pointerTo("Invalid exposed ports, please specify a valid port/protocol. Ex: 80/tcp ")
		}
	}

	return nil
}

func pointerTo[T any](value T) *T {
	return &value
}

func validateExposedPorts(ports []string) bool {
	for _, port := range ports {
		pp := strings.Split(port, "/")
		if len(pp) != 2 {
			return false
		}

		if pp[1] != "tcp" && pp[1] != "udp" && pp[1] != "icmp" {
			return false
		}

		if pp[0] == "" {
			return false
		}

		pn, err := strconv.Atoi(pp[0])
		if err != nil {
			return false
		}

		if pn < 0 || pn > 65535 {
			return false
		}

	}
	return true
}

func validateImportConfigs(configs []string, conf *Config) bool {
	validConfigs := make(map[string]bool)
	for _, configName := range conf.ConfigFiles {
		validConfigs[configName.Name] = true
	}

	for _, config := range configs {
		if _, ok := validConfigs[config]; !ok {
			return false
		}

	}
	return true
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
