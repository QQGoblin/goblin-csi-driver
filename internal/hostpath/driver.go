package hostpath

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	"goblin-csi-driver/internal"
	csi_common "goblin-csi-driver/internal/csi-common"
	"goblin-csi-driver/internal/utils/config"
	"goblin-csi-driver/internal/utils/log"
)

type Driver struct {
	drivername string
	nodeid     string
	endpoint   string
	csiDriver  *csi_common.CSIDriver
	config     *config.HostPathConfig
	servers    csi_common.Servers
}

func NewDriver(drivername, nodeid, endpoint string, c *config.Options) *Driver {
	return &Driver{
		drivername: drivername,
		nodeid:     nodeid,
		endpoint:   endpoint,
		config:     c.HostPathConfig,
		servers:    csi_common.Servers{},
	}
}

func (r *Driver) NewServers() {

	r.csiDriver.AddControllerServiceCapabilities(
		[]csi.ControllerServiceCapability_RPC_Type{
			csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
		},
	)

	r.csiDriver.AddNodeServiceCapabilities(
		[]csi.NodeServiceCapability_RPC_Type{
			csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
		},
	)

	r.servers.IS = &IdentityServer{
		DefaultIdentityServer: csi_common.NewDefaultIdentityServer(r.csiDriver),
	}

	r.servers.CS = &ControllerServer{
		DefaultControllerServer: csi_common.NewDefaultControllerServer(r.csiDriver),
		hostPathConfig:          r.config,
	}

	r.servers.NS = &NodeServer{
		DefaultNodeServer: csi_common.NewDefaultNodeServer(r.csiDriver),
		hostPathConfig:    r.config,
	}

}

func (r *Driver) Run() {

	// 配置节点亲和拓扑
	topology := map[string]string{
		internal.TopologyKeyNode: r.nodeid,
	}

	// 初始化 IdentityServer 和 NodeServer
	r.csiDriver = csi_common.NewCSIDriver(r.drivername, internal.DriverVersion, r.nodeid, topology)

	if r.csiDriver == nil {
		log.FatalLogMsg("Failed to initialize CSI Driver.")
	}

	// 初始化 server
	r.NewServers()

	// 启动 CSI gRPC 服务
	s := csi_common.NewNonBlockingGRPCServer()
	s.Start(r.endpoint, internal.DefaultHistogramOption, r.servers)
	s.Wait()
}
