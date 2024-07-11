package update

import (
	"sort"
	"strings"

	buildsafev1 "github.com/buildsafedev/bsf-apis/go/buildsafe/v1"
	"golang.org/x/mod/semver"
)

const (
	// UpdateTypePatch represents patch update type
	UpdateTypePatch = iota
	// UpdateTypeMinor represents minor update type
	UpdateTypeMinor
	// UpdateTypeDate represents date update type
	UpdateTypeDate
	// UpdateTypePinned represents pinned update type
	UpdateTypePinned
)

// ParseUpdateType parses the update type of the package
func ParseUpdateType(pkg string) int {
	if strings.Contains(pkg, "~") {
		return UpdateTypePatch
	} else if strings.Contains(pkg, "^") {
		return UpdateTypeMinor
	} else if strings.Contains(pkg, "#") {
		return UpdateTypeDate
	} else {
		return UpdateTypePinned
	}
}

// ParsePackage parses the package as given in bsf.hcl and returns the name and version. It removes the update type information.
func ParsePackage(pkg string) (name, version string) {
	nameWithVersion := strings.SplitN(pkg, "@", 2)
	if len(nameWithVersion) < 2 {
		return "", ""
	}

	version = nameWithVersion[1]
	if strings.HasPrefix(version, "~") || strings.HasPrefix(version, "^") {
		version = version[1:]
	}

	return nameWithVersion[0], version
}

// GetDateBasedVersion returns the latest date version for the given version.
func GetDateBasedVersion(v *buildsafev1.FetchPackagesResponse, version string) string {
	if v == nil {
		return ""
	}

	// sorting epochs
	sort.Slice(v.Packages, func(i, j int) bool {
		return v.Packages[i].EpochSeconds > v.Packages[j].EpochSeconds
	})

	return v.Packages[0].Version
}

// GetLatestPatchVersion returns the latest patch version for the given version.
func GetLatestPatchVersion(v *buildsafev1.FetchPackagesResponse, version string) string {
	if v == nil {
		return ""
	}

	desiredMajorMinor := semver.MajorMinor("v" + version)

	validVersions := make([]string, 0, len(v.Packages))
	for _, v := range v.Packages {
		if semver.MajorMinor("v"+v.Version) == desiredMajorMinor {
			validVersions = append(validVersions, "v"+v.Version)
		}
	}

	semver.Sort(validVersions)

	return strings.TrimPrefix(validVersions[len(validVersions)-1], "v")
}

// GetLatestMinorVersion returns the latest minor version for the given version.
func GetLatestMinorVersion(v *buildsafev1.FetchPackagesResponse, version string) string {
	if v == nil {
		return ""
	}

	desiredMajorMinor := semver.Major("v" + version)

	validVersions := make([]string, 0, len(v.Packages))
	for _, v := range v.Packages {
		if semver.Major("v"+v.Version) == desiredMajorMinor {
			validVersions = append(validVersions, "v"+v.Version)
		}
	}
	semver.Sort(validVersions)

	return strings.TrimPrefix(validVersions[len(validVersions)-1], "v")
}

// TrimVersionInfo trims the version information and returns the package name and version.
func TrimVersionInfo(pkg string) (string, string) {
	if pkg == "" {
		return "", ""
	}
	s := strings.Split(pkg, "@")
	name := s[0]
	version := s[1]

	if strings.HasPrefix(version, "~") {
		version = strings.TrimPrefix(version, "~")
	}

	if strings.HasPrefix(version, "^") {
		version = strings.TrimPrefix(version, "^")
	}

	if strings.HasPrefix(version, "#") {
		version = strings.TrimPrefix(version, "#")
	}

	if strings.HasPrefix(version, "v") {
		version = strings.TrimPrefix(version, "v")
	}

	return name, version
}

// ComparePackages compares the devVersions and the runtimeVersions
func ComparePackages(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	counts := make(map[string]bool)
	for _, item := range a {
		counts[item] = true
	}
	for _, item := range b {
		if !counts[item] {
			return false
		}
	}

	return true
}
