
{
	description = "";
	
	inputs = {
		 nixpkgs-a731d0cb71c58f56895f71a5b02eda2962a46746.url = "github:nixos/nixpkgs/a731d0cb71c58f56895f71a5b02eda2962a46746";
		 nixpkgs-7445ccd775d8b892fc56448d17345443a05f7fb4.url = "github:nixos/nixpkgs/7445ccd775d8b892fc56448d17345443a05f7fb4";
		 nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14.url = "github:nixos/nixpkgs/ac5c1886fd9fe49748d7ab80accc4c847481df14";
			
		nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
		 gomod2nix.url = "github:nix-community/gomod2nix";
		gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
		
		
		
		
		
		

		 
	};
	
	outputs = inputs@{ self, nixpkgs, 
	 gomod2nix, 
	
	
	
	
	 nixpkgs-a731d0cb71c58f56895f71a5b02eda2962a46746, 
	 nixpkgs-7445ccd775d8b892fc56448d17345443a05f7fb4, 
	 nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14, 
	 }: let
	  supportedSystems = [ "x86_64-linux" "aarch64-darwin" "x86_64-darwin" "aarch64-linux" ];
	  
	  
	  forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
		inherit system;
		
		 nixpkgs-a731d0cb71c58f56895f71a5b02eda2962a46746-pkgs = import nixpkgs-a731d0cb71c58f56895f71a5b02eda2962a46746 { inherit system; };
		 nixpkgs-7445ccd775d8b892fc56448d17345443a05f7fb4-pkgs = import nixpkgs-7445ccd775d8b892fc56448d17345443a05f7fb4 { inherit system; };
		 nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14-pkgs = import nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14 { inherit system; };
		
		 buildGoApplication = gomod2nix.legacyPackages.${system}.buildGoApplication;
		pkgs = import nixpkgs { inherit system;  };
		
		
	  });
	in {
	  packages = forEachSupportedSystem ({ pkgs,
		 buildGoApplication, 
		
		
		 nixpkgs-a731d0cb71c58f56895f71a5b02eda2962a46746-pkgs, 
		 nixpkgs-7445ccd775d8b892fc56448d17345443a05f7fb4-pkgs, 
		 nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14-pkgs, 
		 ... }: {
		default = pkgs.callPackage ./default.nix {
			 inherit buildGoApplication;
			go = pkgs.go_1_22; 
			
			
			
		};
	  });
	
	  devShells = forEachSupportedSystem ({ pkgs, 
		 buildGoApplication, 
		
		
		 nixpkgs-a731d0cb71c58f56895f71a5b02eda2962a46746-pkgs, 
		 nixpkgs-7445ccd775d8b892fc56448d17345443a05f7fb4-pkgs, 
		 nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14-pkgs, 
		 ... }: {
		devShell = pkgs.mkShell {
		  # The Nix packages provided in the environment
		  packages =  [
			nixpkgs-a731d0cb71c58f56895f71a5b02eda2962a46746-pkgs.delve  
			nixpkgs-a731d0cb71c58f56895f71a5b02eda2962a46746-pkgs.go  
			nixpkgs-7445ccd775d8b892fc56448d17345443a05f7fb4-pkgs.go-task  
			nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14-pkgs.gotools  
			
		  ];
		};
	  });
	
	  runtimeEnvs = forEachSupportedSystem ({ pkgs,
		 buildGoApplication, 
		
		
		 nixpkgs-a731d0cb71c58f56895f71a5b02eda2962a46746-pkgs,  nixpkgs-7445ccd775d8b892fc56448d17345443a05f7fb4-pkgs,  nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14-pkgs,  ... }: {
		runtime = pkgs.buildEnv {
		  name = "runtimeenv";
		  paths = [ 
			nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14-pkgs.cacert   
			
		   ];
		};
	   });

	   devEnvs = forEachSupportedSystem ({ pkgs,
		 buildGoApplication, 
		
		
	    nixpkgs-a731d0cb71c58f56895f71a5b02eda2962a46746-pkgs,  nixpkgs-7445ccd775d8b892fc56448d17345443a05f7fb4-pkgs,  nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14-pkgs,  ... }: {
		development = pkgs.buildEnv {
		  name = "devenv";
		  paths = [ 
			nixpkgs-a731d0cb71c58f56895f71a5b02eda2962a46746-pkgs.delve  
			nixpkgs-a731d0cb71c58f56895f71a5b02eda2962a46746-pkgs.go  
			nixpkgs-7445ccd775d8b892fc56448d17345443a05f7fb4-pkgs.go-task  
			nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14-pkgs.gotools  
			
		   ];
		};
	   });
       
	   
	   
	};
}
