package types

import (
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
)

type StorageClassRequest struct {
	Name string `json:"name"`
	//Namespace string             `json:"namespace"`
	Labels []ListMapItem `json:"labels"`
	//制备器  AWS EBS, GCP PD, Azure Disk/File, vSphere NFS 的nfs-subdir-external-provisioner
	Provisioner string `json:"provisioner"`
	//卷绑定参数配置  挂载选项
	MountOptions []string `json:"mountOptions"`
	//制备器入参  archiveOnDelete: "false" 就是自动回收存储空间
	Parameters []ListMapItem `json:"parameters"`
	//卷回收策略 Delete Retain PV 会转为 Released 状态。
	ReclaimPolicy corev1.PersistentVolumeReclaimPolicy `json:"reclaimPolicy"`
	//是否允许扩充
	AllowVolumeExpansion bool `json:"allowVolumeExpansion"`
	//卷绑定模式
	VolumeBindingMode storagev1.VolumeBindingMode `json:"volumeBindingMode"`
}

type StorageClassReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
	} `json:"data"` // return data
}
type StorageClassRes struct {
	Name string `json:"name"`
	//Namespace string             `json:"namespace"`
	Labels []ListMapItem `json:"labels"`
	//制备器
	Provisioner string `json:"provisioner"`
	//卷绑定参数配置
	MountOptions []string `json:"mountOptions"`
	//制备器入参
	Parameters []ListMapItem `json:"parameters"`
	//卷回收策略
	ReclaimPolicy corev1.PersistentVolumeReclaimPolicy `json:"reclaimPolicy"`
	//是否允许扩充
	AllowVolumeExpansion bool `json:"allowVolumeExpansion"`
	//卷绑定模式
	VolumeBindingMode storagev1.VolumeBindingMode `json:"volumeBindingMode"`
	Age               int64                       `json:"age"`
}

type StorageClassListReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		List []*StorageClassRes `json:"list"`
	} `json:"data"` // return data
}
type DeleteStorageClassReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
	} `json:"data"` // return data
}
