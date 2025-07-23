package configmap

import (
	"github.com/xiaofan193/k8sadmin/internal/pkg/maputils"
	"github.com/xiaofan193/k8sadmin/internal/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Req2K8s struct {
}

func (*Req2K8s) CmReq2K8sConvert(configMapReq *types.ConfigMap) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapReq.Name,
			Namespace: configMapReq.Namespace,
			Labels:    maputils.ToMap(configMapReq.Labels),
		},
		Data: maputils.ToMap(configMapReq.Data),
	}
}
