package app

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/containerd/containerd/platforms"
	"github.com/crazy-max/undock/pkg/config"
	ximage "github.com/crazy-max/undock/pkg/extractor/image"
	"github.com/crazy-max/undock/pkg/image"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
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

	xcli, err := ximage.New(c.ctx, c.meta, c.cli, c.platform)
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
