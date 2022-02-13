# Changelog

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
