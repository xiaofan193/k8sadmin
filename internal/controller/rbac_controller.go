package controller

import (
	"context"
	"github.com/xiaofan193/k8sadmin/internal/pkg/maputils"
	"github.com/xiaofan193/k8sadmin/internal/types/rbac"
	"github.com/xiaofan193/k8sadmin/pkg/global"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
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

func (s *RbacpController) GetRoleDetail(ctx context.Context, namespace string, name string) (*rbac.RoleRes, error) {
	resRole := &rbac.RoleRes{}
	// Role
	if namespace != "" {
		roleK8s, err := s.KubeConfigSet.RbacV1().Roles(namespace).Get(ctx, name, metav1.GetOptions{})

		if err != nil {
			return nil, err
		}
		resRole.Name = roleK8s.Name
		resRole.Namespace = roleK8s.Namespace
		resRole.Labels = maputils.ToList(roleK8s.Labels)
		resRole.Rules = roleK8s.Rules

	} else {
		roleK8s, err := s.KubeConfigSet.RbacV1().ClusterRoles().Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		resRole.Name = roleK8s.Name
		resRole.Namespace = roleK8s.Namespace
		resRole.Labels = maputils.ToList(roleK8s.Labels)
		resRole.Rules = roleK8s.Rules
	}
	return resRole, nil
	// ClusterRoles
}

func (s *RbacpController) GetRoleList(ctx context.Context, namespace string) ([]*rbac.Role, error) {
	resRoleList := make([]*rbac.Role, 0)
	if namespace != "" {
		roleK8sList, err := s.KubeConfigSet.RbacV1().Roles(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return resRoleList, err
		}
		for _, item := range roleK8sList.Items {
			resRoleList = append(resRoleList, &rbac.Role{
				Name:      item.Name,
				Namespace: item.Namespace,
				Age:       item.CreationTimestamp.Unix(),
			})
		}

	} else {
		roleK8sList, err := s.KubeConfigSet.RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{})
		if err != nil {
			return resRoleList, err
		}
		for _, item := range roleK8sList.Items {
			resRoleList = append(resRoleList, &rbac.Role{
				Name:      item.Name,
				Namespace: item.Namespace,
				Age:       item.CreationTimestamp.Unix(),
			})
		}
	}
	return resRoleList, nil
}

func (s *RbacpController) CreateOrUpdateRole(ctx context.Context, reqParam *rbac.RoleRequest) error {
	if reqParam.Namespace == "" {
		// clusterRole
		clusterRoleK8s := &rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name:      reqParam.Name,
				Namespace: reqParam.Namespace,
				Labels:    maputils.ToMap(reqParam.Labels),
			},
			Rules: reqParam.Rules,
		}
		clusterRoleSrc, err := s.KubeConfigSet.RbacV1().ClusterRoles().Get(ctx, reqParam.Name, metav1.GetOptions{})
		if err != nil {
			_, err = s.KubeConfigSet.RbacV1().ClusterRoles().Create(ctx, clusterRoleK8s, metav1.CreateOptions{})
			if err != nil {
				return err
			}
		} else {
			// update
			clusterRoleSrc.ObjectMeta.Labels = clusterRoleK8s.Labels
			clusterRoleSrc.Rules = clusterRoleK8s.Rules
			_, err = s.KubeConfigSet.RbacV1().ClusterRoles().Update(ctx, clusterRoleSrc, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		}
	} else {
		// role
		nsRoleK8sReq := &rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      reqParam.Name,
				Namespace: reqParam.Namespace,
				Labels:    maputils.ToMap(reqParam.Labels),
			},
			Rules: reqParam.Rules,
		}
		roleApi := s.KubeConfigSet.RbacV1().Roles(nsRoleK8sReq.Namespace)
		if nsRoleSrc, err := roleApi.Get(ctx, nsRoleK8sReq.Name, metav1.GetOptions{}); err != nil {
			_, err = roleApi.Create(ctx, nsRoleK8sReq, metav1.CreateOptions{})
			if err != nil {
				return err
			}
		} else {
			nsRoleSrc.Labels = nsRoleK8sReq.Labels
			nsRoleSrc.Rules = nsRoleK8sReq.Rules
			_, err = roleApi.Update(ctx, nsRoleSrc, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *RbacpController) CreateOrUpdateRolebing(ctx context.Context, reqParam *rbac.RoleBindingRequest) {

}
