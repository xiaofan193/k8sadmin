package types

import corev1 "k8s.io/api/core/v1"

type NfsVolumeSource struct {
	NfsPath     string `json:"nfsPath"`
	NfsServer   string `json:"nfsServer"`
	NfsReadOnly bool   `json:"nfsReadOnly"`
}
type VolumeSource struct {
	Type            string          `json:"type"`
	NfsVolumeSource NfsVolumeSource `json:"nfsVolumeSource"`
}
type PersistentVolume struct {
	Name string `json:"name"`
	//ns 不必传
	//Namespace string             `json:"namespace"`
	Labels []ListMapItem `json:"labels"`
	//pv容量
	Capacity int32 `json:"capacity"`
	//数据读写权限
	AccessModes []corev1.PersistentVolumeAccessMode `json:"accessModes"`
	//pv回收策略
	ReClaimPolicy corev1.PersistentVolumeReclaimPolicy `json:"reClaimPolicy"`
	VolumeSource  VolumeSource                         `json:"volumeSource"`
}

type PersistentVolumeRequest struct {
	Name string `json:"name"`
	//ns 不必传
	//Namespace string             `json:"namespace"`
	Labels []ListMapItem `json:"labels"`
	//pv容量
	Capacity int32 `json:"capacity"`
	//数据读写权限
	AccessModes []corev1.PersistentVolumeAccessMode `json:"accessModes"`
	//pv回收策略
	ReClaimPolicy corev1.PersistentVolumeReclaimPolicy `json:"reClaimPolicy"`
	VolumeSource  VolumeSource                         `json:"volumeSource"`
}

type CreatePersistentVolumeRequestReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
	} `json:"data"` // return data
}

type PersistentVolumeRes struct {
	Name string `json:"name"`
	//pv容量
	Capacity int32 `json:"capacity"`
	//ns 不必传
	//Namespace string             `json:"namespace"`
	Labels []ListMapItem `json:"labels"`
	//数据读写权限
	AccessModes []corev1.PersistentVolumeAccessMode `json:"accessModes"`
	//pv回收策略
	ReClaimPolicy corev1.PersistentVolumeReclaimPolicy `json:"reClaimPolicy"`
	//todo 待完善
	Status corev1.PersistentVolumePhase `json:"status"`
	//被具备某个pvc绑定
	Claim string `json:"claim"`
	//创建时间
	Age int64 `json:"age"`
	//状况描述
	Reason string `json:"reason"`
	//sc 名称
	StorageClassName string `json:"storageClassName"`
}

type PersistentVolumeResListReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		List []*PersistentVolumeRes
	} `json:"data"` // return data
}
