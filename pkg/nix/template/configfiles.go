package template

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/buildsafedev/bsf/pkg/hcl2nix"
)

// ConfigFiles holds config files to export
type ConfigFiles struct {
	// Name of the config
	Name string

	// DestinationDir is the directory to copy config files to in the container
	// This directory will be created in the root of the image.
	// default: /
	DestinationDir string

	// Name of files to copy from root of project
	Files []string
}

const (
	confAttrTmpl = `
   configs = forEachSupportedSystem({pkgs,...}:{
        {{range $config := .}}
	config_{{$config.Name}} = pkgs.stdenvNoCC.mkDerivation {
        name = " {{$config.Name}}  ";
        src = ../.;
        dontUnpack = true;
        dontBuild = true;
        phases = [ "installPhase" ];
        installPhase = ''
          mkdir -p $out
          mkdir -p $out/tmp
          mkdir -p $out/{{$config.DestinationDir}}
          {{ range $config.Files}}cp -r $src/{{.}} $out/{{$config.DestinationDir}} 
		  {{end}}
        '';
       };
       {{end}}
     });
 `
)

func quote(s string) string {
	return strings.ReplaceAll(s, "\n", "\\n")
}

// GenerateConfigAttr generates the Nix attribute set for config files
func GenerateConfigAttr(cfg []ConfigFiles) (*string, error) {
	tmpl, err := template.New("confAttr").Funcs(template.FuncMap{
		"quote": quote,
	}).
		Parse(confAttrTmpl)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, cfg)
	if err != nil {
		return nil, err
	}

	result := buf.String()
	return &result, nil
}

func hclConfFilesToConfFiles(hclCfg []hcl2nix.ConfigFiles) []ConfigFiles {
	cfg := make([]ConfigFiles, len(hclCfg))
	for i, f := range hclCfg {
		cfg[i] = ConfigFiles{
			Name:           f.Name,
			DestinationDir: f.DestinationDir,
			Files:          f.Files,
		}
	}
	return cfg
}
