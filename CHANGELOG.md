# Changelog

## 0.3.0 (2022/03/28)

* support `.gz` format (#22)
* `UNDOCK_CACHE_DIR` env var to set cache dir and predefined one in Docker image (#8)
* build(deps): bump github.com/stretchr/testify from 1.7.0 to 1.7.1 (#18)
* build(deps): bump github.com/docker/docker from 20.10.12+incompatible to 20.10.14+incompatible (#16 #20)
* build(deps): bump github.com/containers/image/v5 from 5.19.1 to 5.20.0 (#13)
* build(deps): bump github.com/mholt/archiver/v4 from 4.0.0-alpha.4 to 4.0.0-alpha.5 (#12)
* build(deps): bump github.com/alecthomas/kong from 0.4.0 to 0.5.0 (#10 #17)
* build(deps): bump github.com/containerd/containerd from 1.5.9 to 1.6.2 (#9 #14 #21)

## 0.2.0 (2022/02/13)

* support more sources through specific schemes (#7)
    * `containers-storage://<store>`: image located in a local container storage.
    * `docker://<ref>`: image in a registry implementing the "Docker Registry HTTP API V2". (default)
    * `docker-archive://<path>`: image is stored in the `docker-save` formatted file.
    * `docker-daemon://<ref>`: image stored in the docker daemon's internal storage.
    * `oci://<path>`: image compliant with the "Open Container Image Layout Specification".
    * `oci-archive://<path>`: image compliant with the "Open Container Image Layout Specification" stored as a tar archive.
    * `ostree://<ref>`: image in the local ostree repository.
* ci: e2e workflow (#3)
* docs website with mkdocs (#1)
* build(deps): bump github.com/alecthomas/kong from 0.3.0 to 0.4.0 (#4)
* build(deps): bump github.com/containers/image/v5 from 5.18.0 to 5.19.1 (#2 #5)

## 0.1.0 (2022/01/24)

* Initial version
