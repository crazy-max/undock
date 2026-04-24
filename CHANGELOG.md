# Changelog

## 0.13.0 (2026/04/24)

* Use backoff for image copy retries on Windows by @crazy-max in #436
* Preserve logrus fields in zerolog output by @crazy-max in #432
* Simplify shutdown lifecycle by @crazy-max in #430
* Drop docker client wrapper by @crazy-max in #431
* Go 1.26 by @crazy-max in #434
* MkDocs Materials 9.7.5 by @crazy-max in #416
* Bump github.com/rs/zerolog to 1.35.1 in #427

**Full Changelog**: [`v0.12.0...v0.13.0`](https://github.com/crazy-max/undock/compare/v0.12.0...v0.13.0)

## 0.12.0 (2026/04/10)

* Handle OCI whiteouts correctly during blob extraction by @crazy-max in #413
* Retry cache copy on transient Windows rename failures by @crazy-max in #411
* Bump github.com/alecthomas/kong to 1.15.0 in #396
* Bump github.com/rs/zerolog to 1.35.0 in #395
* Bump github.com/sirupsen/logrus to 1.9.4 in #377
* Bump golang.org/x/sync to 0.20.0 in #392
* Bump golang.org/x/sys to 0.43.0 in #410
* Bump go.podman.io/image/v5 to 5.39.2 in #397

**Full Changelog**: [`v0.11.0...v0.12.0`](https://github.com/crazy-max/undock/compare/v0.11.0...v0.12.0)

## 0.11.0 (2025/12/30)

* Go 1.25 by @crazy-max in #368
* Alpine Linux 3.23 by @crazy-max in #372
* MkDocs Materials 9.6.20 by @crazy-max in #375
* Bump github.com/alecthomas/kong to 1.13.0 in #360
* Bump github.com/docker/docker to 28.5.2+incompatible in #357 #369 #370
* Bump github.com/mholt/archives to 0.1.5 in #354
* Bump go.podman.io/image/v5 to 5.38.0 by @crazy-max in #369 #371
* Bump golang.org/x/crypto to 0.45.0 in #373
* Bump golang.org/x/sync to 0.19.0 in #363
* Bump golang.org/x/sys to 0.39.0 in #362

**Full Changelog**: [`v0.10.0...v0.11.0`](https://github.com/crazy-max/undock/compare/v0.10.0...v0.11.0)

## 0.10.0 (2025/04/18)

* Fix fallback to home for cache dir by @crazy-max in #316
* Bump github.com/alecthomas/kong to 1.10.0 in #286 #311
* Bump github.com/containers/image/v5 to 5.33.1 in #289
* Bump github.com/mholt/archives to 0.1.1 in #285 #312
* Bump github.com/opencontainers/image-spec to 1.1.1 in #303
* Bump github.com/rs/zerolog to 1.34.0 in #307
* Bump golang.org/x/sync to 0.13.0 in #313
* Bump golang.org/x/sys to 0.32.0 in #284 #314

**Full Changelog**: [`v0.9.0...v0.10.0`](https://github.com/crazy-max/undock/compare/v0.9.0...v0.10.0)

## 0.9.0 (2024/12/24)

* Go 1.23 by @crazy-max in #271
* Alpine Linux 3.21 by @crazy-max in #271
* Switch to github.com/containerd/platforms by @crazy-max in #272
* Switch to github.com/mholt/archives by @crazy-max in #282
* Bump github.com/alecthomas/kong to 1.6.0 in #280
* Bump github.com/containers/image/v5 to 5.33.0 in #278
* Bump github.com/docker/docker to 27.3.1+incompatible by @crazy-max in #281
* Bump github.com/stretchr/testify to 1.10.0 in #274
* Bump golang.org/x/sync to 0.10.0 in #276
* Bump golang.org/x/sys to 0.28.0 in #277

**Full Changelog**: [`v0.8.0...v0.9.0`](https://github.com/crazy-max/undock/compare/v0.8.0...v0.9.0)

## 0.8.0 (2024/06/20)

* Enables automatic API version negotiation for Docker client by @crazy-max in #242
* Move `extractor pkg` out of internal by @crazy-max in #224
* Bump github.com/alecthomas/kong to 0.9.0 in #209
* Bump github.com/containerd/containerd from 1.7.11 to 1.7.18 in #222 #236
* Bump github.com/containers/image/v5 to 5.31.1 in #188 #207 #230 #241
* Bump github.com/docker/docker to 26.1.4+incompatible in #223 #237
* Bump github.com/opencontainers/image-spec to 1.1.0 in #204
* Bump github.com/rs/zerolog to 1.33.0 in #200 #234
* Bump github.com/stretchr/testify to 1.9.0 in #205
* Bump golang.org/x/net from 0.22.0 to 0.23.0 in #225
* Bump golang.org/x/sync from 0.5.0 to 0.7.0 in #215
* Bump golang.org/x/sys to 0.21.0 in #214 #235

**Full Changelog**: [`v0.7.0...v0.8.0`](https://github.com/crazy-max/undock/compare/v0.7.0...v0.8.0)

## 0.7.0 (2023/12/21)

* Docker auth config support by @crazy-max in #186
* Go 1.21 by @crazy-max in #179 #185
* Bump github.com/alecthomas/kong to 0.8.1 in #169
* Bump github.com/containerd/containerd to 1.7.11 in #177
* Bump github.com/docker/docker to 24.0.7+incompatible in #171
* Bump github.com/go-jose/go-jose/v3 to 3.0.1 in #184
* Bump github.com/rs/zerolog to 1.31.0 in #165
* Bump golang.org/x/crypto to 0.17.0 in #187
* Bump golang.org/x/sync to 0.5.0 in #172
* Bump golang.org/x/sys to 0.15.0 in #175

**Full Changelog**: [`v0.6.0...v0.7.0`](https://github.com/crazy-max/undock/compare/v0.7.0-rc.1...v0.7.0)

## 0.6.0 (2023/09/15)

* Warn on unknown blob format by @crazy-max in #163
* Use forked module to fix nil pointer dereference by @crazy-max in #164
* Bump github.com/containerd/containerd to 1.7.6 in #146 #159
* Bump github.com/containers/image/v5 to 5.28.0 in #160
* Bump github.com/docker/docker to 24.0.5+incompatible in #138 #140
* Bump github.com/opencontainers/image-spec to 1.1.0-rc5 in #135 #161
* Bump github.com/rs/zerolog to 1.30.0 in #142
* Bump golang.org/x/sys to 0.11.0 in #136 #143
* Bump golang.org/x/sys to 0.11.0 in #143

**Full Changelog**: [`v0.5.0...v0.6.0`](https://github.com/crazy-max/undock/compare/v0.5.0...v0.6.0)

## 0.5.0 (2023/07/02)

* Support `image:tag@digest` format by @crazy-max in #131
* Go 1.20 by @crazy-max in #114 #133
* Alpine Linux 3.18 by @crazy-max in #134
* Bump github.com/alecthomas/kong to 0.8.0 in #129
* Bump github.com/containers/image/v5 to 5.26.1 in #93 #101 #130 #132
* Bump github.com/docker/docker to 24.0.2+incompatible in #106 #115 #123
* Bump github.com/containerd/containerd to 1.7.2 in #97 #126
* Bump github.com/rs/zerolog to 1.29.1 in #102
* Bump github.com/mholt/archiver/v4 to 4.0.0-alpha.8 in #107
* Bump github.com/sirupsen/logrus to 1.9.3 in #125
* Bump github.com/sigstore/rekor to 1.1.1 in #113
* Bump github.com/stretchr/testify to 1.8.4 in #94 #124
* Bump golang.org/x/sync to 0.3.0 in #111 #128
* Bump golang.org/x/sys to 0.7.0 in #96 #100 #110 #127

**Full Changelog**: [`v0.4.0...v0.5.0`](https://github.com/crazy-max/undock/compare/v0.4.0...v0.5.0)

## 0.4.0 (2023/02/14)

* Go 1.19 by @crazy-max in #65 #25
* Alpine Linux 3.17 by @crazy-max in #90 #50
* Enhance workflow by @crazy-max in #66
* Bump github.com/mholt/archiver/v4 to 4.0.0-alpha.7 in #43 #24
* Bump github.com/containers/image/v5 to 5.24.1 in #87 #72 #62 #31 #26
* Bump github.com/containerd/containerd to 1.6.17 in #89 #76 #59 #45 #32 #30
* Bump github.com/docker/docker to 20.10.23+incompatible in #81 #71 #64 #46 #40
* Bump github.com/stretchr/testify to 1.8.1 in #67 #54 #44
* Bump github.com/alecthomas/kong to 0.7.1 in #49 #73
* Bump github.com/rs/zerolog to 1.29.0 in #83 #60 #47
* Bump github.com/sirupsen/logrus to 1.9.0 in #56
* Bump github.com/opencontainers/image-spec to 1.1.0-rc2 in #68
* Bump golang.org/x/sys to 0.5.0 in #91

**Full Changelog**: [`v0.3.0...v0.4.0`](https://github.com/crazy-max/undock/compare/v0.3.0...v0.4.0)

## 0.3.0 (2022/03/28)

* Support `.gz` format by @crazy-max in #22
* `UNDOCK_CACHE_DIR` env var to set cache dir and predefined one in Docker image by @crazy-max in #8
* Bump github.com/stretchr/testify to 1.7.1 in #18
* Bump github.com/docker/docker to 20.10.14+incompatible in #16 #20
* Bump github.com/containers/image/v5 to 5.20.0 in #13
* Bump github.com/mholt/archiver/v4 to 4.0.0-alpha.5 in #12
* Bump github.com/alecthomas/kong to 0.5.0 in #10 #17
* Bump github.com/containerd/containerd to 1.6.2 in #9 #14 #21

**Full Changelog**: [`v0.2.0...v0.3.0`](https://github.com/crazy-max/undock/compare/v0.2.0...v0.3.0)

## 0.2.0 (2022/02/13)

* Support more sources through specific schemes by @crazy-max in #7
    * `containers-storage://<store>`: image located in a local container storage.
    * `docker://<ref>`: image in a registry implementing the "Docker Registry HTTP API V2". (default)
    * `docker-archive://<path>`: image is stored in the `docker-save` formatted file.
    * `docker-daemon://<ref>`: image stored in the docker daemon's internal storage.
    * `oci://<path>`: image compliant with the "Open Container Image Layout Specification".
    * `oci-archive://<path>`: image compliant with the "Open Container Image Layout Specification" stored as a tar archive.
    * `ostree://<ref>`: image in the local ostree repository.
* Docs website with mkdocs by @crazy-max in #1
* Bump github.com/alecthomas/kong to 0.4.0 in #4
* Bump github.com/containers/image/v5 to 5.19.1 in #2 #5

**Full Changelog**: [`v0.1.0...v0.2.0`](https://github.com/crazy-max/undock/compare/v0.1.0...v0.2.0)

## 0.1.0 (2022/01/24)

* Initial version
