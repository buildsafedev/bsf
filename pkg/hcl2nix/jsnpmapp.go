package hcl2nix

// JsNpmApp defines the parameters for a Javascript application.
type JsNpmApp struct {
	// PackageName: Name of the Package
	PackageName string `hcl:"packageName"`
	// PackageRoot: Source path to root file.
	PackageRoot string `hcl:"packageRoot"`
	// PackageJsonPath: Path to package.json file.
	PackageJsonPath string `hcl:"packageJsonPath,optional"`
	// PackageLockPath: Path to package-lock.json file.
	PackageLockPath string `hcl:"packageLockPath,optional"`
}
