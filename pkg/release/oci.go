package release

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	ociv1 "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/credentials"
	"oras.land/oras-go/v2/registry/remote/retry"
)

// OCIParams holds release parameters
type OCIParams struct {
	Name     string
	Version  string
	Platform string
	Dir      string
	ImageLoc string
}

// OciRelease releases the artifacts to OCI
func OciRelease(params OCIParams) error {

	files, err := WalkDir(params.Dir)
	if err != nil {
		return err
	}

	fs, err := file.New("/tmp/")
	if err != nil {
		return err
	}

	defer fs.Close()

	fds, err := createFileDescriptors(fs, params.Dir, files)
	if err != nil {
		return err
	}

	manifest, err := createManifest(fs, fds)
	if err != nil {
		return err
	}

	if err = fs.Tag(context.Background(), manifest, params.Version); err != nil {
		return err
	}

	err = ociPushArtifact(fs, params.Version, params.ImageLoc)
	if err != nil {
		return err
	}

	// ociCreateFlake()

	return nil
}

func createManifest(fs *file.Store, fds []ociv1.Descriptor) (ociv1.Descriptor, error) {
	artifactType := "application/vnd.buildsafe.dev.config"
	opts := oras.PackManifestOptions{
		Layers: fds,
	}
	manifestDescriptor, err := oras.PackManifest(context.Background(), fs, oras.PackManifestVersion1_1, artifactType, opts)
	if err != nil {
		return ociv1.Descriptor{}, err
	}
	return manifestDescriptor, nil

}

// ociPushArtifact pushes the artifact to the OCI registry
func ociPushArtifact(fs *file.Store, tag string, imageLoc string) error {
	repo, err := remote.NewRepository(imageLoc)
	if err != nil {
		return err
	}

	storeOpts := credentials.StoreOptions{}
	credStore, err := credentials.NewStoreFromDocker(storeOpts)
	if err != nil {
		return err
	}
	repo.Client = &auth.Client{
		Client:     retry.DefaultClient,
		Cache:      auth.NewCache(),
		Credential: credentials.Credential(credStore),
	}

	_, err = oras.Copy(context.Background(), fs, tag, repo, tag, oras.DefaultCopyOptions)
	if err != nil {
		return err
	}
	return nil
}

func createFileDescriptors(fs *file.Store, dir string, files []string) ([]ociv1.Descriptor, error) {
	mediaType := "application/octet-stream"
	fileDescriptors := make([]ociv1.Descriptor, 0, len(files))
	for _, name := range files {
		if _, err := os.Stat(name); os.IsNotExist(err) {
			continue // Skip files that don't exist
		}
		// absPath is the absolute path to the file from /
		absPath, err := filepath.Abs(name)
		if err != nil {
			return nil, err
		}

		fileDescriptor, err := fs.Add(context.Background(), strings.TrimPrefix(name, dir), mediaType, absPath)
		if err != nil {
			return nil, err
		}
		fileDescriptors = append(fileDescriptors, fileDescriptor)
	}

	return fileDescriptors, nil
}

func ociCreateFlake(fds []ociv1.Descriptor, params OCIParams) error {
	// temp file
	ff, err := os.CreateTemp("/tmp", "flake.nix")
	if err != nil {
		return err
	}

	defer ff.Close()

	// TODO: write flake to temp file

	return nil
}
