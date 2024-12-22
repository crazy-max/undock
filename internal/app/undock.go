package app

import (
	"context"
	"os"
	"strings"

	"github.com/containerd/containerd/platforms"
	"github.com/crazy-max/undock/internal/config"
	ximage "github.com/crazy-max/undock/pkg/extractor/image"
	"github.com/crazy-max/undock/pkg/image"
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

// Undock represents an active undock object
type Undock struct {
	ctx      context.Context
	meta     config.Meta
	cli      config.Cli
	platform ocispecs.Platform
}

// New creates new undock instance
func New(meta config.Meta, cli config.Cli) (*Undock, error) {
	platform := platforms.DefaultSpec()
	if len(cli.Platform) > 0 {
		var err error
		if platform, err = platforms.Parse(cli.Platform); err != nil {
			return nil, errors.Wrapf(err, "invalid platform %q", cli.Platform)
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

	if ok, err := validateScheme(c.cli.Source); err != nil {
		return err
	} else if !ok {
		return errors.Errorf("unsupported source %q", c.cli.Source)
	}

	xcli, err := ximage.New(c.ctx, ximage.Options{
		Source:   c.cli.Source,
		Platform: c.platform,
		Includes: c.cli.Includes,
		All:      c.cli.All,

		Dist: c.cli.Dist,
		Wrap: c.cli.Wrap,

		RegistryInsecure:  c.cli.Insecure,
		RegistryUserAgent: c.meta.UserAgent,

		CacheDir: c.cli.CacheDir,
	})
	if err != nil {
		return err
	}
	return xcli.Extract()
}

func validateScheme(source string) (bool, error) {
	schemes := []string{"containers-storage", "docker", "docker-archive", "docker-daemon", "oci", "oci-archive", "ostree"}
	for _, scheme := range schemes {
		if strings.HasPrefix(source, scheme+"://") {
			return true, nil
		}
	}
	_, err := image.Reference(source)
	return err == nil, err
}

// Close closes undock
func (c *Undock) Close() {
	// noop
}
