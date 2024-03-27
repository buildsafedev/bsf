package build

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/bom-squad/protobom/pkg/formats"
	"github.com/bom-squad/protobom/pkg/sbom"
	"github.com/spf13/cobra"

	"github.com/buildsafedev/bsf/cmd/configure"
	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/config"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	nixcmd "github.com/buildsafedev/bsf/pkg/nix/cmd"
	"github.com/buildsafedev/bsf/pkg/provenance"
	bsbom "github.com/buildsafedev/bsf/pkg/sbom"
)

var (
	output string
)

func init() {
	BuildCmd.Flags().StringVarP(&output, "output", "o", "", "location of the build artifacts generated")
}

// BuildCmd represents the build command
var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "builds the project",
	Long: `builds the project based on instructions defined in bsf.hcl.
	Build occurs in a sandboxed environment where only current directory is available. 
	It is recommended to check in the files in version control system(ex: Git) before building.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat("bsf.hcl"); err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: Has the project been initialized?"))
			fmt.Println(styles.HintStyle.Render("hint: ", "run `bsf init` to initialize the project"))
			os.Exit(1)
		}
		fmt.Println(styles.HighlightStyle.Render("Building, please be patient..."))

		conf, err := configure.PreCheckConf()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		if output == "" {
			output = "bsf-result"
		}

		err = nixcmd.Build(conf, output+"/result")
		if err != nil {
			if isNoFileError(err.Error()) {
				fmt.Println(styles.ErrorStyle.Render(err.Error() + "\n Please ensure all necessary files are added/committed in your version control system"))
				fmt.Println(styles.HintStyle.Render("hint: run git add .  "))
				os.Exit(1)
			}
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		fmt.Println(styles.HighlightStyle.Render("Generating artifacts..."))

		err = GenerateArtifcats(conf, output)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		fmt.Println(styles.SucessStyle.Render(fmt.Sprintf("Build completed successfully, please check the %s directory", output)))

	},
}

// GenerateSBOM generates the Software Bill of Materials (SBOM)
func GenerateSBOM(conf *config.Config, output string, lockFile *hcl2nix.LockFile, appDetails *nixcmd.App, graph *gographviz.Graph) error {
	appNode := &sbom.Node{
		Id:             bsbom.PurlFromNameVersion(appDetails.Name, appDetails.Version),
		PrimaryPurpose: []sbom.Purpose{sbom.Purpose_APPLICATION},
		Name:           appDetails.Name,
		Hashes: map[int32]string{
			int32(sbom.HashAlgorithm_SHA256): appDetails.Hash,
		},
	}

	bom := bsbom.PackageGraphToSBOM(appNode, lockFile, graph)
	bomSt := bsbom.NewStatement(appDetails)

	spdxBom, err := bomSt.ToJSON(bom, formats.SPDX23JSON)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(output, "spdx.json"), spdxBom, 0644)
	if err != nil {
		return err
	}

	cdxBom, err := bomSt.ToJSON(bom, formats.CDX15JSON)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(output, "cdx.json"), cdxBom, 0644)
	if err != nil {
		return err
	}

	return nil
}

// GenerateProvenance generates the provenance
func GenerateProvenance(conf *config.Config, output string, appDetails *nixcmd.App, graph *gographviz.Graph) error {
	drvPath, err := nixcmd.GetDrvPathFromResult(output)
	if err != nil {
		return err
	}

	drv, err := provenance.GetDerivation(drvPath)
	if err != nil {
		return err
	}

	provSt := provenance.NewStatement(appDetails)
	err = provSt.FromDerivationClosure(drvPath, drv, graph)
	if err != nil {
		return err
	}
	provJ, err := provSt.ToJSON()
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(output, "provenance.json"), provJ, 0644)
	if err != nil {
		return err
	}

	return nil
}

// GenerateArtifcats generates remaining artifacts after build
func GenerateArtifcats(conf *config.Config, output string) error {
	// Read the bsf.lock file
	lockData, err := os.ReadFile("bsf.lock")
	if err != nil {
		fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
		os.Exit(1)
	}

	lockFile := &hcl2nix.LockFile{}
	err = json.Unmarshal(lockData, lockFile)
	if err != nil {
		fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
		os.Exit(1)
	}

	appDetails, graph, err := nixcmd.GetRuntimeClosureGraph(output)
	if err != nil {
		fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
		os.Exit(1)
	}

	err = GenerateSBOM(conf, output, lockFile, appDetails, graph)
	if err != nil {
		fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
		os.Exit(1)
	}

	err = GenerateProvenance(conf, output, appDetails, graph)
	if err != nil {
		fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
		os.Exit(1)
	}

	return nil
}

func isNoFileError(err string) bool {
	return strings.Contains(err, "No such file or directory")
}
