package attestation

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	intoto "github.com/in-toto/in-toto-golang/in_toto"
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
	"https://in-toto.io/attestation/",
	"https://slsa.dev/verification_summary/v1",
	"https://spdx.github.io/spdx-spec/v2.3/",
	"https://cyclonedx.org/specification/overview/",
}

func ValidateInTotoStatement(file []byte) error {

	scanner := bufio.NewScanner(bytes.NewReader(file))
	for scanner.Scan() {
		line := scanner.Text()
		var statement intoto.StatementHeader
		if err := json.Unmarshal([]byte(line), &statement); err != nil {
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
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func validatePredicateType(statement intoto.StatementHeader) error {
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
