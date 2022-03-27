package csi_common

import (
	"context"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"goblin-csi-driver/internal/util/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DefaultControllerServer points to default driver.
type DefaultControllerServer struct {
	Driver *CSIDriver
}

// ControllerPublishVolume publish volume on node.
func (cs *DefaultControllerServer) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ControllerUnpublishVolume unpublish on node.
func (cs *DefaultControllerServer) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ListVolumes lists volumes.
func (cs *DefaultControllerServer) ListVolumes(ctx context.Context, req *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// GetCapacity get volume capacity.
func (cs *DefaultControllerServer) GetCapacity(ctx context.Context, req *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ControllerGetCapabilities implements the default GRPC callout.
// Default supports all capabilities.
func (cs *DefaultControllerServer) ControllerGetCapabilities(ctx context.Context, req *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	log.TraceLog(ctx, "Using default ControllerGetCapabilities")
	if cs.Driver == nil {
		return nil, status.Error(codes.Unimplemented, "Controller server is not enabled")
	}

	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: cs.Driver.capabilities,
	}, nil
}

// ListSnapshots lists snapshots.
func (cs *DefaultControllerServer) ListSnapshots(ctx context.Context, req *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// ControllerGetVolume fetch volume information.
func (cs *DefaultControllerServer) ControllerGetVolume(ctx context.Context, req *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}
