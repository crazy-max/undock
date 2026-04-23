package app

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/containerd/platforms"
	"github.com/crazy-max/undock/internal/config"
	digest "github.com/opencontainers/go-digest"
	specs "github.com/opencontainers/image-spec/specs-go"
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultsPlatform(t *testing.T) {
	app, err := New(config.Meta{}, config.Cli{})
	require.NoError(t, err)

	assert.Equal(t, platforms.DefaultSpec(), app.platform)
}

func TestNewParsesPlatform(t *testing.T) {
	app, err := New(config.Meta{}, config.Cli{Platform: "linux/arm64/v8"})
	require.NoError(t, err)

	assert.Equal(t, "linux", app.platform.OS)
	assert.Equal(t, "arm64", app.platform.Architecture)
	assert.Equal(t, "v8", app.platform.Variant)
}

func TestNewRejectsInvalidPlatform(t *testing.T) {
	_, err := New(config.Meta{}, config.Cli{Platform: "linux/nope/extra/parts"})
	require.Error(t, err)
	require.ErrorContains(t, err, `invalid platform "linux/nope/extra/parts"`)
}

func TestValidateSchemeAcceptsKnownSchemes(t *testing.T) {
	testCases := []string{
		"containers-storage://image",
		"docker://alpine:latest",
		"docker-archive://archive.tar",
		"docker-daemon://alpine:latest",
		"oci://layout",
		"oci-archive://layout.tar",
		"ostree://ref",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			ok, err := validateScheme(tc)
			require.NoError(t, err)
			assert.True(t, ok)
		})
	}
}

func TestValidateSchemeAcceptsImageReference(t *testing.T) {
	ok, err := validateScheme("alpine:latest")
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestValidateSchemeRejectsUnsupportedSource(t *testing.T) {
	ok, err := validateScheme("ftp://example.com/image")
	require.Error(t, err)
	assert.False(t, ok)
}

func TestStartExtractsOCIImage(t *testing.T) {
	root := t.TempDir()
	layoutDir := filepath.Join(root, "layout")
	cacheDir := filepath.Join(root, "cache")
	distDir := filepath.Join(root, "dist")
	require.NoError(t, os.MkdirAll(distDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(distDir, "stale.txt"), []byte("old"), 0o644))

	createOCIImageLayout(t, layoutDir, platforms.DefaultSpec(), []ociLayerEntry{
		{name: "etc/app/config.yaml", body: "wanted"},
		{name: "usr/bin/tool", body: "skip"},
	})

	app, err := New(config.Meta{
		UserAgent: "undock-tests",
	}, config.Cli{
		Source:   "oci://" + layoutDir,
		Dist:     distDir,
		CacheDir: cacheDir,
		Includes: []string{"/etc/app"},
		RmDist:   true,
	})
	require.NoError(t, err)

	require.NoError(t, app.Start(context.Background()))

	require.NoFileExists(t, filepath.Join(distDir, "stale.txt"))
	requireFileContent(t, filepath.Join(distDir, "etc", "app", "config.yaml"), "wanted")
	require.NoFileExists(t, filepath.Join(distDir, "usr", "bin", "tool"))

	cacheEntries, err := os.ReadDir(cacheDir)
	require.NoError(t, err)
	var foundDigestDir bool
	for _, entry := range cacheEntries {
		if !entry.IsDir() || entry.Name() == "blobs" {
			continue
		}
		foundDigestDir = true
		require.DirExists(t, filepath.Join(cacheDir, entry.Name(), "blobs", "sha256"))
	}
	assert.True(t, foundDigestDir)
	require.DirExists(t, filepath.Join(cacheDir, "blobs"))
}

type ociLayerEntry struct {
	name string
	body string
}

func createOCIImageLayout(t *testing.T, dir string, platform ocispecs.Platform, entries []ociLayerEntry) {
	t.Helper()

	require.NoError(t, os.MkdirAll(filepath.Join(dir, ocispecs.ImageBlobsDir, "sha256"), 0o755))

	layerPayload, layerDigest := marshalLayer(t, entries)
	layerDesc := writeOCIObject(t, dir, layerPayload, ocispecs.MediaTypeImageLayer)

	configPayload, err := json.Marshal(map[string]any{
		"architecture": platform.Architecture,
		"os":           platform.OS,
		"variant":      platform.Variant,
		"rootfs": map[string]any{
			"type":     "layers",
			"diff_ids": []string{layerDigest.String()},
		},
	})
	require.NoError(t, err)
	configDesc := writeOCIObject(t, dir, configPayload, ocispecs.MediaTypeImageConfig)

	manifestPayload, manifestDesc := marshalManifest(t, dir, configDesc, []ocispecs.Descriptor{layerDesc})
	manifestDesc.Platform = &platform

	indexPayload, err := json.Marshal(ocispecs.Index{
		Versioned: specs.Versioned{SchemaVersion: 2},
		MediaType: ocispecs.MediaTypeImageIndex,
		Manifests: []ocispecs.Descriptor{manifestDesc},
	})
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(filepath.Join(dir, ocispecs.ImageLayoutFile), []byte(`{"imageLayoutVersion":"`+ocispecs.ImageLayoutVersion+`"}`), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, ocispecs.ImageIndexFile), indexPayload, 0o644))
	_ = manifestPayload
}

func marshalLayer(t *testing.T, entries []ociLayerEntry) ([]byte, digest.Digest) {
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

	payload := buf.Bytes()
	return payload, digest.FromBytes(payload)
}

func marshalManifest(t *testing.T, dir string, config ocispecs.Descriptor, layers []ocispecs.Descriptor) ([]byte, ocispecs.Descriptor) {
	t.Helper()

	payload, err := json.Marshal(ocispecs.Manifest{
		Versioned: specs.Versioned{SchemaVersion: 2},
		MediaType: ocispecs.MediaTypeImageManifest,
		Config:    config,
		Layers:    layers,
	})
	require.NoError(t, err)

	return payload, writeOCIObject(t, dir, payload, ocispecs.MediaTypeImageManifest)
}

func writeOCIObject(t *testing.T, dir string, payload []byte, mediaType string) ocispecs.Descriptor {
	t.Helper()

	dgst := digest.FromBytes(payload)
	blobPath := filepath.Join(dir, ocispecs.ImageBlobsDir, dgst.Algorithm().String(), dgst.Hex())
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
