package sbom

import (
	"time"

	"github.com/awalterschulze/gographviz"
	"github.com/bom-squad/protobom/pkg/sbom"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func sbomTools() []*sbom.Tool {
	return []*sbom.Tool{
		{
			Name: "bsf",
			// TODO: this version should be picked from ldFlags.
			Version: "0.0.0",
			Vendor:  "buildsafe",
		},
	}
}

// PackageGraphToSBOM converts the package graph to a SBOM
func PackageGraphToSBOM(appNode *sbom.Node, lockFile *hcl2nix.LockFile, graph *gographviz.Graph) *sbom.Document {
	document := sbom.NewDocument()

	document.Metadata.Tools = sbomTools()

	document.Metadata.Name = "SBOM for " + appNode.Name

	document.NodeList.AddRootNode(appNode)

	parseDotGraph(document, appNode, graph)

	parseLockfileToSBOMNodes(document, appNode, lockFile)

	return document

}

func parseLockfileToSBOMNodes(document *sbom.Document, appNode *sbom.Node, lf *hcl2nix.LockFile) {
	for _, pkg := range lf.Packages {
		snode := sbom.Node{
			Id:               PurlFromNameVersion(pkg.Package.Name, pkg.Package.Version),
			Type:             sbom.Node_PACKAGE,
			Name:             pkg.Package.Name,
			Version:          pkg.Package.Version,
			UrlHome:          pkg.Package.Homepage,
			UrlDownload:      pkg.Package.Homepage,
			Licenses:         []string{pkg.Package.SpdxId},
			LicenseConcluded: pkg.Package.SpdxId,
			Description:      pkg.Package.Description,
			ReleaseDate:      timestamppb.New(time.Unix(int64(pkg.Package.EpochSeconds), 0)),
		}

		document.NodeList.AddNode(&snode)
		if pkg.Runtime {
			document.NodeList.RelateNodeAtID(&snode, appNode.Id, sbom.Edge_runtimeDependency)
		} else {
			document.NodeList.RelateNodeAtID(&snode, appNode.Id, sbom.Edge_devDependency)
			document.NodeList.RelateNodeAtID(&snode, appNode.Id, sbom.Edge_devTool)
		}
	}

	return
}

func parseDotGraph(document *sbom.Document, appNode *sbom.Node, graph *gographviz.Graph) {
	for _, node := range graph.Nodes.Nodes {
		name := node.Attrs["name"]
		version := node.Attrs["version"]
		if name == appNode.Name {
			continue
		}

		snode := sbom.Node{
			Name:           name,
			Type:           sbom.Node_PACKAGE,
			Id:             PurlFromNameVersion(name, version),
			Version:        version,
			PrimaryPurpose: []sbom.Purpose{sbom.Purpose_DATA},
			Hashes: map[int32]string{
				int32(sbom.HashAlgorithm_SHA256): node.Attrs["hash"],
			},
		}
		document.NodeList.AddNode(&snode)
		document.NodeList.RelateNodeAtID(&snode, appNode.Id, sbom.Edge_contains)
	}

	return
}

// PurlFromNameVersion returns a package url for the given name and version
func PurlFromNameVersion(name, version string) string {
	return "pkg:" + "nix/" + name + "@v" + version
}
