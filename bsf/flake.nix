
{
	description = "";
	
	inputs = {
		 nixpkgs-96d1259aefb7350ebc4fbcc0718447fe30321151.url = "github:nixos/nixpkgs/96d1259aefb7350ebc4fbcc0718447fe30321151";
		 nixpkgs-c5e85c459830b30d1e54ca4ae0d4d37fc23adbe2.url = "github:nixos/nixpkgs/c5e85c459830b30d1e54ca4ae0d4d37fc23adbe2";
		 nixpkgs-eeee184c00a7e542d2a837252a0ed4e74dd27dc5.url = "github:nixos/nixpkgs/eeee184c00a7e542d2a837252a0ed4e74dd27dc5";
		 nixpkgs-7fbe081b14e1363801a6a60e105b403f37048ea6.url = "github:nixos/nixpkgs/7fbe081b14e1363801a6a60e105b403f37048ea6";
			
		nixpkgs.url = "github:nixos/nixpkgs/nixos-23.11";
	};
	
	outputs = { self, nixpkgs,  nixpkgs-96d1259aefb7350ebc4fbcc0718447fe30321151, 
	 nixpkgs-c5e85c459830b30d1e54ca4ae0d4d37fc23adbe2, 
	 nixpkgs-eeee184c00a7e542d2a837252a0ed4e74dd27dc5, 
	 nixpkgs-7fbe081b14e1363801a6a60e105b403f37048ea6, 
	 }: let
	  supportedSystems = [ "x86_64-linux" "aarch64-darwin" "x86_64-darwin" "aarch64-linux" ];
	  forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
		 nixpkgs-96d1259aefb7350ebc4fbcc0718447fe30321151-pkgs = import nixpkgs-96d1259aefb7350ebc4fbcc0718447fe30321151 { inherit system; };
		 nixpkgs-c5e85c459830b30d1e54ca4ae0d4d37fc23adbe2-pkgs = import nixpkgs-c5e85c459830b30d1e54ca4ae0d4d37fc23adbe2 { inherit system; };
		 nixpkgs-eeee184c00a7e542d2a837252a0ed4e74dd27dc5-pkgs = import nixpkgs-eeee184c00a7e542d2a837252a0ed4e74dd27dc5 { inherit system; };
		 nixpkgs-7fbe081b14e1363801a6a60e105b403f37048ea6-pkgs = import nixpkgs-7fbe081b14e1363801a6a60e105b403f37048ea6 { inherit system; };
		
		pkgs = import nixpkgs { inherit system; };
	  });
	in {
	  packages = forEachSupportedSystem ({ pkgs,  nixpkgs-96d1259aefb7350ebc4fbcc0718447fe30321151-pkgs, 
		 nixpkgs-c5e85c459830b30d1e54ca4ae0d4d37fc23adbe2-pkgs, 
		 nixpkgs-eeee184c00a7e542d2a837252a0ed4e74dd27dc5-pkgs, 
		 nixpkgs-7fbe081b14e1363801a6a60e105b403f37048ea6-pkgs, 
		 }: {
		default = pkgs.callPackage ./default.nix {
		  
		};
	  });
	
	  devShells = forEachSupportedSystem ({ pkgs,  nixpkgs-96d1259aefb7350ebc4fbcc0718447fe30321151-pkgs, 
		 nixpkgs-c5e85c459830b30d1e54ca4ae0d4d37fc23adbe2-pkgs, 
		 nixpkgs-eeee184c00a7e542d2a837252a0ed4e74dd27dc5-pkgs, 
		 nixpkgs-7fbe081b14e1363801a6a60e105b403f37048ea6-pkgs, 
		 }: {
		devShell = pkgs.mkShell {
		  # The Nix packages provided in the environment
		  packages =  [
			nixpkgs-96d1259aefb7350ebc4fbcc0718447fe30321151-pkgs.delve  
			nixpkgs-eeee184c00a7e542d2a837252a0ed4e74dd27dc5-pkgs.go  
			nixpkgs-7fbe081b14e1363801a6a60e105b403f37048ea6-pkgs.goreleaser  
			nixpkgs-c5e85c459830b30d1e54ca4ae0d4d37fc23adbe2-pkgs.gotools  
			
		  ];
		};
	  });
	
	  runtimeEnvs = forEachSupportedSystem ({ pkgs,  nixpkgs-96d1259aefb7350ebc4fbcc0718447fe30321151-pkgs,  nixpkgs-c5e85c459830b30d1e54ca4ae0d4d37fc23adbe2-pkgs,  nixpkgs-eeee184c00a7e542d2a837252a0ed4e74dd27dc5-pkgs,  nixpkgs-7fbe081b14e1363801a6a60e105b403f37048ea6-pkgs,  }: {
		runtime = pkgs.buildEnv {
		  name = "runtimeenv";
		  paths = [ 
			nixpkgs-96d1259aefb7350ebc4fbcc0718447fe30321151-pkgs.cacert   
			
		   ];
		};
	});
	};
}
