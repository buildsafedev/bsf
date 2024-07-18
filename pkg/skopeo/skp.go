package skopeo

import (
    "context"
    "fmt"
    "io"
    "github.com/containers/image/v5/copy"
    "github.com/containers/image/v5/signature"
    "github.com/containers/image/v5/transports/alltransports"
)

// ParseImageName parses the image reference
func ParseImageName(ref string) (alltransports.ImageReference, error) {
    return alltransports.ParseImageName(ref)
}

// LoadImageToDocker loads the image to the Docker daemon
func LoadImageToDocker(dir, imageName string, out io.Writer) error {
    ctx := context.Background()

    // Parsed the source and destination image references
    srcRef, err := ParseImageName("dir:" + dir)
    if err != nil {
        return fmt.Errorf("parsing source image reference: %w", err)
    }
    destRef, err := ParseImageName("docker-daemon:" + imageName)
    if err != nil {
        return fmt.Errorf("parsing destination image reference: %w", err)
    }

    // Created a policy context
    policyContext, err := signature.NewPolicyContext(&signature.Policy{})
    if err != nil {
        return fmt.Errorf("creating policy context: %w", err)
    }
    defer policyContext.Destroy()

    // Copied the image
    _, err = copy.Image(ctx, policyContext, destRef, srcRef, &copy.Options{
        ReportWriter: out,
    })
    if err != nil {
        return fmt.Errorf("copying image to Docker: %w", err)
    }

    return nil
}

// LoadImageToPodman loads the image to Podman storage
func LoadImageToPodman(dir, imageName string, out io.Writer) error {
    ctx := context.Background()

    srcRef, err := ParseImageName("dir:" + dir)
    if err != nil {
        return fmt.Errorf("parsing source image reference: %w", err)
    }
    destRef, err := ParseImageName("containers-storage:" + imageName)
    if err != nil {
        return fmt.Errorf("parsing destination image reference: %w", err)
    }

    policyContext, err := signature.NewPolicyContext(&signature.Policy{})
    if err != nil {
        return fmt.Errorf("creating policy context: %w", err)
    }
    defer policyContext.Destroy()

    _, err = copy.Image(ctx, policyContext, destRef, srcRef, &copy.Options{
        ReportWriter: out,
    })
    if err != nil {
        return fmt.Errorf("copying image to Podman: %w", err)
    }

    return nil
}

// PushImage pushes an image to a Docker registry
func PushImage(dir, imageName string, out io.Writer) error {
    ctx := context.Background()

    srcRef, err := ParseImageName("dir:" + dir)
    if err != nil {
        return fmt.Errorf("parsing source image reference: %w", err)
    }
    destRef, err := ParseImageName("docker://" + imageName)
    if err != nil {
        return fmt.Errorf("parsing destination image reference: %w", err)
    }

    policyContext, err := signature.NewPolicyContext(&signature.Policy{})
    if err != nil {
        return fmt.Errorf("creating policy context: %w", err)
    }
    defer policyContext.Destroy()

    _, err = copy.Image(ctx, policyContext, destRef, srcRef, &copy.Options{
        ReportWriter: out,
    })
    if err != nil {
        return fmt.Errorf("pushing image to Docker registry: %w", err)
    }

    return nil
}
