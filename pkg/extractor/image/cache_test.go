package image

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"

	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.podman.io/image/v5/types"
)

func TestSrcCtxUsesRegistrySettings(t *testing.T) {
	c := &Client{
		opts: Options{
			Platform: ocispecs.Platform{
				OS:           "linux",
				Architecture: "arm64",
				Variant:      "v8",
			},
			CacheDir: filepath.Join("cache", "root"),

			RegistryUserAgent: "undock-tests",
		},
	}

	sysCtx, err := c.srcCtx(&types.DockerAuthConfig{Username: "alice"}, true)
	require.NoError(t, err)

	assert.Equal(t, "alice", sysCtx.DockerAuthConfig.Username)
	assert.True(t, sysCtx.DockerDaemonInsecureSkipTLSVerify)
	assert.Equal(t, types.OptionalBoolTrue, sysCtx.DockerInsecureSkipTLSVerify)
	assert.Equal(t, "undock-tests", sysCtx.DockerRegistryUserAgent)
	assert.Equal(t, "linux", sysCtx.OSChoice)
	assert.Equal(t, "arm64", sysCtx.ArchitectureChoice)
	assert.Equal(t, "v8", sysCtx.VariantChoice)
	assert.Equal(t, filepath.Join("cache", "root", "blobs"), sysCtx.BlobInfoCacheDir)
}

func TestDstCtxForcesDecompression(t *testing.T) {
	c := &Client{opts: Options{CacheDir: filepath.Join("cache", "root")}}

	sysCtx, err := c.dstCtx("ignored")
	require.NoError(t, err)

	assert.True(t, sysCtx.DirForceDecompress)
	assert.Equal(t, filepath.Join("cache", "root", "blobs"), sysCtx.BlobInfoCacheDir)
}

func TestProgressWriterTrimsTrailingWhitespace(t *testing.T) {
	var buf bytes.Buffer
	writer := progressWriter{logger: zerolog.New(&buf)}

	n, err := writer.Write([]byte("copying blob\n"))
	require.NoError(t, err)
	assert.Equal(t, len("copying blob\n"), n)

	var event map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &event))
	assert.Equal(t, "copying blob", event["message"])
}
