package hcl2nix

// PoetryApp defines the parameters for a poetry application.
type PoetryApp struct {
	// ProjectDir: path to the root of the project.
	ProjectDir string `hcl:"projectDir"`
	// Src: project source (defaults to projectDir).
	Src string `hcl:"src,optional"`
	// Pyproject: path to pyproject.toml (defaults to projectDir + "/pyproject.toml").
	Pyproject string `hcl:"pyproject,optional"`
	// Poetrylock: poetry.lock file path (defaults to projectDir + "/poetry.lock").
	Poetrylock string `hcl:"poetrylock,optional"`
	// PreferWheels: Use wheels rather than sdist as much as possible (defaults to false).
	PreferWheels bool `hcl:"preferWheels,optional"`
	// CheckGroups: Which Poetry 1.2.0+ dependency groups to run unit tests (defaults to [ "dev" ]).
	CheckGroups []string `hcl:"checkGroups,optional"`

	// TODO: allow customsing python version
	// PythonVersion     string `hcl:"pythonVersion,optional"`

}
