package hcl2nix

// RustApp defines the parameters for a rust application.
type RustApp struct {
	// CrateName: name of the project.
	CrateName string `hcl:"projectName"`
	// Release: To enable or disable the release mode, defaults to "true".
	Release bool `hcl:"release"`
	// RustVersion: Version of Rust
	RustVersion string `hcl:"rustVersion,optional"`
	// RustToolChain: Used to override the toolchain
	RustToolChain string `hcl:"rustToolchain,optional"`
	// RustChannel: To support legacy use, this can be a version when supplied alone. Defaults to "stable".
	RustChannel string `hcl:"rustChannel,optional"`
	// RustProfile: Can be set to "minimal" or "default".
	RustProfile string `hcl:"rustProfile,optional"`
	// ExtraRustComponents: Extra rust components to be added with the build process
	ExtraRustComponents []string `hcl:"extraRustComponents,optional"`
}
