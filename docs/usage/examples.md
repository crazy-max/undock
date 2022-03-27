# Examples

In this page, we quickly go over the basic use cases that will allow you to
make the most of Undock.

## Simple

=== "Command"

    ```shell
    $ undock --rm-dist crazymax/buildx-pkg:latest ./dist
    ```

=== "Output result"

    ```text
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

## Extract a multi-platform image

You can extract all architectures for a source image if this one is a
manifest list:

=== "Command"

    ```shell
    $ undock --rm-dist --all crazymax/buildx-pkg:latest ./dist
    ```

=== "Output result"

    ```text
    ./dist/
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
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-rhel7.aarch64.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-rocky8.aarch64.rpm
    │   ├── docker-buildx_0.7.0-r0~53-gb265f1cf.m_aarch64.apk
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian10_arm64.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian11_arm64.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian10_arm64.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian11_arm64.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu1804_arm64.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2004_arm64.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2104_arm64.deb
    │   └── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2110_arm64.deb
    ├── linux_armv6
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-centos7.armv6hl.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-centos8.armv6hl.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora33.armv6hl.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora34.armv6hl.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora35.armv6hl.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-ol8.armv6hl.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-rhel7.armv6hl.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-rocky8.armv6hl.rpm
    │   ├── docker-buildx_0.7.0-r0~53-gb265f1cf.m_armhf.apk
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian10_armel.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian11_armel.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian10_armel.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian11_armel.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu1804_armel.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2004_armel.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2104_armel.deb
    │   └── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2110_armel.deb
    ├── linux_armv7
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-centos7.armv7hl.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-centos8.armv7hl.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora33.armv7hl.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora34.armv7hl.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora35.armv7hl.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-ol8.armv7hl.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-rhel7.armv7hl.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-rocky8.armv7hl.rpm
    │   ├── docker-buildx_0.7.0-r0~53-gb265f1cf.m_armv7.apk
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian10_armhf.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian11_armhf.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian10_armhf.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian11_armhf.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu1804_armhf.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2004_armhf.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2104_armhf.deb
    │   └── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2110_armhf.deb
    ├── linux_ppc64le
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-centos7.ppc64le.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-centos8.ppc64le.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora33.ppc64le.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora34.ppc64le.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora35.ppc64le.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-ol8.ppc64le.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-rhel7.ppc64le.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-rocky8.ppc64le.rpm
    │   ├── docker-buildx_0.7.0-r0~53-gb265f1cf.m_ppc64le.apk
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian10_ppc64el.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian11_ppc64el.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian10_ppc64el.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian11_ppc64el.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu1804_ppc64el.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2004_ppc64el.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2104_ppc64el.deb
    │   └── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2110_ppc64el.deb
    ├── linux_riscv64
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-centos7.riscv64.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-centos8.riscv64.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora33.riscv64.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora34.riscv64.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora35.riscv64.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-ol8.riscv64.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-rhel7.riscv64.rpm
    │   ├── docker-buildx-0.7.0~53-gb265f1cf.m-rocky8.riscv64.rpm
    │   ├── docker-buildx_0.7.0-r0~53-gb265f1cf.m_riscv64.apk
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian10_riscv64.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian11_riscv64.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian10_riscv64.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian11_riscv64.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu1804_riscv64.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2004_riscv64.deb
    │   ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2104_riscv64.deb
    │   └── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2110_riscv64.deb
    └── linux_s390x
        ├── docker-buildx-0.7.0~53-gb265f1cf.m-centos7.s390x.rpm
        ├── docker-buildx-0.7.0~53-gb265f1cf.m-centos8.s390x.rpm
        ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora33.s390x.rpm
        ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora34.s390x.rpm
        ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora35.s390x.rpm
        ├── docker-buildx-0.7.0~53-gb265f1cf.m-ol8.s390x.rpm
        ├── docker-buildx-0.7.0~53-gb265f1cf.m-rhel7.s390x.rpm
        ├── docker-buildx-0.7.0~53-gb265f1cf.m-rocky8.s390x.rpm
        ├── docker-buildx_0.7.0-r0~53-gb265f1cf.m_s390x.apk
        ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian10_s390x.deb
        ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian11_s390x.deb
        ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian10_s390x.deb
        ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian11_s390x.deb
        ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu1804_s390x.deb
        ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2004_s390x.deb
        ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2104_s390x.deb
        └── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2110_s390x.deb
    ```

## Wrap contents in a single directory

For a manifest list, you can merge the output in the dist folder:

=== "Command"

    ```shell
    $ undock --wrap --rm-dist --all crazymax/buildx-pkg:latest ./dist
    ```

=== "Output result"

    ```text
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
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora33.x86_64.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora34.aarch64.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora34.armv6hl.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora34.armv7hl.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora34.ppc64le.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora34.riscv64.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora34.s390x.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora34.x86_64.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora35.aarch64.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora35.armv6hl.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora35.armv7hl.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora35.ppc64le.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora35.riscv64.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora35.s390x.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-fedora35.x86_64.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-ol8.aarch64.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-ol8.armv6hl.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-ol8.armv7hl.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-ol8.ppc64le.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-ol8.riscv64.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-ol8.s390x.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-ol8.x86_64.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-rhel7.aarch64.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-rhel7.armv6hl.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-rhel7.armv7hl.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-rhel7.ppc64le.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-rhel7.riscv64.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-rhel7.s390x.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-rhel7.x86_64.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-rocky8.aarch64.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-rocky8.armv6hl.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-rocky8.armv7hl.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-rocky8.ppc64le.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-rocky8.riscv64.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-rocky8.s390x.rpm
    ├── docker-buildx-0.7.0~53-gb265f1cf.m-rocky8.x86_64.rpm
    ├── docker-buildx_0.7.0-r0~53-gb265f1cf.m_aarch64.apk
    ├── docker-buildx_0.7.0-r0~53-gb265f1cf.m_armhf.apk
    ├── docker-buildx_0.7.0-r0~53-gb265f1cf.m_armv7.apk
    ├── docker-buildx_0.7.0-r0~53-gb265f1cf.m_ppc64le.apk
    ├── docker-buildx_0.7.0-r0~53-gb265f1cf.m_riscv64.apk
    ├── docker-buildx_0.7.0-r0~53-gb265f1cf.m_s390x.apk
    ├── docker-buildx_0.7.0-r0~53-gb265f1cf.m_x86_64.apk
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian10_amd64.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian10_arm64.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian10_armel.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian10_armhf.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian10_ppc64el.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian10_riscv64.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian10_s390x.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian11_amd64.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian11_arm64.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian11_armel.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian11_armhf.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian11_ppc64el.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian11_riscv64.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-debian11_s390x.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian10_amd64.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian10_arm64.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian10_armel.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian10_armhf.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian10_ppc64el.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian10_riscv64.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian10_s390x.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian11_amd64.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian11_arm64.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian11_armel.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian11_armhf.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian11_ppc64el.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian11_riscv64.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-raspbian11_s390x.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu1804_amd64.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu1804_arm64.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu1804_armel.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu1804_armhf.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu1804_ppc64el.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu1804_riscv64.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu1804_s390x.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2004_amd64.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2004_arm64.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2004_armel.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2004_armhf.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2004_ppc64el.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2004_riscv64.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2004_s390x.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2104_amd64.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2104_arm64.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2104_armel.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2104_armhf.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2104_ppc64el.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2104_riscv64.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2104_s390x.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2110_amd64.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2110_arm64.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2110_armel.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2110_armhf.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2110_ppc64el.deb
    ├── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2110_riscv64.deb
    └── docker-buildx_0.7.0~53-gb265f1cf.m-ubuntu2110_s390x.deb
    ```

## Extract a subset of files/dirs

It can be useful to extract contents of a specific subset of files/dirs, so
you don't extract the whole file system if the source image is not a scratch
one.

=== "Command"

    ```shell
    $ undock --include /usr/local/bin --rm-dist --all crazymax/diun:latest ./dist
    ```

=== "Output result"

    ```text
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
    │   └── usr
    │       └── local
    │           └── bin
    │               └── diun
    ```

## Using the Docker image

You can also use the [official Docker image](../install/docker.md):

```shell
$ docker run --rm -t \
  -v $(pwd)/dist:/dist \
  crazymax/undock:latest \
    --include /usr/local/bin --all crazymax/diun:latest /dist
```

With this command, the cache will be deleted because the container is removed
as soon as the command ends, but you can define a volume to keep it:

```shell
$ docker run --rm -t \
  -v $(pwd)/dist:/dist \
  -v $(pwd)/cache:/var/cache/undock \
  crazymax/undock:latest \
    --include /usr/local/bin --all crazymax/diun:latest /dist
```
