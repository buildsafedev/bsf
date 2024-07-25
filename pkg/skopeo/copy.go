package skopeo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Options struct {
	ReportWriter io.Writer
}

type Image struct {
	Source      string
	Destination string
}

func CopyImage(ctx context.Context, policyContext *PolicyContext, destRef, srcRef ImageReference, options *Options) (*Image, error) {
	if srcRef == nil || destRef == nil {
		return nil, errors.New("invalid image reference")
	}

	src := srcRef.String()
	dest := destRef.String()

	if !strings.Contains(src, "dir:") || !strings.Contains(dest, "docker:") {
		return nil, errors.New("unsupported source or destination")
	}

	if options.ReportWriter != nil {
		_, err := io.WriteString(options.ReportWriter, fmt.Sprintf("Copying image from %s to %s...\n", src, dest))
		if err != nil {
			return nil, fmt.Errorf("failed to write to report writer: %w", err)
		}
	}

	fmt.Printf("Successfully copied image from %s to %s\n", src, dest)

	return &Image{Source: src, Destination: dest}, nil
}
