# Docker image

## About

Undock provides automatically updated Docker :whale: images within several registries:

| Registry                                                                                           | Image                       |
|----------------------------------------------------------------------------------------------------|-----------------------------|
| [Docker Hub](https://hub.docker.com/r/crazymax/undock/)                                            | `crazymax/undock`           |
| [GitHub Container Registry](https://github.com/users/crazy-max/packages/container/package/undock)  | `ghcr.io/crazy-max/undock`  |

It is possible to always use the latest stable tag or to use another service that handles updating Docker images.

!!! note
    Want to be notified of new releases? Check out :bell: [Diun (Docker Image Update Notifier)](https://github.com/crazy-max/diun) project!

Following platforms for this image are available:

```
$ docker run --rm mplatform/mquery crazymax/undock:latest
Image: crazymax/undock:latest
 * Manifest List: Yes
 * Supported platforms:
   - linux/amd64
   - linux/arm/v6
   - linux/arm/v7
   - linux/arm64
   - linux/ppc64le
   - linux/s390x
```
