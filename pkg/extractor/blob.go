package extractor

import (
	"context"
	"io"
	"io/fs"
	"os"
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
		if !fileIsIncluded(pathsInArchive, f.NameInArchive) {
			return nil
		}

		if f.FileInfo.IsDir() {
			opts.Logger.Trace().Msgf("Extracting %s", f.NameInArchive)
		} else {
			opts.Logger.Debug().Msgf("Extracting %s", f.NameInArchive)
		}

		path := filepath.Join(dest, f.NameInArchive)
		if err = os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			return err
		}

		switch {
		case f.FileInfo.IsDir():
			return os.MkdirAll(path, f.Mode())
		case f.FileInfo.Mode().IsRegular():
			return writeFile(ctx, path, f)
		case f.FileInfo.Mode()&fs.ModeSymlink != 0:
			return writeSymlink(ctx, path, f)
		default:
			return errors.Errorf("cannot handle file mode: %v", f.FileInfo.Mode())
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
