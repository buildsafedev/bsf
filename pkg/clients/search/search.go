package search

import (
	"crypto/tls"
	"log"
	"sort"

	"golang.org/x/mod/semver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	buildsafev1 "github.com/buildsafedev/bsf-apis/go/buildsafe/v1"
)

// NewClientWithAddr initializes a Client with a specific API address
func NewClientWithAddr(addr string, tlsSkip bool) (buildsafev1.SearchServiceClient, error) {
	var creds credentials.TransportCredentials
	if tlsSkip {
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		creds = credentials.NewTLS(tlsConfig)
	} else {
		creds = insecure.NewCredentials()
	}
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return buildsafev1.NewSearchServiceClient(conn), nil
}

// SortPackagesWithTimestamp sorts packages with timestamp with latest being the first element
func SortPackagesWithTimestamp(packageVersions []*buildsafev1.Package) []*buildsafev1.Package {
	if packageVersions == nil {
		return nil
	}
	sort.Slice(packageVersions, func(i, j int) bool {
		return packageVersions[i].EpochSeconds > packageVersions[j].EpochSeconds
	})
	return packageVersions
}

// SortPackagesWithSamver sorts packages with samver with latest being the first element
func SortPackagesWithVersion(packageVersions []*buildsafev1.Package) []*buildsafev1.Package {
	if packageVersions == nil {
		return nil
	}

	for i := range packageVersions {
		packageVersions[i].Version = "v" + packageVersions[i].Version
	}
	sort.Slice(packageVersions, func(i, j int) bool {
		if semver.Compare(packageVersions[i].Version, packageVersions[j].Version) > 0 {
			return true
		}
		return false
	})
	for i := range packageVersions {
		packageVersions[i].Version = packageVersions[i].Version[1:]
	}

	return packageVersions
}

// SortPackages sorts the pkg based on  Semantic Versioning
func SortPackages(packageVersions []*buildsafev1.Package) []*buildsafev1.Package {

	semverPkgs := make([]*buildsafev1.Package, 0)
	nonsemverPkgs := make([]*buildsafev1.Package, 0)

	for i := range packageVersions {
		packageVersions[i].Version = "v" + packageVersions[i].Version
		if semver.IsValid(packageVersions[i].Version) {
			packageVersions[i].Version = packageVersions[i].Version[1:]
			semverPkgs = append(semverPkgs, packageVersions[i])
			continue
		} else {
			packageVersions[i].Version = packageVersions[i].Version[1:]
			nonsemverPkgs = append(nonsemverPkgs, packageVersions[i])
			continue
		}
	}

	if nonsemverPkgs != nil {
		semverPkgs := SortPackagesWithVersion(semverPkgs)
		nonsemverPkgs := SortPackagesWithTimestamp(nonsemverPkgs)
		semverPkgs = append(semverPkgs, nonsemverPkgs...)
		return semverPkgs
	}

	return SortPackagesWithVersion(packageVersions)
}
