package image

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/containerd/platforms"
	"github.com/containers/image/v5/manifest"
	"github.com/crazy-max/undock/pkg/extractor"
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

// Client represents an active image extractor object
type Client struct {
	*extractor.Client
	ctx    context.Context
	opts   Options
	logger zerolog.Logger
}

// Options represents image extractor options
type Options struct {
	// Source image reference
	Source string
	// Platform to enforce for Source image
	Platform ocispecs.Platform
	// Includes a subset of files/dirs from the Source image
	Includes []string
	// All extracts all architectures if Source image is a manifest list
	All bool

	// Dist folder
	Dist string
	// Wrap merges output in Dist folder for a manifest list
	Wrap bool

	// RegistryInsecure allows contacting the registry or docker daemon over
	// HTTP, or HTTPS with failed TLS verification
	RegistryInsecure bool
	// RegistryUserAgent is the User-Agent string to send to the registry
	RegistryUserAgent string

	// CacheDir is the directory where the cache is stored
	CacheDir string
}

// New creates new image extractor instance
func New(ctx context.Context, opts Options) (*extractor.Client, error) {
	logger := log.With().Str("src", opts.Source).Logger()

	if opts.Platform.OS == "" || opts.Platform.Architecture == "" {
		opts.Platform = platforms.DefaultSpec()
		logger.Warn().Msgf("platform not set, using %s", platforms.Format(opts.Platform))
	}

	datadir := opts.CacheDir
	if len(datadir) == 0 {
		datadir = os.Getenv("XDG_DATA_HOME")
		if len(datadir) == 0 {
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, errors.Wrap(err, "failed to get home directory")
			}
			datadir = filepath.Join(home, ".local", "share")
		}
	}

	opts.CacheDir = filepath.Join(datadir, "undock", "cache")
	if err := os.MkdirAll(opts.CacheDir, 0700); err != nil {
		return nil, errors.Wrapf(err, "failed to create cache directory %q", opts.CacheDir)
	}

	return &extractor.Client{
		Handler: &Client{
			ctx:    ctx,
			opts:   opts,
			logger: logger,
		},
	}, nil
}

// Type returns the extractor type
func (c *Client) Type() string {
	return "image"
}

// Extract extracts a registry image
func (c *Client) Extract() error {
	c.logger.Info().Msg("Extracting source")

	manblob, cachedir, err := c.cacheSource(c.opts.Source)
	if err != nil {
		return errors.Wrap(err, "cannot cache source")
	}

	type manifestEntry struct {
		platform ocispecs.Platform
		manifest *manifest.OCI1
	}

	var mans []manifestEntry

	mtype := manifest.GuessMIMEType(manblob)
	if mtype == ocispecs.MediaTypeImageManifest {
		man, err := manifest.OCI1FromManifest(manblob)
		if err != nil {
			return errors.Wrap(err, "cannot create OCI manifest instance from blob")
		}
		mans = append(mans, manifestEntry{
			platform: c.opts.Platform,
			manifest: man,
		})
	} else if mtype == ocispecs.MediaTypeImageIndex {
		ocindex, err := manifest.OCI1IndexFromManifest(manblob)
		if err != nil {
			return errors.Wrap(err, "cannot create OCI manifest index instance from blob")
		}
		for _, m := range ocindex.Manifests {
			mblob, err := os.ReadFile(path.Join(cachedir, "blobs", m.Digest.Algorithm().String(), m.Digest.Hex()))
			if err != nil {
				return errors.Wrapf(err, "cannot read OCI manifest JSON for platform %s", platforms.Format(*m.Platform))
			}
			man, err := manifest.OCI1FromManifest(mblob)
			if err != nil {
				return errors.Wrap(err, "cannot create OCI manifest instance from blob")
			}
			mans = append(mans, manifestEntry{
				platform: *m.Platform,
				manifest: man,
			})
		}
	}

	eg, _ := errgroup.WithContext(c.ctx)
	for _, me := range mans {
		func(me manifestEntry) {
			eg.Go(func() error {
				dest := c.opts.Dist
				if !c.opts.Wrap && len(mans) > 1 {
					dest = path.Join(c.opts.Dist, fmt.Sprintf("%s_%s%s", me.platform.OS, me.platform.Architecture, me.platform.Variant))
				}
				for _, layer := range me.manifest.LayerInfos() {
					sublogger := c.logger.With().
						Str("platform", platforms.Format(me.platform)).
						Str("media-type", layer.MediaType).
						Str("blob", layer.Digest.String()).Logger()
					if err = extractor.ExtractBlob(path.Join(cachedir, "blobs", layer.Digest.Algorithm().String(), layer.Digest.Hex()), dest, extractor.ExtractBlobOpts{
						Context:  c.ctx,
						Logger:   sublogger,
						Includes: c.opts.Includes,
					}); err != nil {
						return err
					}
				}
				return nil
			})
		}(me)
	}

	return eg.Wait()
}

//nolint:unused
func formatReference(source string) (string, string) {
	scheme := sourceScheme(source)
	switch scheme {
	case "":
		return "docker://" + source, "docker"
	case "docker":
		return source, scheme
	default:
		return strings.Replace(source, scheme+"://", scheme+":", 1), scheme
	}
}

//nolint:unused
func sourceScheme(source string) string {
	schemes := []string{"containers-storage", "docker", "docker-archive", "docker-daemon", "oci", "oci-archive", "ostree"}
	for _, scheme := range schemes {
		if strings.HasPrefix(source, scheme+"://") {
			return scheme
		}
	}
	return ""
}
