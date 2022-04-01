package csi_common

import (
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"goblin-csi-driver/internal/utils/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
)

// CSIDriver stores hostpath information.
type CSIDriver struct {
	name    string
	nodeID  string
	version string
	// topology constraints that this nodeserver will advertise
	topology               map[string]string
	controllerCapabilities []*csi.ControllerServiceCapability
	nodeCapabilities       []*csi.NodeServiceCapability
	vc                     []*csi.VolumeCapability_AccessMode
}

// NewCSIDriver Creates a NewCSIDriver object. Assumes vendor
// version is equal to hostpath version &  does not support optional
// hostpath plugin info manifest field. Refer to CSI spec for more details.
func NewCSIDriver(name, v, nodeID string, topology map[string]string) *CSIDriver {
	if name == "" {
		klog.Errorf("Driver name missing")

		return nil
	}

	if nodeID == "" {
		klog.Errorf("NodeID missing")

		return nil
	}
	// TODO version format and validation
	if v == "" {
		klog.Errorf("Version argument missing")

		return nil
	}

	driver := CSIDriver{
		name:     name,
		version:  v,
		nodeID:   nodeID,
		topology: topology,
	}

	return &driver
}

// ValidateControllerServiceRequest validates the controller
// plugin controllerCapabilities.
func (d *CSIDriver) ValidateControllerServiceRequest(c csi.ControllerServiceCapability_RPC_Type) error {
	if c == csi.ControllerServiceCapability_RPC_UNKNOWN {
		return nil
	}

	for _, capability := range d.controllerCapabilities {
		if c == capability.GetRpc().GetType() {
			return nil
		}
	}

	return status.Error(codes.InvalidArgument, fmt.Sprintf("%s", c)) //nolint
}

// AddControllerServiceCapabilities stores the controller controllerCapabilities
// in hostpath object.
func (d *CSIDriver) AddControllerServiceCapabilities(cl []csi.ControllerServiceCapability_RPC_Type) {
	csc := make([]*csi.ControllerServiceCapability, 0, len(cl))

	for _, c := range cl {
		log.DefaultLog("Enabling controller service capability: %v", c.String())
		csc = append(csc, NewControllerServiceCapability(c))
	}

	d.controllerCapabilities = csc
}

func (d *CSIDriver) AddNodeServiceCapabilities(cl []csi.NodeServiceCapability_RPC_Type) {
	csc := make([]*csi.NodeServiceCapability, 0, len(cl))

	for _, c := range cl {
		log.DefaultLog("Enabling controller service capability: %v", c.String())
		csc = append(csc, NewNodeServiceCapability(c))
	}

	d.nodeCapabilities = csc
}

// AddVolumeCapabilityAccessModes stores volume access modes.
func (d *CSIDriver) AddVolumeCapabilityAccessModes(
	vc []csi.VolumeCapability_AccessMode_Mode) []*csi.VolumeCapability_AccessMode {
	vca := make([]*csi.VolumeCapability_AccessMode, 0, len(vc))
	for _, c := range vc {
		log.DefaultLog("Enabling volume access mode: %v", c.String())
		vca = append(vca, NewVolumeCapabilityAccessMode(c))
	}
	d.vc = vca

	return vca
}

// GetVolumeCapabilityAccessModes returns access modes.
func (d *CSIDriver) GetVolumeCapabilityAccessModes() []*csi.VolumeCapability_AccessMode {
	return d.vc
}
