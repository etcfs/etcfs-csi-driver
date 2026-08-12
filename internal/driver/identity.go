package driver

import (
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type identityServer struct {
	csi.UnimplementedIdentityServer
	cfg Config
}

func (s *identityServer) GetPluginInfo(context.Context, *csi.GetPluginInfoRequest) (
	*csi.GetPluginInfoResponse, error) {
	return &csi.GetPluginInfoResponse{
		Name:          s.cfg.Name,
		VendorVersion: s.cfg.Version,
	}, nil
}

func (s *identityServer) GetPluginCapabilities(context.Context, *csi.GetPluginCapabilitiesRequest) (
	*csi.GetPluginCapabilitiesResponse, error) {
	caps := []*csi.PluginCapability{{
		Type: &csi.PluginCapability_Service_{Service: &csi.PluginCapability_Service{
			Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
		}},
	}}
	return &csi.GetPluginCapabilitiesResponse{Capabilities: caps}, nil
}

// Probe reports readiness. The driver has no asynchronous initialisation: once
// it is serving, it is ready, and a controller that has lost etcd should fail
// the individual call with the reason rather than take the whole plugin out of
// service.
func (s *identityServer) Probe(context.Context, *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	return &csi.ProbeResponse{Ready: wrapperspb.Bool(true)}, nil
}
