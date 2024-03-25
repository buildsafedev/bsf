package provenance

import (
	"context"
	"encoding/json"

	"github.com/awalterschulze/gographviz"
	slsav1 "github.com/buildsafedev/bsf/pkg/slsa/v1"
	intoto "github.com/in-toto/in-toto-golang/in_toto"
	intotoCom "github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/common"
	"github.com/nix-community/go-nix/pkg/derivation"
	"github.com/nix-community/go-nix/pkg/derivation/store"
	"google.golang.org/protobuf/types/known/structpb"

	nixcmd "github.com/buildsafedev/bsf/pkg/nix/cmd"
)

// Statement is a struct to hold the provenance statement
type Statement struct {
	intoto.StatementHeader
	Predicate slsav1.Provenance
}

// NewStatement creates a new provenance statement
func NewStatement(appDetails *nixcmd.App) *Statement {
	st := Statement{}
	st.Type = "https://in-toto.io/Statement/v1"
	st.PredicateType = "https://slsa.dev/provenance/v1"
	st.Subject = []intoto.Subject{
		{
			Name: appDetails.Name,
			Digest: intotoCom.DigestSet{
				"sha256": appDetails.Hash,
			},
		},
	}
	return &st
}

// GetDerivation gets the derivation from the store
func GetDerivation(dir string) (*derivation.Derivation, error) {
	store, err := store.NewFromURI("")
	if err != nil {
		return nil, err
	}

	drv, err := store.Get(context.Background(), dir)
	if err != nil {
		return nil, err
	}

	return drv, err
}

// FromDerivationClosure gets the provenance from the graph
func (s *Statement) FromDerivationClosure(drvPath string, drv *derivation.Derivation, graph *gographviz.Graph) error {
	rds := graphToResourceDes(graph)
	prov := slsav1.Provenance{
		BuildDefinition: &slsav1.BuildDefinition{
			// TODO: this should link to a doc like https://slsa-framework.github.io/github-actions-buildtypes/workflow/v1
			BuildType: "nix",
			ExternalParameters: &structpb.Struct{
				// TODO: we should have git information in here.
				Fields: map[string]*structpb.Value{
					"derivation": {
						Kind: &structpb.Value_StringValue{
							StringValue: drvPath,
						},
					},
				},
			},
			ResolvedDependencies: rds,
			InternalParameters: &structpb.Struct{
				Fields: make(map[string]*structpb.Value),
			},
		},
		RunDetails: &slsav1.RunDetails{
			Builder: &slsav1.Builder{
				// TODO: discuss with SLSA community of what the ID should be in case of Nix
				Id: drv.Builder,
				Version: map[string]string{
					"nix": ">=2.18.0",
				},
			},
		},
	}

	s.Predicate = prov

	return nil
}

func graphToResourceDes(graph *gographviz.Graph) []*slsav1.ResourceDescriptor {
	rds := make([]*slsav1.ResourceDescriptor, 0, len(graph.Nodes.Nodes))

	for _, node := range graph.Nodes.Nodes {
		rds = append(rds, &slsav1.ResourceDescriptor{
			Uri:  "/nix/store/" + nixcmd.CleanNameFromGraph(node.Name),
			Name: node.Attrs["name"],
			Digest: map[string]string{
				"sha256": node.Attrs["hash"],
			},
			Annotations: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"version": {
						Kind: &structpb.Value_StringValue{
							StringValue: node.Attrs["version"],
						},
					},
				},
			},
		})
	}

	return rds
}

// ToJSON converts the provenance statement to JSON
func (s *Statement) ToJSON() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}
