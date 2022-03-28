package image

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/containerd/containerd/platforms"
	"github.com/containers/image/v5/manifest"
	"github.com/crazy-max/undock/internal/config"
	"github.com/crazy-max/undock/internal/extractor"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

// Client represents an active image extractor object
type Client struct {
	*extractor.Client
	ctx      context.Context
	meta     config.Meta
	cli      config.Cli
	platform specs.Platform
	logger   zerolog.Logger
}

// New creates new image extractor instance
func New(ctx context.Context, meta config.Meta, cli config.Cli, platform specs.Platform) (*extractor.Client, error) {
	return &extractor.Client{
		Handler: &Client{
			ctx:      ctx,
			meta:     meta,
			cli:      cli,
			platform: platform,
			logger:   log.With().Str("src", cli.Source).Logger(),
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

	manblob, cachedir, err := c.cacheSource(c.cli.Source)
	if err != nil {
		return errors.Wrap(err, "cannot cache source")
	}

	type manifestEntry struct {
		platform specs.Platform
		manifest *manifest.OCI1
	}

	var mans []manifestEntry

	mtype := manifest.GuessMIMEType(manblob)
	if mtype == specs.MediaTypeImageManifest {
		man, err := manifest.OCI1FromManifest(manblob)
		if err != nil {
			return errors.Wrap(err, "cannot create OCI manifest instance from blob")
		}
		mans = append(mans, manifestEntry{
			platform: c.platform,
			manifest: man,
		})
	} else if mtype == specs.MediaTypeImageIndex {
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
	for _, mane := range mans {
		func(mane manifestEntry) {
			eg.Go(func() error {
				dest := c.cli.Dist
				if !c.cli.Wrap && len(mans) > 1 {
					dest = path.Join(c.cli.Dist, fmt.Sprintf("%s_%s%s", mane.platform.OS, mane.platform.Architecture, mane.platform.Variant))
				}
				if err := os.MkdirAll(dest, 0o700); err != nil {
					return errors.Wrapf(err, "failed to create destination folder %q", dest)
				}
				for _, layer := range mane.manifest.LayerInfos() {
					sublogger := c.logger.With().Str("platform", platforms.Format(mane.platform)).Str("blob", layer.Digest.String()).Logger()
					if err = extractor.ExtractBlob(path.Join(cachedir, "blobs", layer.Digest.Algorithm().String(), layer.Digest.Hex()), dest, extractor.ExtractBlobOpts{
						Context:  c.ctx,
						Logger:   sublogger,
						Includes: c.cli.Includes,
					}); err != nil {
						return err
					}
				}
				return nil
			})
		}(mane)
	}

	return eg.Wait()
}

//nolint:deadcode
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

func sourceScheme(source string) string {
	schemes := []string{"containers-storage", "docker", "docker-archive", "docker-daemon", "oci", "oci-archive", "ostree"}
	for _, scheme := range schemes {
		if strings.HasPrefix(source, scheme+"://") {
			return scheme
		}
	}
	return ""
}
