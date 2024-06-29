package template

import (
	"os"
	"testing"
)

func TestGenerateRemoteFlake(t *testing.T) {
	testRemoteFile := RemoteFile{
		Name:           "testName",
		Version:        "1.0.0",
		PlatformURLs:   map[string]string{"linux/amd64": "http://example.com/linux/amd64.tar.gz"},
		PlatformHashes: map[string]string{"linux/amd64": "sha256-abcdef1234567890"},
		Binaries:       []string{"bin1", "bin2"},
	}

	err := GenerateRemoteFlake(testRemoteFile, os.Stdout)
	if err != nil {
		t.Errorf("GenerateRemoteFlake returned an error: %v", err)
		return
	}

}
