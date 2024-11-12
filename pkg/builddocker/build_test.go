package builddocker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEditDockerFile(t *testing.T) {
	tests := []struct {
		ociArtifactName string
		name            string
		lines           []string
		tag             string
		expectedRes     []string
		expectError     bool
	}{
		{
			name:            "FROM Command with bsf tag",
			ociArtifactName: "docker.io/buildsafe/bsf",
			lines: []string{
				"FROM docker.io/buildsafe/bsf:v1.1.1 AS build",
				"RUN apt-get update",
			},
			tag: "latest",
			expectedRes: []string{
				"FROM docker.io/buildsafe/bsf:latest AS build",
				"RUN apt-get update",
			},
			expectError: false,
		},
		{
			name:            "No FROM Command with bsf tag",
			ociArtifactName: "docker.io/buildsafe/bsf",
			lines: []string{
				"FROM ubuntu:latest",
				"RUN apt-get update",
			},
			tag:         "latest",
			expectedRes: nil,
			expectError: true,
		},
		{
			name:            "No FROM Command",
			ociArtifactName: "docker.io/buildsafe/bsf",
			lines: []string{
				"RUN apt-get update",
			},
			tag:         "latest",
			expectedRes: nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := editDockerfile(tt.lines, tt.ociArtifactName, tt.tag)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRes, res)
			}
		})
	}
}
