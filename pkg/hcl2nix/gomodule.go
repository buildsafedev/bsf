package hcl2nix

// GoModule is a struct to hold Nix Go module parameters
type GoModule struct {
	Name       string `hcl:"name"`
	SourcePath string `hcl:"src"`
	VendorHash string `hcl:"vendorHash"`
	DoCheck    bool   `hcl:"doCheck,optional"`
	Meta       *Meta  `hcl:"meta"`
}

// Meta holds Nix lib meta parameters
type Meta struct {
	Description string `hcl:"description"`
}
