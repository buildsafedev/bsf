package hcl2nix

// RustApp defines the parameters for a rust application.
type RustApp struct {
	// WorkspaceSrc: Source of the workspace.
	WorkspaceSrc string `hcl:"workspaceSrc"`
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
	// RootFeatures: A list of activated features on your workspace's crates.
	RootFeatures []string `hcl:"rootFeatures,optional"`
	// FetchCrateAlternativeRegistry: A fetcher for crates on alternative registries.
	FetchCrateAlternativeRegistry string `hcl:"fetchCrateAlternativeRegistry,optional"`
	// HostPlatformCpu: Equivalent to rust's target-cpu codegen option.
	HostPlatformCpu string `hcl:"hostPlatformCpu,optional"`
	// HostPlatformFeatures: Equivalent to rust's target-feature codegen option.
	HostPlatformFeatures []string `hcl:"hostPlatformFeatures,optional"`
	// ExtraRustComponents: Extra rust components to be added with the build process
	ExtraRustComponents []string `hcl:"extraRustComponents,optional"`
	// CargoUnstableFlags: Flags that affect cargo unstable features.
	CargoUnstableFlags []string `hcl:"cargoUnstableFlags,optional"`
	// RustcLinkFlags: Pass extra flags directly to rustc during non-build invocations
	RustcLinkFlags []string `hcl:"rustcLinkFlags,optional"`
	// RustcBuildFlags: Pass extra flags directly to Rustc during build invocations
	RustcBuildFlags []string `hcl:"rustcBuildFlags,optional"`
}
