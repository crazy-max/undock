//go:build integration

package integration

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/crazy-max/undock/internal/app"
	"github.com/crazy-max/undock/internal/config"
	"github.com/google/go-containerregistry/pkg/registry"
	digest "github.com/opencontainers/go-digest"
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/require"
	"go.podman.io/image/v5/copy"
	"go.podman.io/image/v5/manifest"
	"go.podman.io/image/v5/signature"
	"go.podman.io/image/v5/transports/alltransports"
	"go.podman.io/image/v5/types"
)

func integrationRegistry(t *testing.T) string {
	t.Helper()
	if endpoint := os.Getenv("TEST_REGISTRY"); endpoint != "" {
		return endpoint
	}
	server := httptest.NewServer(registry.New())
	t.Cleanup(server.Close)
	return server.URL
}

func integrationCountingRegistry(t *testing.T) (string, *atomic.Int64) {
	t.Helper()
	var downloads atomic.Int64
	endpoint := integrationRegistryProxy(t, func(_ *http.Response) error {
		downloads.Add(1)
		return nil
	})
	return endpoint, &downloads
}

// Intercept layer responses while leaving registry discovery and configs intact.
func integrationRegistryProxy(t *testing.T, modifyLayer func(*http.Response) error) string {
	t.Helper()
	upstream, err := url.Parse(integrationRegistry(t))
	require.NoError(t, err)
	proxy := httputil.NewSingleHostReverseProxy(upstream)
	var layers sync.Map
	proxy.ModifyResponse = func(resp *http.Response) error {
		_, isLayer := layers.Load(filepath.Base(resp.Request.URL.Path))
		if resp.Request.Method == http.MethodGet && strings.Contains(resp.Request.URL.Path, "/blobs/") && isLayer {
			return modifyLayer(resp)
		}
		return nil
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut && strings.Contains(r.URL.Path, "/manifests/") {
			payload, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var manifest ocispecs.Manifest
			if err := json.Unmarshal(payload, &manifest); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			for _, layer := range manifest.Layers {
				layers.Store(layer.Digest.String(), struct{}{})
			}
			r.Body = io.NopCloser(bytes.NewReader(payload))
		}
		proxy.ServeHTTP(w, r)
	}))
	t.Cleanup(server.Close)
	return server.URL
}

func integrationRequest(t *testing.T, method, endpoint, mediaType string, data []byte, status int) http.Header {
	t.Helper()
	req, err := http.NewRequestWithContext(t.Context(), method, endpoint, bytes.NewReader(data)) // #nosec G704 -- Endpoint is the controlled test registry.
	require.NoError(t, err)
	req.Header.Set("Content-Type", mediaType)
	resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(req) // #nosec G704 -- Requests upload fixtures to the test registry.
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, status, resp.StatusCode, "%s %s", method, endpoint)
	return resp.Header
}

// Publish real layers without pulling or running a public base image.
func integrationImage(t *testing.T, endpoint, repo, revision string, platforms ...ocispecs.Platform) string {
	t.Helper()
	if len(platforms) == 0 {
		platforms = []ocispecs.Platform{
			{
				OS:           "linux",
				Architecture: "amd64",
			},
			{
				OS:           "linux",
				Architecture: "arm64",
				Variant:      "v8",
			},
		}
	}
	var manifests []ocispecs.Descriptor
	for _, platform := range platforms {
		arch := platform.Architecture
		dir := t.TempDir()
		lower, _ := marshalLayer(t, []ociLayerEntry{
			{
				name: "fixture/remove/old.txt",
				body: "old",
			},
			{
				name: "fixture/opaque/old.txt",
				body: "old",
			},
			{
				name: "excluded.txt",
				body: "excluded",
			},
		})
		upper, upperDiffID := marshalLayer(t, []ociLayerEntry{
			{
				name: "fixture/.wh.remove",
			},
			{
				name: "fixture/opaque/.wh..wh..opq",
			},
			{
				name: "fixture/remove/new.txt",
				body: revision,
			},
			{
				name: "fixture/opaque/new.txt",
				body: revision,
			},
			{
				name: "fixture/" + arch + ".txt",
				body: arch,
			},
		})
		var compressed bytes.Buffer
		zw := gzip.NewWriter(&compressed)
		_, err := zw.Write(upper)
		require.NoError(t, err)
		require.NoError(t, zw.Close())
		upper = compressed.Bytes()
		layers := []ocispecs.Descriptor{
			writeOCIObject(t, dir, lower, ocispecs.MediaTypeImageLayer),
			writeOCIObject(t, dir, upper, ocispecs.MediaTypeImageLayerGzip),
		}
		cfg, err := json.Marshal(ocispecs.Image{
			OS:           platform.OS,
			Architecture: arch,
			Variant:      platform.Variant,
			RootFS:       ocispecs.RootFS{Type: "layers", DiffIDs: []digest.Digest{layers[0].Digest, upperDiffID}},
		})
		require.NoError(t, err)
		configDesc := writeOCIObject(t, dir, cfg, ocispecs.MediaTypeImageConfig)
		for _, blob := range [][]byte{lower, upper, cfg} {
			integrationUploadBlob(t, endpoint, repo, blob)
		}
		payload, desc := marshalManifest(t, dir, configDesc, layers)
		integrationRequest(t, http.MethodPut, endpoint+"/v2/"+repo+"/manifests/"+desc.Digest.String(), desc.MediaType, payload, http.StatusCreated)
		desc.Platform = &platform
		manifests = append(manifests, desc)
	}
	index, err := json.Marshal(ocispecs.Index{SchemaVersion: 2, MediaType: ocispecs.MediaTypeImageIndex, Manifests: manifests})
	require.NoError(t, err)
	integrationRequest(t, http.MethodPut, endpoint+"/v2/"+repo+"/manifests/latest", ocispecs.MediaTypeImageIndex, index, http.StatusCreated)
	return strings.TrimPrefix(endpoint, "http://") + "/" + repo + ":latest"
}

func integrationUploadBlob(t *testing.T, endpoint, repo string, blob []byte) {
	t.Helper()
	headers := integrationRequest(t, http.MethodPost, endpoint+"/v2/"+repo+"/blobs/uploads/", "", nil, http.StatusAccepted)
	base, err := url.Parse(endpoint)
	require.NoError(t, err)
	location, err := base.Parse(headers.Get("Location"))
	require.NoError(t, err)
	query := location.Query()
	query.Set("digest", digest.FromBytes(blob).String())
	location.RawQuery = query.Encode()
	integrationRequest(t, http.MethodPut, location.String(), "application/octet-stream", blob, http.StatusCreated)
}

func integrationLayerImage(t *testing.T, entries []ociLayerEntry, dockerSchema2 bool) string {
	t.Helper()
	endpoint := integrationRegistry(t)
	repo := "undock-" + digest.FromString(t.TempDir()).Encoded()[:16]
	dir := t.TempDir()
	layer, diffID := marshalLayer(t, entries)
	var compressed bytes.Buffer
	zw := gzip.NewWriter(&compressed)
	_, err := zw.Write(layer)
	require.NoError(t, err)
	require.NoError(t, zw.Close())
	layer = compressed.Bytes()
	mediaType, configType, layerType := ocispecs.MediaTypeImageManifest, ocispecs.MediaTypeImageConfig, ocispecs.MediaTypeImageLayerGzip
	if dockerSchema2 {
		mediaType, configType, layerType = manifest.DockerV2Schema2MediaType, manifest.DockerV2Schema2ConfigMediaType, manifest.DockerV2Schema2LayerMediaType
	}
	cfg, err := json.Marshal(ocispecs.Image{
		OS:           "linux",
		Architecture: "amd64",
		RootFS:       ocispecs.RootFS{Type: "layers", DiffIDs: []digest.Digest{diffID}},
	})
	require.NoError(t, err)
	payload, err := json.Marshal(ocispecs.Manifest{
		SchemaVersion: 2,
		MediaType:     mediaType,
		Config:        writeOCIObject(t, dir, cfg, configType),
		Layers:        []ocispecs.Descriptor{writeOCIObject(t, dir, layer, layerType)},
	})
	require.NoError(t, err)
	integrationUploadBlob(t, endpoint, repo, layer)
	integrationUploadBlob(t, endpoint, repo, cfg)
	integrationRequest(t, http.MethodPut, endpoint+"/v2/"+repo+"/manifests/latest", mediaType, payload, http.StatusCreated)
	return strings.TrimPrefix(endpoint, "http://") + "/" + repo + ":latest"
}

type ociLayerEntry struct {
	name     string
	body     string
	mode     int64
	typeflag byte
	linkname string
}

func marshalLayer(t *testing.T, entries []ociLayerEntry) ([]byte, digest.Digest) {
	t.Helper()

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for _, entry := range entries {
		hdr := &tar.Header{
			Name:     entry.name,
			Mode:     0o644,
			Size:     int64(len(entry.body)),
			Typeflag: entry.typeflag,
			Linkname: entry.linkname,
		}
		if entry.mode != 0 {
			hdr.Mode = entry.mode
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
		SchemaVersion: 2,
		MediaType:     ocispecs.MediaTypeImageManifest,
		Config:        config,
		Layers:        layers,
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

func integrationExtract(t *testing.T, cli config.Cli) {
	t.Helper()
	require.NoError(t, integrationRun(t, cli))
}

func integrationRun(t *testing.T, cli config.Cli) error {
	t.Helper()
	application, err := app.New(config.Meta{UserAgent: "undock-integration"}, cli)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeoutCause(t.Context(), time.Minute, errors.New("extraction timed out"))
	defer cancel()
	return application.Start(ctx)
}

func integrationContents(t *testing.T, dir, revision string) {
	t.Helper()
	requireFileContent(t, filepath.Join(dir, "fixture/remove/new.txt"), revision)
	requireFileContent(t, filepath.Join(dir, "fixture/opaque/new.txt"), revision)
	require.NoFileExists(t, filepath.Join(dir, "fixture/remove/old.txt"))
	require.NoFileExists(t, filepath.Join(dir, "fixture/opaque/old.txt"))
	require.NoFileExists(t, filepath.Join(dir, "fixture/.wh.remove"))
	require.NoFileExists(t, filepath.Join(dir, "fixture/opaque/.wh..wh..opq"))
}

func integrationArchive(t *testing.T) (string, string) {
	t.Helper()
	endpoint := integrationRegistry(t)
	repo := "undock-" + digest.FromString(t.TempDir()).Encoded()[:16]

	source := integrationImage(t, endpoint, repo, "second")
	// Relative paths avoid confusing Windows drive letters with transport separators.
	t.Chdir(t.TempDir())
	archive := "image.tar"
	src, err := alltransports.ParseImageName("docker://" + source)
	require.NoError(t, err)
	imageName := repo + ":latest"
	dst, err := alltransports.ParseImageName("docker-archive:" + archive + ":" + imageName)
	require.NoError(t, err)
	policy, err := signature.NewPolicyContext(&signature.Policy{Default: []signature.PolicyRequirement{signature.NewPRInsecureAcceptAnything()}})
	require.NoError(t, err)
	defer policy.Destroy() //nolint:errcheck
	_, err = copy.Image(t.Context(), policy, dst, src, &copy.Options{SourceCtx: &types.SystemContext{DockerInsecureSkipTLSVerify: types.OptionalBoolTrue, OSChoice: "linux", ArchitectureChoice: "amd64"}})
	require.NoError(t, err)

	return archive, imageName
}

func requirePublicRegistry(t *testing.T) {
	t.Helper()
	if os.Getenv("TEST_PUBLIC_REGISTRY") != "1" {
		t.Skip("live Docker Hub test requires TEST_PUBLIC_REGISTRY=1")
	}
	// Use anonymous access even when the developer has registry credentials.
	home := t.TempDir()
	for _, key := range []string{"HOME", "USERPROFILE", "DOCKER_CONFIG", "XDG_CONFIG_HOME", "XDG_RUNTIME_DIR"} {
		t.Setenv(key, home)
	}
}

func requireFileSHA256(t *testing.T, filename, expected string) {
	t.Helper()
	f, err := os.Open(filename)
	require.NoError(t, err)
	defer f.Close()
	hash := sha256.New()
	_, err = io.Copy(hash, f)
	require.NoError(t, err)
	require.Equal(t, expected, hex.EncodeToString(hash.Sum(nil)), filename)
}

func integrationDownloadRecovery(t *testing.T, fault func(*http.Response) error, message string) {
	t.Helper()
	var enabled atomic.Bool
	var faults atomic.Int64
	endpoint := integrationRegistryProxy(t, func(resp *http.Response) error {
		if enabled.Load() {
			faults.Add(1)
			return fault(resp)
		}
		return nil
	})
	repo := "undock-" + digest.FromString(t.TempDir()).Encoded()[:16]
	source := integrationImage(t, endpoint, repo, "recovered")
	cli := config.Cli{Source: source, Platform: "linux/amd64", Insecure: true, Dist: t.TempDir(), CacheDir: t.TempDir()}
	enabled.Store(true)
	err := integrationRun(t, cli)
	require.ErrorContains(t, err, message)
	require.Positive(t, faults.Load(), "the failure must reach a layer download")
	t.Logf("injected download failure: %v", err)
	entries, err := os.ReadDir(cli.Dist)
	require.NoError(t, err)
	require.Empty(t, entries, "failed caching must not start extraction")
	enabled.Store(false)
	// Retain the failed download's cache and destination for the recovery attempt.
	integrationExtract(t, cli)
	integrationContents(t, cli.Dist, "recovered")
}
