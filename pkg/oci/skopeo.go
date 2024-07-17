package oci

import (
    "context"
    "fmt"
    "os"
    "github.com/containers/image/v5/copy"
    "github.com/containers/image/v5/signature"
    "github.com/containers/image/v5/transports/alltransports"
)

// LoadDocker loads the image to the docker daemon
func LoadDocker(daemon, dir, imageName string) error {
    ctx := context.Background()

    // Parsed the source and destination image references
    srcRef, err := alltransports.ParseImageName("dir:" + dir)
    if err != nil {
        return fmt.Errorf("parsing source image reference: %w", err)
    }
    destRef, err := alltransports.ParseImageName("docker-daemon:" + imageName)
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
        ReportWriter: os.Stdout,
    })
    if err != nil {
        return fmt.Errorf("copying image: %w", err)
    }

    return nil
}

// LoadPodman loads the image to the podman
func LoadPodman(dir, imageName string) error {
    ctx := context.Background()

    srcRef, err := alltransports.ParseImageName("dir:" + dir)
    if err != nil {
        return fmt.Errorf("parsing source image reference: %w", err)
    }
    destRef, err := alltransports.ParseImageName("containers-storage:" + imageName)
    if err != nil {
        return fmt.Errorf("parsing destination image reference: %w", err)
    }

    policyContext, err := signature.NewPolicyContext(&signature.Policy{})
    if err != nil {
        return fmt.Errorf("creating policy context: %w", err)
    }
    defer policyContext.Destroy()

    _, err = copy.Image(ctx, policyContext, destRef, srcRef, &copy.Options{
        ReportWriter: os.Stdout,
    })
    if err != nil {
        return fmt.Errorf("copying image: %w", err)
    }

    return nil
}

// Push image to registry
func Push(dir, imageName string) error {
    ctx := context.Background()

    srcRef, err := alltransports.ParseImageName("dir:" + dir)
    if err != nil {
        return fmt.Errorf("parsing source image reference: %w", err)
    }
    destRef, err := alltransports.ParseImageName("docker://" + imageName)
    if err != nil {
        return fmt.Errorf("parsing destination image reference: %w", err)
    }

    policyContext, err := signature.NewPolicyContext(&signature.Policy{})
    if err != nil {
        return fmt.Errorf("creating policy context: %w", err)
    }
    defer policyContext.Destroy()

    _, err = copy.Image(ctx, policyContext, destRef, srcRef, &copy.Options{
        ReportWriter: os.Stdout,
    })
    if err != nil {
        return fmt.Errorf("copying image: %w", err)
    }

    return nil
}
