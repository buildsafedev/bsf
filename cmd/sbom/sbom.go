package sbom

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bom-squad/protobom/pkg/formats"
	"github.com/bom-squad/protobom/pkg/sbom"
	"github.com/bom-squad/protobom/pkg/writer"
	"github.com/spf13/cobra"

	"github.com/buildsafedev/bsf/cmd/configure"
	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	nixcmd "github.com/buildsafedev/bsf/pkg/nix/cmd"
	bsbom "github.com/buildsafedev/bsf/pkg/sbom"
)

var (
	output string
)

func init() {
	SBOMCmd.Flags().StringVarP(&output, "output", "o", "", "location of the generated SBOM file")
}

// SBOMCmd represents the sbom command that generates the SBOM for the project
var SBOMCmd = &cobra.Command{
	Use:   "sbom",
	Short: "sbom generates the SBOM for the project",
	Long: `sbom generates the SBOM for the project. It should be run after a sucesfull "bsf build".
		Example:
		bsf sbom spdx
		bsf sbom cdx
		bsf sbom spdx -o spdx.json
		bsf sbom cdx -o cdx.json

		Currently, spdx 2.3 and cdx 1.5 are supported.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println(styles.HintStyle.Render("hint:", "run `bsf sbom spdx or bsf sbom cdx to generate the sbom"))
			os.Exit(1)
		}

		if args[0] != "spdx" && args[0] != "cdx" {
			fmt.Println(styles.ErrorStyle.Render("error:", "Invalid SBOM type. Supported types are: spdx, cdx"))
			fmt.Println(styles.HintStyle.Render("hint:", "Run  bsf sbom spdx  or  bsf sbom cdx"))
			os.Exit(1)
		}

		var format formats.Format
		if args[0] == "spdx" {
			format = formats.SPDX23JSON
		} else {
			format = formats.CDX15JSON
		}

		conf, err := configure.PreCheckConf()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		fmt.Println(styles.HighlightStyle.Render("Generating SBOM..."))

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

		// re-generating to make sure we have the latest data.
		err = nixcmd.Build(conf)
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

		appDetails, graph, err := nixcmd.GetRuntimeClosureGraph()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}
		appNode := &sbom.Node{
			Id:             bsbom.PurlFromNameVersion(appDetails.Name, appDetails.Version),
			PrimaryPurpose: []sbom.Purpose{sbom.Purpose_APPLICATION},
			Name:           appDetails.Name,
			Hashes: map[int32]string{
				int32(sbom.HashAlgorithm_SHA256): appDetails.Hash,
			},
		}
		w := writer.New()
		doc := bsbom.PackageGraphToSBOM(appNode, lockFile, graph)

		if output != "" {
			sbomFile, err := os.Create(output)
			if err != nil {
				fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
				os.Exit(1)
			}
			defer sbomFile.Close()

			err = w.WriteStreamWithOptions(doc, sbomFile, &writer.Options{
				Format: format,
			})
			if err != nil {
				fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
				os.Exit(1)
			}
			fmt.Println(styles.SucessStyle.Render("SBOM generated successfully"))
			return
		}

		err = w.WriteStreamWithOptions(doc, os.Stdout, &writer.Options{
			Format: format,
		})
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		fmt.Println(styles.SucessStyle.Render("SBOM generated successfully"))

	},
}
