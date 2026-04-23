package extractor

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/mholt/archives"
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

func TestExtractBlobKeepsSameLayerFilesWhenOpaqueWhiteoutComesLater(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "dist")

	layer1 := filepath.Join(root, "layer1.tar")
	writeTarFile(t, layer1, []tarEntry{
		{name: "dir/lower.txt", body: "lower"},
	})

	layer2 := filepath.Join(root, "layer2.tar")
	writeTarFile(t, layer2, []tarEntry{
		{name: "dir/keep.txt", body: "keep"},
		{name: "dir/.wh..wh..opq"},
	})

	opts := ExtractBlobOpts{
		Context: context.Background(),
		Logger:  zerolog.New(io.Discard),
	}

	require.NoError(t, ExtractBlob(layer1, dest, opts))
	require.NoError(t, ExtractBlob(layer2, dest, opts))

	require.NoFileExists(t, filepath.Join(dest, "dir", "lower.txt"))
	require.FileExists(t, filepath.Join(dest, "dir", "keep.txt"))
}

func TestExtractBlobKeepsSameLayerFilesWhenWhiteoutComesLater(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "dist")

	layer1 := filepath.Join(root, "layer1.tar")
	writeTarFile(t, layer1, []tarEntry{
		{name: "dir/sub/lower.txt", body: "lower"},
	})

	layer2 := filepath.Join(root, "layer2.tar")
	writeTarFile(t, layer2, []tarEntry{
		{name: "dir/sub/keep.txt", body: "keep"},
		{name: "dir/.wh.sub"},
	})

	opts := ExtractBlobOpts{
		Context: context.Background(),
		Logger:  zerolog.New(io.Discard),
	}

	require.NoError(t, ExtractBlob(layer1, dest, opts))
	require.NoError(t, ExtractBlob(layer2, dest, opts))

	require.NoFileExists(t, filepath.Join(dest, "dir", "sub", "lower.txt"))
	require.FileExists(t, filepath.Join(dest, "dir", "sub", "keep.txt"))
}

func TestExtractBlobRejectsBreakoutArchivePath(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "dist")
	outside := filepath.Join(root, "escape.txt")

	layer := filepath.Join(root, "layer.tar")
	writeTarFile(t, layer, []tarEntry{
		{name: "../escape.txt", body: "nope"},
	})

	err := ExtractBlob(layer, dest, ExtractBlobOpts{
		Context: context.Background(),
		Logger:  zerolog.New(io.Discard),
	})
	require.ErrorContains(t, err, "resolves outside destination")
	require.NoFileExists(t, outside)
}

func TestExtractBlobRejectsBreakoutWhiteoutPath(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "dist")
	outside := filepath.Join(root, "escape.txt")
	require.NoError(t, os.WriteFile(outside, []byte("keep"), 0o644))

	layer := filepath.Join(root, "layer.tar")
	writeTarFile(t, layer, []tarEntry{
		{name: "..\\.wh.escape.txt"},
	})

	err := ExtractBlob(layer, dest, ExtractBlobOpts{
		Context: context.Background(),
		Logger:  zerolog.New(io.Discard),
	})
	require.ErrorContains(t, err, "resolves outside destination")
	require.FileExists(t, outside)
}

func TestExtractBlobIgnoresReservedWhiteoutMetadataFile(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "dist")
	require.NoError(t, os.MkdirAll(filepath.Join(dest, "dir"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dest, "dir", ".wh.keep"), []byte("keep"), 0o644))

	layer2 := filepath.Join(root, "layer2.tar")
	writeTarFile(t, layer2, []tarEntry{
		{name: "dir/.wh..wh.keep"},
	})

	opts := ExtractBlobOpts{
		Context: context.Background(),
		Logger:  zerolog.New(io.Discard),
	}

	require.NoError(t, ExtractBlob(layer2, dest, opts))

	require.FileExists(t, filepath.Join(dest, "dir", ".wh.keep"))
}

func TestExtractBlobSkipsWhiteoutMetadataDirectory(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "dist")

	layer := filepath.Join(root, "layer.tar")
	writeTarFile(t, layer, []tarEntry{
		{name: ".wh..wh.plnk/link", body: "metadata"},
		{name: "real.txt", body: "real"},
	})

	err := ExtractBlob(layer, dest, ExtractBlobOpts{
		Context: context.Background(),
		Logger:  zerolog.New(io.Discard),
	})
	require.NoError(t, err)

	require.NoDirExists(t, filepath.Join(dest, ".wh..wh.plnk"))
	require.FileExists(t, filepath.Join(dest, "real.txt"))
}

func TestExtractBlobExtractsTarGz(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "dist")

	layer := filepath.Join(root, "layer.tar.gz")
	writeTarGzFile(t, layer, []tarEntry{
		{name: "usr/bin/tool", body: "binary"},
	})

	err := ExtractBlob(layer, dest, ExtractBlobOpts{
		Context: context.Background(),
		Logger:  zerolog.New(io.Discard),
	})
	require.NoError(t, err)

	require.FileExists(t, filepath.Join(dest, "usr", "bin", "tool"))
}

func TestExtractBlobDoesNotOvermatchIncludes(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "dist")

	layer := filepath.Join(root, "layer.tar")
	writeTarFile(t, layer, []tarEntry{
		{name: "etc/config/app.yaml", body: "wanted"},
		{name: "etcetera/config/app.yaml", body: "unwanted"},
	})

	err := ExtractBlob(layer, dest, ExtractBlobOpts{
		Context:  context.Background(),
		Includes: []string{"/etc/"},
		Logger:   zerolog.New(io.Discard),
	})
	require.NoError(t, err)

	require.FileExists(t, filepath.Join(dest, "etc", "config", "app.yaml"))
	require.NoFileExists(t, filepath.Join(dest, "etcetera", "config", "app.yaml"))
}

func TestExtractBlobKeepsRecreatedDirWhenWhiteoutComesLater(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "dist")

	layer1 := filepath.Join(root, "layer1.tar")
	writeTarFile(t, layer1, []tarEntry{
		{name: "dir/sub/lower.txt", body: "lower"},
	})

	layer2 := filepath.Join(root, "layer2.tar")
	writeTarFile(t, layer2, []tarEntry{
		{name: "dir/sub/", typeflag: tar.TypeDir},
		{name: "dir/.wh.sub"},
	})

	opts := ExtractBlobOpts{
		Context: context.Background(),
		Logger:  zerolog.New(io.Discard),
	}

	require.NoError(t, ExtractBlob(layer1, dest, opts))
	require.NoError(t, ExtractBlob(layer2, dest, opts))

	require.DirExists(t, filepath.Join(dest, "dir", "sub"))
	require.NoFileExists(t, filepath.Join(dest, "dir", "sub", "lower.txt"))
}

func TestWriteFileHonorsCanceledContext(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source.txt")
	dest := filepath.Join(root, "dest.txt")
	require.NoError(t, os.WriteFile(source, []byte("payload"), 0o644))

	info, err := os.Stat(source)
	require.NoError(t, err)

	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(nil)

	err = writeFile(ctx, dest, archives.FileInfo{
		FileInfo: info,
		Open: func() (fs.File, error) {
			return os.Open(source)
		},
	})
	require.ErrorIs(t, err, context.Canceled)

	data, readErr := os.ReadFile(dest)
	require.NoError(t, readErr)
	require.Empty(t, data)
}

type tarEntry struct {
	name     string
	body     string
	mode     int64
	typeflag byte
}

func writeTarFile(t *testing.T, filename string, entries []tarEntry) {
	t.Helper()

	f, err := os.Create(filename)
	require.NoError(t, err)
	defer f.Close()

	tw := tar.NewWriter(f)
	defer tw.Close()

	writeTarEntries(t, tw, entries)
}

func writeTarGzFile(t *testing.T, filename string, entries []tarEntry) {
	t.Helper()

	f, err := os.Create(filename)
	require.NoError(t, err)
	defer f.Close()

	gw := gzip.NewWriter(f)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	writeTarEntries(t, tw, entries)
}

func writeTarEntries(t *testing.T, tw *tar.Writer, entries []tarEntry) {
	t.Helper()

	for _, entry := range entries {
		mode := entry.mode
		if mode == 0 {
			mode = 0o644
			if entry.typeflag == tar.TypeDir {
				mode = 0o755
			}
		}

		typeflag := entry.typeflag
		if typeflag == 0 {
			typeflag = tar.TypeReg
		}

		size := int64(len(entry.body))
		if typeflag == tar.TypeDir {
			size = 0
		}

		hdr := &tar.Header{
			Name:     entry.name,
			Mode:     mode,
			Size:     size,
			Typeflag: typeflag,
		}
		require.NoError(t, tw.WriteHeader(hdr))
		if entry.body == "" || typeflag == tar.TypeDir {
			continue
		}
		_, err := tw.Write([]byte(entry.body))
		require.NoError(t, err)
	}
}
