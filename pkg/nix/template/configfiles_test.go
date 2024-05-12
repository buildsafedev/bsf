package template

import (
	"fmt"
	"strings"
	"testing"
)

func TestGenerateConfigAttr(t *testing.T) {
	cfg := []ConfigFiles{
		{
			Name:           "Test1",
			DestinationDir: "dir1",
			Files:          []string{"file1", "file2"},
		},
		{
			Name:           "Test2",
			DestinationDir: "dir2",
			Files:          []string{"file3", "file4"},
		},
	}

	result, err := GenerateConfigAttr(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(*result, "config_Test1") || !strings.Contains(*result, "config_Test2") {
		t.Errorf("Generated template does not contain expected config names")
	}

	if !strings.Contains(*result, "cp -r $src/file1 $out/dir1") || !strings.Contains(*result, "cp -r $src/file2 $out/dir1") {
		t.Errorf("Generated template does not contain expected file copy commands for Test1")
	}

	if !strings.Contains(*result, "cp -r $src/file3 $out/dir2") || !strings.Contains(*result, "cp -r $src/file4 $out/dir2") {
		t.Errorf("Generated template does not contain expected file copy commands for Test2")
	}

	fmt.Println(*result)
}
