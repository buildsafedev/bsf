package attestation

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/buildsafedev/bsf/pkg/provenance"
)

var PredicateTypes = []string{
	"SLSA Provenance",
	"Link",
	"SCAI Report",
	"Runtime Traces",
	"SLSA Verification Summary",
	"SPDX",
	"CycloneDX",
	"Vulnerability",
	"Release",
	"Test Result",
}

var predicateURIs = []string{
	"https://slsa.dev/provenance/",
	"https://in-toto.io/attestation/link/v0.3",
	"https://in-toto.io/attestation/scai/attribute-report",
	"https://in-toto.io/attestation/runtime-trace/v0.1",
	"https://slsa.dev/verification_summary/v1",
	"https://spdx.github.io/spdx-spec/v2.3/",
	"https://cyclonedx.org/specification/overview/",
	"https://in-toto.io/attestation/vulns",
	"https://in-toto.io/attestation/release",
	"https://in-toto.io/attestation/test-result/v0.1",
}

func ValidateInTotoStatement(line []byte) error {
	var statement provenance.Statement
	if err := json.Unmarshal(line, &statement); err != nil {
		return fmt.Errorf("invalid JSON: %v", err)
	}

	if statement.Type != "https://in-toto.io/Statement/v1" {
		return fmt.Errorf("invalid _type: %s", statement.Type)
	}

	if err := validatePredicateType(statement); err != nil {
		return err
	}

	if len(statement.Subject) == 0 {
		return fmt.Errorf("subject is empty")
	}

	for _, subject := range statement.Subject {
		if subject.Name == "" {
			return fmt.Errorf("subject name is empty")
		}
	}

	return nil
}

func validatePredicateType(statement provenance.Statement) error {
	if statement.PredicateType == "" {
		return fmt.Errorf("predicateType is empty")
	}

	for _, uri := range predicateURIs {
		if strings.Contains(statement.PredicateType, uri) {
			return nil
		}
	}

	return fmt.Errorf("predicateType %s is invalid", statement.PredicateType)
}
