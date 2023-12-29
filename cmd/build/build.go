package build

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/buildsafedev/bsf/cmd/styles"
	nixcmd "github.com/buildsafedev/bsf/pkg/nix/cmd"
)

// BuildCmd represents the build command
var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "builds the project",
	Long: `builds the project based on instructions defined in bsf.hcl.
	Build occurs in a sandboxed environment where only current directory is available. 
	It is recommended to check in the files before building.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		stdOut, err := nixcmd.Build()
		if err != nil {
			gotHash := isHashMismatchError(err.Error())
			if gotHash == "" {
				fmt.Println(styles.ErrorStyle.Render("error: %v", err.Error()))
			} else {
				fmt.Println(fmt.Sprintf(styles.ErrorStyle.Render("Hash mismatch detected. Please insert the following hash in the build app/module section of bsf.hcl : ", gotHash)))
			}
			os.Exit(1)
		}

		fmt.Println(styles.SucessStyle.Render(stdOut))

	},
}

func isHashMismatchError(err string) string {
	if strings.Contains(err, "hash mismatch") {
		re := regexp.MustCompile(`got:\s+(sha256-.*)`)
		matches := re.FindStringSubmatch(err)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}
