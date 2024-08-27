package template

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateOCIAttr(t *testing.T) {
	artifacts := []OCIArtifact{
		{
			Artifact:      "Test1",
			Name:          "artifact1",
			EnvVars:       []string{"VAR1=value1", "VAR2=value2"},
			ExposedPorts:  []string{"8080", "8081"},
			ImportConfigs: []string{"config1", "config2"},
		},
		{
			Artifact:     "Test2",
			Name:         "artifact2",
			EnvVars:      []string{"VAR3=value3", "VAR4=value4"},
			ExposedPorts: []string{"8082", "8083"},
		},
	}

	result, err := GenerateOCIAttr(artifacts)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(*result, "ociImage_Test1") || !strings.Contains(*result, "ociImage_Test2") {
		t.Errorf("Generated template does not contain expected OCI image names")
	}

	if !strings.Contains(*result, "VAR1=value1") || !strings.Contains(*result, "VAR2=value2") {
		t.Errorf("Generated template does not contain expected environment variables for Test1")
	}

	if !strings.Contains(*result, "VAR3=value3") || !strings.Contains(*result, "VAR4=value4") {
		t.Errorf("Generated template does not contain expected environment variables for Test2")
	}

	if !strings.Contains(*result, "8080") || !strings.Contains(*result, "8081") {
		t.Errorf("Generated template does not contain expected exposed ports for Test1")
	}

	if !strings.Contains(*result, "8082") || !strings.Contains(*result, "8083") {
		t.Errorf("Generated template does not contain expected exposed ports for Test2")
	}

	fmt.Println(*result)
}

func TestGetReqPkgs(t *testing.T) {
	tests := []struct {
		name     string
		layers   []string
		fl       Flake
		expected [][]string
	}{
		{
			name:   "Split runtime packages",
			layers: []string{"split(packages.runtime)"},
			fl: Flake{
				RuntimePackages: map[string]string{
					"go":    "d919897915f0f91216d2501b617d670deee993a0",
					"nginx": "e3f1b7d7e09f8f5371b2cb1e3a0bc6c3b03f78a0",
				},
			},
			expected: [][]string{
				{"nixpkgs-d919897915f0f91216d2501b617d670deee993a0-pkgs.go"},
				{"nixpkgs-e3f1b7d7e09f8f5371b2cb1e3a0bc6c3b03f78a0-pkgs.nginx"},
			},
		},
		{
			name:   "Split dev packages",
			layers: []string{"split(packages.dev)"},
			fl: Flake{
				DevPackages: map[string]string{
					"bash": "f2c55c8e7d3d843f75e2f18c8bf707b8a77c8a0f",
					"zsh":  "a5b69c8e5b6d78364c4f938ac142c8e6a6b2d3a0",
				},
			},
			expected: [][]string{
				{"nixpkgs-f2c55c8e7d3d843f75e2f18c8bf707b8a77c8a0f-pkgs.bash"},
				{"nixpkgs-a5b69c8e5b6d78364c4f938ac142c8e6a6b2d3a0-pkgs.zsh"},
			},
		},
		{
			name:   "Specific dev package",
			layers: []string{"packages.dev.bash"},
			fl: Flake{
				DevPackages: map[string]string{
					"bash": "f2c55c8e7d3d843f75e2f18c8bf707b8a77c8a0f",
					"zsh":  "a5b69c8e5b6d78364c4f938ac142c8e6a6b2d3a0",
				},
			},
			expected: [][]string{
				{"nixpkgs-f2c55c8e7d3d843f75e2f18c8bf707b8a77c8a0f-pkgs.bash"},
			},
		},
		{
			name:   "Specific runtime package",
			layers: []string{"packages.runtime.go"},
			fl: Flake{
				RuntimePackages: map[string]string{
					"go":    "d919897915f0f91216d2501b617d670deee993a0",
					"nginx": "e3f1b7d7e09f8f5371b2cb1e3a0bc6c3b03f78a0",
				},
			},
			expected: [][]string{
				{"nixpkgs-d919897915f0f91216d2501b617d670deee993a0-pkgs.go"},
			},
		},
		{
			name:   "Combined dev and runtime packages",
			layers: []string{"packages.dev.go + packages.runtime.nginx"},
			fl: Flake{
				DevPackages: map[string]string{
					"go": "d919897915f0f91216d2501b617d670deee993a0",
				},
				RuntimePackages: map[string]string{
					"nginx": "e3f1b7d7e09f8f5371b2cb1e3a0bc6c3b03f78a0",
				},
			},
			expected: [][]string{
				{
					"nixpkgs-d919897915f0f91216d2501b617d670deee993a0-pkgs.go",
					"nixpkgs-e3f1b7d7e09f8f5371b2cb1e3a0bc6c3b03f78a0-pkgs.nginx",
				},
			},
		},
		{
			name:   "Combined dev pkg and whole runtime",
			layers: []string{"packages.dev.go + packages.runtime"},
			fl: Flake{
				DevPackages: map[string]string{
					"go": "d919897915f0f91216d2501b617d670deee993a0",
				},
				RuntimePackages: map[string]string{
					"nginx": "e3f1b7d7e09f8f5371b2cb1e3a0bc6c3b03f78a0",
				},
			},
			expected: [][]string{
				{
					"nixpkgs-d919897915f0f91216d2501b617d670deee993a0-pkgs.go",
					"inputs.self.runtimeEnvs.${system}.runtime",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getReqPkgs(tt.layers, tt.fl)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetLayers(t *testing.T) {
	tests := []struct {
		name     string
		layers   []string
		fl       Flake
		expected []string
	}{
		{
			name:   "Basic runtime and dev split",
			layers: []string{"split(packages.runtime)", "split(packages.dev)"},
			fl: Flake{
				DevPackages: map[string]string{
					"bash": "f2c55c8e7d3d843f75e2f18c8bf707b8a77c8a0f",
				},
				RuntimePackages: map[string]string{
					"go": "d919897915f0f91216d2501b617d670deee993a0",
				},
			},
			expected: []string{
				`(nix2containerPkgs.nix2container.buildLayer { 
					copyToRoot = [
						nixpkgs-d919897915f0f91216d2501b617d670deee993a0-pkgs.go
					];
				})`,
				`(nix2containerPkgs.nix2container.buildLayer { 
					copyToRoot = [
						nixpkgs-f2c55c8e7d3d843f75e2f18c8bf707b8a77c8a0f-pkgs.bash
					];
				})`,
			},
		},
		{
			name:   "Combined dev pkg and split runtime",
			layers: []string{"packages.dev.go + split(packages.runtime)"},
			fl: Flake{
				DevPackages: map[string]string{
					"go": "7445ccd775d8b892fc56448d17345443a05f7fb4",
				},
				RuntimePackages: map[string]string{
					"cacert": "ac5c1886fd9fe49748d7ab80accc4c847481df14",
				},
			},
			expected: []string{
				`(nix2containerPkgs.nix2container.buildLayer { 
					copyToRoot = [
						nixpkgs-7445ccd775d8b892fc56448d17345443a05f7fb4-pkgs.go
						nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14-pkgs.cacert
					];
				})`,
			},
		},
		{
			name:   "Single dev package",
			layers: []string{"packages.dev.go"},
			fl: Flake{
				DevPackages: map[string]string{
					"go": "7445ccd775d8b892fc56448d17345443a05f7fb4",
				},
				RuntimePackages: map[string]string{},
			},
			expected: []string{
				`(nix2containerPkgs.nix2container.buildLayer { 
					copyToRoot = [
						nixpkgs-7445ccd775d8b892fc56448d17345443a05f7fb4-pkgs.go
					];
				})`,
			},
		},
		{
			name:   "Combined dev and runtime packages",
			layers: []string{"packages.dev.go + packages.runtime.cacert"},
			fl: Flake{
				DevPackages: map[string]string{
					"go": "7445ccd775d8b892fc56448d17345443a05f7fb4",
				},
				RuntimePackages: map[string]string{
					"cacert": "ac5c1886fd9fe49748d7ab80accc4c847481df14",
				},
			},
			expected: []string{
				`(nix2containerPkgs.nix2container.buildLayer { 
					copyToRoot = [
						nixpkgs-7445ccd775d8b892fc56448d17345443a05f7fb4-pkgs.go
						nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14-pkgs.cacert
					];
				})`,
			},
		},
		{
			name:   "Combined dev and whole runtime",
			layers: []string{"packages.dev.go + packages.runtime"},
			fl: Flake{
				DevPackages: map[string]string{
					"go": "7445ccd775d8b892fc56448d17345443a05f7fb4",
				},
				RuntimePackages: map[string]string{},
			},
			expected: []string{
				`(nix2containerPkgs.nix2container.buildLayer { 
					copyToRoot = [
						nixpkgs-7445ccd775d8b892fc56448d17345443a05f7fb4-pkgs.go
						inputs.self.runtimeEnvs.${system}.runtime
					];
				})`,
			},
		},
		{
			name:   "Split dev and runtime",
			layers: []string{"split(packages.dev)", "split(packages.runtime)"},
			fl: Flake{
				DevPackages: map[string]string{
					"go":   "7445ccd775d8b892fc56448d17345443a05f7fb4",
					"bash": "f2c55c8e7d3d843f75e2f18c8bf707b8a77c8a0f",
				},
				RuntimePackages: map[string]string{
					"cacert": "ac5c1886fd9fe49748d7ab80accc4c847481df14",
				},
			},
			expected: []string{
				`(nix2containerPkgs.nix2container.buildLayer { 
					copyToRoot = [
						nixpkgs-7445ccd775d8b892fc56448d17345443a05f7fb4-pkgs.go
					];
				})`,
				`(nix2containerPkgs.nix2container.buildLayer { 
					copyToRoot = [
						nixpkgs-f2c55c8e7d3d843f75e2f18c8bf707b8a77c8a0f-pkgs.bash
					];
				})`,
				`(nix2containerPkgs.nix2container.buildLayer { 
					copyToRoot = [
						nixpkgs-ac5c1886fd9fe49748d7ab80accc4c847481df14-pkgs.cacert
					];
				})`,
			},
		},
		{
			name:   "No split, just raw layers",
			layers: []string{"packages.dev", "packages.runtime"},
			fl: Flake{
				DevPackages: map[string]string{
					"bash": "f2c55c8e7d3d843f75e2f18c8bf707b8a77c8a0f",
				},
				RuntimePackages: map[string]string{
					"go": "d919897915f0f91216d2501b617d670deee993a0",
				},
			},
			expected: []string{
				`(nix2containerPkgs.nix2container.buildLayer { 
					copyToRoot = [
						inputs.self.devEnvs.${system}.development
					];
				})`,
				`(nix2containerPkgs.nix2container.buildLayer { 
					copyToRoot = [
						inputs.self.runtimeEnvs.${system}.runtime
					];
				})`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getLayers(tt.layers, tt.fl)
			for i := range result {
				result[i] = normalizeString(result[i])
			}
			for i := range tt.expected {
				tt.expected[i] = normalizeString(tt.expected[i])
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

func normalizeString(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
