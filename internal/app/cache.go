package app

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/docker"
	dockerconfig "github.com/containers/image/v5/pkg/docker/config"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
	"github.com/crazy-max/undock/pkg/image"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func (c *Undock) cacheSource(logger zerolog.Logger, src string) ([]byte, string, error) {
	srcCtx, dgst, err := c.srcCtx(src, c.cli.Insecure)
	if err != nil {
		return nil, "", err
	}
	srcRef, err := alltransports.ParseImageName(fmt.Sprintf("docker://%s", src))
	if err != nil {
		return nil, "", errors.Wrapf(err, "invalid container image source %s", src)
	}

	cachedir := filepath.Join(c.cli.CacheDir, dgst.Encoded())

	dstRef, err := alltransports.ParseImageName(fmt.Sprintf("oci:%s", cachedir))
	if err != nil {
		return nil, "", errors.Wrapf(err, "invalid oci cache folder %s", cachedir)
	}
	dstCtx, err := c.dstCtx(cachedir)
	if err != nil {
		return nil, "", err
	}

	imageSelection := copy.CopySystemImage
	if c.cli.All {
		imageSelection = copy.CopyAllImages
	}

	policyContext, err := signature.NewPolicyContext(&signature.Policy{Default: []signature.PolicyRequirement{signature.NewPRInsecureAcceptAnything()}})
	if err != nil {
		return nil, "", err
	}
	defer policyContext.Destroy()

	manblob, err := copy.Image(c.ctx, policyContext, dstRef, srcRef, &copy.Options{
		ReportWriter:                          &progressWriter{logger: logger},
		SourceCtx:                             srcCtx,
		DestinationCtx:                        dstCtx,
		ImageListSelection:                    imageSelection,
		OptimizeDestinationImageAlreadyExists: true,
	})

	return manblob, cachedir, err
}

func (c *Undock) srcCtx(name string, insecure bool) (*types.SystemContext, *digest.Digest, error) {
	img, err := image.Parse(name)
	if err != nil {
		return nil, nil, err
	}

	auth, err := dockerconfig.GetCredentials(nil, img.Domain)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "cannot find docker credentials")
	}

	sysCtx := &types.SystemContext{
		DockerAuthConfig:                  &auth,
		DockerDaemonInsecureSkipTLSVerify: insecure,
		DockerInsecureSkipTLSVerify:       types.NewOptionalBool(insecure),
		DockerRegistryUserAgent:           c.meta.UserAgent,
		OSChoice:                          c.platform.OS,
		ArchitectureChoice:                c.platform.Architecture,
		VariantChoice:                     c.platform.Variant,
		BlobInfoCacheDir:                  filepath.Join(c.cli.CacheDir, "blobs"),
	}

	ref, err := image.Reference(img.String())
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot parse reference")
	}

	dgst, err := docker.GetDigest(c.ctx, sysCtx, ref)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot get image digest from HEAD request")
	}

	return sysCtx, &dgst, nil
}

func (c *Undock) dstCtx(name string) (*types.SystemContext, error) {
	return &types.SystemContext{
		DirForceDecompress: true,
		BlobInfoCacheDir:   filepath.Join(c.cli.CacheDir, "blobs"),
	}, nil
}

type progressWriter struct {
	logger zerolog.Logger
	writer io.Writer
}

func (w *progressWriter) Write(p []byte) (n int, err error) {
	w.logger.Info().Msgf("%s", strings.TrimSpace(string(p)))
	return len(p), nil
}
