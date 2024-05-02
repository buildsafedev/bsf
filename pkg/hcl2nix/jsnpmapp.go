package hcl2nix

// RustApp defines the parameters for a rust application.
type JsNpmApp struct {
	// PackageName: Name of the Package
	PackageName string `hcl:"packageName"`
	// PackageJsonPath: Source path to the package.json and package-lock.json file.
	PackageRoot string `hcl:"packageRoot"`
}
