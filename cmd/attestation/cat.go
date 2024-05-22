package attestation

import (
	"fmt"
	"os"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/attestation"
	"github.com/spf13/cobra"
)

var (
	fileName      string
	predicateType string
	subject       string
)

func init() {
	listCmd.Flags().StringVarP(&fileName, "filePath", "f", "", "path to the JSONL file")
	listCmd.Flags().StringVarP(&predicateType, "predicate-type", "pred", "", "type of the predicate")
	listCmd.Flags().StringVarP(&subject, "subject", "sub", "", "subject of the predicate")
}

// AttCmd represents the attestation command
var catCmd = &cobra.Command{
	Use:   "cat",
	Short: "prints out the predicate type in JSON",
	Long: `
	bsf att cat -f <path-to-file> --predicate-type <predicate-type>
	bsf att cat -f <path-to-file> --predicate-type <predicate-type> --subject <subject>
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if fileName == "" || predicateType == "" {
			fmt.Println(styles.HintStyle.Render("hint: bsf att cat -f <path-to-file> --predicate-type <predicate-type"))
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

		pred := attestation.GetPredicate(psMap, predicateType, subject)
		fmt.Println(pred)
	},
}
