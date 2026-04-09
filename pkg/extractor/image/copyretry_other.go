//go:build !windows

package image

import (
	"go.podman.io/image/v5/copy"
	"go.podman.io/image/v5/signature"
	"go.podman.io/image/v5/types"
)

func (c *Client) copyCachedImage(policyContext *signature.PolicyContext, dstRef, srcRef types.ImageReference, opts *copy.Options) ([]byte, error) {
	return copy.Image(c.ctx, policyContext, dstRef, srcRef, opts)
}
