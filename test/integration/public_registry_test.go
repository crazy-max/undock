//go:build integration

package integration

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/crazy-max/undock/internal/app"
	"github.com/crazy-max/undock/internal/config"
	"github.com/stretchr/testify/require"
)

func TestPublicRegistry(t *testing.T) {
	requirePublicRegistry(t)
	cli := config.Cli{
		// The short name exercises Docker Hub's default registry and namespace.
		Source:   "alpine:3.24@sha256:f59abe656b692d901e184be3445e084786602108599d2a51d229c0f5a48c71bf",
		Platform: "linux/amd64",
		Dist:     t.TempDir(),
		CacheDir: t.TempDir(),
		Includes: []string{"/etc/alpine-release"},
	}
	// Insecure remains false: token acquisition and blob downloads use HTTPS.
	integrationExtract(t, cli)
	requireFileContent(t, filepath.Join(cli.Dist, "etc", "alpine-release"), "3.24.4\n")
	require.NoDirExists(t, filepath.Join(cli.Dist, "bin"))
}

func TestPublicRegistryMultiPlatform(t *testing.T) {
	requirePublicRegistry(t)
	packages := []struct {
		directory string
		path      string
		checksum  string
	}{
		{
			directory: "darwin_amd64",
			path:      "static/darwin/amd64/docker-buildx-plugin_0.39.1.tgz",
			checksum:  "bd901128e703e8902e4a6ea43fd3ccb91010a1070877aa1df89f4361a7a0575a",
		},
		{
			directory: "darwin_arm64",
			path:      "static/darwin/arm64/docker-buildx-plugin_0.39.1.tgz",
			checksum:  "27739c107130475f4a0d11597398c764c6044ac3d7f98a2860c679f0bb92ca22",
		},
		{
			directory: "linux_amd64",
			path:      "static/linux/amd64/docker-buildx-plugin_0.39.1.tgz",
			checksum:  "ba8c820dd24ee731f80721dad268a3e09b743e72abbdbd923ca7d0eca8d5a019",
		},
		{
			directory: "linux_arm64",
			path:      "static/linux/arm64/docker-buildx-plugin_0.39.1.tgz",
			checksum:  "a1006b9513cea4d4529145a29f76ef9a390d0749ce9f1eef1745865b2c2dbdd0",
		},
		{
			directory: "linux_armv6",
			path:      "static/linux/arm/v6/docker-buildx-plugin_0.39.1.tgz",
			checksum:  "6d140d52bce4cbb7e40d5a2d1392d33647975e17a5126b70668dd4f848708999",
		},
		{
			directory: "linux_armv7",
			path:      "static/linux/arm/v7/docker-buildx-plugin_0.39.1.tgz",
			checksum:  "11d599d265dab79a7b39dd3966345c40e46d8196d1420c658732f7d44dc898e2",
		},
		{
			directory: "linux_ppc64le",
			path:      "static/linux/ppc64le/docker-buildx-plugin_0.39.1.tgz",
			checksum:  "7a6a4a6e29bae10bc955b797176b9fe67fb2eefd34d0df2bbd0b8e3734178732",
		},
		{
			directory: "linux_riscv64",
			path:      "static/linux/riscv64/docker-buildx-plugin_0.39.1.tgz",
			checksum:  "b5222890d8ea034316ab5e0fb3de73d488e80fbc3a034a77954dc6c3f658d8c3",
		},
		{
			directory: "linux_s390x",
			path:      "static/linux/s390x/docker-buildx-plugin_0.39.1.tgz",
			checksum:  "4c3b56fa279dc4ef705d7d22c1552fdf20e9eb9b6035d65812caa959762f0d73",
		},
		{
			directory: "windows_amd64",
			path:      "static/windows/amd64/docker-buildx-plugin_0.39.1.zip",
			checksum:  "c451face44272d646a6514f3f6e4ec6d194811ebeb4132beb27647ef823ec80f",
		},
		{
			directory: "windows_arm64",
			path:      "static/windows/arm64/docker-buildx-plugin_0.39.1.zip",
			checksum:  "68ecbb1d3a9d4f48fc69ea2e2cf83e5f42da089a5334dde43228a151f2141e9c",
		},
	}
	cache := t.TempDir()
	for _, tc := range []struct {
		name     string
		wrap     bool
		includes []string
	}{
		{
			name: "split",
		},
		{
			name: "wrap",
			wrap: true,
		},
		{
			name:     "split-include",
			includes: []string{"/static"},
		},
		{
			name:     "wrap-include",
			wrap:     true,
			includes: []string{"/static"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cli := config.Cli{
				Source:   "dockereng/packaging:buildx-v0.39.1@sha256:0dcfe9931281584037025d778761675b9ee1ea41c38c05c4ec652adf37e3e838",
				Dist:     t.TempDir(),
				CacheDir: cache,
				All:      true,
				Wrap:     tc.wrap,
				Includes: tc.includes,
				RmDist:   true,
			}
			for range 2 {
				require.NoError(t, os.WriteFile(filepath.Join(cli.Dist, "stale.txt"), []byte("stale"), 0o600))
				application, err := app.New(config.Meta{UserAgent: "undock-integration"}, cli)
				require.NoError(t, err)
				// This packaging index includes distribution packages as well as binaries.
				ctx, cancel := context.WithTimeoutCause(t.Context(), 5*time.Minute, errors.New("public image extraction timed out"))
				err = application.Start(ctx)
				cancel()
				require.NoError(t, err)
				require.NoFileExists(t, filepath.Join(cli.Dist, "stale.txt"))
				var directories []string
				for _, pkg := range packages {
					dir := cli.Dist
					if !tc.wrap {
						dir = filepath.Join(dir, pkg.directory)
						directories = append(directories, pkg.directory)
					} else {
						require.NoDirExists(t, filepath.Join(dir, pkg.directory))
					}
					requireFileSHA256(t, filepath.Join(dir, filepath.FromSlash(pkg.path)), pkg.checksum)
					if len(tc.includes) > 0 {
						entries, err := os.ReadDir(dir)
						require.NoError(t, err)
						require.Len(t, entries, 1)
						require.Equal(t, "static", entries[0].Name())
					} else {
						data, err := os.ReadFile(filepath.Join(dir, "metadata.env"))
						require.NoError(t, err)
						require.Contains(t, string(data), "VERSION=v0.39.1\n")
					}
				}
				if !tc.wrap {
					entries, err := os.ReadDir(cli.Dist)
					require.NoError(t, err)
					var actual []string
					for _, entry := range entries {
						actual = append(actual, entry.Name())
					}
					require.ElementsMatch(t, directories, actual)
				}
			}
		})
	}
}
