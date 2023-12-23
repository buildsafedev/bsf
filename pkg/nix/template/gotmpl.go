package template

// GoPackage holds go package parameters
type GoPackage struct {
	Name       string
	SourcePath string
	VendorHash string
	Meta       Meta
}

const (
	golangTmpl = `
	{
	   lib,
	   stdenv,
	   buildGoModule,
	   ... 
	 }: buildGoModule {
	   name = "";
	 
	   src = ../.;  
	 
	   vendorHash = lib.fakeHash;
	 
	   meta = with lib; {
		 description = "";
	   };
	 }
	`
)
