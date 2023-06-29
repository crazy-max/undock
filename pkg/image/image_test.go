package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		desc     string
		name     string
		expected Image
	}{
		{
			desc: "bintray artifactory-oss",
			name: "jfrog-docker-reg2.bintray.io/jfrog/artifactory-oss:4.0.0",
			expected: Image{
				Domain: "jfrog-docker-reg2.bintray.io",
				Path:   "jfrog/artifactory-oss",
				Tag:    "4.0.0",
			},
		},
		{
			desc: "bintray xray-server",
			name: "docker.bintray.io/jfrog/xray-server:2.8.6",
			expected: Image{
				Domain: "docker.bintray.io",
				Path:   "jfrog/xray-server",
				Tag:    "2.8.6",
			},
		},
		{
			desc: "dockerhub alpine",
			name: "alpine",
			expected: Image{
				Domain: "docker.io",
				Path:   "library/alpine",
				Tag:    "latest",
			},
		},
		{
			desc: "dockerhub crazymax/nextcloud",
			name: "docker.io/crazymax/nextcloud:latest",
			expected: Image{
				Domain: "docker.io",
				Path:   "crazymax/nextcloud",
				Tag:    "latest",
			},
		},
		{
			desc: "gcr busybox",
			name: "gcr.io/google-containers/busybox:latest",
			expected: Image{
				Domain: "gcr.io",
				Path:   "google-containers/busybox",
				Tag:    "latest",
			},
		},
		{
			desc: "gcr busybox tag/digest",
			name: "gcr.io/google-containers/busybox:latest" + sha256digest,
			expected: Image{
				Domain: "gcr.io",
				Path:   "google-containers/busybox",
				Tag:    "latest",
				Digest: sha256digest,
			},
		},
		{
			desc: "github ddns-route53",
			name: "docker.pkg.github.com/crazy-max/ddns-route53/ddns-route53:latest",
			expected: Image{
				Domain: "docker.pkg.github.com",
				Path:   "crazy-max/ddns-route53/ddns-route53",
				Tag:    "latest",
			},
		},
		{
			desc: "gitlab meltano",
			name: "registry.gitlab.com/meltano/meltano",
			expected: Image{
				Domain: "registry.gitlab.com",
				Path:   "meltano/meltano",
				Tag:    "latest",
			},
		},
		{
			desc: "quay hypercube",
			name: "quay.io/coreos/hyperkube",
			expected: Image{
				Domain: "quay.io",
				Path:   "coreos/hyperkube",
				Tag:    "latest",
			},
		},
		{
			desc: "ghcr ddns-route53",
			name: "ghcr.io/crazy-max/ddns-route53",
			expected: Image{
				Domain: "ghcr.io",
				Path:   "crazy-max/ddns-route53",
				Tag:    "latest",
			},
		},
		{
			desc: "ghcr radarr",
			name: "ghcr.io/linuxserver/radarr",
			expected: Image{
				Domain: "ghcr.io",
				Path:   "linuxserver/radarr",
				Tag:    "latest",
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			img, err := Parse(tt.name)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tt.expected.Domain, img.Domain)
			assert.Equal(t, tt.expected.Path, img.Path)
			assert.Equal(t, tt.expected.Tag, img.Tag)
		})
	}
}
