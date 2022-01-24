package app

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func (c *Undock) extract(logger zerolog.Logger, filename string, dest string) error {
	var r io.ReadCloser
	dt, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer dt.Close()

	format, err := archiver.Identify(filename, dt)
	if err != nil {
		return err
	}
	logger.Debug().Msgf("Blob format %s detected", format.Name())

	var u archiver.Extractor
	var d archiver.Decompressor
	switch format.Name() {
	case ".zip":
		u = archiver.Zip{}
	case ".tar.gz":
		u = archiver.Tar{}
		d = archiver.Gz{}
	case ".tar.xz":
		u = archiver.Tar{}
		d = archiver.Xz{}
	case ".zst":
		u = archiver.Tar{}
		d = archiver.Zstd{}
	default:
		return errors.Errorf("blob format not supported: %s", format.Name())
	}

	if d == nil {
		r = dt
	} else {
		r, err = d.OpenReader(dt)
		if err != nil {
			return err
		}
	}

	var pathsInArchive []string
	if len(c.cli.Includes) > 0 {
		for _, inc := range c.cli.Includes {
			inc = strings.TrimPrefix(inc, "/")
			if len(inc) > 0 {
				pathsInArchive = append(pathsInArchive, inc)
			}
		}
	}
	if len(pathsInArchive) == 0 {
		pathsInArchive = nil
	}

	err = u.Extract(c.ctx, r, pathsInArchive, func(ctx context.Context, f archiver.File) error {
		if f.FileInfo.IsDir() {
			logger.Trace().Msgf("Extracting %s", f.NameInArchive)
		} else {
			logger.Debug().Msgf("Extracting %s", f.NameInArchive)
		}

		path := filepath.Join(dest, f.NameInArchive)
		if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
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
			return fmt.Errorf("cannot handle file mode: %v", f.FileInfo.Mode())
		}
	})

	return err
}

func writeFile(ctx context.Context, path string, f archiver.File) error {
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

func writeSymlink(ctx context.Context, path string, f archiver.File) error {
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
