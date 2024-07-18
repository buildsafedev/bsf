package skopeo

import (
	"context"
	"errors"
	"fmt"
	"io"
)

type Options struct {
	ReportWriter io.Writer
}

type Image struct{}

func CopyImage(ctx context.Context, policyContext *PolicyContext, destRef, srcRef ImageReference, options *Options) (*Image, error) {
	if srcRef == nil || destRef == nil {
		return nil, errors.New("invalid image reference")
	}

	if options.ReportWriter != nil {
		_, err := io.WriteString(options.ReportWriter, "Copying image...\n")
		if err != nil {
			return nil, fmt.Errorf("failed to write to report writer: %w", err)
		}
	}

	return &Image{}, nil
}
