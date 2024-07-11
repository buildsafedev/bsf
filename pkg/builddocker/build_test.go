package builddocker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEditDockerFile(t *testing.T) {
	tests := []struct {
		name        string
		lines       []string
		isDev       bool
		tag         string
		expectedRes []string
		expectError bool
	}{
		{
			name: "Dev",
			lines: []string{
				"FROM ubuntu:18.04 # bsfimage:dev",
				"RUN apt-get update",
			},
			isDev: true,
			tag:   "latest",
			expectedRes: []string{
				"FROM ubuntu:latest # bsfimage:dev",
				"RUN apt-get update",
			},
			expectError: false,
		},
		{
			name: "Runtime",
			lines: []string{
				"FROM ubuntu:18.04 # bsfimage:runtime",
				"RUN apt-get update",
			},
			isDev: false,
			tag:   "latest",
			expectedRes: []string{
				"FROM ubuntu:latest # bsfimage:runtime",
				"RUN apt-get update",
			},
			expectError: false,
		},
		{
			name: "No FROM Command with bsf tag",
			lines: []string{
				"FROM ubuntu:latest",
				"RUN apt-get update",
			},
			isDev:       true,
			tag:         "latest",
			expectedRes: nil,
			expectError: true,
		},
		{
			name: "No FROM Command",
			lines: []string{
				"RUN apt-get update",
			},
			isDev:       true,
			tag:         "latest",
			expectedRes: nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := editDockerfile(tt.lines, tt.isDev, tt.tag)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRes, res)
			}
		})
	}
}
