package extractor

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mholt/archives"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// ExtractBlobOpts holds extract blob options
type ExtractBlobOpts struct {
	Context  context.Context
	Logger   zerolog.Logger
	Includes []string
}

func ExtractBlob(filename string, dest string, opts ExtractBlobOpts) error {
	opts.Logger.Info().Msgf("Extracting blob")

	dt, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer dt.Close()

	format, input, err := archives.Identify(opts.Context, filename, dt)
	if err != nil {
		opts.Logger.Warn().Msg("Blob format not recognized")
		return nil
	}
	opts.Logger.Debug().Msgf("Blob format %s detected", format.Extension())

	extractor, ok := format.(archives.Extractor)
	if !ok {
		// .gz is a special case, as it is a compressed tarball
		if format.Extension() == ".gz" {
			extractor = archives.Tar{}
			input, err = archives.Gz{}.OpenReader(input)
			if err != nil {
				return err
			}
		} else {
			return errors.Errorf("blob format not supported: %s", format.Extension())
		}
	}

	var pathsInArchive []string
	for _, inc := range opts.Includes {
		inc = strings.TrimPrefix(inc, "/")
		if len(inc) > 0 {
			pathsInArchive = append(pathsInArchive, inc)
		}
	}

	return extractor.Extract(opts.Context, input, func(ctx context.Context, f archives.FileInfo) error {
		if target, opaque, ok := whiteoutTarget(f.NameInArchive); ok {
			if !pathIntersects(pathsInArchive, target) {
				return nil
			}
			if opaque {
				opts.Logger.Debug().Msgf("Applying opaque whiteout %s", f.NameInArchive)
				return applyOpaqueWhiteout(dest, target)
			}
			opts.Logger.Debug().Msgf("Applying whiteout %s", f.NameInArchive)
			return removePath(filepath.Join(dest, filepath.FromSlash(target)))
		}

		if !fileIsIncluded(pathsInArchive, f.NameInArchive) {
			return nil
		}

		if f.IsDir() {
			opts.Logger.Trace().Msgf("Extracting %s", f.NameInArchive)
		} else {
			opts.Logger.Debug().Msgf("Extracting %s", f.NameInArchive)
		}

		outPath := filepath.Join(dest, filepath.FromSlash(f.NameInArchive))
		if err = os.MkdirAll(filepath.Dir(outPath), 0o700); err != nil {
			return err
		}

		switch {
		case f.IsDir():
			return os.MkdirAll(outPath, f.Mode())
		case f.Mode().IsRegular():
			return writeFile(ctx, outPath, f)
		case f.Mode()&fs.ModeSymlink != 0:
			return writeSymlink(ctx, outPath, f)
		default:
			return errors.Errorf("cannot handle file mode: %v", f.Mode())
		}
	})
}

func fileIsIncluded(filenameList []string, filename string) bool {
	// include all files if there is no specific list
	if len(filenameList) == 0 {
		return true
	}
	for _, fn := range filenameList {
		// exact matches are of course included
		if filename == fn {
			return true
		}
		// also consider the file included if its parent folder/path is in the list
		if strings.HasPrefix(filename, strings.TrimSuffix(fn, "/")+"/") {
			return true
		}
	}
	return false
}

func pathIntersects(filenameList []string, filename string) bool {
	// include all paths if there is no specific list
	if len(filenameList) == 0 {
		return true
	}
	for _, fn := range filenameList {
		trimmed := strings.TrimSuffix(fn, "/")
		if filename == trimmed {
			return true
		}
		if strings.HasPrefix(filename, trimmed+"/") {
			return true
		}
		if strings.HasPrefix(trimmed, filename+"/") {
			return true
		}
	}
	return false
}

func whiteoutTarget(filename string) (target string, opaque bool, ok bool) {
	dir, base := path.Split(path.Clean(filename))
	switch {
	case base == ".wh..wh..opq":
		return strings.TrimSuffix(dir, "/"), true, true
	case strings.HasPrefix(base, ".wh."):
		return path.Join(dir, strings.TrimPrefix(base, ".wh.")), false, true
	default:
		return "", false, false
	}
}

func applyOpaqueWhiteout(dest string, target string) error {
	targetPath := filepath.Join(dest, filepath.FromSlash(target))
	entries, err := os.ReadDir(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, entry := range entries {
		if err := removePath(filepath.Join(targetPath, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

func removePath(path string) error {
	err := os.RemoveAll(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func writeFile(ctx context.Context, path string, f archives.FileInfo) error {
	r, err := f.Open()
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(path)
	if err != nil {
		return err
	}
	defer w.Close()

	err = w.Chmod(f.Mode())
	if err != nil {
		return err
	}

	_, err = io.Copy(w, readerContext(ctx, r))
	return err
}

func writeSymlink(_ context.Context, path string, f archives.FileInfo) error {
	if f.LinkTarget == "" {
		return errors.Errorf("symlink target is empty for %s", f.Name())
	}

	_, err := os.Lstat(path)
	if err == nil {
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}

	return os.Symlink(f.LinkTarget, path)
}

type reader struct {
	ctx context.Context
	r   io.Reader
}

func readerContext(ctx context.Context, r io.Reader) io.Reader {
	return reader{ctx, r}
}

func (r reader) Read(p []byte) (int, error) {
	err := r.ctx.Err()
	if err != nil {
		return 0, err
	}
	n, err := r.r.Read(p)
	if err != nil {
		return n, err
	}
	return n, r.ctx.Err()
}
