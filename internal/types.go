package internal

const (

	// 默认 CSI。sock 目录（容器内的路径）
	DefaultEndPoint        string = "/csi/csi.sock"
	DefaultHistogramOption string = ""
	TopologyKeyNode               = "topology.goblin.hostpath.csi/node"

	// HostPath plugin Default
	DefaultHostPathState = "/var/lib/goblin/hostpath"
)

type CSIPluginType string

const (
	CSINodeServer       CSIPluginType = "node"
	CSIControllerServer CSIPluginType = "controller"
)

var CSIPluginTypeMap = map[string]CSIPluginType{
	"node":       CSINodeServer,
	"controller": CSIControllerServer,
}
