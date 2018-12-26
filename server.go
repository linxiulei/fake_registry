package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"
)

const (
	destLayerDigest  = "sha256:a5a6f2f73cd8abbdc55d0df0d8834f7262713e87d6c8800ea3851f103025e0f0"
	destLayerDiffIds = "sha256:a5a6f2f73cd8abbdc55d0df0d8834f7262713e87d6c8800ea3851f103025e0f0"
)

type Fake struct {
	LayerJson  string
	ConfigJson string
}

func NewFake() *Fake {
	configJson := `{
    "architecture": "amd64",
    "config": {
        "ArgsEscaped": true,
        "AttachStderr": false,
        "AttachStdin": false,
        "AttachStdout": false,
        "Cmd": [
            "python3"
        ],
        "Domainname": "",
        "Entrypoint": null,
        "Env": [
            "PATH=/usr/local/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
            "LANG=C.UTF-8",
            "GPG_KEY=0D96DF4D4110E5C43FBFB17F2D347EA6AA65421D",
            "PYTHON_VERSION=3.7.1",
            "PYTHON_PIP_VERSION=18.1"
        ],
        "Hostname": "",
        "Image": "sha256:5da82d22e1b08312469fd4e2662b820f398b59ac7c4f8876b88c5de04f66b5a2",
        "Labels": null,
        "OnBuild": [],
        "OpenStdin": false,
        "StdinOnce": false,
        "Tty": false,
        "User": "",
        "Volumes": null,
        "WorkingDir": ""
    },
    "container": "6b3bd4182966ad2e7f567e278e21e845cb3b82dce5d3bb1a6b1cc871ebfded57",
    "container_config": {
        "ArgsEscaped": true,
        "AttachStderr": false,
        "AttachStdin": false,
        "AttachStdout": false,
        "Cmd": [
            "/bin/sleep"
        ],
        "Domainname": "",
        "Entrypoint": null,
        "Env": [
            "PATH=/usr/local/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
            "LANG=C.UTF-8",
            "GPG_KEY=0D96DF4D4110E5C43FBFB17F2D347EA6AA65421D",
            "PYTHON_VERSION=3.7.1",
            "PYTHON_PIP_VERSION=18.1"
        ],
        "Hostname": "6b3bd4182966",
        "Image": "sha256:5da82d22e1b08312469fd4e2662b820f398b59ac7c4f8876b88c5de04f66b5a2",
        "Labels": {},
        "OnBuild": [],
        "OpenStdin": false,
        "StdinOnce": false,
        "Tty": false,
        "User": "",
        "Volumes": null,
        "WorkingDir": ""
    },
    "created": "2018-11-16T06:26:33.671893564Z",
    "docker_version": "17.06.2-ce",
    "history": [
        {
            "created": "2018-11-16T06:26:33.671893564Z",
            "created_by": "/bin/sh -c #(nop)  CMD [\"python3\"]",
            "empty_layer": true
        }
    ],
    "os": "linux",
    "rootfs": {
        "diff_ids": [
            "%s"
        ],
        "type": "layers"
    }
}`

	configJson = fmt.Sprintf(configJson, destLayerDiffIds)
	layerJson := `{
   "schemaVersion": 2,
   "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
   "config": {
      "mediaType": "application/vnd.docker.container.image.v1+json",
      "size": %d,
      "digest": "%s"
   },
   "layers": [
      {
         "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
         "size": 0,
         "digest": "%s"
      }
	]
}`
	layerJson = fmt.Sprintf(
		layerJson,
		len(configJson),
		fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(configJson))),
		destLayerDigest,
	)

	return &Fake{
		LayerJson:  layerJson,
		ConfigJson: configJson,
	}
}

func (f *Fake) GetLayerDigest() string {
	digest := fmt.Sprintf("sha256:%x",
		sha256.Sum256([]byte(f.LayerJson)))
	return digest
}

func (f *Fake) GetConfigDigest() string {
	digest := fmt.Sprintf("sha256:%x",
		sha256.Sum256([]byte(f.ConfigJson)))
	return digest
}

var fake *Fake

func rootManifestHandler(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Set(
		"Content-Type",
		"application/vnd.docker.distribution.manifest.v2+json")
	header.Set(
		"Docker-Content-Digest", fake.GetLayerDigest(),
	)
	fmt.Fprintf(w, fake.LayerJson)
}

func layerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, fake.LayerJson)
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, fake.ConfigJson)
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	fake = NewFake()
	http.HandleFunc("/v2/library/test/manifests/latest",
		rootManifestHandler)
	http.HandleFunc("/v2/library/test/manifests/"+fake.GetLayerDigest(),
		layerHandler)
	http.HandleFunc("/v2/library/test/blobs/"+fake.GetConfigDigest(),
		configHandler)
	fmt.Print(http.ListenAndServe("127.0.0.1:8084",
		logRequest(http.DefaultServeMux)))
}
