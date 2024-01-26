
{
	description = "";
	
	inputs = {
		 nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14.url = "github:nixos/nixpkgs/ac5c1886fd9fe49748d7ab80accc4c847481df14";
		 nixpkgs-1ebb7d7bba2953a4223956cfb5f068b0095f84a7.url = "github:nixos/nixpkgs/1ebb7d7bba2953a4223956cfb5f068b0095f84a7";
		 nixpkgs-a6515b40d18282e1deee8129209998b4b62e4bec.url = "github:nixos/nixpkgs/a6515b40d18282e1deee8129209998b4b62e4bec";
			
		nixpkgs.url = "github:nixos/nixpkgs/nixos-23.11";
	};
	
	outputs = { self, nixpkgs,  nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14, 
	 nixpkgs-1ebb7d7bba2953a4223956cfb5f068b0095f84a7, 
	 nixpkgs-a6515b40d18282e1deee8129209998b4b62e4bec, 
	 }: let
	  supportedSystems = [ "x86_64-linux" "aarch64-darwin" "x86_64-darwin" "aarch64-linux" ];
	  forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
		 nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14-pkgs = import nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14 { inherit system; };
		 nixpkgs-1ebb7d7bba2953a4223956cfb5f068b0095f84a7-pkgs = import nixpkgs-1ebb7d7bba2953a4223956cfb5f068b0095f84a7 { inherit system; };
		 nixpkgs-a6515b40d18282e1deee8129209998b4b62e4bec-pkgs = import nixpkgs-a6515b40d18282e1deee8129209998b4b62e4bec { inherit system; };
		
		pkgs = import nixpkgs { inherit system; };
	  });
	in {
	  packages = forEachSupportedSystem ({ pkgs,  nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14-pkgs, 
		 nixpkgs-1ebb7d7bba2953a4223956cfb5f068b0095f84a7-pkgs, 
		 nixpkgs-a6515b40d18282e1deee8129209998b4b62e4bec-pkgs, 
		 }: {
		default = pkgs.callPackage ./default.nix {
		  
		};
	  });
	
	  devShells = forEachSupportedSystem ({ pkgs,  nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14-pkgs, 
		 nixpkgs-1ebb7d7bba2953a4223956cfb5f068b0095f84a7-pkgs, 
		 nixpkgs-a6515b40d18282e1deee8129209998b4b62e4bec-pkgs, 
		 }: {
		devShell = pkgs.mkShell {
		  # The Nix packages provided in the environment
		  packages =  [
			nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14-pkgs.delve  
			nixpkgs-1ebb7d7bba2953a4223956cfb5f068b0095f84a7-pkgs.go  
			nixpkgs-a6515b40d18282e1deee8129209998b4b62e4bec-pkgs.goreleaser  
			nixpkgs-1ebb7d7bba2953a4223956cfb5f068b0095f84a7-pkgs.gotools  
			nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14-pkgs.ngrok  
			
		  ];
		};
	  });
	
	  runtimeEnvs = forEachSupportedSystem ({ pkgs,  nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14-pkgs,  nixpkgs-1ebb7d7bba2953a4223956cfb5f068b0095f84a7-pkgs,  nixpkgs-a6515b40d18282e1deee8129209998b4b62e4bec-pkgs,  }: {
		runtime = pkgs.buildEnv {
		  name = "runtimeenv";
		  paths = [ 
			nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14-pkgs.cacert   
			
		   ];
		};
	});
	};
}
