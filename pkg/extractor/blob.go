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

const (
	whiteoutPrefix     = ".wh."
	whiteoutMetaPrefix = ".wh..wh."
	whiteoutOpaqueDir  = ".wh..wh..opq"
	whiteoutLinkDir    = ".wh..wh.plnk"
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

	createdInLayer := map[string]struct{}{}

	return extractor.Extract(opts.Context, input, func(ctx context.Context, f archives.FileInfo) error {
		entryName, err := normalizeArchivePath(f.NameInArchive)
		if err != nil {
			return err
		}
		if shouldSkipReservedWhiteoutPath(entryName) {
			opts.Logger.Debug().Msgf("Skipping reserved whiteout metadata %s", f.NameInArchive)
			return nil
		}

		if target, opaque, ok := whiteoutTarget(entryName); ok {
			if !pathIntersects(pathsInArchive, target) {
				return nil
			}
			if opaque {
				opts.Logger.Debug().Msgf("Applying opaque whiteout %s", f.NameInArchive)
				return applyOpaqueWhiteout(dest, target, createdInLayer)
			}
			opts.Logger.Debug().Msgf("Applying whiteout %s", f.NameInArchive)
			return applyWhiteout(dest, target, createdInLayer)
		}

		if !fileIsIncluded(pathsInArchive, entryName) {
			return nil
		}

		if f.IsDir() {
			opts.Logger.Trace().Msgf("Extracting %s", f.NameInArchive)
		} else {
			opts.Logger.Debug().Msgf("Extracting %s", f.NameInArchive)
		}

		outPath := filepath.Join(dest, filepath.FromSlash(entryName))
		if err = os.MkdirAll(filepath.Dir(outPath), 0o700); err != nil {
			return err
		}

		switch {
		case f.IsDir():
			err = os.MkdirAll(outPath, f.Mode())
		case f.Mode().IsRegular():
			err = writeFile(ctx, outPath, f)
		case f.Mode()&fs.ModeSymlink != 0:
			err = writeSymlink(ctx, outPath, f)
		default:
			return errors.Errorf("cannot handle file mode: %v", f.Mode())
		}

		if err != nil {
			return err
		}
		createdInLayer[entryName] = struct{}{}
		return nil
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

func normalizeArchivePath(filename string) (string, error) {
	filename = strings.ReplaceAll(filename, "\\", "/")
	cleaned := path.Clean(strings.TrimPrefix(filename, "/"))
	if cleaned == "." {
		return ".", nil
	}
	if cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return "", errors.Errorf("archive path %q resolves outside destination", filename)
	}
	return cleaned, nil
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
	case base == whiteoutOpaqueDir:
		return strings.TrimSuffix(dir, "/"), true, true
	case strings.HasPrefix(base, whiteoutMetaPrefix):
		return "", false, false
	case strings.HasPrefix(base, whiteoutPrefix):
		return path.Join(dir, strings.TrimPrefix(base, whiteoutPrefix)), false, true
	default:
		return "", false, false
	}
}

func shouldSkipReservedWhiteoutPath(filename string) bool {
	for _, segment := range strings.Split(filename, "/") {
		switch {
		case segment == whiteoutOpaqueDir:
			return false
		case segment == whiteoutLinkDir:
			return true
		case strings.HasPrefix(segment, whiteoutMetaPrefix):
			return true
		}
	}
	return false
}

func applyOpaqueWhiteout(dest string, target string, createdInLayer map[string]struct{}) error {
	return removePathPreservingCurrentLayer(dest, target, createdInLayer, false)
}

func applyWhiteout(dest string, target string, createdInLayer map[string]struct{}) error {
	return removePathPreservingCurrentLayer(dest, target, createdInLayer, true)
}

func removePathPreservingCurrentLayer(dest string, target string, createdInLayer map[string]struct{}, removeSelf bool) error {
	targetPath := filepath.Join(dest, filepath.FromSlash(target))
	info, err := os.Lstat(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	protectSelf, protectedChildren := protectedPathsInCurrentLayer(createdInLayer, target)

	if !info.IsDir() {
		if protectSelf {
			return nil
		}
		return removePath(targetPath)
	}

	if removeSelf && !protectSelf && len(protectedChildren) == 0 {
		return removePath(targetPath)
	}

	entries, err := os.ReadDir(targetPath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if _, ok := protectedChildren[entry.Name()]; ok {
			continue
		}
		if err := removePath(filepath.Join(targetPath, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

func protectedPathsInCurrentLayer(createdInLayer map[string]struct{}, target string) (bool, map[string]struct{}) {
	protectedChildren := make(map[string]struct{})
	targetPrefix := strings.TrimSuffix(target, "/") + "/"
	var protectSelf bool

	for created := range createdInLayer {
		switch {
		case created == target:
			protectSelf = true
		case strings.HasPrefix(created, targetPrefix):
			rest := strings.TrimPrefix(created, targetPrefix)
			child, _, _ := strings.Cut(rest, "/")
			if child != "" {
				protectedChildren[child] = struct{}{}
			}
		}
	}

	return protectSelf, protectedChildren
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
