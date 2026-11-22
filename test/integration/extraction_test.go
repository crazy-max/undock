//go:build integration

package integration

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/containerd/platforms"
	"github.com/crazy-max/undock/internal/config"
	digest "github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/require"
)

func TestDefaultPlatform(t *testing.T) {
	endpoint := integrationRegistry(t)
	repo := "undock-" + digest.FromString(t.TempDir()).Encoded()[:16]

	native := integrationImage(t, endpoint, repo+"-native", "native", platforms.DefaultSpec())
	cli := config.Cli{Source: native, Insecure: true, Dist: t.TempDir(), CacheDir: t.TempDir()}
	integrationExtract(t, cli)
	integrationContents(t, cli.Dist, "native")
}

func TestRegistryExtraction(t *testing.T) {
	endpoint := integrationRegistry(t)
	repo := "undock-" + digest.FromString(t.TempDir()).Encoded()[:16]
	source := integrationImage(t, endpoint, repo, "first")
	for _, tc := range []struct {
		name     string
		all      bool
		wrap     bool
		includes []string
	}{
		{
			name: "platform",
		},
		{
			name: "all",
			all:  true,
		},
		{
			name:     "wrap",
			all:      true,
			wrap:     true,
			includes: []string{"/fixture/amd64.txt", "/fixture/arm64.txt"},
		},
		{
			name:     "include",
			includes: []string{"/fixture"},
		},
		{
			name:     "all-include",
			all:      true,
			includes: []string{"/fixture"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cli := config.Cli{Source: source, Platform: "linux/amd64", Insecure: true, Dist: t.TempDir(), CacheDir: t.TempDir(), All: tc.all, Wrap: tc.wrap, Includes: tc.includes, RmDist: true}
			for range 2 {
				require.NoError(t, os.WriteFile(filepath.Join(cli.Dist, "stale.txt"), []byte("stale"), 0o600))
				integrationExtract(t, cli)
				require.NoFileExists(t, filepath.Join(cli.Dist, "stale.txt"))
				if tc.all && !tc.wrap {
					for _, arch := range []string{"amd64", "arm64"} {
						dir := filepath.Join(cli.Dist, "linux_"+arch)
						if arch == "arm64" {
							dir += "v8"
						}
						integrationContents(t, dir, "first")
						requireFileContent(t, filepath.Join(dir, "fixture", arch+".txt"), arch)
						if len(tc.includes) > 0 {
							require.NoFileExists(t, filepath.Join(dir, "excluded.txt"))
						}
					}
				} else {
					if !tc.wrap {
						integrationContents(t, cli.Dist, "first")
					}
					requireFileContent(t, filepath.Join(cli.Dist, "fixture/amd64.txt"), "amd64")
					if tc.wrap {
						requireFileContent(t, filepath.Join(cli.Dist, "fixture/arm64.txt"), "arm64")
					} else {
						require.NoFileExists(t, filepath.Join(cli.Dist, "fixture/arm64.txt"))
					}
					if len(tc.includes) > 0 {
						require.NoFileExists(t, filepath.Join(cli.Dist, "excluded.txt"))
					} else {
						requireFileContent(t, filepath.Join(cli.Dist, "excluded.txt"), "excluded")
					}
				}
			}
		})
	}
}

func TestRegistryDigest(t *testing.T) {
	endpoint := integrationRegistry(t)
	repo := "undock-" + digest.FromString(t.TempDir()).Encoded()[:16]
	source := integrationImage(t, endpoint, repo, "first")

	headers := integrationRequest(t, http.MethodHead, endpoint+"/v2/"+repo+"/manifests/latest", "", nil, http.StatusOK)
	dgst := headers.Get("Docker-Content-Digest")
	require.NotEmpty(t, dgst)
	cli := config.Cli{Source: source + "@" + dgst, Platform: "linux/amd64", Insecure: true, Dist: t.TempDir(), CacheDir: t.TempDir()}
	integrationExtract(t, cli)
	integrationContents(t, cli.Dist, "first")
}

func TestRegistryTagUpdate(t *testing.T) {
	endpoint := integrationRegistry(t)
	repo := "undock-" + digest.FromString(t.TempDir()).Encoded()[:16]
	source := integrationImage(t, endpoint, repo, "first")

	cli := config.Cli{Source: source, Platform: "linux/amd64", Insecure: true, Dist: t.TempDir(), CacheDir: t.TempDir(), RmDist: true}
	integrationExtract(t, cli)
	integrationContents(t, cli.Dist, "first")
	integrationImage(t, endpoint, repo, "second")
	integrationExtract(t, cli)
	integrationContents(t, cli.Dist, "second")
}

func TestDockerArchive(t *testing.T) {
	archive, _ := integrationArchive(t)
	cli := config.Cli{Source: "docker-archive://" + archive, Platform: "linux/amd64", Dist: t.TempDir(), CacheDir: t.TempDir(), Includes: []string{"/fixture"}}
	integrationExtract(t, cli)
	integrationContents(t, cli.Dist, "second")
	require.NoFileExists(t, filepath.Join(cli.Dist, "excluded.txt"))
}

func TestDockerDaemon(t *testing.T) {
	if os.Getenv("TEST_DOCKER") != "1" {
		t.Skip("Docker daemon tests require TEST_DOCKER=1; run with test/compose.yaml")
	}
	archive, imageName := integrationArchive(t)
	command := func(args ...string) {
		t.Helper()
		ctx, cancel := context.WithTimeoutCause(context.Background(), time.Minute, errors.New("Docker command timed out"))
		defer cancel()
		output, err := exec.CommandContext(ctx, "docker", args...).CombinedOutput()
		require.NoError(t, err, "%s", output)
	}
	command("load", "--input", archive)
	t.Cleanup(func() { command("image", "rm", imageName) })
	saved := "saved.tar"
	command("save", "--output", saved, imageName)
	for _, transport := range []string{"docker-archive://" + archive, "docker-archive://" + saved, "docker-daemon://" + imageName} {
		cli := config.Cli{Source: transport, Platform: "linux/amd64", Dist: t.TempDir(), CacheDir: t.TempDir(), Includes: []string{"/fixture"}}
		integrationExtract(t, cli)
		integrationContents(t, cli.Dist, "second")
		require.NoFileExists(t, filepath.Join(cli.Dist, "excluded.txt"))
	}
}

func TestCLI(t *testing.T) {
	endpoint := integrationRegistry(t)
	repo := "undock-" + digest.FromString(t.TempDir()).Encoded()[:16]

	source := integrationImage(t, endpoint, repo, "cli")
	dist := t.TempDir()
	ctx, cancel := context.WithTimeoutCause(t.Context(), time.Minute, errors.New("CLI extraction timed out"))
	defer cancel()
	binary := os.Getenv("UNDOCK_BIN")
	if binary == "" {
		binary = "undock"
	}
	output, err := exec.CommandContext(ctx, binary, "--version").CombinedOutput() // #nosec G702 -- Test binary supplied by the runner.
	require.NoError(t, err, "%s", output)
	require.NotEmpty(t, output)
	for _, tc := range []struct {
		name  string
		flags []string
	}{
		{
			name:  "include",
			flags: []string{"--include", "/fixture"},
		},
		{
			name:  "all",
			flags: []string{"--all"},
		},
		{
			name:  "wrap",
			flags: []string{"--all", "--wrap", "--include", "/fixture/amd64.txt", "--include", "/fixture/arm64.txt"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cache := t.TempDir()
			for range 2 {
				require.NoError(t, os.WriteFile(filepath.Join(dist, "stale.txt"), []byte("stale"), 0o600))
				args := []string{"--insecure", "--platform", "linux/amd64", "--cachedir", cache, "--rm-dist"}
				args = append(args, tc.flags...)
				args = append(args, source, dist)
				output, err := exec.CommandContext(ctx, binary, args...).CombinedOutput() // #nosec G702 -- Test executable and controlled fixture arguments.
				require.NoError(t, err, "%s", output)
				require.NoFileExists(t, filepath.Join(dist, "stale.txt"))
				switch tc.name {
				case "include":
					integrationContents(t, dist, "cli")
					require.NoFileExists(t, filepath.Join(dist, "excluded.txt"))
				case "all":
					integrationContents(t, filepath.Join(dist, "linux_amd64"), "cli")
					integrationContents(t, filepath.Join(dist, "linux_arm64v8"), "cli")
				case "wrap":
					requireFileContent(t, filepath.Join(dist, "fixture/amd64.txt"), "amd64")
					requireFileContent(t, filepath.Join(dist, "fixture/arm64.txt"), "arm64")
				}
			}
		})
	}
}
