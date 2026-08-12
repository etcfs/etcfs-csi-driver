package driver

import (
	"context"
	"errors"
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type nodeServer struct {
	csi.UnimplementedNodeServer
	cfg Config
}

func (s *nodeServer) NodeGetInfo(context.Context, *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	return &csi.NodeGetInfoResponse{NodeId: s.cfg.NodeID}, nil
}

func (s *nodeServer) NodeGetCapabilities(context.Context, *csi.NodeGetCapabilitiesRequest) (
	*csi.NodeGetCapabilitiesResponse, error) {
	return &csi.NodeGetCapabilitiesResponse{Capabilities: []*csi.NodeServiceCapability{{
		Type: &csi.NodeServiceCapability_Rpc{Rpc: &csi.NodeServiceCapability_RPC{
			Type: csi.NodeServiceCapability_RPC_GET_VOLUME_STATS,
		}},
	}}}, nil
}

// NodePublishVolume bind mounts the volume's directory at the target path.
//
// There is no NodeStageVolume: staging exists to mount a device once per node
// before sharing it between pods, and the EtcFS mount is already there,
// maintained by the daemon rather than by this driver.
func (s *nodeServer) NodePublishVolume(_ context.Context, req *csi.NodePublishVolumeRequest) (
	*csi.NodePublishVolumeResponse, error) {
	target := req.GetTargetPath()
	if target == "" {
		return nil, status.Error(codes.InvalidArgument, "target path is required")
	}
	if err := validateCapability(req.GetVolumeCapability()); err != nil {
		return nil, err
	}
	source, err := volumeDir(s.cfg.MountPath, req.GetVolumeId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	mounted, err := isMountPoint(s.cfg.MountPath)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition,
			"EtcFS mount %s is not usable on this node: %v", s.cfg.MountPath, err)
	}
	if !mounted {
		return nil, status.Errorf(codes.FailedPrecondition,
			"%s is not a mount point on this node: the EtcFS daemon is not running, and publishing "+
				"here would give the pod a local directory that looks shared but is not",
			s.cfg.MountPath)
	}

	// Static provisioning names a directory that may not exist yet; dynamic
	// provisioning already created it. Either way this is idempotent.
	if err := os.MkdirAll(source, 0o777); err != nil {
		return nil, status.Errorf(codes.Internal, "create volume directory: %v", err)
	}
	if err := os.MkdirAll(target, 0o750); err != nil {
		return nil, status.Errorf(codes.Internal, "create target path: %v", err)
	}

	if alreadyMounted, err := isMountPoint(target); err == nil && alreadyMounted {
		return &csi.NodePublishVolumeResponse{}, nil
	}

	flags := uintptr(unix.MS_BIND)
	if err := unix.Mount(source, target, "", flags, ""); err != nil {
		return nil, status.Errorf(codes.Internal, "bind mount %s at %s: %v", source, target, err)
	}
	if req.GetReadonly() {
		// A read-only bind mount takes two calls: the flag is ignored on the
		// initial MS_BIND and only takes effect on a remount.
		if err := unix.Mount("", target, "", unix.MS_BIND|unix.MS_REMOUNT|unix.MS_RDONLY, ""); err != nil {
			_ = unix.Unmount(target, unix.MNT_DETACH)
			return nil, status.Errorf(codes.Internal, "remount %s read-only: %v", target, err)
		}
	}
	return &csi.NodePublishVolumeResponse{}, nil
}

func (s *nodeServer) NodeUnpublishVolume(_ context.Context, req *csi.NodeUnpublishVolumeRequest) (
	*csi.NodeUnpublishVolumeResponse, error) {
	target := req.GetTargetPath()
	if target == "" {
		return nil, status.Error(codes.InvalidArgument, "target path is required")
	}
	// EINVAL means the path is not a mount and ENOENT that it is already gone;
	// both are the state this call is trying to reach, and kubelet retries the
	// call until it succeeds.
	if err := unix.Unmount(target, 0); err != nil &&
		!errors.Is(err, unix.EINVAL) && !errors.Is(err, unix.ENOENT) {
		return nil, status.Errorf(codes.Internal, "unmount %s: %v", target, err)
	}
	if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
		return nil, status.Errorf(codes.Internal, "remove target path: %v", err)
	}
	return &csi.NodeUnpublishVolumeResponse{}, nil
}

// NodeGetVolumeStats reports the filesystem's own numbers. EtcFS quotas are
// computed by a sweep rather than reflected in statfs, so a volume with a
// quota still reports cluster-wide capacity here; `etcfsctl quota` is what
// reports usage against the limit.
func (s *nodeServer) NodeGetVolumeStats(_ context.Context, req *csi.NodeGetVolumeStatsRequest) (
	*csi.NodeGetVolumeStatsResponse, error) {
	path := req.GetVolumePath()
	if path == "" {
		return nil, status.Error(codes.InvalidArgument, "volume path is required")
	}
	var st unix.Statfs_t
	if err := unix.Statfs(path, &st); err != nil {
		return nil, status.Errorf(codes.NotFound, "statfs %s: %v", path, err)
	}
	bsize := st.Bsize
	return &csi.NodeGetVolumeStatsResponse{Usage: []*csi.VolumeUsage{
		{
			Unit:      csi.VolumeUsage_BYTES,
			Total:     int64(st.Blocks) * bsize,
			Available: int64(st.Bavail) * bsize,
			Used:      int64(st.Blocks-st.Bfree) * bsize,
		},
		{
			Unit:      csi.VolumeUsage_INODES,
			Total:     int64(st.Files),
			Available: int64(st.Ffree),
			Used:      int64(st.Files - st.Ffree),
		},
	}}, nil
}
