package template

// OCIArtifact holds parameters for OCI artifacts
type OCIArtifact struct {
	Environment  string
	Name         string
	Cmd          []string
	Entrypoint   []string
	Platform     string
	EnvVars      []string
	ConfigFiles  []ConfigFiles
	ExposedPorts []string
}
