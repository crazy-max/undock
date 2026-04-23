package image

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/containerd/platforms"
	digest "github.com/opencontainers/go-digest"
	specs "github.com/opencontainers/image-spec/specs-go"
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultsCacheDir(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", t.TempDir())

	cli, err := New(context.Background(), Options{Source: "alpine:latest"})
	require.NoError(t, err)

	handler, ok := cli.Handler.(*Client)
	require.True(t, ok)

	assert.Equal(t, platforms.DefaultSpec(), handler.opts.Platform)
	assert.Equal(t, filepath.Join(os.Getenv("XDG_DATA_HOME"), "undock", "cache"), handler.opts.CacheDir)
	assert.DirExists(t, handler.opts.CacheDir)
}

func TestNewKeepsExplicitCacheDir(t *testing.T) {
	cacheDir := filepath.Join(t.TempDir(), "cache")

	cli, err := New(context.Background(), Options{
		Source:   "alpine:latest",
		Platform: platforms.DefaultSpec(),
		CacheDir: cacheDir,
	})
	require.NoError(t, err)

	handler, ok := cli.Handler.(*Client)
	require.True(t, ok)

	assert.Equal(t, cacheDir, handler.opts.CacheDir)
	assert.DirExists(t, cacheDir)
}

func TestExtractCachedSourceSingleManifest(t *testing.T) {
	root := t.TempDir()
	cachedir := filepath.Join(root, "cache")
	dest := filepath.Join(root, "dist")

	layer := writeLayerBlob(t, cachedir, []layerEntry{
		{name: "etc/app/config.yaml", body: "wanted"},
		{name: "usr/bin/tool", body: "skip"},
	})
	manblob := writeManifestBlob(t, cachedir, []ocispecs.Descriptor{layer})

	c := &Client{
		ctx: context.Background(),
		opts: Options{
			Dist:     dest,
			Includes: []string{"/etc/app"},
		},
		logger: zerolog.New(io.Discard),
	}

	require.NoError(t, c.extractCachedSource(manblob, cachedir))
	require.FileExists(t, filepath.Join(dest, "etc", "app", "config.yaml"))
	require.NoFileExists(t, filepath.Join(dest, "usr", "bin", "tool"))
}

func TestExtractCachedSourceSplitsPlatforms(t *testing.T) {
	root := t.TempDir()
	cachedir := filepath.Join(root, "cache")
	dest := filepath.Join(root, "dist")

	amd64Layer := writeLayerBlob(t, cachedir, []layerEntry{{name: "arch.txt", body: "amd64"}})
	_, amd64Desc := writeManifestBlobWithDescriptor(t, cachedir, []ocispecs.Descriptor{amd64Layer})
	amd64Desc.Platform = &ocispecs.Platform{OS: "linux", Architecture: "amd64"}

	arm64Layer := writeLayerBlob(t, cachedir, []layerEntry{{name: "arch.txt", body: "arm64"}})
	_, arm64Desc := writeManifestBlobWithDescriptor(t, cachedir, []ocispecs.Descriptor{arm64Layer})
	arm64Desc.Platform = &ocispecs.Platform{OS: "linux", Architecture: "arm64", Variant: "v8"}

	indexBlob := writeIndexBlob(t, []ocispecs.Descriptor{amd64Desc, arm64Desc})

	c := &Client{
		ctx:    context.Background(),
		opts:   Options{Dist: dest},
		logger: zerolog.New(io.Discard),
	}

	require.NoError(t, c.extractCachedSource(indexBlob, cachedir))

	requireFileContent(t, filepath.Join(dest, "linux_amd64", "arch.txt"), "amd64")
	requireFileContent(t, filepath.Join(dest, "linux_arm64v8", "arch.txt"), "arm64")
}

func TestExtractCachedSourceWrapsPlatforms(t *testing.T) {
	root := t.TempDir()
	cachedir := filepath.Join(root, "cache")
	dest := filepath.Join(root, "dist")

	amd64Layer := writeLayerBlob(t, cachedir, []layerEntry{{name: "amd64.txt", body: "amd64"}})
	_, amd64Desc := writeManifestBlobWithDescriptor(t, cachedir, []ocispecs.Descriptor{amd64Layer})
	amd64Desc.Platform = &ocispecs.Platform{OS: "linux", Architecture: "amd64"}

	arm64Layer := writeLayerBlob(t, cachedir, []layerEntry{{name: "arm64.txt", body: "arm64"}})
	_, arm64Desc := writeManifestBlobWithDescriptor(t, cachedir, []ocispecs.Descriptor{arm64Layer})
	arm64Desc.Platform = &ocispecs.Platform{OS: "linux", Architecture: "arm64", Variant: "v8"}

	indexBlob := writeIndexBlob(t, []ocispecs.Descriptor{amd64Desc, arm64Desc})

	c := &Client{
		ctx: context.Background(),
		opts: Options{
			Dist: dest,
			Wrap: true,
		},
		logger: zerolog.New(io.Discard),
	}

	require.NoError(t, c.extractCachedSource(indexBlob, cachedir))
	requireFileContent(t, filepath.Join(dest, "amd64.txt"), "amd64")
	requireFileContent(t, filepath.Join(dest, "arm64.txt"), "arm64")
}

func TestExtractCachedSourceFailsOnMissingManifestBlob(t *testing.T) {
	root := t.TempDir()
	cachedir := filepath.Join(root, "cache")

	indexBlob := writeIndexBlob(t, []ocispecs.Descriptor{{
		MediaType: ocispecs.MediaTypeImageManifest,
		Digest:    digest.FromString("missing-manifest"),
		Size:      123,
		Platform:  &ocispecs.Platform{OS: "linux", Architecture: "amd64"},
	}})

	c := &Client{
		ctx:    context.Background(),
		opts:   Options{Dist: filepath.Join(root, "dist")},
		logger: zerolog.New(io.Discard),
	}

	err := c.extractCachedSource(indexBlob, cachedir)
	require.ErrorContains(t, err, "cannot read OCI manifest JSON for platform linux/amd64")
}

func TestExtractCachedSourceFailsOnInvalidManifestBlob(t *testing.T) {
	root := t.TempDir()
	cachedir := filepath.Join(root, "cache")

	manifestDesc := writeOCIObject(t, cachedir, []byte(`not-json`), ocispecs.MediaTypeImageManifest)
	manifestDesc.Platform = &ocispecs.Platform{OS: "linux", Architecture: "amd64"}
	indexBlob := writeIndexBlob(t, []ocispecs.Descriptor{manifestDesc})

	c := &Client{
		ctx:    context.Background(),
		opts:   Options{Dist: filepath.Join(root, "dist")},
		logger: zerolog.New(io.Discard),
	}

	err := c.extractCachedSource(indexBlob, cachedir)
	require.ErrorContains(t, err, "cannot create OCI manifest instance from blob")
}

type layerEntry struct {
	name string
	body string
}

func writeLayerBlob(t *testing.T, cachedir string, entries []layerEntry) ocispecs.Descriptor {
	t.Helper()

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for _, entry := range entries {
		hdr := &tar.Header{
			Name: entry.name,
			Mode: 0o644,
			Size: int64(len(entry.body)),
		}
		require.NoError(t, tw.WriteHeader(hdr))
		_, err := tw.Write([]byte(entry.body))
		require.NoError(t, err)
	}
	require.NoError(t, tw.Close())

	return writeOCIObject(t, cachedir, buf.Bytes(), ocispecs.MediaTypeImageLayer)
}

func writeManifestBlob(t *testing.T, cachedir string, layers []ocispecs.Descriptor) []byte {
	t.Helper()
	blob, _ := writeManifestBlobWithDescriptor(t, cachedir, layers)
	return blob
}

func writeManifestBlobWithDescriptor(t *testing.T, cachedir string, layers []ocispecs.Descriptor) ([]byte, ocispecs.Descriptor) {
	t.Helper()

	configDesc := writeOCIObject(t, cachedir, []byte(`{}`), ocispecs.MediaTypeImageConfig)
	payload, err := json.Marshal(ocispecs.Manifest{
		Versioned: specs.Versioned{SchemaVersion: 2},
		MediaType: ocispecs.MediaTypeImageManifest,
		Config:    configDesc,
		Layers:    layers,
	})
	require.NoError(t, err)

	return payload, writeOCIObject(t, cachedir, payload, ocispecs.MediaTypeImageManifest)
}

func writeIndexBlob(t *testing.T, manifests []ocispecs.Descriptor) []byte {
	t.Helper()

	payload, err := json.Marshal(ocispecs.Index{
		Versioned: specs.Versioned{SchemaVersion: 2},
		MediaType: ocispecs.MediaTypeImageIndex,
		Manifests: manifests,
	})
	require.NoError(t, err)

	return payload
}

func writeOCIObject(t *testing.T, cachedir string, payload []byte, mediaType string) ocispecs.Descriptor {
	t.Helper()

	dgst := digest.FromBytes(payload)
	blobPath := filepath.Join(cachedir, "blobs", dgst.Algorithm().String(), dgst.Hex())
	require.NoError(t, os.MkdirAll(filepath.Dir(blobPath), 0o755))
	require.NoError(t, os.WriteFile(blobPath, payload, 0o644))

	return ocispecs.Descriptor{
		MediaType: mediaType,
		Digest:    dgst,
		Size:      int64(len(payload)),
	}
}

func requireFileContent(t *testing.T, filename string, expected string) {
	t.Helper()

	data, err := os.ReadFile(filename)
	require.NoError(t, err)
	require.Equal(t, expected, string(data))
}
