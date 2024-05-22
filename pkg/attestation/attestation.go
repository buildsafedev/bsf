package attestation

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	intoto "github.com/in-toto/in-toto-golang/in_toto"
)

var ValidPreds = map[string][]string{
	"https://slsa.dev/provenance/":                         {"SLSA Provenance", "Provenance"},
	"https://in-toto.io/attestation/vulns":                 {"Vulnerability", "vuln"},
	"https://slsa.dev/verification_summary/v1":             {"SLSA Verification Summary", "vsa"},
	"https://in-toto.io/attestation/test-result/v0.1":      {"Test Result", "test-result"},
	"https://spdx.dev/Document":                            {"SPDX"},
	"https://spdx.github.io/spdx-spec":                     {"SPDX"},
	"https://in-toto.io/attestation/scai/attribute-report": {"SCAI Report", "scai"},
	"https://in-toto.io/attestation/runtime-trace/":        {"Runtime Traces", "runtime-trace"},
	"https://in-toto.io/attestation/release":               {"Release"},
	"https://in-toto.io/attestation/link":                  {"Link"},
	"https://cyclonedx.org/bom":                            {"CycloneDX", "CDX"},
	"https://cyclonedx.org/specification/overview/":        {"CycloneDX", "CDX"},
}

var predicateURIs = []string{
	"https://slsa.dev/provenance/",
	"https://in-toto.io/attestation/",
	"https://slsa.dev/verification_summary/v1",
	"https://spdx.github.io/spdx-spec/v2.3/",
	"https://spdx.dev/Document",
	"https://cyclonedx.org/specification/overview/",
	"https://cyclonedx.org/bom",
}

func ValidateInTotoStatement(file []byte) (map[string][]intoto.Statement, error) {
	var predStatementMap = make(map[string][]intoto.Statement)

	scanner := bufio.NewScanner(bytes.NewReader(file))
	for scanner.Scan() {
		line := scanner.Text()
		var statement intoto.Statement
		if err := json.Unmarshal([]byte(line), &statement); err != nil {
			return nil, fmt.Errorf("invalid JSON: %v", err)
		}

		if err := validatePredicateType(statement.StatementHeader); err != nil {
			return nil, err
		}

		if len(statement.Subject) == 0 {
			return nil, fmt.Errorf("subject is empty")
		}

		for _, subject := range statement.Subject {
			if subject.Name == "" {
				return nil, fmt.Errorf("subject name is empty")
			}
		}
		predStatementMap[statement.PredicateType] = append(predStatementMap[statement.PredicateType], statement)

	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return predStatementMap, nil
}

func validatePredicateType(statement intoto.StatementHeader) error {

	if !strings.Contains("https://in-toto.io/Statement/v1", statement.Type) {
		return fmt.Errorf("invalid _type: %s", statement.Type)
	}

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

func GetPredicate(psMap map[string][]intoto.Statement, predtype string, subject string) ([]intoto.Statement) {
	var keysToFind []string
	var predsFound []intoto.Statement

	// Find the key in ValidPreds based on the predtype
	for key, values := range ValidPreds {
		for _, value := range values {
			if strings.EqualFold(value, predtype) {
				keysToFind = append(keysToFind, key)
				break
			}
		}
	}

	// Loop over the psMap and find the matching key
	for key, statements := range psMap {
		for _, keyToFind := range keysToFind {
			if strings.Contains(key, keyToFind) {
				for _, statement := range statements {
					if subject != "" {
						for _, predSub := range statement.StatementHeader.Subject {
							if predSub.Name == subject {
								predsFound = append(predsFound, statement)
							}
						}
					} else {
						predsFound = append(predsFound, statement)
					}
				}
			}
		}
	}
	return predsFound
}
