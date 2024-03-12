
{
	description = "";
	
	inputs = {
		 nixpkgs-1ebb7d7bba2953a4223956cfb5f068b0095f84a7.url = "github:nixos/nixpkgs/1ebb7d7bba2953a4223956cfb5f068b0095f84a7";
		 nixpkgs-a731d0cb71c58f56895f71a5b02eda2962a46746.url = "github:nixos/nixpkgs/a731d0cb71c58f56895f71a5b02eda2962a46746";
		 nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14.url = "github:nixos/nixpkgs/ac5c1886fd9fe49748d7ab80accc4c847481df14";
		 nixpkgs-5a7b241264578c55cd25aa7422121aef072ce588.url = "github:nixos/nixpkgs/5a7b241264578c55cd25aa7422121aef072ce588";
			
		nixpkgs.url = "github:nixos/nixpkgs/nixos-23.11";
		 gomod2nix.url = "github:nix-community/gomod2nix";
		gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
	};
	
	outputs = { self, nixpkgs, 
	 gomod2nix, 
	 nixpkgs-1ebb7d7bba2953a4223956cfb5f068b0095f84a7, 
	 nixpkgs-a731d0cb71c58f56895f71a5b02eda2962a46746, 
	 nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14, 
	 nixpkgs-5a7b241264578c55cd25aa7422121aef072ce588, 
	 }: let
	  supportedSystems = [ "x86_64-linux" "aarch64-darwin" "x86_64-darwin" "aarch64-linux" ];
	  forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
		 nixpkgs-1ebb7d7bba2953a4223956cfb5f068b0095f84a7-pkgs = import nixpkgs-1ebb7d7bba2953a4223956cfb5f068b0095f84a7 { inherit system; };
		 nixpkgs-a731d0cb71c58f56895f71a5b02eda2962a46746-pkgs = import nixpkgs-a731d0cb71c58f56895f71a5b02eda2962a46746 { inherit system; };
		 nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14-pkgs = import nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14 { inherit system; };
		 nixpkgs-5a7b241264578c55cd25aa7422121aef072ce588-pkgs = import nixpkgs-5a7b241264578c55cd25aa7422121aef072ce588 { inherit system; };
		
		 buildGoApplication = gomod2nix.legacyPackages.${system}.buildGoApplication;
		pkgs = import nixpkgs { inherit system; };
	  });
	in {
	  packages = forEachSupportedSystem ({ pkgs,
		 buildGoApplication, 
		 nixpkgs-1ebb7d7bba2953a4223956cfb5f068b0095f84a7-pkgs, 
		 nixpkgs-a731d0cb71c58f56895f71a5b02eda2962a46746-pkgs, 
		 nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14-pkgs, 
		 nixpkgs-5a7b241264578c55cd25aa7422121aef072ce588-pkgs, 
		 }: {
		default = pkgs.callPackage ./default.nix {
			 inherit buildGoApplication; 
		};
	  });
	
	  devShells = forEachSupportedSystem ({ pkgs, 
		 buildGoApplication, 
		 nixpkgs-1ebb7d7bba2953a4223956cfb5f068b0095f84a7-pkgs, 
		 nixpkgs-a731d0cb71c58f56895f71a5b02eda2962a46746-pkgs, 
		 nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14-pkgs, 
		 nixpkgs-5a7b241264578c55cd25aa7422121aef072ce588-pkgs, 
		 }: {
		devShell = pkgs.mkShell {
		  # The Nix packages provided in the environment
		  packages =  [
			nixpkgs-a731d0cb71c58f56895f71a5b02eda2962a46746-pkgs.delve  
			nixpkgs-5a7b241264578c55cd25aa7422121aef072ce588-pkgs.go_1_22  
			nixpkgs-5a7b241264578c55cd25aa7422121aef072ce588-pkgs.goreleaser  
			nixpkgs-1ebb7d7bba2953a4223956cfb5f068b0095f84a7-pkgs.gotools  
			
		  ];
		};
	  });
	
	  runtimeEnvs = forEachSupportedSystem ({ pkgs,
		 buildGoApplication, 
		 nixpkgs-1ebb7d7bba2953a4223956cfb5f068b0095f84a7-pkgs,  nixpkgs-a731d0cb71c58f56895f71a5b02eda2962a46746-pkgs,  nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14-pkgs,  nixpkgs-5a7b241264578c55cd25aa7422121aef072ce588-pkgs,  }: {
		runtime = pkgs.buildEnv {
		  name = "runtimeenv";
		  paths = [ 
			nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14-pkgs.cacert   
			
		   ];
		};
	});
	};
}
