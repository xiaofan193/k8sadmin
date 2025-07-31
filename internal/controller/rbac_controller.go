package controller

import (
	"context"
	"github.com/xiaofan193/k8sadmin/internal/pkg/maputils"
	"github.com/xiaofan193/k8sadmin/internal/types/rbac"
	"github.com/xiaofan193/k8sadmin/pkg/global"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sync"
)

var (
	rbacInstance       *RbacpController
	rbacControllerOnce sync.Once
)

type RbacpController struct {
	KubeConfigSet *kubernetes.Clientset
	CONF          global.Server
}

func NewRbacController() *RbacpController {
	rbacControllerOnce.Do(func() {
		rbacInstance = &RbacpController{
			KubeConfigSet: global.GlobalKubeConfigSet,
		}
	})
	return rbacInstance
}

func (s *RbacpController) ServiceAccounts(ctx context.Context, namespace string, name string) ([]*rbac.ServiceAccount, error) {
	list, err := s.KubeConfigSet.CoreV1().ServiceAccounts(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	resList := make([]*rbac.ServiceAccount, 0)

	for _, item := range list.Items {
		resList = append(resList, &rbac.ServiceAccount{
			Name:      item.Name,
			Namespace: item.Namespace,
			Age:       item.CreationTimestamp.Unix(),
		})
	}

	return resList, err
}

func (s *RbacpController) CreateServiceAccount(ctx context.Context, reqParam *rbac.ServiceAccountRequest) error {
	saK8s := corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      reqParam.Name,
			Namespace: reqParam.Namespace,
			Labels:    maputils.ToMap(reqParam.Labels),
		},
	}

	_, err := s.KubeConfigSet.CoreV1().ServiceAccounts(reqParam.Namespace).Create(ctx, &saK8s, metav1.CreateOptions{})

	return err
}

func (s *RbacpController) DeleteServiceAccount(ctx context.Context, namespace string, name string) error {
	return s.KubeConfigSet.CoreV1().ServiceAccounts(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
