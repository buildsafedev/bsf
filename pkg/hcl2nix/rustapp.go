package hcl2nix

// RustApp defines the parameters for a rust application.
type RustApp struct {
	// ProjectName: name of the project.
	ProjectName string `hcl:"projectName"`
	// Version of Rust
	RustVersion string `hcl:"rustVersion"`
}
