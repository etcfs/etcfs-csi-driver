# etcfs-csi-driver

Kubernetes CSI driver for [EtcFS](https://github.com/etcfs/etcfs) — provisions
and mounts EtcFS volumes as `ReadWriteMany` `Block` volumes, with fencing
handled through the core EtcFS fencing protocol.

Split out of the main [etcfs/etcfs](https://github.com/etcfs/etcfs) repository
(formerly the `csi/` nested Go module there); full commit history for this
driver was preserved in the split.

See [github.com/etcfs/etcfs-docs](https://github.com/etcfs/etcfs-docs) for the
Kubernetes deployment guide, and `deploy/helm/etcfs-csi` here for the Helm
chart.

## Layout

- `cmd/etcfs-csi` — driver entrypoint
- `internal/driver` — CSI controller/node/identity implementation
- `internal/etcdtest` — etcd test helper (vendored from the core repo's
  integration test suite)
- `deploy/helm/etcfs-csi` — Helm chart
- `examples/` — static and dynamic provisioning examples

## Building

```bash
go build ./cmd/etcfs-csi
```

Requires [github.com/etcfs/etcfs](https://github.com/etcfs/etcfs) as a
published Go module dependency (`pkg/metadata`).
