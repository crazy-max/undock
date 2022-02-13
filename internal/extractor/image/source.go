package image

import (
	"strings"

	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
)

var sourceSchemes = []string{"containers-storage", "docker", "docker-archive", "docker-daemon", "oci", "oci-archive", "ostree"}

type Source struct {
	name string
}

func NewSource(name string) *Source {
	return &Source{
		name: name,
	}
}

func (s *Source) Scheme() string {
	for _, scheme := range sourceSchemes {
		if strings.HasPrefix(s.name, scheme+"://") {
			return scheme
		}
	}
	return "docker"
}

func (s *Source) HasScheme() bool {
	for _, scheme := range sourceSchemes {
		if strings.HasPrefix(s.name, scheme+"://") {
			return true
		}
	}
	return false
}

func (s *Source) Reference() (types.ImageReference, error) {
	return alltransports.ParseImageName(s.String())
}

func (s *Source) String() string {
	if !s.HasScheme() {
		return "docker://" + s.name
	}
	scheme := s.Scheme()
	switch scheme {
	case "docker":
		return s.name
	default:
		return strings.Replace(s.name, scheme+"://", scheme+":", 1)
	}
}
