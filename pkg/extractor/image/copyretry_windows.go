//go:build windows

package image

import (
	"errors"
	"math/rand/v2"
	"os"
	"time"

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
	var (
		manblob []byte
		err     error
	)

	for attempt := 1; attempt <= copyImageAttempts; attempt++ {
		manblob, err = copy.Image(c.ctx, policyContext, dstRef, srcRef, opts)
		if err == nil || !isRetryableCopyImageError(err) || attempt == copyImageAttempts {
			return manblob, err
		}

		delay := copyImageBackoff << min(attempt-1, 3)
		delay = min(delay, copyImageMaxDelay)
		delay = delay/2 + time.Duration(rand.Int64N(int64(delay/2)+1))

		c.logger.Warn().Err(err).Int("attempt", attempt).Int("max-attempts", copyImageAttempts).Dur("backoff", delay).Msg("Retrying cache copy after transient Windows filesystem error")

		timer := time.NewTimer(delay)
		select {
		case <-c.ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return nil, c.ctx.Err()
		case <-timer.C:
		}
	}

	return manblob, err
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
