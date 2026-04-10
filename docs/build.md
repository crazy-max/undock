# Build

Everything is dockerized and handled by [Buildx Bake](https://docs.docker.com/build/bake/) for an agnostic usage of
this project:

```shell
git clone https://github.com/crazy-max/undock.git undock
cd undock

# build docker image and output to docker with undock:local tag (default)
docker buildx bake

# build binary in ./bin/build/undock
docker buildx bake binary

# build artifact
docker buildx bake artifact

# build artifact for many platforms
docker buildx bake artifact-all

# build multi-platform image
docker buildx bake image-all
```
