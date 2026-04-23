//go:build windows

package image

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/windows"
)

func TestIsRetryableWindowsRenameError(t *testing.T) {
	assert.True(t, isRetryableWindowsRenameError(windows.ERROR_ACCESS_DENIED))
	assert.True(t, isRetryableWindowsRenameError(windows.ERROR_SHARING_VIOLATION))
	assert.True(t, isRetryableWindowsRenameError(windows.ERROR_LOCK_VIOLATION))
	assert.False(t, isRetryableWindowsRenameError(windows.ERROR_FILE_NOT_FOUND))
}

func TestIsRetryableCopyImageError(t *testing.T) {
	assert.True(t, isRetryableCopyImageError(&os.LinkError{Err: windows.ERROR_ACCESS_DENIED}))
	assert.True(t, isRetryableCopyImageError(&os.PathError{Err: windows.ERROR_SHARING_VIOLATION}))
	assert.False(t, isRetryableCopyImageError(&os.PathError{Err: windows.ERROR_FILE_NOT_FOUND}))
	assert.False(t, isRetryableCopyImageError(errors.New("plain error")))
}
