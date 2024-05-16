package attestation

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	filePath string
)

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

		fmt.Println("Valid")
	},
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
