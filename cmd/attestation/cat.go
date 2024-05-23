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
	predicateType string
	predicate     bool
	subject       string
	output        string
)

func init() {
	catCmd.Flags().StringVarP(&predicateType, "predicate-type", "t", "", "type of the predicate")
	catCmd.Flags().StringVarP(&subject, "subject", "s", "", "subject of the predicate")
	catCmd.Flags().StringVarP(&output, "output", "o", "", "name of the output file")
	catCmd.Flags().BoolVarP(&predicate, "predicate", "p", false, "print predicate")
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
	bsf att cat <path-to-file> --predicate-type <predicate-type>
	bsf att cat <path-to-file> --predicate-type <predicate-type> --subject <subject>
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "" || predicateType == "" {
			fmt.Println(styles.HintStyle.Render("hint: bsf att cat <path-to-file> --predicate-type <predicate-type"))
			os.Exit(1)
		}

		if !slices.Contains(validPredArgs, predicateType) {
			fmt.Print(styles.HintStyle.Render("Hint: validate predicate types:", strings.Join(validPredArgs, ", ")))
			os.Exit(1)
		}

		fileName := args[0]

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

		relSts := attestation.GetRelevantStatements(psMap, predicateType, subject)

		if len(relSts) == 0 {
			fmt.Println(styles.ErrorStyle.Render("no relevant statements found"))
			os.Exit(1)
		}

		for _, relSt := range relSts {
			var data interface{}
			if predicate {
				data = relSt.Predicate
			} else {
				data = relSt
			}

			predJSON, err := json.MarshalIndent(data, " ", "  ")
			if err != nil {
				fmt.Println(err)
				continue
			}

			if output != "" {
				if err := os.WriteFile(output, predJSON, 0644); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			} else {
				fmt.Println(string(predJSON))
			}
		}
	},
}
