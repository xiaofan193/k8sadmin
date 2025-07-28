package secrete

import (
	"github.com/xiaofan193/k8sadmin/internal/pkg/maputils"
	"github.com/xiaofan193/k8sadmin/internal/types"
	corev1 "k8s.io/api/core/v1"
)

type K8s2Res struct {
}

func (K8s2Res) SecretK8s2ResItemConvert(secret corev1.Secret) *types.Secret {
	return &types.Secret{
		Name:      secret.Name,
		Namespace: secret.Namespace,
		Type:      secret.Type,
		DataNum:   len(secret.Data),
		Age:       secret.CreationTimestamp.Unix(),
	}
}

func (K8s2Res) SecretK8s2ResDetailConvert(secret corev1.Secret) *types.Secret {
	return &types.Secret{
		Name:      secret.Name,
		Namespace: secret.Namespace,
		Type:      secret.Type,
		DataNum:   len(secret.Data),
		Age:       secret.CreationTimestamp.Unix(),
		Data:      maputils.ToListWithMapByte(secret.Data),
		Labels:    maputils.ToList(secret.Labels),
	}
}
