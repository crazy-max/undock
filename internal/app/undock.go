package app

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/containerd/containerd/platforms"
	"github.com/containers/image/v5/manifest"
	"github.com/crazy-max/undock/internal/config"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

// Undock represents an active undock object
type Undock struct {
	ctx      context.Context
	meta     config.Meta
	cli      config.Cli
	platform specs.Platform
}

// New creates new undock instance
func New(meta config.Meta, cli config.Cli) (*Undock, error) {
	datadir := cli.CacheDir
	if len(datadir) == 0 {
		datadir = os.Getenv("XDG_DATA_HOME")
		if len(datadir) == 0 {
			home := os.Getenv("HOME")
			if len(home) == 0 {
				return nil, errors.New("neither XDG_DATA_HOME nor HOME was set non-empty")
			}
			datadir = filepath.Join(home, ".local", "share")
		}
	}
	cli.CacheDir = filepath.Join(datadir, "undock", "cache")
	if err := os.MkdirAll(cli.CacheDir, 0700); err != nil {
		return nil, errors.Wrapf(err, "failed to create cache directory %q", cli.CacheDir)
	}

	platform := platforms.DefaultSpec()
	if len(cli.Platform) > 0 {
		if p, err := platforms.Parse(cli.Platform); err != nil {
			return nil, errors.Wrapf(err, "invalid platform %q", cli.Platform)
		} else {
			platform = p
		}
	}

	return &Undock{
		ctx:      context.Background(),
		meta:     meta,
		cli:      cli,
		platform: platform,
	}, nil
}

// Start starts undock
func (c *Undock) Start() error {
	if _, err := os.Stat(c.cli.Dist); err == nil && c.cli.RmDist {
		if err := os.RemoveAll(c.cli.Dist); err != nil {
			return errors.Wrapf(err, "failed to remove dist folder %q", c.cli.Dist)
		}
	}
	if err := os.MkdirAll(c.cli.Dist, 0700); err != nil {
		return errors.Wrapf(err, "failed to create dist folder %q", c.cli.Dist)
	}

	logger := log.With().Str("src", c.cli.Source).Logger()

	logger.Info().Msg("Extracting source image")
	manblob, cachedir, err := c.cacheSource(logger, c.cli.Source)
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
				if err := os.MkdirAll(dest, 0700); err != nil {
					return errors.Wrapf(err, "failed to create destination folder %q", dest)
				}
				for _, layer := range mane.manifest.LayerInfos() {
					sublogger := logger.With().Str("platform", platforms.Format(mane.platform)).Str("blob", layer.Digest.String()).Logger()
					sublogger.Info().Msgf("Extracting blob")
					if err = c.extract(sublogger, path.Join(cachedir, "blobs", layer.Digest.Algorithm().String(), layer.Digest.Hex()), dest); err != nil {
						return err
					}
				}
				return nil
			})
		}(mane)
	}

	return eg.Wait()
}

// Close closes undock
func (c *Undock) Close() {
	// noop
}
