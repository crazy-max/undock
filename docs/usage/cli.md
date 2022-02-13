# Command Line

## Usage

```shell
undock [options] <source> <dist>
```

## Options

```
$ undock --help
Usage: undock <source> <dist>

Extract contents of a container image in a local folder. More info: https://github.com/crazy-max/undock

Arguments:
  <source>    Source image. (eg. alpine:latest)
  <dist>      Dist folder. (eg. ./dist)

Flags:
  -h, --help                   Show context-sensitive help.
      --version
      --log-level="info"       Set log level ($LOG_LEVEL).
      --log-json               Enable JSON logging output ($LOG_JSON).
      --log-caller             Add file:line of the caller to log output ($LOG_CALLER).
      --log-nocolor            Disable colorized output ($LOG_NOCOLOR).
      --cachedir=STRING        Set cache path. (eg. ~/.local/share/undock/cache)
      --platform=STRING        Enforce platform for source image. (eg. linux/amd64)
      --all                    Extract all architectures if source is a manifest list.
      --include=INCLUDE,...    Include a subset of files/dirs from the source image.
      --insecure               Allow contacting the registry or docker daemon over HTTP, or HTTPS with failed TLS verification.
      --rm-dist                Removes dist folder.
      --wrap                   For a manifest list, merge output in dist folder.
```

### Source image

`source` argument can be a container image from a registry, a local docker
image, a container store reference, etc. Following schemes can be used:

* `containers-storage://<store>`: image located in a local container storage[^1].
* `docker://<ref>`: image in a registry implementing the "Docker Registry HTTP API V2"[^1].
* `docker-archive://<path>`: image is stored in the `docker-save` formatted file[^1].
* `docker-daemon://<ref>`: image stored in the docker daemon's internal storage[^1].
* `oci://<path>`: image compliant with the "Open Container Image Layout Specification"[^1].
* `oci-archive://<path>`: image compliant with the "Open Container Image Layout Specification" stored as a tar archive[^1].
* `ostree://<ref>`: image in the local ostree repository[^1].

!!! note
    `docker://` is used by default if scheme unset.

```console
# registry image
undock --rm-dist crazymax/buildx-pkg:latest ./dist
# or
undock --rm-dist docker://crazymax/buildx-pkg:latest ./dist

# archive docker image
docker pull crazymax/buildx-pkg:latest
docker save crazymax/buildx-pkg:latest > archive.tar
undock --rm-dist docker-archive://archive.tar ./dist

# local docker image
docker build -t myimage:local .
undock --rm-dist docker-daemon://myimage:local ./dist
```

## Environment variables

Following environment variables can be used in place:

| Name               | Default       | Description   |
|--------------------|---------------|---------------|
| `LOG_LEVEL`        | `info`        | Log level output |
| `LOG_JSON`         | `false`       | Enable JSON logging output |
| `LOG_CALLER`       | `false`       | Enable to add `file:line` of the caller |
| `LOG_NOCOLOR`      | `false`       | Disable the colorized output |

[^1]: See [containers image transport page](https://github.com/containers/image/blob/main/docs/containers-transports.5.md) for more info.
