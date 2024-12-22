package image

import (
	"fmt"
	"strings"

	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/docker/reference"
	"github.com/containers/image/v5/types"
	"github.com/pkg/errors"
)

func DockerReference(name string) (types.ImageReference, error) {
	ref, err := Reference(name)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse reference")
	}
	refStr := ref.String()
	if !strings.HasPrefix(refStr, "//") {
		refStr = fmt.Sprintf("//%s", refStr)
	}
	return docker.ParseReference(refStr)
}

func Reference(name string) (reference.Named, error) {
	name = strings.TrimPrefix(name, "//")

	ref, err := reference.ParseNormalizedNamed(name)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing normalized named %q", name)
	} else if ref == nil {
		return nil, errors.Errorf("%q is not a named reference", name)
	}

	if _, hasTag := ref.(reference.NamedTagged); hasTag {
		ref, err = normalizeTaggedDigestedNamed(ref)
		if err != nil {
			return nil, errors.Wrapf(err, "normalizing tagged digested name %q", name)
		}
		return ref, nil
	}
	if _, hasDigest := ref.(reference.Digested); hasDigest {
		return ref, nil
	}

	return reference.TagNameOnly(ref), nil
}

// normalizeTaggedDigestedNamed strips the tag off the specified named
// reference iff it is tagged and digested. Note that the tag is entirely
// ignored.
func normalizeTaggedDigestedNamed(named reference.Named) (reference.Named, error) {
	_, isTagged := named.(reference.NamedTagged)
	if !isTagged {
		return named, nil
	}
	digested, isDigested := named.(reference.Digested)
	if !isDigested {
		return named, nil
	}
	// Now strip off the tag.
	newNamed := reference.TrimNamed(named)
	// And re-add the digest.
	newNamed, err := reference.WithDigest(newNamed, digested.Digest())
	if err != nil {
		return named, err
	}
	return newNamed, nil
}
