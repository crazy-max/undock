# Command Line

## Usage

```shell
undock [options]
```

## Options

```
$ undock --help
Usage: undock <source> <dist>

Extract contents of a container image in a local folder. More info: https://github.com/crazy-max/undock

Arguments:
  <source>    Source image from a registry. (eg. alpine:latest)
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
      --insecure               Allow contacting the registry over HTTP, or HTTPS with failed TLS verification.
      --rm-dist                Removes dist folder.
      --wrap                   For a manifest list, merge output in dist folder.
```

## Environment variables

Following environment variables can be used in place:

| Name               | Default       | Description   |
|--------------------|---------------|---------------|
| `LOG_LEVEL`        | `info`        | Log level output |
| `LOG_JSON`         | `false`       | Enable JSON logging output |
| `LOG_CALLER`       | `false`       | Enable to add `file:line` of the caller |
| `LOG_NOCOLOR`      | `false`       | Disable the colorized output |
