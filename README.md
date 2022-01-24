<p align="center"><a href="https://github.com/crazy-max/undock" target="_blank"><img height="256" src="https://raw.githubusercontent.com/crazy-max/undock/master/.github/undock.png"></a></p>

<p align="center">
  <a href="https://github.com/crazy-max/undock/releases/latest"><img src="https://img.shields.io/github/release/crazy-max/undock.svg?style=flat-square" alt="GitHub release"></a>
  <a href="https://github.com/crazy-max/undock/releases/latest"><img src="https://img.shields.io/github/downloads/crazy-max/undock/total.svg?style=flat-square" alt="Total downloads"></a>
  <a href="https://github.com/crazy-max/undock/actions?workflow=build"><img src="https://img.shields.io/github/workflow/status/crazy-max/undock/build?label=build&logo=github&style=flat-square" alt="Build Status"></a>
  <a href="https://hub.docker.com/r/crazymax/undock/"><img src="https://img.shields.io/docker/stars/crazymax/undock.svg?style=flat-square&logo=docker" alt="Docker Stars"></a>
  <a href="https://hub.docker.com/r/crazymax/undock/"><img src="https://img.shields.io/docker/pulls/crazymax/undock.svg?style=flat-square&logo=docker" alt="Docker Pulls"></a>
  <br /><a href="https://goreportcard.com/report/github.com/crazy-max/undock"><img src="https://goreportcard.com/badge/github.com/crazy-max/undock?style=flat-square" alt="Go Report"></a>
  <a href="https://codecov.io/gh/crazy-max/undock"><img src="https://img.shields.io/codecov/c/github/crazy-max/undock?logo=codecov&style=flat-square" alt="Codecov"></a>
  <a href="https://github.com/sponsors/crazy-max"><img src="https://img.shields.io/badge/sponsor-crazy--max-181717.svg?logo=github&style=flat-square" alt="Become a sponsor"></a>
  <a href="https://www.paypal.me/crazyws"><img src="https://img.shields.io/badge/donate-paypal-00457c.svg?logo=paypal&style=flat-square" alt="Donate Paypal"></a>
</p>

## About

**Undock** is a CLI application that allows you to extract contents of a
container image in a local folder.

___

* [Usage](#usage)
  * [Minimal](#minimal)
  * [Extract for all architectures](#extract-for-all-architectures)
  * [Wrap contents in a single directory](#wrap-contents-in-a-single-directory)
  * [Extract a subset of files/dirs](#extract-a-subset-of-filesdirs)
* [Build](#build)
* [Contributing](#contributing)
* [License](#license)

## Usage

```shell
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

### Minimal

```shell
$ undock --rm-dist crazymax/buildx-pkg:latest ./dist
$ tree ./dist
./dist
├── docker-buildx-0.7.0~53-gb265f1cf.m-centos7.x86_64.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-centos8.x86_64.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora33.x86_64.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora34.x86_64.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora35.x86_64.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-ol8.x86_64.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-rhel7.x86_64.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-rocky8.x86_64.rpm
├── docker-buildx_0.7.0-r0~53-gb265f1cf.m_x86_64.apk
├── docker-buildx_0.7.0~53-gb265f1cf.m-debian10_amd64.deb
├── docker-buildx_0.7.0~53-gb265f1cf.m-debian11_amd64.deb
├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian10_amd64.deb
├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian11_amd64.deb
├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu1804_amd64.deb
├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2004_amd64.deb
├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2104_amd64.deb
└── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2110_amd64.deb
```

### Extract for all architectures

You can extract for all architectures if source image is a manifest list:

```shell
$ undock --rm-dist --all crazymax/buildx-pkg:latest ./dist
$ tree ./dist
./dist
├── linux_amd64
│   ├── docker-buildx-0.7.0~53-gb265f1cf.m-centos7.x86_64.rpm
│   ├── docker-buildx-0.7.0~53-gb265f1cf.m-centos8.x86_64.rpm
│   ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora33.x86_64.rpm
│   ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora34.x86_64.rpm
│   ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora35.x86_64.rpm
│   ├── docker-buildx-0.7.0~53-gb265f1cf.m-ol8.x86_64.rpm
│   ├── docker-buildx-0.7.0~53-gb265f1cf.m-rhel7.x86_64.rpm
│   ├── docker-buildx-0.7.0~53-gb265f1cf.m-rocky8.x86_64.rpm
│   ├── docker-buildx_0.7.0-r0~53-gb265f1cf.m_x86_64.apk
│   ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian10_amd64.deb
│   ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian11_amd64.deb
│   ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian10_amd64.deb
│   ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian11_amd64.deb
│   ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu1804_amd64.deb
│   ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2004_amd64.deb
│   ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2104_amd64.deb
│   └── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2110_amd64.deb
├── linux_arm64
│   ├── docker-buildx-0.7.0~53-gb265f1cf.m-centos7.aarch64.rpm
│   ├── docker-buildx-0.7.0~53-gb265f1cf.m-centos8.aarch64.rpm
│   ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora33.aarch64.rpm
│   ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora34.aarch64.rpm
│   ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora35.aarch64.rpm
│   ├── docker-buildx-0.7.0~53-gb265f1cf.m-ol8.aarch64.rpm
...
```

### Wrap contents in a single directory

For a manifest list, merge output in dist folder:

```shell
$ undock --wrap --rm-dist --all crazymax/buildx-pkg:latest ./dist
$ tree ./dist
./dist
├── docker-buildx-0.7.0~53-gb265f1cf.m-centos7.aarch64.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-centos7.armv6hl.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-centos7.armv7hl.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-centos7.ppc64le.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-centos7.riscv64.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-centos7.s390x.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-centos7.x86_64.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-centos8.aarch64.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-centos8.armv6hl.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-centos8.armv7hl.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-centos8.ppc64le.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-centos8.riscv64.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-centos8.s390x.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-centos8.x86_64.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora33.aarch64.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora33.armv6hl.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora33.armv7hl.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora33.ppc64le.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora33.riscv64.rpm
├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora33.s390x.rpm
...
```

### Extract a subset of files/dirs

It can be useful to extract contents of a specific subset of files/dirs if the
source image is not a scratch one.

```shell
$ undock --include /usr/local/bin --rm-dist --all crazymax/diun:latest ./dist
$ tree ./dist
./dist
├── linux_386
│   └── usr
│       └── local
│           └── bin
│               └── diun
├── linux_amd64
│   └── usr
│       └── local
│           └── bin
│               └── diun
├── linux_arm64
│   └── usr
│       └── local
│           └── bin
│               └── diun
├── linux_armv6
│   └── usr
│       └── local
│           └── bin
│               └── diun
├── linux_armv7
│   └── usr
│       └── local
│           └── bin
│               └── diun
└── linux_ppc64le
    └── usr
        └── local
            └── bin
                └── diun
```

## Build

Everything is dockerized and handled by [buildx bake](docker-bake.hcl) for an
agnostic usage of this repo:

```shell
git clone https://github.com/crazy-max/undock.git undock
cd undock

# build docker image and output to docker with undock:local tag (default)
docker buildx bake

# build binary in ./bin/undock
docker buildx bake binary

# build artifact
docker buildx bake artifact

# build artifact for many platforms
docker buildx bake artifact-all

# build multi-platform image
docker buildx bake image-all
```

## Contributing

Want to contribute? Awesome! The most basic way to show your support is to star the project, or to raise issues. If
you want to open a pull request, please read the [contributing guidelines](.github/CONTRIBUTING.md).

You can also support this project by [**becoming a sponsor on GitHub**](https://github.com/sponsors/crazy-max) or by
making a [Paypal donation](https://www.paypal.me/crazyws) to ensure this journey continues indefinitely!

Thanks again for your support, it is much appreciated! :pray:

## License

MIT. See `LICENSE` for more details.
