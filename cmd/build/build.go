package build

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/awalterschulze/gographviz"
	"github.com/bom-squad/protobom/pkg/formats"
	"github.com/bom-squad/protobom/pkg/sbom"
	"github.com/spf13/cobra"

	binit "github.com/buildsafedev/bsf/cmd/init"
	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/generate"
	bgit "github.com/buildsafedev/bsf/pkg/git"
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
		sc, fh, err := binit.GetBSFInitializers()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		err = generate.Generate(fh, sc)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		if output == "" {
			output = "bsf-result"
		}

		err = bgit.Add("bsf/")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		err = bgit.Ignore(output + "/")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		err = nixcmd.Build(output + "/result")
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

		err = GenerateArtifcats(output)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
			os.Exit(1)
		}

		fmt.Println(styles.SucessStyle.Render(fmt.Sprintf("Build completed successfully, please check the %s directory", output)))

	},
}

// GenerateSBOM generates the Software Bill of Materials (SBOM)
func GenerateSBOM(w io.Writer, output string, lockFile *hcl2nix.LockFile, appDetails *nixcmd.App, graph *gographviz.Graph) error {
	appNode := &sbom.Node{
		Id:             bsbom.PurlFromNameVersion(appDetails.Name, appDetails.Version),
		PrimaryPurpose: []sbom.Purpose{sbom.Purpose_APPLICATION},
		Name:           appDetails.Name,
		Hashes: map[int32]string{
			int32(sbom.HashAlgorithm_SHA256): appDetails.BinaryHash,
		},
	}

	bom := bsbom.PackageGraphToSBOM(appNode, lockFile, graph)
	bomSt := bsbom.NewStatement(appDetails)

	spdxBom, err := bomSt.ToJSON(bom, formats.SPDX23JSON)
	if err != nil {
		return err
	}
	_, err = w.Write(append(spdxBom, []byte("\n")...))
	if err != nil {
		return err
	}

	cdxBom, err := bomSt.ToJSON(bom, formats.CDX15JSON)
	if err != nil {
		return err
	}
	_, err = w.Write(append(cdxBom, []byte("\n")...))
	if err != nil {
		return err
	}
	return nil
}

// GenerateProvenance generates the provenance
func GenerateProvenance(w io.Writer, output string, appDetails *nixcmd.App, graph *gographviz.Graph) error {
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

	_, err = w.Write(provJ)
	if err != nil {
		return err
	}

	return nil
}

// GenerateArtifcats generates remaining artifacts after build
func GenerateArtifcats(output string) error {
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

	attestationsPath := filepath.Join(output, "attestations.intoto.jsonl")
	attFile, err := os.Create(attestationsPath)
	if err != nil {
		fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
		os.Exit(1)
	}
	defer attFile.Close()

	err = GenerateSBOM(attFile, output, lockFile, appDetails, graph)
	if err != nil {
		fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
		os.Exit(1)
	}

	err = GenerateProvenance(attFile, output, appDetails, graph)
	if err != nil {
		fmt.Println(styles.ErrorStyle.Render("error:", err.Error()))
		os.Exit(1)
	}

	return nil
}

func isNoFileError(err string) bool {
	return strings.Contains(err, "No such file or directory") || strings.Contains(err, "does not contain a 'bsf/flake.nix' file")
}
