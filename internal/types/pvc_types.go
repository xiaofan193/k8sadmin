package types

import corev1 "k8s.io/api/core/v1"

type PersistentVolumeClaimRequest struct {
	Name             string                              `json:"name"`
	Namespace        string                              `json:"namespace"`
	Labels           []ListMapItem                       `json:"labels"`
	AccessModes      []corev1.PersistentVolumeAccessMode `json:"accessModes"`
	Capacity         int32                               `json:"capacity"`
	Selector         []ListMapItem                       `json:"selector"`
	StorageClassName string                              `json:"storageClassName"`
}

type CreatePersistentVolumeClaimReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
	} `json:"data"` // return data
}
type PersistentVolumeClaimRes struct {
	Name             string                              `json:"name"`
	Namespace        string                              `json:"namespace"`
	Labels           []ListMapItem                       `json:"labels"`
	AccessModes      []corev1.PersistentVolumeAccessMode `json:"accessModes"`
	Capacity         int32                               `json:"capacity"`
	Selector         []ListMapItem                       `json:"selector"`
	StorageClassName string                              `json:"storageClassName"`
	Age              int64                               `json:"age"`
	Volume           string                              `json:"volume"`
	//pvc 状态
	Status corev1.PersistentVolumeClaimPhase `json:"status"`
}
type PersistentVolumeClaimResListReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		List []*PersistentVolumeClaimRes
	} `json:"data"` // return data
}
