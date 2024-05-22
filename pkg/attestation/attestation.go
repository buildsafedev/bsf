package attestation

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	intoto "github.com/in-toto/in-toto-golang/in_toto"
)

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
