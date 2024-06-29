package template

import (
	"io"
	"text/template"
)

// RemoteFile holds remote flake parameters
type RemoteFile struct {
	Name           string
	Version        string
	PlatformURLs   map[string]string
	PlatformHashes map[string]string
	Binaries       []string
}

var remoteFlakeTmpl = `{
  description = "{{.Name}}";

  inputs = { };

  outputs = { self, nixpkgs }:
    let
      supportedSystems = [ {{range $platform, $url := .PlatformURLs}}"{{$platform}}" {{end}}];
      forEachSystem = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
        inherit system;
        pkgs = import nixpkgs { inherit system; };
      });

      getBinaryUrl = { system, version }: 
        let
          urls = {
            {{range $platform, $url := .PlatformURLs}}"{{$platform}}" = "{{$url}}";
            {{end}}
          };
        in urls.${system};

      getBinaryHash = { system }:
        let 
          hashes = {
            {{range $platform, $hash := .PlatformHashes}}"{{$platform}}" = "{{$hash}}";
            {{end}}
          };
        in hashes.${system};
        
    in
    {
      packages = forEachSystem ({ pkgs, ... }: {
        default = pkgs.stdenvNoCC.mkDerivation rec {
          name = "{{.Name}}";
          version = "{{.Version}}";

          src = pkgs.fetchurl {
            url = getBinaryUrl { system = pkgs.system; version = version; };
            sha256 = getBinaryHash {system = pkgs.system;};
          };

          dontUnpack = false;
          phases = [ "unpackPhase" "installPhase" ];
          dontBuild = true;

          unpackPhase = ''
            tar -xzf $src
          '';

          installPhase = ''
            mkdir -p $out/bin
            cp -r ./bin/* $out/bin/
            chmod +x $out/bin/*
          '';
        };
      });
    };
}`

// GenerateRemoteFlake generates a flake to fetch remote files
func GenerateRemoteFlake(fl RemoteFile, wr io.Writer) error {
	data := RemoteFile{
		Name:           fl.Name,
		Version:        fl.Version,
		PlatformURLs:   fl.PlatformURLs,
		PlatformHashes: fl.PlatformHashes,
		Binaries:       fl.Binaries,
	}

	tmpl, err := template.New("remoteFlake").Parse(remoteFlakeTmpl)
	if err != nil {
		return err
	}

	err = tmpl.Execute(wr, data)
	if err != nil {
		return err
	}

	return nil
}
