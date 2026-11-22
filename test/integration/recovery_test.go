//go:build integration

package integration

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/crazy-max/undock/internal/config"
	digest "github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/require"
)

func TestMissingPlatformRecovery(t *testing.T) {
	endpoint := integrationRegistry(t)
	repo := "undock-" + digest.FromString(t.TempDir()).Encoded()[:16]
	source := integrationImage(t, endpoint, repo, "recovered")
	cli := config.Cli{Source: source, Platform: "linux/s390x", Insecure: true, Dist: t.TempDir(), CacheDir: t.TempDir()}
	err := integrationRun(t, cli)
	require.ErrorContains(t, err, "no image found")
	require.ErrorContains(t, err, "s390x")
	entries, err := os.ReadDir(cli.Dist)
	require.NoError(t, err)
	require.Empty(t, entries)
	cli.Platform = "linux/amd64"
	integrationExtract(t, cli)
	integrationContents(t, cli.Dist, "recovered")
}

func TestInterruptedDownloadRecovery(t *testing.T) {
	integrationDownloadRecovery(t, func(resp *http.Response) error {
		// Keep the advertised size but close the stream halfway through the body.
		resp.Body = struct {
			io.Reader
			io.Closer
		}{io.LimitReader(resp.Body, resp.ContentLength/2), resp.Body}
		return nil
	}, "unexpected EOF")
}

func TestCorruptBlobRecovery(t *testing.T) {
	integrationDownloadRecovery(t, func(resp *http.Response) error {
		payload, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return err
		}
		// Preserve the size and digest URL while supplying different content.
		resp.Body = io.NopCloser(bytes.NewReader(bytes.Repeat([]byte{'x'}, len(payload))))
		return nil
	}, "Digest did not match")
}

func TestRegistryErrorRecovery(t *testing.T) {
	integrationDownloadRecovery(t, func(resp *http.Response) error {
		resp.Body.Close()
		body := `{"errors":[{"code":"UNKNOWN","message":"injected registry failure"}]}`
		resp.StatusCode = http.StatusServiceUnavailable
		resp.Status = "503 Service Unavailable"
		resp.Header.Set("Content-Type", "application/json")
		resp.ContentLength = int64(len(body))
		resp.Header.Set("Content-Length", strconv.Itoa(len(body)))
		resp.Body = io.NopCloser(strings.NewReader(body))
		return nil
	}, "503")
}
