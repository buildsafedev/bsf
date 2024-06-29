package platformutils

import "testing"

func TestDetermineFormat(t *testing.T) {
	tests := []struct {
		name     string
		platform string
		want     FormatType
	}{
		{"Linux AMD64", "linux/amd64", OSArchFormat},
		{"Darwin ARM64", "darwin/arm64", OSArchFormat},
		{"X86_64 Linux", "x86_64-linux", ArchOSFormat},
		{"Aarch64 Darwin", "aarch64-darwin", ArchOSFormat},
		{"Unknown Format", "windows/amd64", UnknownFormat},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetermineFormat(tt.platform); got != tt.want {
				t.Errorf("DetermineFormat(%s) = %v, want %v", tt.platform, got, tt.want)
			}
		})
	}
}
