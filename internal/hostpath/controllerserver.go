package hostpath

import (
	"context"
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/pborman/uuid"
	csi_common "goblin-csi-driver/internal/csi-common"
	"goblin-csi-driver/internal/utils/config"
	"goblin-csi-driver/internal/utils/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ControllerServer struct {
	*csi_common.DefaultControllerServer
	hostPathConfig *config.HostPathConfig
}

func (cs *ControllerServer) CreateVolume(ctx context.Context, request *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	log.DebugObjMessage(request)

	// Check arguments
	if len(request.GetName()) == 0 {
		// csi-provisioner 生成的 volumeName 是 pvc-<pvc uid>
		return nil, status.Error(codes.InvalidArgument, "Name missing in request")
	}

	// Create Volume 的主要任务只是生成 volumeID,提供给后续接口使用
	// volumeID 可以考虑直接从 VolumeName 中获取 PVC UID，方便后续管理
	volumeID := uuid.NewUUID().String()

	// 对于 accessibilityRequirements.GetRequisite() 属性，csi-provisioner:v2.0.4 代码的逻辑如下：
	//
	// 如果 csi-provisioner 不开启 strict-topology 那么生成 AccessibilityRequirements 的方式如下：
	// 1. 随机选择一个节点的 CSINode，获取对应的 topologyKeys
	// 2. 在所有 Node 上获取 topologyKeys 对应的 LabelKey-LabelValue
	// 3. 根据获取的 LabelKey-LabelValue 形成 accessibilityRequirements

	// 如果 csi-provisioner 开启 strict-topology 那么生成 AccessibilityRequirements 的方式如下：
	// 1. 根据 volume.kubernetes.io/selected-node 的值，获取对应 CSINode 的 topologyKeys
	// 2. 从 selectNode 上获取 topologyKeys 对应的 LabelKey-LabelValue
	// 3. 根据获取的 LabelKey-LabelValue 形成 selectedTopology
	// 4. 验证 StorageClass 中指定的 allowedTopologies 是否是 selectedTopology 的子集（即 Selected-Node 是否满足 SC 的限制条件）
	// 5. 输出 selectedTopology 形成 accessibilityRequirements

	// PS：所以这里需要开启 strict-topology 配置

	accessibilityRequirements := request.GetAccessibilityRequirements()

	// // TODO: 需要通过 API 或者其他手段，使 ControllerServer 进程能够访问对应的节点，并创建目录

	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:           volumeID,
			CapacityBytes:      request.GetCapacityRange().GetRequiredBytes(),
			VolumeContext:      request.GetParameters(),
			ContentSource:      request.GetVolumeContentSource(),
			AccessibleTopology: accessibilityRequirements.GetRequisite(),
		},
	}, nil
}

func (cs *ControllerServer) DeleteVolume(ctx context.Context, request *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	log.DebugObjMessage(request)

	// Check arguments
	if len(request.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}

	// TODO: 需要根据 VolumeId 判断 PV 位于哪个节点，并且删除该节点的指定目录（VolumeID 记录在 PV.spec.csi.volumeHandle 字段）
	// TODO: 需要通过 API 或者其他手段，使 ControllerServer 进程能够访问对应的节点

	return &csi.DeleteVolumeResponse{}, nil
}
