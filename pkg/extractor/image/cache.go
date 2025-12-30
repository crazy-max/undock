package image

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	dockercli "github.com/crazy-max/undock/pkg/docker"
	"github.com/crazy-max/undock/pkg/image"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.podman.io/image/v5/copy"
	"go.podman.io/image/v5/docker"
	"go.podman.io/image/v5/docker/reference"
	"go.podman.io/image/v5/pkg/docker/config"
	"go.podman.io/image/v5/signature"
	"go.podman.io/image/v5/transports/alltransports"
	"go.podman.io/image/v5/types"
)

func (c *Client) cacheSource(src string) ([]byte, string, error) {
	srcObj := NewSource(src)
	srcRef, err := srcObj.Reference()
	if err != nil {
		return nil, "", errors.Wrapf(err, "cannot parse reference '%s'", srcObj.String())
	}

	var dockerRef types.ImageReference
	var dockerAuth types.DockerAuthConfig
	if srcObj.Scheme() == "docker" {
		dockerRef, err = image.DockerReference(strings.TrimPrefix(src, "docker://"))
		if err != nil {
			return nil, "", errors.Wrap(err, "cannot parse docker reference")
		}
		dockerAuth, err = config.GetCredentials(nil, reference.Domain(dockerRef.DockerReference()))
		if err != nil {
			c.logger.Warn().Msg("cannot retrieve Docker credentials")
		}
	}

	srcCtx, err := c.srcCtx(&dockerAuth, c.opts.RegistryInsecure)
	if err != nil {
		return nil, "", errors.Wrap(err, "cannot create source context")
	}

	var cacheDigest string
	switch srcObj.Scheme() {
	case "docker":
		if dgst, err := docker.GetDigest(c.ctx, srcCtx, dockerRef); err == nil {
			cacheDigest = srcObj.Scheme() + "-" + dgst.Encoded()
		} else {
			return nil, "", errors.Wrap(err, "cannot get docker reference digest")
		}
	case "docker-daemon":
		dcli, err := dockercli.New(c.ctx)
		if err != nil {
			return nil, "", err
		}
		if img, err := dcli.ImageInspect(strings.TrimPrefix(src, "docker-daemon://")); err == nil {
			cacheDigest = srcObj.Scheme() + "-" + strings.TrimPrefix(img.ID, "sha256:")
		} else {
			return nil, "", err
		}
	default:
		// TODO: Find a proper way to create a cache fingerprint. Best effort atm with transport.
		srcHash := sha256.New()
		srcHash.Write([]byte(srcRef.StringWithinTransport()))
		cacheDigest = srcObj.Scheme() + "-" + hex.EncodeToString(srcHash.Sum(nil))
	}

	cachedir := filepath.Join(c.opts.CacheDir, cacheDigest)
	c.logger.Info().Msgf("Computed cache digest %s", cacheDigest)

	dstRef, err := alltransports.ParseImageName(fmt.Sprintf("oci:%s", cachedir))
	if err != nil {
		return nil, "", errors.Wrapf(err, "invalid oci cache folder %s", cachedir)
	}
	dstCtx, err := c.dstCtx(cachedir)
	if err != nil {
		return nil, "", err
	}

	imageSelection := copy.CopySystemImage
	if c.opts.All {
		imageSelection = copy.CopyAllImages
	}

	policyContext, err := signature.NewPolicyContext(&signature.Policy{Default: []signature.PolicyRequirement{signature.NewPRInsecureAcceptAnything()}})
	if err != nil {
		return nil, "", err
	}
	defer policyContext.Destroy() //nolint:errcheck

	manblob, err := copy.Image(c.ctx, policyContext, dstRef, srcRef, &copy.Options{
		ReportWriter:                          &progressWriter{logger: c.logger},
		SourceCtx:                             srcCtx,
		DestinationCtx:                        dstCtx,
		ImageListSelection:                    imageSelection,
		OptimizeDestinationImageAlreadyExists: true,
	})

	return manblob, cachedir, err
}

func (c *Client) srcCtx(auth *types.DockerAuthConfig, insecure bool) (*types.SystemContext, error) {
	sysCtx := &types.SystemContext{
		DockerAuthConfig:                  auth,
		DockerDaemonInsecureSkipTLSVerify: insecure,
		DockerInsecureSkipTLSVerify:       types.NewOptionalBool(insecure),
		DockerRegistryUserAgent:           c.opts.RegistryUserAgent,
		OSChoice:                          c.opts.Platform.OS,
		ArchitectureChoice:                c.opts.Platform.Architecture,
		VariantChoice:                     c.opts.Platform.Variant,
		BlobInfoCacheDir:                  filepath.Join(c.opts.CacheDir, "blobs"),
	}
	return sysCtx, nil
}

func (c *Client) dstCtx(_ string) (*types.SystemContext, error) {
	return &types.SystemContext{
		DirForceDecompress: true,
		BlobInfoCacheDir:   filepath.Join(c.opts.CacheDir, "blobs"),
	}, nil
}

type progressWriter struct {
	logger zerolog.Logger
	writer io.Writer //nolint:unused
}

func (w *progressWriter) Write(p []byte) (n int, err error) {
	w.logger.Info().Msgf("%s", strings.TrimSpace(string(p)))
	return len(p), nil
}
