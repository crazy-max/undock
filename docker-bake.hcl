variable "GO_VERSION" {
  default = null
}

variable "DESTDIR" {
  default = "./bin"
}

# GITHUB_REF is the actual ref that triggers the workflow and used as version
# when tag is pushed! https://docs.github.com/en/actions/learn-github-actions/environment-variables#default-environment-variables
variable "GITHUB_REF" {
  default = ""
}

target "_common" {
  args = {
    BUILDKIT_CONTEXT_KEEP_GIT_DIR = 1
    GO_VERSION = GO_VERSION
    GIT_REF = GITHUB_REF
  }
}

# Special target: https://github.com/docker/metadata-action#bake-definition
target "docker-metadata-action" {
  tags = ["undock:local"]
}

group "default" {
  targets = ["image-local"]
}

target "binary" {
  inherits = ["_common"]
  target = "binary"
  output = ["${DESTDIR}/build"]
}

target "artifact" {
  inherits = ["_common"]
  target = "artifact"
  output = ["${DESTDIR}/artifact"]
}

target "artifact-all" {
  inherits = ["artifact"]
  platforms = [
    "darwin/amd64",
    "darwin/arm64",
    "linux/amd64",
    "linux/arm/v5",
    "linux/arm/v6",
    "linux/arm/v7",
    "linux/arm64",
    "linux/ppc64le",
    "linux/riscv64",
    "linux/s390x",
    "windows/amd64",
    "windows/arm64"
  ]
}

target "release" {
  target = "release"
  output = ["${DESTDIR}/release"]
  contexts = {
    artifacts = "${DESTDIR}/artifact"
  }
}

target "image" {
  inherits = ["_common", "docker-metadata-action"]
}

target "image-local" {
  inherits = ["image"]
  output = ["type=docker"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/amd64",
    "linux/arm/v6",
    "linux/arm/v7",
    "linux/arm64",
    "linux/ppc64le",
    "linux/s390x"
  ]
}

target "test" {
  inherits = ["_common"]
  target = "test-coverage"
  output = ["${DESTDIR}/coverage"]
}

target "vendor" {
  inherits = ["_common"]
  dockerfile = "./hack/vendor.Dockerfile"
  target = "update"
  output = ["."]
}

target "gomod-outdated" {
  inherits = ["_common"]
  dockerfile = "./hack/vendor.Dockerfile"
  target = "outdated"
  output = ["type=cacheonly"]
}

target "docs" {
  dockerfile = "./hack/docs.Dockerfile"
  target = "release"
  output = ["${DESTDIR}/site"]
}

group "validate" {
  targets = ["lint", "vendor-validate"]
}

target "lint" {
  inherits = ["_common"]
  dockerfile = "./hack/lint.Dockerfile"
  target = "lint"
  output = ["type=cacheonly"]
}

target "vendor-validate" {
  inherits = ["_common"]
  dockerfile = "./hack/vendor.Dockerfile"
  target = "validate"
  output = ["type=cacheonly"]
}
