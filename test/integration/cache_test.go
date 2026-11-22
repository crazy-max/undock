//go:build integration

package integration

import (
	"path/filepath"
	"testing"

	"github.com/crazy-max/undock/internal/config"
	digest "github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/require"
)

func TestRegistryCacheAvoidsBlobDownloads(t *testing.T) {
	endpoint, downloads := integrationCountingRegistry(t)
	repo := "undock-" + digest.FromString(t.TempDir()).Encoded()[:16]
	source := integrationImage(t, endpoint, repo, "cached")
	cli := config.Cli{Source: source, Insecure: true, All: true, CacheDir: t.TempDir()}

	for run := range 3 {
		// A fresh destination proves cached blobs still produce the complete output.
		cli.Dist = t.TempDir()
		if run == 2 {
			// A fresh cache is the control: it must download blobs again.
			cli.CacheDir = t.TempDir()
		}
		before := downloads.Load()
		integrationExtract(t, cli)
		after := downloads.Load()
		if run == 1 {
			require.Equal(t, before, after, "warm cache must not issue layer blob GET requests")
		} else {
			require.Greater(t, after, before, "cold cache must issue layer blob GET requests")
		}
		t.Logf("run %d: %d layer blob GET requests", run+1, after-before)
		for _, platform := range []string{"linux_amd64", "linux_arm64v8"} {
			integrationContents(t, filepath.Join(cli.Dist, platform), "cached")
		}
		requireFileContent(t, filepath.Join(cli.Dist, "linux_amd64/fixture/amd64.txt"), "amd64")
		requireFileContent(t, filepath.Join(cli.Dist, "linux_arm64v8/fixture/arm64.txt"), "arm64")
	}
}
