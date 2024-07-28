package oci

import (
    "context"
    "fmt"
    "os"
	"github.com/BalaadityaPatanjali/image/v5/copy"
    "github.com/BalaadityaPatanjali/image/v5/signature"
    "github.com/BalaadityaPatanjali/image/v5/transports/alltransports"
)

// Handles the image operation based on the provided source and destination types
func handleImageOperation(srcType, destType, dir, imageName string) error {
    ctx := context.Background()

    // Parses the source and destination image references
    srcRef, err := alltransports.ParseImageName(fmt.Sprintf("%s:%s", srcType, dir))
    if err != nil {
        return fmt.Errorf("parsing source image reference: %w", err)
    }
    destRef, err := alltransports.ParseImageName(fmt.Sprintf("%s:%s", destType, imageName))
    if err != nil {
        return fmt.Errorf("parsing destination image reference: %w", err)
    }

    // Creates a policy context
    policyContext, err := signature.NewPolicyContext(&signature.Policy{})
    if err != nil {
        return fmt.Errorf("creating policy context: %w", err)
    }
    defer policyContext.Destroy()

    // Copies the image
    _, err = copy.Image(ctx, policyContext, destRef, srcRef, &copy.Options{
        ReportWriter: os.Stdout,
    })
    if err != nil {
        return fmt.Errorf("copying image: %w", err)
    }

    return nil
}

// LoadDocker loads the image to the Docker daemon
func LoadDocker(daemon, dir, imageName string) error {
    return handleImageOperation("dir", "docker-daemon", dir, imageName)
}

// LoadPodman loads the image to Podman
func LoadPodman(dir, imageName string) error {
    return handleImageOperation("dir", "containers-storage", dir, imageName)
}

// Push pushes the image to a registry
func Push(dir, imageName string) error {
    return handleImageOperation("dir", "docker", dir, imageName)
}
