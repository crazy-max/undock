# Download

Undock binaries are available on [releases]({{ config.repo_url }}releases/latest) page.

Choose the archive matching the destination platform:

* [`undock_{{ git.tag | trim('v') }}_darwin_arm64.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/undock_{{ git.tag | trim('v') }}_darwin_arm64.tar.gz)
* [`undock_{{ git.tag | trim('v') }}_darwin_amd64.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/undock_{{ git.tag | trim('v') }}_darwin_amd64.tar.gz)
* [`undock_{{ git.tag | trim('v') }}_linux_amd64.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/undock_{{ git.tag | trim('v') }}_linux_amd64.tar.gz)
* [`undock_{{ git.tag | trim('v') }}_linux_arm64.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/undock_{{ git.tag | trim('v') }}_linux_arm64.tar.gz)
* [`undock_{{ git.tag | trim('v') }}_linux_armv5.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/undock_{{ git.tag | trim('v') }}_linux_armv5.tar.gz)
* [`undock_{{ git.tag | trim('v') }}_linux_armv6.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/undock_{{ git.tag | trim('v') }}_linux_armv6.tar.gz)
* [`undock_{{ git.tag | trim('v') }}_linux_armv7.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/undock_{{ git.tag | trim('v') }}_linux_armv7.tar.gz)
* [`undock_{{ git.tag | trim('v') }}_linux_ppc64le.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/undock_{{ git.tag | trim('v') }}_linux_ppc64le.tar.gz)
* [`undock_{{ git.tag | trim('v') }}_linux_riscv64.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/undock_{{ git.tag | trim('v') }}_linux_riscv64.tar.gz)
* [`undock_{{ git.tag | trim('v') }}_linux_s390x.tar.gz`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/undock_{{ git.tag | trim('v') }}_linux_s390x.tar.gz)
* [`undock_{{ git.tag | trim('v') }}_windows_amd64.zip`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/undock_{{ git.tag | trim('v') }}_windows_amd64.zip)
* [`undock_{{ git.tag | trim('v') }}_windows_arm64.zip`]({{ config.repo_url }}/releases/download/v{{ git.tag | trim('v') }}/undock_{{ git.tag | trim('v') }}_windows_arm64.zip)

And extract Undock:

```shell
wget -qO- {{ config.repo_url }}releases/download/v{{ git.tag | trim('v') }}/undock_{{ git.tag | trim('v') }}_linux_amd64.tar.gz | tar -zxvf - undock
```

After getting the binary, it can be tested with [`./undock --help`](../usage/cli.md) command and moved to a
permanent location.
