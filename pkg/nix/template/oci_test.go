package template

import (
	"fmt"
	"strings"
	"testing"
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

	if !strings.Contains(*result, "inputs.self.devEnvs.${system}.development") {
		t.Errorf("Generated template does not contain expected devEnvs reference")
	}
	fmt.Println(*result)
}
