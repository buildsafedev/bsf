package attestation

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

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
	catCmd.Flags().StringVarP(&fileName, "filePath", "f", "", "path to the JSONL file")
	catCmd.Flags().StringVarP(&predicateType, "predicate-type", "p", "", "type of the predicate")
	catCmd.Flags().StringVarP(&subject, "subject", "s", "", "subject of the predicate")
}

var validPredArgs = []string{
	"provenance",
	"vulnerability",
	"vsa",
	"test-result",
	"spdx",
	"scai",
	"runtime-trace",
	"release",
	"link",
	"cdx",
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

		if !slices.Contains(validPredArgs, predicateType) {
			fmt.Print(styles.HintStyle.Render("Hint: validate predicate types:", strings.Join(validPredArgs, ", ")))
			os.Exit(1)
		}

		isValidJSONL, _, err := validateFile(fileName, "JSON")
		if !isValidJSONL {
			fmt.Println(styles.ErrorStyle.Render("error parsing JSONL:", err.Error()))
			os.Exit(1)
		}

		isValidInToto, psMap, err := validateFile(fileName, "inToto")
		if !isValidInToto {
			fmt.Println(styles.ErrorStyle.Render("error validating intoto attestation:", err.Error()))
			os.Exit(1)
		}

		preds := attestation.GetPredicate(psMap, predicateType, subject)
		for _, pred := range preds {
			jsonData, err := json.MarshalIndent(pred, "", "  ")
			if err != nil {
				fmt.Println(styles.ErrorStyle.Render("error marshalling predicate to JSON:", err.Error()))
				os.Exit(1)
			}
			fmt.Println(string(jsonData))

		}
	},
}
