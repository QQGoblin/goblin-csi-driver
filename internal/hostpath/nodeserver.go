package hostpath

import (
	"context"
	"fmt"
	"github.com/container-storage-interface/spec/lib/go/csi"
	csi_common "goblin-csi-driver/internal/csi-common"
	"goblin-csi-driver/internal/utils/config"
	"goblin-csi-driver/internal/utils/log"
	"k8s.io/utils/mount"
	"os"
)

type NodeServer struct {
	*csi_common.DefaultNodeServer
	hostPathConfig *config.HostPathConfig
}

func (n *NodeServer) NodePublishVolume(ctx context.Context, request *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {

	log.DebugObjMessage(request)

	// 获取 targetPath 路径
	volID := request.GetVolumeId()
	targetPath := request.GetTargetPath()
	path := fmt.Sprintf("%s/%s", n.hostPathConfig.State, volID)

	// 判断是否是 ephemeral 卷
	ephemeralVolume := request.GetVolumeContext()["csi.storage.k8s.io/ephemeral"] == "true"

	// 特殊处理 ephemeral vol
	if ephemeralVolume {
		// 对于 ephemeral 卷 k8s只会调用 NodePublishVolume 和 NodeUnPublishVolume，所以需要在这两个接口中进行创建和清理本地目录
		path := fmt.Sprintf("%s/ephemeral-%s", n.hostPathConfig.State, volID)
		err := os.MkdirAll(path, 0777)
		if err != nil {
			return nil, err
		}
	}

	// IsNotMountPoint 接口判断 targetPath 是否被挂在
	// 如果该目录不存在返回异常 ErrNotExist
	// 该接口通过检查 /proc/mounts 文件的内容，判断 targetPath 是否被挂载
	notMnt, err := mount.IsNotMountPoint(mount.New(""), targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(targetPath, 0750); err != nil {
				return nil, fmt.Errorf("create target path: %w", err)
			}
			notMnt = true
		} else {
			return nil, fmt.Errorf("check target path: %w", err)
		}
	}

	// 已经被挂载过则返回
	if !notMnt {
		return &csi.NodePublishVolumeResponse{}, nil
	}

	// 创建挂载参数
	readOnly := request.GetReadonly()
	options := []string{"bind"}
	if readOnly {
		options = append(options, "ro")
	}

	// 挂载目录，底层实际上执行的是 mount --bind path targetPath/mount
	mounter := mount.New("")
	if err := mounter.Mount(path, targetPath, "", options); err != nil {
		if ephemeralVolume {
			if rmErr := os.RemoveAll(path); rmErr != nil && !os.IsNotExist(rmErr) {
				log.ErrorLog(ctx, "clean %s failed, err: %v", path, rmErr)
			}
		}
		return nil, fmt.Errorf("failed to mount device: %s at %s: %s", path, targetPath, err)
	}

	return &csi.NodePublishVolumeResponse{}, nil
}

func (n *NodeServer) NodeUnpublishVolume(ctx context.Context, request *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {

	log.DebugObjMessage(request)

	// 检查卷是否有挂载
	notMnt, err := mount.IsNotMountPoint(mount.New(""), request.GetTargetPath())
	if err != nil {
		if os.IsNotExist(err) {
			notMnt = true
		} else {
			return nil, fmt.Errorf("check target path: %w", err)
		}
	}

	// 没有被挂载过则返回
	if notMnt {
		return &csi.NodeUnpublishVolumeResponse{}, nil
	}
	// 卸载目录
	mounter := mount.New("")
	if err := mounter.Unmount(request.GetTargetPath()); err != nil {
		return nil, fmt.Errorf("failed to unmount device: %s: %s", request, err)
	}

	// 删除 ephemeral 卷的目录
	ephemeralVolumePath := fmt.Sprintf("%s/ephemeral-%s", n.hostPathConfig.State, request.GetVolumeId())
	if err := os.RemoveAll(ephemeralVolumePath); err != nil {
		if !os.IsNotExist(err) {
			log.ErrorLog(ctx, "remove %s path err: %v", ephemeralVolumePath, err)
		}
	}

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

func (n *NodeServer) NodeStageVolume(ctx context.Context, request *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	log.DebugObjMessage(request)
	log.DebugLogMsg("没有实现 NodeStageVolume 操作逻辑")
	return &csi.NodeStageVolumeResponse{}, nil
}

func (n *NodeServer) NodeUnstageVolume(ctx context.Context, request *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	log.DebugObjMessage(request)
	log.DebugLogMsg("没有实现 NodeUnStageVolume 操作逻辑")
	return &csi.NodeUnstageVolumeResponse{}, nil
}
