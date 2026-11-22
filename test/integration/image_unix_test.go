//go:build integration && !windows

package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/crazy-max/undock/internal/config"
	"github.com/stretchr/testify/require"
)

func TestImageExecutable(t *testing.T) {
	source := integrationLayerImage(t, []ociLayerEntry{
		{
			name: "bin/tool",
			body: "#!/bin/sh\nexit 0\n",
			mode: 0o755,
		},
	}, false)
	dist := t.TempDir()
	integrationExtract(t, config.Cli{Source: source, Platform: "linux/amd64", Insecure: true, Dist: dist, CacheDir: t.TempDir()})
	info, err := os.Stat(filepath.Join(dist, "bin/tool"))
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0o755), info.Mode().Perm())
	requireFileContent(t, filepath.Join(dist, "bin/tool"), "#!/bin/sh\nexit 0\n")
}
