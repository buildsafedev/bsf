package release

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/buildsafedev/bsf/cmd/styles"
	"github.com/buildsafedev/bsf/pkg/hcl2nix"
	"github.com/buildsafedev/bsf/pkg/platformutils"
	"github.com/buildsafedev/bsf/pkg/release"
)

var (
	platform string
)

func init() {
	ReleaseCmd.Flags().StringVarP(&platform, "platform", "p", "", "platform to release to")
}

// ReleaseCmd represents the release command
var ReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "release the artifacts ",
	Long: `
	
	GITHUB_ACCESS_TOKEN=xxx bsf release app v0.1
	GITHUB_ACCESS_TOKEN=xxx bsf release app v0.1 --platform linux/amd64
	`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(styles.TextStyle.Render("Validating configuration..." + args[0]))

		if len(args) < 2 {
			fmt.Println(styles.ErrorStyle.Render("error: ", "missing arguments"))
			fmt.Println(styles.HintStyle.Render("hint:", "run `bsf release app v0.1` to release the artifact"))
			os.Exit(1)
		}

		accessToken := os.Getenv("GITHUB_ACCESS_TOKEN")
		if accessToken == "" {
			fmt.Println(styles.ErrorStyle.Render("error: ", "GITHUB_ACCESS_TOKEN is not set"))
			fmt.Println(styles.HintStyle.Render("hint:", "run `export GITHUB_ACCESS_TOKEN=xxx` to set the access token"))
			os.Exit(1)
		}

		conf, err := hcl2nix.ReadHclFile("bsf.hcl")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", "failed to read config"))
			fmt.Println(styles.HintStyle.Render("hint:", "run `bsf init` "))
			os.Exit(1)
		}

		ghHCL, err := hcl2nix.ReadGitHubReleaseParams(conf, args[0])
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			fmt.Println(styles.HintStyle.Render("hint:", "run `bsf init` "))
			os.Exit(1)
		}

		version := args[1]

		var plat string
		if platform == "" {
			tos, tarch := platformutils.FindPlatform(platform)
			platform = tos + "/" + tarch
		}

		switch platformutils.DetermineFormat(platform) {
		case platformutils.ArchOSFormat:
			plat = platformutils.ArchOSToOSArch(platform)

		case platformutils.OSArchFormat:
			plat = platform

		case platformutils.UnknownFormat:
			fmt.Println(styles.ErrorStyle.Render("error: ", "unknown platform format"))
			os.Exit(1)
		}

		fmt.Println(styles.TextStyle.Render("Creating release for " + args[0] + " version " + version))
		releaseParams := release.GHParams{
			Owner:       ghHCL.Owner,
			Repo:        ghHCL.Repo,
			Version:     version,
			Platform:    plat,
			AccessToken: accessToken,
			Dir:         ghHCL.Dir + "/result",
		}

		ghr := release.NewGHRelease(releaseParams)
		err = ghr.GHReleaseCreate(releaseParams)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}
		fmt.Println(styles.TextStyle.Render("Release created for " + args[0] + " version " + version))

		fmt.Println(styles.TextStyle.Render("Preparing artifact... "))
		tmpDir, err := os.MkdirTemp("/tmp", "bsf-release")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		archiveName := ghHCL.App + "-" + strings.ReplaceAll(plat, "/", "-") + ".tar.gz"
		archivePath := tmpDir + "/" + archiveName
		artifactPath, err := filepath.Abs(ghHCL.Dir + "result")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		err = release.CreateTarGzFromSymlink(artifactPath, archivePath)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		fmt.Println(styles.TextStyle.Render("Uploading artifact... "))

		err = ghr.UploadFileToRelease(context.Background(), archiveName, archivePath)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}
		fmt.Println(styles.TextStyle.Render("Artifact uploaded"))

		fmt.Println(styles.TextStyle.Render("Uploading attestations... "))

		attPath, err := filepath.Abs(ghHCL.Dir + "attestations.intoto.jsonl")
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		err = ghr.UploadFileToRelease(context.Background(), ghHCL.App+"-"+strings.ReplaceAll(plat, "/", "-")+"-attestations.intoto.jsonl", attPath)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

		fmt.Println(styles.TextStyle.Render("Attestations uploaded"))
		fmt.Println(styles.SucessStyle.Render("Release created for " + args[0] + " version " + version + " for " + plat))

		err = ghr.Flake(tmpDir, archiveName)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render("error: ", err.Error()))
			os.Exit(1)
		}

	},
}
