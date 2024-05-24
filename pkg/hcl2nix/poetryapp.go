package hcl2nix

import "github.com/BurntSushi/toml"

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

// Poetry defines structs to match the structure of the TOML file
type Poetry struct {
	Name         string            `toml:"name"`
	Version      string            `toml:"version"`
	Description  string            `toml:"description"`
	Authors      []string          `toml:"authors"`
	License      string            `toml:"license"`
	Readme       string            `toml:"readme"`
	Dependencies map[string]string `toml:"dependencies"`
	Scripts      map[string]string `toml:"scripts"`
}

// PoetryBuildSystem defines structs to match the structure of the TOML file
type PoetryBuildSystem struct {
	Requires     []string `toml:"requires"`
	BuildBackend string   `toml:"build-backend"`
}

// PoetryTool defines structs to match the structure of the TOML file
type PoetryTool struct {
	Poetry Poetry `toml:"poetry"`
}

// PoetryConfig defines structs to match the structure of the TOML file
type PoetryConfig struct {
	Tool        PoetryTool        `toml:"tool"`
	BuildSystem PoetryBuildSystem `toml:"build-system"`
}

func parsePyProject(path string) (*PoetryConfig, error) {
	var config PoetryConfig

	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
