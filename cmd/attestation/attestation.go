package attestation

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	filePath string
)

// InTotoStatement represents the structure of an in-toto attestation
type InTotoStatement struct {
	Type          string      `json:"type"`
	PredicateType string      `json:"predicateType"`
	Subject       []Subject   `json:"subject"`
	Predicate     interface{} `json:"Predicate"`
}

type Subject struct {
	Name   string            `json:"name"`
	Digest map[string]string `json:"digest"`
}

func init() {
	AttCmd.Flags().StringVarP(&filePath, "filePath", "f", "", "path to the JSONL file")
}

// AttCmd represents the attestation command
var AttCmd = &cobra.Command{
	Use:   "att",
	Short: "",
	Long: `
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if filePath == "" {
			fmt.Println("Usage: bsf att -f <path.to.JSONL_file>")
			return
		}

		if !isJSONL(filePath) {
			fmt.Println("Invalid JSONL format.")
			return
		}

		if !isValidInToto(filePath) {
			fmt.Println("Invalid intoto attestation.")
			return
		}

		fmt.Println("Valid")
	},
}

func isValidInToto(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Text()
		if err := validateInTotoStatement(line); err != nil {
			fmt.Printf("Line %d is not a valid in-toto attestation: %v\n", lineNumber, err)
		} else {
			fmt.Printf("Line %d is a valid in-toto attestation\n", lineNumber)
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return true
}

func isJSONL(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return false
	}
	defer file.Close()

	// Read the file line by line and validate each line as JSON
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var data interface{}
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			fmt.Println("Error parsing JSON:", err)
			return false
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return false
	}

	return true
}

func validateInTotoStatement(line string) error {
	var statement InTotoStatement
	if err := json.Unmarshal([]byte(line), &statement); err != nil {
		return fmt.Errorf("invalid JSON: %v", err)
	}

	if statement.Type != "https://in-toto.io/Statement/v1" {
		return fmt.Errorf("invalid _type: %s", statement.Type)
	}

	if statement.PredicateType == "" {
		return fmt.Errorf("predicateType is empty")
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
