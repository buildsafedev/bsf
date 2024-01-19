package search

import (
	"log"
	"sort"

	buildsafev1 "github.com/buildsafedev/cloud-api/apis/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewClientWithAddr initializes a Client with a specific API address
func NewClientWithAddr(addr string) (buildsafev1.SearchServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
