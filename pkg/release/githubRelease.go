package release

import (
	"context"
	"os"
	"strings"

	"github.com/google/go-github/v62/github"
	"golang.org/x/oauth2"

	"github.com/buildsafedev/bsf/pkg/crypto"
	nixtmpl "github.com/buildsafedev/bsf/pkg/nix/template"
	"github.com/buildsafedev/bsf/pkg/platformutils"
)

// GHParams holds release parameters for GitHub release
type GHParams struct {
	Repo        string
	Owner       string
	Version     string
	Platform    string
	Dir         string
	AccessToken string
}

// GHRelease holds the Github release
type GHRelease struct {
	client    *github.Client
	params    GHParams
	releaseID int64
}

// NewGHRelease creates a new GHRelease
func NewGHRelease(params GHParams) *GHRelease {
	return &GHRelease{
		client: createGHClient(params.AccessToken),
		params: params,
	}
}

// GHReleaseCreate uploads the artifacts to Github
func (gh *GHRelease) GHReleaseCreate(params GHParams) error {

	ghRelease, err := gh.getReleaseByTag(context.Background())
	if err != nil {
		return err
	}
	if ghRelease == nil {
		_, err = gh.createRelease(context.Background())
		if err != nil {
			return err
		}
	}

	return nil
}

func createGHClient(accessToken string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

func (gh *GHRelease) getReleaseByTag(ctx context.Context) (*github.RepositoryRelease, error) {
	release, _, err := gh.client.Repositories.GetReleaseByTag(ctx, gh.params.Owner, gh.params.Repo, gh.params.Version)
	if err != nil {
		if _, ok := err.(*github.ErrorResponse); ok {
			// Release does not exist
			return nil, nil
		}
		return nil, err
	}
	return release, nil
}

func (gh *GHRelease) createRelease(ctx context.Context) (*github.RepositoryRelease, error) {
	release := &github.RepositoryRelease{
		TagName: github.String(gh.params.Version),
		Name:    github.String(gh.params.Version),
	}
	release, _, err := gh.client.Repositories.CreateRelease(ctx, gh.params.Owner, gh.params.Repo, release)
	if err != nil {
		return nil, err
	}
	gh.releaseID = *release.ID
	return release, nil
}

// UploadFileToRelease uploads a file to a GitHub release
func (gh *GHRelease) UploadFileToRelease(ctx context.Context, name, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, _, err = gh.client.Repositories.UploadReleaseAsset(ctx, gh.params.Owner, gh.params.Repo, gh.releaseID, &github.UploadOptions{Name: name}, file)
	return err
}

// Flake generates a flake.nix file for the release
func (gh *GHRelease) Flake(tmpDir, artifactName string) error {
	archiveHash, err := crypto.FileSHA256(tmpDir + "/" + artifactName)
	if err != nil {
		return err
	}

	nixHash, err := crypto.HexToBase64(archiveHash)
	if err != nil {
		return err
	}

	params := nixtmpl.RemoteFile{
		Name:    gh.params.Repo,
		Version: gh.params.Version,
		PlatformURLs: map[string]string{
			platformutils.OSArchToArchOS(gh.params.Platform): gh.getArtifactURL(artifactName),
		},
		PlatformHashes: map[string]string{
			platformutils.OSArchToArchOS(gh.params.Platform): "sha256-" + nixHash,
		},
	}

	fh, err := os.Create(tmpDir + "/flake.nix")
	if err != nil {
		return err
	}
	defer fh.Close()

	err = nixtmpl.GenerateRemoteFlake(params, fh)
	if err != nil {
		return err
	}

	flakeName := strings.ReplaceAll(gh.params.Platform, "/", "-") + "-flake.nix"

	err = gh.UploadFileToRelease(context.Background(), flakeName, fh.Name())
	if err != nil {
		return err
	}

	return nil
}

func (gh *GHRelease) getArtifactURL(fileName string) string {

	return "https://github.com/" + gh.params.Owner + "/" + gh.params.Repo + "/releases/download/" + gh.params.Version + "/" + fileName
}
