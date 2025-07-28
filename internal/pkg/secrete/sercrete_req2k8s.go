package secrete

import (
	"github.com/xiaofan193/k8sadmin/internal/pkg/maputils"
	"github.com/xiaofan193/k8sadmin/internal/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Req2K8s struct {
}

func (r *Req2K8s) SecretReq2K8sConvert(secret *types.Secret) corev1.Secret {
	return corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secret.Name,
			Namespace: secret.Namespace,
			Labels:    maputils.ToMap(secret.Labels),
		},
		Type:       secret.Type,
		StringData: maputils.ToMap(secret.Data),
	}
}
