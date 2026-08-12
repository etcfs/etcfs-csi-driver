// The CSI driver is a nested module so the container-storage-interface and
// gRPC dependency graph stays out of the root module, which every other binary
// in the repository builds from.  The replace keeps the two halves in lockstep:
// the fencing semantics the controller relies on are the ones in this working
// tree, not those of an older published tag.
module github.com/MHS-20/EtcFS/csi

go 1.24.0

require (
	github.com/MHS-20/EtcFS v0.9.0
	github.com/container-storage-interface/spec v1.9.0
	go.etcd.io/etcd/client/v3 v3.5.18
	golang.org/x/sys v0.29.0
	google.golang.org/grpc v1.66.0
	google.golang.org/protobuf v1.36.1
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_golang v1.21.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	go.etcd.io/etcd/api/v3 v3.5.18 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.18 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.17.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240903143218-8af14fe29dc1 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
)

replace github.com/MHS-20/EtcFS => ../
