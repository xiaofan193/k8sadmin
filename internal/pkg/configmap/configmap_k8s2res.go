package configmap

import (
	"github.com/xiaofan193/k8sadmin/internal/pkg/maputils"
	"github.com/xiaofan193/k8sadmin/internal/types"
	corev1 "k8s.io/api/core/v1"
)

type K82Res struct {
}

func (K82Res) GeCmReqItem(configMap *corev1.ConfigMap) types.ConfigMapRes {
	return types.ConfigMapRes{
		Name:      configMap.Name,
		Namespace: configMap.Namespace,
		DataNum:   len(configMap.Data),
		Age:       configMap.CreationTimestamp.Unix(),
	}
}

func (this K82Res) GeCmReqDetail(configMap *corev1.ConfigMap) *types.ConfigMapRes {
	detail := this.GeCmReqItem(configMap)
	detail.Labels = maputils.ToList(configMap.Labels)
	detail.Data = maputils.ToList(configMap.Data)
	return &detail
}
