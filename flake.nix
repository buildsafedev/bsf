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
            "aarch64-darwin" = "sha256-W553BRHGgRZXZHtcRxPZa/QNBjk5gyrWHOAKPApEhwQ=";
            "x86_64-linux" = "sha256-p7o//av/y3qrDupdKnCi/YUy+pr1GsSRth6TFci4zEs=";
            "x86_64-darwin" = "sha256-VP55fwVCg63oQjuS5ttQfUUx4MsTeDSYFPcJVSVZ9OY=";
            "aarch64-linux" = "sha256-0WR3AkpnQ/qcrRJ5HXzjVYZcF5IMYhyNhFDd9GfaiDk=";
          };
          in hash.${system};
        
    in
    {
      packages = forEachSystem ({ pkgs, ... }: {
        default = pkgs.stdenvNoCC.mkDerivation  rec {
          name = "bsf";
          version = "0.3.0";

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
