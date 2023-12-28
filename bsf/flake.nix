
{
	description = "bsf flake";
	
	inputs = {
		 nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688.url = "github.com/nixos/nixpkgs/a89ba043dda559ebc57fc6f1fa8cf3a0b207f688";
			
	};
	
	outputs = { self,  nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688, 
	 }: let
	  supportedSystems = [ "x86_64-linux" "aarch64-darwin" "x86_64-darwin" "aarch64-linux" ];
	  forEachSupportedSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
		 nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688-pkgs = import nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688 { inherit system; };
		
	  });
	in {
	  # packages = forEachSupportedSystem ({  nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688-pkgs, 
		#  }: {
		# default = pkgs.callPackage ./default.nix {
		  
		# };
	  # });
	
	  devShells = forEachSupportedSystem ({  nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688-pkgs, 
		 }: {
		devShell = pkgs.mkShell {
		  # The Nix packages provided in the environment
		  packages =  [
			nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688-pkgs.delve  
			nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688-pkgs.go  
			nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688-pkgs.gotools  
			
		  ];
		};
	  });
	
	  runtimeEnvs = forEachSupportedSystem ({  nixpkgs-a89ba043dda559ebc57fc6f1fa8cf3a0b207f688-pkgs,  }: {
		runtime = pkgs.buildEnv {
		  name = "runtimeenv";
		  paths = [ 
			
		   ];
		};
	});
	};
}
