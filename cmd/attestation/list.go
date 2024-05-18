package attestation

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/attestation"
	"github.com/spf13/cobra"
)

var (
	filePath string
)

func init() {
	listCmd.Flags().StringVarP(&filePath, "filePath", "f", "", "path to the JSONL file")
}

// AttCmd represents the attestation command
var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "lists predicate types",
	Long: `validates if JSONL is valid and if all JSON blocks
	are intoto attestation and lists all the predicate types 
	available.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if filePath == "" {
			fmt.Println(styles.ErrorStyle.Render("Usage: bsf att ls -f <path.to.JSONL_file>"))
			os.Exit(1)
		}

		isValidJSONL, err := isJSONL(filePath)
		if !isValidJSONL {
			fmt.Println(styles.ErrorStyle.Render("error parsing JSONL:", err.Error()))
			os.Exit(1)
		}
		fmt.Println(styles.SucessStyle.Render("✅ JSONL is valid"))

		isValidInToto, err := attestation.IsInToto(filePath)
		if !isValidInToto {
			fmt.Println(styles.ErrorStyle.Render("error validating intoto attestation:", err.Error()))
			os.Exit(1)
		}
		fmt.Println(styles.SucessStyle.Render("✅ intoto attestations are valid"))

		fmt.Println(styles.TextStyle.Render("List of available predicates:"))
		for _, predicateType := range attestation.PredicateTypes {
			fmt.Println(styles.TextStyle.Render(predicateType))
		}
	},
}

func isJSONL(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read the file line by line and validate each line as JSON
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var data interface{}
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			return false, err
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return true, nil
}
