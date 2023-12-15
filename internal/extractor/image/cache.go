package image

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
	dockercli "github.com/crazy-max/undock/pkg/docker"
	"github.com/crazy-max/undock/pkg/image"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func (c *Client) cacheSource(src string) ([]byte, string, error) {
	srcCtx, srcObj, err := c.srcCtx(src, c.cli.Insecure)
	if err != nil {
		return nil, "", errors.Wrap(err, "cannot create source context")
	}
	srcRef, err := srcObj.Reference()
	if err != nil {
		return nil, "", errors.Wrapf(err, "cannot parse reference '%s'", srcObj.String())
	}

	var cacheDigest string
	switch srcObj.Scheme() {
	case "docker":
		dockerRef, err := image.DockerReference(strings.TrimPrefix(src, "docker://"))
		if err != nil {
			return nil, "", errors.Wrap(err, "cannot parse docker reference")
		}
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
		if img, err := dcli.ImageInspectWithRaw(strings.TrimPrefix(src, "docker-daemon://")); err == nil {
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

	cachedir := filepath.Join(c.cli.CacheDir, cacheDigest)
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
	if c.cli.All {
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

func (c *Client) srcCtx(name string, insecure bool) (*types.SystemContext, *Source, error) {
	sysCtx := &types.SystemContext{
		DockerDaemonInsecureSkipTLSVerify: insecure,
		DockerInsecureSkipTLSVerify:       types.NewOptionalBool(insecure),
		DockerRegistryUserAgent:           c.meta.UserAgent,
		OSChoice:                          c.platform.OS,
		ArchitectureChoice:                c.platform.Architecture,
		VariantChoice:                     c.platform.Variant,
		BlobInfoCacheDir:                  filepath.Join(c.cli.CacheDir, "blobs"),
	}
	return sysCtx, NewSource(name), nil
}

func (c *Client) dstCtx(cacheDir string) (*types.SystemContext, error) {
	return &types.SystemContext{
		DirForceDecompress: true,
		BlobInfoCacheDir:   filepath.Join(cacheDir, "blobs"),
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
