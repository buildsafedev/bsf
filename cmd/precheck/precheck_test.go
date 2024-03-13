package precheck

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
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

	tempDir := os.TempDir()

	tmpFile, err := ioutil.TempFile(tempDir, "daemon.json")
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	data := struct {
		Features struct {
			ContainerdSnapshotter bool `json:"containerd-snapshotter"`
		} `json:"features"`
	}{
		Features: struct {
			ContainerdSnapshotter bool `json:"containerd-snapshotter"`
		}{
			ContainerdSnapshotter: true,
		},
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	if _, err := tmpFile.Write(jsonData); err != nil {
		fmt.Println("Error writing to temporary file:", err)
		return
	}

	dockerdaemon_json = tmpFile.Name()
	resp, err := IsSnapshotterEnabled()
	if err != nil {
		t.Error(err)
	}

	if !resp {
		t.Errorf(" ⚠️  containerd image store is not enabled [ https://docs.docker.com/storage/containerd/ ]")
	}

}
