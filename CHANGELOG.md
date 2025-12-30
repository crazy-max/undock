# Changelog

## 0.11.0 (2025/12/30)

* Go 1.25 (#368)
* Alpine Linux 3.23 (#372)
* MkDocs Materials 9.6.20 (#375)
* Bump github.com/alecthomas/kong to 1.13.0 (#360)
* Bump github.com/docker/docker to 28.5.2+incompatible (#357 #369 #370)
* Bump github.com/mholt/archives to 0.1.5 (#354)
* Bump go.podman.io/image/v5 to 5.38.0 (#369 #371)
* Bump golang.org/x/crypto to 0.45.0 (#373)
* Bump golang.org/x/sync to 0.19.0 (#363)
* Bump golang.org/x/sys to 0.39.0 (#362)

## 0.10.0 (2025/04/18)

* Fix fallback to home for cache dir by @crazy-max (#316)
* Bump github.com/alecthomas/kong to 1.10.0 (#286 #311)
* Bump github.com/containers/image/v5 to 5.33.1 (#289)
* Bump github.com/mholt/archives to 0.1.1 (#285 #312)
* Bump github.com/opencontainers/image-spec to 1.1.1 (#303)
* Bump github.com/rs/zerolog to 1.34.0 (#307)
* Bump golang.org/x/sync to 0.13.0 (#313)
* Bump golang.org/x/sys to 0.32.0 (#284 #314)

## 0.9.0 (2024/12/24)

* Go 1.23 (#271)
* Alpine Linux 3.21 (#271)
* Switch to github.com/containerd/platforms (#272)
* Switch to github.com/mholt/archives (#282)
* Bump github.com/alecthomas/kong to 1.6.0 (#280)
* Bump github.com/containers/image/v5 to 5.33.0 (#278)
* Bump github.com/docker/docker to 27.3.1+incompatible (#281)
* Bump github.com/stretchr/testify to 1.10.0 (#274)
* Bump golang.org/x/sync to 0.10.0 (#276)
* Bump golang.org/x/sys to 0.28.0 (#277)

## 0.8.0 (2024/06/20)

* Enables automatic API version negotiation for Docker client (#242)
* Move `extractor pkg` out of internal (#224)
* Bump github.com/alecthomas/kong to 0.9.0 (#209)
* Bump github.com/containerd/containerd from 1.7.11 to 1.7.18 (#222 #236)
* Bump github.com/containers/image/v5 to 5.31.1 (#188 #207 #230 #241)
* Bump github.com/docker/docker to 26.1.4+incompatible (#223 #237)
* Bump github.com/opencontainers/image-spec to 1.1.0 (#204)
* Bump github.com/rs/zerolog to 1.33.0 (#200 #234)
* Bump github.com/stretchr/testify to 1.9.0 (#205)
* Bump golang.org/x/net from 0.22.0 to 0.23.0 (#225)
* Bump golang.org/x/sync from 0.5.0 to 0.7.0 (#215)
* Bump golang.org/x/sys to 0.21.0 (#214 #235)

## 0.7.0 (2023/12/21)

* Docker auth config support (#186)
* Go 1.21 (#179 #185)
* Bump github.com/alecthomas/kong to 0.8.1 (#169)
* Bump github.com/containerd/containerd to 1.7.11 (#177)
* Bump github.com/docker/docker to 24.0.7+incompatible (#171)
* Bump github.com/go-jose/go-jose/v3 to 3.0.1 (#184)
* Bump github.com/rs/zerolog to 1.31.0 (#165)
* Bump golang.org/x/crypto to 0.17.0 (#187)
* Bump golang.org/x/sync to 0.5.0 (#172)
* Bump golang.org/x/sys to 0.15.0 (#175)

## 0.7.0-rc.1 (2023/12/19)

* Docker auth config support (#186)
* Go 1.21 (#179 #185)
* Bump github.com/alecthomas/kong to 0.8.1 (#169)
* Bump github.com/containerd/containerd to 1.7.11 (#177)
* Bump github.com/docker/docker to 24.0.7+incompatible (#171)
* Bump github.com/go-jose/go-jose/v3 to 3.0.1 (#184)
* Bump github.com/rs/zerolog to 1.31.0 (#165)
* Bump golang.org/x/crypto to 0.17.0 (#187)
* Bump golang.org/x/sync to 0.5.0 (#172)
* Bump golang.org/x/sys to 0.15.0 (#175)

## 0.6.0 (2023/09/15)

* Warn on unknown blob format (#163)
* Use forked module to fix nil pointer dereference (#164)
* Bump github.com/containerd/containerd to 1.7.6 (#146 #159)
* Bump github.com/containers/image/v5 to 5.28.0 (#160)
* Bump github.com/docker/docker to 24.0.5+incompatible  (#138 #140)
* Bump github.com/opencontainers/image-spec to 1.1.0-rc5 (#135 #161)
* Bump github.com/rs/zerolog to 1.30.0 (#142)
* Bump golang.org/x/sys to 0.11.0 (#136 #143)
* Bump golang.org/x/sys to 0.11.0 (#143)

## 0.5.0 (2023/07/02)

* Support `image:tag@digest` format (#131)
* Go 1.20 (#114 #133)
* Alpine Linux 3.18 (#134)
* Bump github.com/alecthomas/kong to 0.8.0 (#129)
* Bump github.com/containers/image/v5 to 5.26.1 (#93 #101 #130 #132)
* Bump github.com/docker/docker to 24.0.2+incompatible (#106 #115 #123)
* Bump github.com/containerd/containerd to 1.7.2 (#97 #126)
* Bump github.com/rs/zerolog to 1.29.1 (#102)
* Bump github.com/mholt/archiver/v4 to 4.0.0-alpha.8 (#107)
* Bump github.com/sirupsen/logrus to 1.9.3 (#125)
* Bump github.com/sigstore/rekor to 1.1.1 (#113)
* Bump github.com/stretchr/testify to 1.8.4 (#94 #124)
* Bump golang.org/x/sync to 0.3.0 (#111 #128)
* Bump golang.org/x/sys to 0.7.0 (#96 #100 #110 #127)

## 0.4.0 (2023/02/14)

* Go 1.19 (#65 #25)
* Alpine Linux 3.17 (#90 #50)
* Enhance workflow (#66)
* Bump github.com/mholt/archiver/v4 to 4.0.0-alpha.7 (#43 #24)
* Bump github.com/containers/image/v5 to 5.24.1 (#87 #72 #62 #31 #26)
* Bump github.com/containerd/containerd to 1.6.17 (#89 #76 #59 #45 #32 #30)
* Bump github.com/docker/docker to 20.10.23+incompatible (#81 #71 #64 #46 #40)
* Bump github.com/stretchr/testify to 1.8.1 (#67 #54 #44)
* Bump github.com/alecthomas/kong to 0.7.1 (#49 #73)
* Bump github.com/rs/zerolog to 1.29.0 (#83 #60 #47)
* Bump github.com/sirupsen/logrus to 1.9.0 (#56)
* Bump github.com/opencontainers/image-spec to 1.1.0-rc2 (#68)
* Bump golang.org/x/sys to 0.5.0 (#91)

## 0.3.0 (2022/03/28)

* Support `.gz` format (#22)
* `UNDOCK_CACHE_DIR` env var to set cache dir and predefined one in Docker image (#8)
* Bump github.com/stretchr/testify to 1.7.1 (#18)
* Bump github.com/docker/docker to 20.10.14+incompatible (#16 #20)
* Bump github.com/containers/image/v5 to 5.20.0 (#13)
* Bump github.com/mholt/archiver/v4 to 4.0.0-alpha.5 (#12)
* Bump github.com/alecthomas/kong to 0.5.0 (#10 #17)
* Bump github.com/containerd/containerd to 1.6.2 (#9 #14 #21)

## 0.2.0 (2022/02/13)

* Support more sources through specific schemes (#7)
    * `containers-storage://<store>`: image located in a local container storage.
    * `docker://<ref>`: image in a registry implementing the "Docker Registry HTTP API V2". (default)
    * `docker-archive://<path>`: image is stored in the `docker-save` formatted file.
    * `docker-daemon://<ref>`: image stored in the docker daemon's internal storage.
    * `oci://<path>`: image compliant with the "Open Container Image Layout Specification".
    * `oci-archive://<path>`: image compliant with the "Open Container Image Layout Specification" stored as a tar archive.
    * `ostree://<ref>`: image in the local ostree repository.
* CI e2e workflow (#3)
* Docs website with mkdocs (#1)
* Bump github.com/alecthomas/kong to 0.4.0 (#4)
* Bump github.com/containers/image/v5 to 5.19.1 (#2 #5)

## 0.1.0 (2022/01/24)

* Initial version
