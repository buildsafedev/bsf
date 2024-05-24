package sbom

import (
	"encoding/json"
	"time"

	"github.com/awalterschulze/gographviz"
	"github.com/bom-squad/protobom/pkg/formats"
	"github.com/bom-squad/protobom/pkg/sbom"
	"github.com/bom-squad/protobom/pkg/writer"
	intoto "github.com/in-toto/in-toto-golang/in_toto"
	intotoCom "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/common"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	bio "github.com/buildsafedev/bsf/pkg/io"
	nixcmd "github.com/buildsafedev/bsf/pkg/nix/cmd"
)

// Statement is a struct to hold the SBOM in SPDX format
type Statement struct {
	intoto.StatementHeader
	// Predicate can be SPDX or CDX format
	Predicate interface{}
}

// NewStatement creates a new SBOM
func NewStatement(appDetails *nixcmd.App) *Statement {
	st := Statement{}
	st.Type = "https://in-toto.io/Statement/v1"
	st.Subject = []intoto.Subject{
		{
			Name: appDetails.Name,
			Digest: intotoCom.DigestSet{
				"sha256": appDetails.BinaryHash,
			},
		},
		{
			Name: "result-" + appDetails.Name,
			Digest: intotoCom.DigestSet{
				"sha256": appDetails.ResultHash,
			},
		},
	}
	return &st
}

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

// ToJSON returns the statement in JSON format
func (s *Statement) ToJSON(bom *sbom.Document, format formats.Format) ([]byte, error) {
	s.PredicateType = "https://spdx.github.io/spdx-spec/v2.3/"
	if format == formats.CDX15JSON {
		s.PredicateType = "https://cyclonedx.org/specification/overview/"
	}

	w := writer.New()
	bomBytes := bio.NewBufferCloser()
	err := w.WriteStreamWithOptions(bom, bomBytes, &writer.Options{
		Format: format,
	})
	if err != nil {
		return nil, err
	}

	// Unmarshal the Predicate into an interface{} to prettify it
	var pred interface{}
	err = json.Unmarshal(bomBytes.Bytes(), &pred)
	if err != nil {
		return nil, err
	}
	s.Predicate = pred

	return json.Marshal(s)
}

func parseLockfileToSBOMNodes(document *sbom.Document, appNode *sbom.Node, lf *hcl2nix.LockFile) {
	for _, pkg := range lf.Packages {
		snode := sbom.Node{
			Id: GeneratePurl(pkg.Package.Name, pkg.Package.Version, "", ""),
			Identifiers: map[int32]string{
				int32(sbom.SoftwareIdentifierType_CPE23): pkg.Package.Cpe,
				int32(sbom.SoftwareIdentifierType_PURL):  GeneratePurl(pkg.Package.Name, pkg.Package.Version, "", ""),
			},
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
		if name == "" || name == appNode.Name {
			continue
		}

		snode := sbom.Node{
			Name:           name,
			Type:           sbom.Node_PACKAGE,
			Id:             GeneratePurl(name, version, "", ""),
			Version:        version,
			PrimaryPurpose: []sbom.Purpose{sbom.Purpose_DATA},
			Identifiers: map[int32]string{
				int32(sbom.SoftwareIdentifierType_PURL): GeneratePurl(name, version, "", ""),
			},
			Hashes: map[int32]string{
				int32(sbom.HashAlgorithm_SHA256): node.Attrs["hash"],
			},
		}
		document.NodeList.AddNode(&snode)
		document.NodeList.RelateNodeAtID(&snode, appNode.Id, sbom.Edge_contains)
	}

	return
}

// GeneratePurl returns a package url for the given name and version
func GeneratePurl(name, version, os, arch string) string {
	purl := "pkg:" + "nix/" + name + "@v" + version
	if os != "" && arch != "" {
		purl += "?os=" + os + "&arch=" + arch
	}

	return purl
}
