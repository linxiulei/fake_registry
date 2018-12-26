# Intro

A simple, fake container image registry, only used to demostrate
secure issues among containerd

# Usage

Setup a host running containerd and there are images pulled you want to
steal.

Get the content digest and modify var `destLayerDigest` in the server.go
Also get the uncompressed content disgest with modifying var
`destLayerDiffIds`

Run the fake register server

```
go run server.go
```

Try to steal the content

```
ctr images pull http://localhost:8084/library/test:latest
```

The content will not need to be downloaded again, but shared from the
local.
