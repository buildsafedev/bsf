package attestation

import (
	"fmt"
	"os"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/attestation"
	"github.com/buildsafedev/bsf/pkg/jsonl"
	tea "github.com/charmbracelet/bubbletea"
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

		isValidJSONL, err := validateFile(filePath, jsonl.ValidateIsJSONL)
		if !isValidJSONL {
			fmt.Println(styles.ErrorStyle.Render("error parsing JSONL:", err.Error()))
			os.Exit(1)
		}
		fmt.Println(styles.SucessStyle.Render("✅ JSONL is valid"))

		isValidInToto, err := validateFile(filePath, attestation.ValidateInTotoStatement)
		if !isValidInToto {
			fmt.Println(styles.ErrorStyle.Render("error validating intoto attestation:", err.Error()))
			os.Exit(1)
		}
		fmt.Println(styles.SucessStyle.Render("✅ intoto attestations are valid"))

		// Get the predicate-subject map
		psMap := attestation.GetPredicateSubjectMap()

		m := initPredSubTable(psMap)
		if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
			fmt.Println(styles.ErrorStyle.Render(fmt.Errorf("error: %v", err).Error()))
			os.Exit(1)
		}
	},
}

func validateFile(filePath string, validateFunc func([]byte) error) (bool, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}

	if err := validateFunc([]byte(file)); err != nil {
		return false, err
	}

	return true, nil
}
