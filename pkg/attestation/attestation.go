package attestation

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	intoto "github.com/in-toto/in-toto-golang/in_toto"
)

// PredicateURIType is a map of predicate URIs to the type of predicate they are
var PredicateURIType = map[string]string{
	"https://slsa.dev/provenance/":                         "provenance",
	"https://in-toto.io/attestation/vulns":                 "vuln",
	"https://slsa.dev/verification_summary/v1":             "vsa",
	"https://in-toto.io/attestation/test-result/":          "test-result",
	"https://spdx.dev/Document":                            "spdx",
	"https://spdx.github.io/spdx-spec":                     "spdx",
	"https://in-toto.io/attestation/scai/attribute-report": "scai",
	"https://in-toto.io/attestation/runtime-trace/":        "runtime-trace",
	"https://in-toto.io/attestation/release":               "release",
	"https://in-toto.io/attestation/link":                  "link",
	"https://cyclonedx.org/bom":                            "cdx",
	"https://cyclonedx.org/specification/overview/":        "cdx",
}

// ValidateInTotoStatement validates the in-toto statement in the byte array
func ValidateInTotoStatement(file []byte) (map[string][]intoto.Statement, error) {
	var predStatementMap = make(map[string][]intoto.Statement)

	scanner := bufio.NewScanner(bytes.NewReader(file))
	for scanner.Scan() {
		line := scanner.Text()
		var statement intoto.Statement
		if err := json.Unmarshal([]byte(line), &statement); err != nil {
			return nil, fmt.Errorf("invalid JSON: %v", err)
		}

		gotPredType, err := GetPredicateType(statement.StatementHeader)
		if err != nil {
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
		predStatementMap[gotPredType] = append(predStatementMap[gotPredType], statement)

	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return predStatementMap, nil
}

// GetPredicateType returns the predicate type for the given statement
func GetPredicateType(statement intoto.StatementHeader) (string, error) {

	if !strings.Contains("https://in-toto.io/Statement/v1", statement.Type) {
		return "", fmt.Errorf("invalid _type: %s", statement.Type)
	}

	if statement.PredicateType == "" {
		return "", fmt.Errorf("predicateType is empty")
	}

	for pred, shortName := range PredicateURIType {
		if strings.Contains(statement.PredicateType, pred) {
			return shortName, nil
		}
	}

	return "", fmt.Errorf("predicateType %s is invalid", statement.PredicateType)
}

// GetRelevantStatements returns the predicate for the given predicate type and subject
func GetRelevantStatements(psMap map[string][]intoto.Statement, predType string, subject string) []intoto.Statement {
	// Filter out the allSts based on the predicate type
	allSts := psMap[predType]

	if subject == "" {
		return allSts
	}

	subSts := make([]intoto.Statement, 0, len(allSts))
	for _, stmt := range allSts {
		for _, subj := range stmt.Subject {
			if subj.Name == subject {
				subSts = append(subSts, stmt)
			}
		}
	}

	return subSts
}
