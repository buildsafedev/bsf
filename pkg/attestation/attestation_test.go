package attestation

import (
	"testing"

	intoto "github.com/in-toto/in-toto-golang/in_toto"
	"github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/common"
	"github.com/stretchr/testify/assert"
)

func TestGetRelevantStatementsForOnePred(t *testing.T) {
	tests := []struct {
		name        string
		psMap       map[string][]intoto.Statement
		predType    string
		subject     string
		expectedRes []intoto.Statement
	}{
		{
			name: "get relevant statements for one pred",
			psMap: map[string][]intoto.Statement{
				"provenance": {
					{
						StatementHeader: intoto.StatementHeader{
							Type:          "https://in-toto.io/Statement/v1",
							PredicateType: "https://slsa.dev/provenance/v1",
							Subject: []intoto.Subject{
								{
									Name: "0.2.0",
									Digest: common.DigestSet{
										"sha256": "",
									},
								},
							},
						},
						Predicate: struct{}{}, // Replace struct{}{} with the actual type of Predicate
					},
				},
			},
			predType: "provenance",
			subject:  "",
			expectedRes: []intoto.Statement{
				{
					StatementHeader: intoto.StatementHeader{
						Type:          "https://in-toto.io/Statement/v1",
						PredicateType: "https://slsa.dev/provenance/v1",
						Subject: []intoto.Subject{
							{
								Name: "0.2.0",
								Digest: common.DigestSet{
									"sha256": "",
								},
							},
						},
					},
					Predicate: struct{}{}, // Replace struct{}{} with the actual type of Predicate
				},
			},
		},
		{
			name: "get relevant statements for multiple pred",
			psMap: map[string][]intoto.Statement{
				"provenance": {
					{
						StatementHeader: intoto.StatementHeader{
							Type:          "https://in-toto.io/Statement/v1",
							PredicateType: "https://slsa.dev/provenance/v1",
							Subject: []intoto.Subject{
								{
									Name: "0.1.0",
									Digest: common.DigestSet{
										"sha256": "",
									},
								},
							},
						},
						Predicate: struct{}{},
					},
					{
						StatementHeader: intoto.StatementHeader{
							Type:          "https://in-toto.io/Statement/v1",
							PredicateType: "https://slsa.dev/provenance/v1",
							Subject: []intoto.Subject{
								{
									Name: "0.2.0",
									Digest: common.DigestSet{
										"sha256": "",
									},
								},
							},
						},
						Predicate: struct{}{},
					},
				},
				"spdx": {
					{
						StatementHeader: intoto.StatementHeader{
							Type:          "https://in-toto.io/Statement/v1",
							PredicateType: "https://spdx.github.io/spdx-spec/v2.3/",
							Subject: []intoto.Subject{
								{
									Name: "0.2.0",
									Digest: common.DigestSet{
										"sha256": "",
									},
								},
							},
						},
						Predicate: struct{}{},
					},
				},
			},
			predType: "provenance",
			subject:  "",
			expectedRes: []intoto.Statement{
				{
					StatementHeader: intoto.StatementHeader{
						Type:          "https://in-toto.io/Statement/v1",
						PredicateType: "https://slsa.dev/provenance/v1",
						Subject: []intoto.Subject{
							{
								Name: "0.1.0",
								Digest: common.DigestSet{
									"sha256": "",
								},
							},
						},
					},
					Predicate: struct{}{},
				},
				{
					StatementHeader: intoto.StatementHeader{
						Type:          "https://in-toto.io/Statement/v1",
						PredicateType: "https://slsa.dev/provenance/v1",
						Subject: []intoto.Subject{
							{
								Name: "0.2.0",
								Digest: common.DigestSet{
									"sha256": "",
								},
							},
						},
					},
					Predicate: struct{}{},
				},
			},
		},
		{
			name: "get relevant statements for one pred with subject",
			psMap: map[string][]intoto.Statement{
				"provenance": {
					{
						StatementHeader: intoto.StatementHeader{
							Type:          "https://in-toto.io/Statement/v1",
							PredicateType: "https://slsa.dev/provenance/v1",
							Subject: []intoto.Subject{
								{
									Name: "0.1.0",
									Digest: common.DigestSet{
										"sha256": "",
									},
								},
							},
						},
						Predicate: struct{}{},
					},
					{
						StatementHeader: intoto.StatementHeader{
							Type:          "https://in-toto.io/Statement/v1",
							PredicateType: "https://slsa.dev/provenance/v1",
							Subject: []intoto.Subject{
								{
									Name: "0.2.0",
									Digest: common.DigestSet{
										"sha256": "",
									},
								},
							},
						},
						Predicate: struct{}{},
					},
				},
				"spdx": {
					{
						StatementHeader: intoto.StatementHeader{
							Type:          "https://in-toto.io/Statement/v1",
							PredicateType: "https://spdx.github.io/spdx-spec/v2.3/",
							Subject: []intoto.Subject{
								{
									Name: "0.2.0",
									Digest: common.DigestSet{
										"sha256": "",
									},
								},
							},
						},
						Predicate: struct{}{},
					},
				},
			},
			predType: "provenance",
			subject:  "0.1.0",
			expectedRes: []intoto.Statement{
				{
					StatementHeader: intoto.StatementHeader{
						Type:          "https://in-toto.io/Statement/v1",
						PredicateType: "https://slsa.dev/provenance/v1",
						Subject: []intoto.Subject{
							{
								Name: "0.1.0",
								Digest: common.DigestSet{
									"sha256": "",
								},
							},
						},
					},
					Predicate: struct{}{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := GetRelevantStatements(tt.psMap, tt.predType, tt.subject)
			assert.Equal(t, tt.expectedRes, res)
		})
	}
}
