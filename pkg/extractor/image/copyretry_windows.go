//go:build windows

package image

import (
	"os"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/pkg/errors"
	"go.podman.io/image/v5/copy"
	"go.podman.io/image/v5/signature"
	"go.podman.io/image/v5/types"
	"golang.org/x/sys/windows"
)

const (
	copyImageAttempts = 5
	copyImageBackoff  = 100 * time.Millisecond
	copyImageMaxDelay = 2 * time.Second
)

func (c *Client) copyCachedImage(policyContext *signature.PolicyContext, dstRef, srcRef types.ImageReference, opts *copy.Options) ([]byte, error) {
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = copyImageBackoff
	bo.RandomizationFactor = 0.5
	bo.Multiplier = 2
	bo.MaxInterval = copyImageMaxDelay
	var attempt int
	return backoff.Retry(c.ctx, func() ([]byte, error) {
		attempt++
		manblob, err := copy.Image(c.ctx, policyContext, dstRef, srcRef, opts)
		if err == nil || isRetryableCopyImageError(err) {
			return manblob, err
		}
		return nil, backoff.Permanent(err)
	},
		backoff.WithBackOff(bo),
		backoff.WithMaxTries(copyImageAttempts),
		backoff.WithNotify(func(err error, delay time.Duration) {
			c.logger.Warn().Err(err).Int("attempt", attempt).Int("max-attempts", copyImageAttempts).Dur("backoff", delay).Msg("Retrying cache copy after transient Windows filesystem error")
		}),
	)
}

func isRetryableCopyImageError(err error) bool {
	var linkErr *os.LinkError
	if errors.As(err, &linkErr) {
		return isRetryableWindowsRenameError(linkErr.Err)
	}
	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		return isRetryableWindowsRenameError(pathErr.Err)
	}
	return false
}

func isRetryableWindowsRenameError(err error) bool {
	return errors.Is(err, windows.ERROR_ACCESS_DENIED) || errors.Is(err, windows.ERROR_SHARING_VIOLATION) || errors.Is(err, windows.ERROR_LOCK_VIOLATION)
}
