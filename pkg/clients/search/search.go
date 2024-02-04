package search

import (
	"crypto/tls"
	"log"
	"sort"

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
