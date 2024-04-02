package hcl2nix

// PoetryApp defines the parameters for a poetry application.
type RustApp struct {
	// ProjectName: name of the project.
	ProjectName string `hcl:"projectName"`
	// Version of Rust
	RustVersion string `hcl:"rustVersion"`
}
