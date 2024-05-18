package attestation

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// InTotoStatement represents the structure of an in-toto attestation
type InTotoStatement struct {
	Type          string      `json:"_type"`
	PredicateType string      `json:"predicateType"`
	Subject       []Subject   `json:"subject"`
	Predicate     interface{} `json:"Predicate"`
}

type Subject struct {
	Name   string            `json:"name"`
	Digest map[string]string `json:"digest"`
}

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

func IsInToto(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if err := validateInTotoStatement(line); err != nil {
			return false, err
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return true, nil
}

func validateInTotoStatement(line string) error {
	var statement InTotoStatement
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
		if len(subject.Digest) == 0 {
			return fmt.Errorf("subject digest is empty")
		}
		for algo, digest := range subject.Digest {
			if strings.TrimSpace(algo) == "" || strings.TrimSpace(digest) == "" {
				return fmt.Errorf("subject digest has invalid algorithm or value")
			}
		}
	}

	return nil
}

func validatePredicateType(statement InTotoStatement) error {
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
