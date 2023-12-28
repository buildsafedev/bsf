
{
	description = "";
	
	inputs = {
		 nixpkgs-eeee184c00a7e542d2a837252a0ed4e74dd27dc5.url = "github:nixos/nixpkgs/eeee184c00a7e542d2a837252a0ed4e74dd27dc5";
		 nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688.url = "github:nixos/nixpkgs/a89ba043dda559ebc57fc6f1fa8cf3a0b207f688";
			
		nixpkgs.url = "github:nixos/nixpkgs/nixos-23.11";
	};
	
	outputs = { self, nixpkgs,  nixpkgs-eeee184c00a7e542d2a837252a0ed4e74dd27dc5, 
	 nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688, 
	 }: let
	  supportedSystems = [ "x86_64-linux" "aarch64-darwin" "x86_64-darwin" "aarch64-linux" ];
	  forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
		 nixpkgs-eeee184c00a7e542d2a837252a0ed4e74dd27dc5-pkgs = import nixpkgs-eeee184c00a7e542d2a837252a0ed4e74dd27dc5 { inherit system; };
		 nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688-pkgs = import nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688 { inherit system; };
		
		pkgs = import nixpkgs { inherit system; };
	  });
	in {
	  packages = forEachSupportedSystem ({ pkgs,  nixpkgs-eeee184c00a7e542d2a837252a0ed4e74dd27dc5-pkgs, 
		 nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688-pkgs, 
		 }: {
		default = pkgs.callPackage ./default.nix {
		  
		};
	  });
	
	  devShells = forEachSupportedSystem ({ pkgs,  nixpkgs-eeee184c00a7e542d2a837252a0ed4e74dd27dc5-pkgs, 
		 nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688-pkgs, 
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
	
	  runtimeEnvs = forEachSupportedSystem ({ pkgs,  nixpkgs-eeee184c00a7e542d2a837252a0ed4e74dd27dc5-pkgs,  nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688-pkgs,  }: {
		runtime = pkgs.buildEnv {
		  name = "runtimeenv";
		  paths = [ 
			nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688-pkgs.cacert   
			
		   ];
		};
	});
	};
}
