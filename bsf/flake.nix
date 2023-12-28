
{
	description = "";
	
	inputs = {
		 nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688.url = "github:nixos/nixpkgs/a89ba043dda559ebc57fc6f1fa8cf3a0b207f688";
		 nixpkgs-eeee184c00a7e542d2a837252a0ed4e74dd27dc5.url = "github:nixos/nixpkgs/eeee184c00a7e542d2a837252a0ed4e74dd27dc5";
			
		nixpkgs.url = "github:nixos/nixpkgs/nixos-23.11";
	};
	
	outputs = { self, nixpkgs,  nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688, 
	 nixpkgs-eeee184c00a7e542d2a837252a0ed4e74dd27dc5, 
	 }: let
	  supportedSystems = [ "x86_64-linux" "aarch64-darwin" "x86_64-darwin" "aarch64-linux" ];
	  forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
		 nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688-pkgs = import nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688 { inherit system; };
		 nixpkgs-eeee184c00a7e542d2a837252a0ed4e74dd27dc5-pkgs = import nixpkgs-eeee184c00a7e542d2a837252a0ed4e74dd27dc5 { inherit system; };
		
		pkgs = import nixpkgs { inherit system; };
	  });
	in {
	  packages = forEachSupportedSystem ({ pkgs,  nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688-pkgs, 
		 nixpkgs-eeee184c00a7e542d2a837252a0ed4e74dd27dc5-pkgs, 
		 }: {
		default = pkgs.callPackage ./default.nix {
		  
		};
	  });
	
	  devShells = forEachSupportedSystem ({ pkgs,  nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688-pkgs, 
		 nixpkgs-eeee184c00a7e542d2a837252a0ed4e74dd27dc5-pkgs, 
		 }: {
		devShell = pkgs.mkShell {
		  # The Nix packages provided in the environment
		  packages =  [
			nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688-pkgs.delve  
			nixpkgs-eeee184c00a7e542d2a837252a0ed4e74dd27dc5-pkgs.go  
			nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688-pkgs.gotools  
			
		  ];
		};
	  });
	
	  runtimeEnvs = forEachSupportedSystem ({ pkgs,  nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688-pkgs,  nixpkgs-eeee184c00a7e542d2a837252a0ed4e74dd27dc5-pkgs,  }: {
		runtime = pkgs.buildEnv {
		  name = "runtimeenv";
		  paths = [ 
			nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688-pkgs.cacert   
			
		   ];
		};
	});
	};
}
