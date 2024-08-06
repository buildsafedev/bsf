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
            "aarch64-darwin" = "sha256-+YhwgHcwqB4XMZvRRFQ/H4/NSBZqHyUo5Yk0OfF7mEA=";
            "x86_64-linux" = "sha256-2nSZbhCxtNK/YnJ7COUHmzRX7gVUV5UITuDmND1vdtM=";
            "x86_64-darwin" = "sha256-c11fPY5+8FcYOGDt2LAcuDy/Sx7hLO7Y7IKrRhIjQio=";
            "aarch64-linux" = "sha256-64ATYM0u5kMnqudvY8iC0NjaHscp2fKPiRmGKN917A8=";
          };
          in hash.${system};
        
    in
    {
      packages = forEachSystem ({ pkgs, ... }: {
        default = pkgs.stdenvNoCC.mkDerivation  rec {
          name = "bsf";
          version = "0.2.4";

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
