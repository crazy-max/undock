package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSourceDefaultsToDocker(t *testing.T) {
	src := NewSource("alpine:latest")

	assert.Equal(t, "docker", src.Scheme())
	assert.False(t, src.HasScheme())
	assert.Equal(t, "docker://alpine:latest", src.String())

	ref, err := src.Reference()
	require.NoError(t, err)
	assert.Equal(t, "docker", ref.Transport().Name())
	require.NotNil(t, ref.DockerReference())
	assert.Equal(t, "docker.io/library/alpine:latest", ref.DockerReference().String())
}

func TestSourceRewritesDaemonScheme(t *testing.T) {
	src := NewSource("docker-daemon://alpine:latest")

	assert.Equal(t, "docker-daemon", src.Scheme())
	assert.True(t, src.HasScheme())
	assert.Equal(t, "docker-daemon:alpine:latest", src.String())
}

func TestSourceParsesArchiveTransport(t *testing.T) {
	src := NewSource("oci-archive://image.tar")

	ref, err := src.Reference()
	require.NoError(t, err)

	assert.Equal(t, "oci-archive", ref.Transport().Name())
	assert.Contains(t, ref.StringWithinTransport(), "image.tar")
}
