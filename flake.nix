{
  description = "bsf";

  inputs = {
    # nixpkgs.url = "github:NixOS/nixpkgs";
  };

  outputs = { self, nixpkgs }:
    let
      supportedSystems = [ "x86_64-linux" "aarch64-darwin" "x86_64-darwin" "aarch64-linux" ];
      forEachSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
        inherit system;
        pkgs = import nixpkgs { inherit system; };
      });

      getBinaryUrl = { system, version }: 
        let
          urls = {
            "x86_64-linux" = "https://github.com/buildsafedev/bsf/releases/download/v${version}/bsf_linux_amd64";
            "aarch64-darwin" = "https://github.com/buildsafedev/bsf/releases/download/v${version}/bsf_darwin_arm64";
            "x86_64-darwin" = "https://github.com/buildsafedev/bsf/releases/download/v${version}/bsf_darwin_amd64";
            "aarch64-linux" = "https://github.com/buildsafedev/bsf/releases/download/v${version}/bsf_linux_arm64";
          };
        in urls.${system};

        getBinaryHash = { system }:
          let 
          hash = {
            "aarch64-darwin" = "sha256-Uz/Iu9uH6Pqu1DFnjvksqgvxfH8hy+yLGDCoSgR+Lxg=";
            "x86_64-linux" = "sha256-I2NRGrkBqOpAYA2QlFOZMW0RKCwPS5CWvXWQjMCNU6U=";
            "x86_64-darwin" = "sha256-m+6kEGJ8OaIJ2UNdmM8ccmdrlJfAGrMq7L4geZNQKsg=";
            "aarch64-linux" = "sha256-l6LN722i6J0chbN4LqQ557DMiih+CWjBkhwfeGgoY20=";
          };
          in hash.${system};
        
    in
    {
      packages = forEachSystem ({ pkgs, ... }: {
        default = pkgs.stdenvNoCC.mkDerivation  rec {
          name = "bsf";
          version = "0.2.1";

          src = pkgs.fetchurl {
            url = getBinaryUrl { system = pkgs.system; version = version; };
            sha256 = getBinaryHash {system = pkgs.system;};
          };

          dontUnpack = true;
          phases = [ "installPhase" ];
          dontBuild = true;

          installPhase = ''
            mkdir -p $out/bin
            cp $src $out/bin/bsf
            chmod +x $out/bin/bsf
          '';

          shellHook = ''
            echo $src
            echo "bsf is available: $out/bin/bsf"
            $out/bin/bsf 
          '';
        };
      });
    };
}
