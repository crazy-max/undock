package image

import (
	"fmt"
	"strings"

	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/types"
)

func Reference(name string) (types.ImageReference, error) {
	if !strings.HasPrefix(name, "//") {
		name = fmt.Sprintf("//%s", name)
	}
	return docker.ParseReference(name)
}
