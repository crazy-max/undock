//go:build integration

package integration

import (
	"archive/tar"
	"os"
	"path/filepath"
	"testing"

	"github.com/crazy-max/undock/internal/config"
	"github.com/stretchr/testify/require"
)

func TestImageSymlink(t *testing.T) {
	if err := os.Symlink("target", filepath.Join(t.TempDir(), "probe")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	source := integrationLayerImage(t, []ociLayerEntry{
		{
			name: "fixture/target",
			body: "linked content",
		},
		{
			name:     "fixture/link",
			typeflag: tar.TypeSymlink,
			linkname: "target",
		},
	}, false)
	dist := t.TempDir()
	integrationExtract(t, config.Cli{Source: source, Platform: "linux/amd64", Insecure: true, Dist: dist, CacheDir: t.TempDir()})
	link := filepath.Join(dist, "fixture/link")
	target, err := os.Readlink(link)
	require.NoError(t, err)
	require.Equal(t, "target", target)
	requireFileContent(t, link, "linked content")
}

func TestImageHardlink(t *testing.T) {
	source := integrationLayerImage(t, []ociLayerEntry{
		{
			name: "fixture/target",
			body: "linked content",
		},
		{
			name:     "fixture/link",
			typeflag: tar.TypeLink,
			linkname: "fixture/target",
		},
	}, false)
	dist := t.TempDir()
	integrationExtract(t, config.Cli{Source: source, Platform: "linux/amd64", Insecure: true, Dist: dist, CacheDir: t.TempDir()})
	requireFileContent(t, filepath.Join(dist, "fixture/link"), "linked content")
	target, err := os.Stat(filepath.Join(dist, "fixture/target"))
	require.NoError(t, err)
	link, err := os.Stat(filepath.Join(dist, "fixture/link"))
	require.NoError(t, err)
	require.True(t, os.SameFile(target, link), "hardlink must share the target's file identity")
}

func TestDockerSchema2Registry(t *testing.T) {
	source := integrationLayerImage(t, []ociLayerEntry{
		{
			name: "fixture/content",
			body: "Docker schema 2",
		},
	}, true)
	dist := t.TempDir()
	integrationExtract(t, config.Cli{Source: source, Platform: "linux/amd64", Insecure: true, Dist: dist, CacheDir: t.TempDir()})
	requireFileContent(t, filepath.Join(dist, "fixture/content"), "Docker schema 2")
}
