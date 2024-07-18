package skopeo

import (
	"errors"
)

type ImageReference interface{}

func ParseImageName(ref string) (ImageReference, error) {
	if ref == "" {
		return nil, errors.New("invalid image reference")
	}
	return ref, nil
}
