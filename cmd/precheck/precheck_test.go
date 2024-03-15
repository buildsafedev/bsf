package precheck

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	bsfbuild "github.com/buildsafedev/bsf/pkg/build"
)

func TestCheckVersionConstraint(t *testing.T) {
	tests := []struct {
		currentVer   string
		requiredVer  string
		expectResult bool
	}{
		{"v1.2.3", "v1.0.0", true},
		{"v1.2.3", "v1.2.3", true},
		{"v1.2.3", "v1.2.4", false},
		{"v1.2.3", "v1.3.0", false},
	}

	for _, test := range tests {
		result := checkVersionGreater(test.currentVer, test.requiredVer)
		if result != test.expectResult {
			t.Errorf("checkVersionConstraint(%s, %s) = %t, want %t", test.currentVer, test.requiredVer, result, test.expectResult)
		}
	}
}

func TestIsSnapshotterEnabled(t *testing.T) {

	tests := []struct {
		name     string
		data     interface{}
		expected interface{}
	}{
		{
			name: "ContainerdSnapshotter is true",
			data: struct {
				Features struct {
					ContainerdSnapshotter bool `json:"containerd-snapshotter"`
				} `json:"features"`
			}{
				Features: struct {
					ContainerdSnapshotter bool `json:"containerd-snapshotter"`
				}{
					ContainerdSnapshotter: true,
				},
			},
			expected: map[string]interface{}{
				"features": map[string]interface{}{
					"containerd-snapshotter": true,
				},
			},
		}, {
			name: "ContainerdSnapshotter is false",
			data: struct {
				Features struct {
					ContainerdSnapshotter bool `json:"containerd-snapshotter"`
				} `json:"features"`
			}{
				Features: struct {
					ContainerdSnapshotter bool `json:"containerd-snapshotter"`
				}{
					ContainerdSnapshotter: false,
				},
			},
			expected: map[string]interface{}{
				"features": map[string]interface{}{
					"containerd-snapshotter": false,
				},
			},
		},
	}

	for _, testCases := range tests {
		t.Run(testCases.name, func(t *testing.T) {

			tempDir := os.TempDir()

			tmpFile, err := os.CreateTemp(tempDir, "daemon.json")
			if err != nil {
				fmt.Println("Error creating temporary file:", err)
				return
			}
			defer os.Remove(tmpFile.Name())

			bsfbuild.DockerDaemonJSON = tmpFile.Name()

			jsonData, err := json.MarshalIndent(testCases.data, "", "  ")
			if err != nil {
				fmt.Println("Error marshaling JSON:", err)
				return
			}
			if _, err := tmpFile.Write(jsonData); err != nil {
				fmt.Println("Error writing to temporary file:", err)
				return
			}

			result, err := bsfbuild.ReadDockerfile()
			if err != nil {
				t.Fail()
			}

			if !reflect.DeepEqual(testCases.expected, result) {
				t.Fail()
			}
		})
	}
}
