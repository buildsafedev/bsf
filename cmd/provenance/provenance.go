package provenance

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/buildsafedev/bsf/cmd/configure"
	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/clients/search"
	"github.com/buildsafedev/bsf/pkg/generate"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	nixcmd "github.com/buildsafedev/bsf/pkg/nix/cmd"
	"github.com/buildsafedev/bsf/pkg/provenance"
)

var (
	output string
)

func init() {
	ProvenanceCMD.Flags().StringVarP(&output, "output", "o", "", "location of the generated provenance file")
}

// ProvenanceCMD represents the provenance command that generates the provenance for the project
var ProvenanceCMD = &cobra.Command{
	Use:   "provenance",
	Short: "provenance generates the provenance for the project",
	Long: `provenance generates a SLSA compaitable provenance for the project. It should be run after a sucesfull "bsf build".
		Example:
		bsf provenance
		bsf sbom provenance -o provenance.json
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(styles.HighlightStyle.Render("Generating Provenance..."))

		// Generating to make sure we have a lock file to work with.
		fh, err := hcl2nix.NewFileHandlers(true)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}
		defer fh.ModFile.Close()
		defer fh.LockFile.Close()
		defer fh.FlakeFile.Close()
		defer fh.DefFlakeFile.Close()

		conf, err := configure.PreCheckConf()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}
		sc, err := search.NewClientWithAddr(conf.BuildSafeAPI, conf.BuildSafeAPITLS)
		if err != nil {
			os.Exit(1)
		}

		// re-generating to make sure we have the latest data.
		err = generate.Generate(fh, sc)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		// Read the bsf.lock file
		data, err := os.ReadFile("bsf.lock")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		lockFile := &hcl2nix.LockFile{}
		err = json.Unmarshal(data, lockFile)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		drvPath, err := nixcmd.GetDrvPathFromResult("")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		drv, err := provenance.GetDerivation(drvPath)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		appDetails, graph, err := nixcmd.GetRuntimeClosureGraph("")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		provSt := provenance.NewStatement(appDetails)
		err = provSt.FromDerivationClosure(drvPath, drv, graph)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}
		provJ, err := provSt.ToJSON()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		// write json to file
		if output != "" {
			err = os.WriteFile(output, provJ, 0644)
			if err != nil {
				fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
				os.Exit(1)
			}
		} else {
			fmt.Println(string(provJ))
		}

		fmt.Println(styles.SucessStyle.Render("Provenance generated successfully"))

	},
}
