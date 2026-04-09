package extractor

import (
	"archive/tar"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestExtractBlobAppliesWhiteoutWithinIncludedTree(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "dist")

	layer1 := filepath.Join(root, "layer1.tar")
	writeTarFile(t, layer1, []tarEntry{
		{name: ".git/lfs/cache_stats/old-a.json", body: "old-a"},
		{name: ".git/lfs/cache_stats/old-b.json", body: "old-b"},
	})

	layer2 := filepath.Join(root, "layer2.tar")
	writeTarFile(t, layer2, []tarEntry{
		{name: ".git/lfs/.wh.cache_stats"},
		{name: ".git/lfs/cache_stats/new.json", body: "new"},
	})

	opts := ExtractBlobOpts{
		Context:  context.Background(),
		Includes: []string{"/.git/lfs"},
		Logger:   zerolog.New(io.Discard),
	}

	require.NoError(t, ExtractBlob(layer1, dest, opts))
	require.NoError(t, ExtractBlob(layer2, dest, opts))

	require.NoFileExists(t, filepath.Join(dest, ".git", "lfs", "cache_stats", "old-a.json"))
	require.NoFileExists(t, filepath.Join(dest, ".git", "lfs", "cache_stats", "old-b.json"))
	require.FileExists(t, filepath.Join(dest, ".git", "lfs", "cache_stats", "new.json"))
}

func TestExtractBlobAppliesWhiteoutForIncludedDescendant(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "dist")

	layer1 := filepath.Join(root, "layer1.tar")
	writeTarFile(t, layer1, []tarEntry{
		{name: "a/b/c.txt", body: "old"},
	})

	layer2 := filepath.Join(root, "layer2.tar")
	writeTarFile(t, layer2, []tarEntry{
		{name: "a/.wh.b"},
	})

	opts := ExtractBlobOpts{
		Context:  context.Background(),
		Includes: []string{"/a/b"},
		Logger:   zerolog.New(io.Discard),
	}

	require.NoError(t, ExtractBlob(layer1, dest, opts))
	require.NoError(t, ExtractBlob(layer2, dest, opts))

	require.NoDirExists(t, filepath.Join(dest, "a", "b"))
}

func TestExtractBlobAppliesOpaqueWhiteoutForIncludedDescendant(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "dist")

	layer1 := filepath.Join(root, "layer1.tar")
	writeTarFile(t, layer1, []tarEntry{
		{name: "dir/sub/old.txt", body: "old"},
	})

	layer2 := filepath.Join(root, "layer2.tar")
	writeTarFile(t, layer2, []tarEntry{
		{name: "dir/.wh..wh..opq"},
	})

	opts := ExtractBlobOpts{
		Context:  context.Background(),
		Includes: []string{"/dir/sub"},
		Logger:   zerolog.New(io.Discard),
	}

	require.NoError(t, ExtractBlob(layer1, dest, opts))
	require.NoError(t, ExtractBlob(layer2, dest, opts))

	require.NoDirExists(t, filepath.Join(dest, "dir", "sub"))
}

type tarEntry struct {
	name string
	body string
}

func writeTarFile(t *testing.T, filename string, entries []tarEntry) {
	t.Helper()

	f, err := os.Create(filename)
	require.NoError(t, err)
	defer f.Close()

	tw := tar.NewWriter(f)
	defer tw.Close()

	for _, entry := range entries {
		hdr := &tar.Header{
			Name: entry.name,
			Mode: 0o644,
			Size: int64(len(entry.body)),
		}
		require.NoError(t, tw.WriteHeader(hdr))
		if entry.body == "" {
			continue
		}
		_, err = tw.Write([]byte(entry.body))
		require.NoError(t, err)
	}
}
