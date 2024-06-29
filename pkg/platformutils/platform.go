package platformutils

import (
	"runtime"
	"strings"
)

// FindPlatform returns OS and arch. Ex: if linux/amd64 is passed, it will return linux amd64.
// If no platform is passed, it will return the current OS and arch.
func FindPlatform(platform string) (string, string) {
	if platform == "" {
		return runtime.GOOS, runtime.GOARCH
	}
	osarch := strings.Split(platform, "/")
	if len(osarch) != 2 {
		return runtime.GOOS, runtime.GOARCH
	}

	return osarch[0], osarch[1]
}

// OSArchToArchOS converts "OS/Arch" format to "Arch-OS" format.
func OSArchToArchOS(input string) string {
	mappings := map[string]string{
		"linux/amd64":  "x86_64-linux",
		"linux/arm64":  "aarch64-linux",
		"darwin/amd64": "x86_64-darwin",
		"darwin/arm64": "aarch64-darwin",
	}

	if val, ok := mappings[input]; ok {
		return val
	}
	return input
}

// ArchOSToOSArch converts "Arch-OS" format to "OS/Arch" format.
func ArchOSToOSArch(input string) string {
	mappings := map[string]string{
		"x86_64-linux":   "linux/amd64",
		"aarch64-linux":  "linux/arm64",
		"x86_64-darwin":  "darwin/amd64",
		"aarch64-darwin": "darwin/arm64",
	}

	if val, ok := mappings[input]; ok {
		return val
	}
	return input
}

// FormatType represents the format of a platform string.
type FormatType string

const (
	// OSArchFormat represents the "OS/Arch" format.
	OSArchFormat FormatType = "OS/Arch"
	// ArchOSFormat represents the "Arch-OS" format.
	ArchOSFormat FormatType = "Arch-OS"
	// UnknownFormat represents an unknown format.
	UnknownFormat FormatType = "Unknown"
)

// DetermineFormat takes a platform string and returns a FormatType indicating its format.
func DetermineFormat(platform string) FormatType {
	// Known prefixes and suffixes for "OS/Arch" and "Arch-OS" formats
	osArchPrefixes := []string{"linux/", "darwin/"}
	archOSPostfixes := []string{"-linux", "-darwin"}

	// Check for "OS/Arch" format by looking for "/"
	if strings.Contains(platform, "/") {
		for _, prefix := range osArchPrefixes {
			if strings.HasPrefix(platform, prefix) {
				return OSArchFormat
			}
		}
	}

	// Check for "Arch-OS" format by looking for known postfixes
	for _, postfix := range archOSPostfixes {
		if strings.HasSuffix(platform, postfix) {
			return ArchOSFormat
		}
	}

	// If neither, return "Unknown"
	return UnknownFormat
}
