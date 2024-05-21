package attestation

import (
	"fmt"
	"os"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/attestation"
	"github.com/buildsafedev/bsf/pkg/jsonl"
	intoto "github.com/in-toto/in-toto-golang/in_toto"
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
			fmt.Println(styles.HintStyle.Render("hint: bsf att ls -f <path.to.JSONL_file>"))
			os.Exit(1)
		}

		isValidJSONL, _, err := validateFile(filePath, "JSON")
		if !isValidJSONL {
			fmt.Println(styles.ErrorStyle.Render("error parsing JSONL:", err.Error()))
			os.Exit(1)
		}
		fmt.Println(styles.SucessStyle.Render("✅ JSONL is valid"))

		isValidInToto, psMap, err := validateFile(filePath, "inToto")
		if !isValidInToto {
			fmt.Println(styles.ErrorStyle.Render("error validating intoto attestation:", err.Error()))
			os.Exit(1)
		}
		fmt.Println(styles.SucessStyle.Render("✅ intoto attestations are valid"))

		// Print the predicate-subject map
		printPredSubjTable(psMap)
	},
}

func validateFile(filePath string, fileType string) (bool, map[string][]intoto.Statement, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return false, nil, err
	}
	// Check the fileType
	switch fileType {
	case "JSON":
		if err := jsonl.ValidateIsJSONL([]byte(file)); err != nil {
			return false, nil, err
		} else {
			return true, nil, nil
		}
	case "inToto":
		if psMap, err := attestation.ValidateInTotoStatement([]byte(file)); err != nil {
			return false, nil, err
		} else {
			return true, psMap, nil
		}
	}
	return false, nil, nil
}
