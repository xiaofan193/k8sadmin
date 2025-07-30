package controller

import (
	"context"
	"github.com/go-dev-frame/sponge/pkg/copier"
	"github.com/xiaofan193/k8sadmin/internal/pkg/secrete"
	"github.com/xiaofan193/k8sadmin/internal/types"
	"github.com/xiaofan193/k8sadmin/pkg/global"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
	"sync"
)

var (
	secreteInstance       *SecreteController
	secreteControllerOnce sync.Once
)

type SecreteController struct {
	KubeConfigSet *kubernetes.Clientset
	CONF          global.Server
}

func NewSecreteController() *SecreteController {
	secreteControllerOnce.Do(func() {
		secreteInstance = &SecreteController{
			KubeConfigSet: global.GlobalKubeConfigSet,
		}
	})
	return secreteInstance
}

func (s *SecreteController) CreateOrUpdateSecret(ctx context.Context, reqParam *types.CreateOrUpadteSecreteRequest) error {
	secretteConvert := &secrete.Req2K8s{}
	// 转换
	secrete := &types.Secret{}
	_ = copier.Copy(secrete, &reqParam)
	secretK8s := secretteConvert.SecretReq2K8sConvert(secrete)
	// 查询是否存在
	_, err := s.KubeConfigSet.CoreV1().Secrets(reqParam.Namespace).Get(ctx, reqParam.Name, metav1.GetOptions{})
	if err != nil {
		_, err = s.KubeConfigSet.CoreV1().Secrets(reqParam.Namespace).Create(ctx, &secretK8s, metav1.CreateOptions{})
	} else {
		_, err = s.KubeConfigSet.CoreV1().Secrets(reqParam.Namespace).Update(ctx, &secretK8s, metav1.UpdateOptions{})
	}
	return err
}

func (s *SecreteController) GetSecretList(ctx context.Context, namespace string, keyword string) ([]*types.Secret, error) {
	list, err := s.KubeConfigSet.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	secreteList := make([]*types.Secret, 0)
	secretConvert := &secrete.K8s2Res{}
	for _, item := range list.Items {

		if !strings.Contains(item.Name, keyword) {
			continue
		}

		secretRes := secretConvert.SecretK8s2ResItemConvert(item)
		secreteList = append(secreteList, secretRes)
	}
	return secreteList, err
}

func (s *SecreteController) GetSecretDetail(ctx context.Context, namespace string, name string) (*types.Secret, error) {
	secretK8s, err := s.KubeConfigSet.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})

	if err != nil {
		return nil, err
	}
	secretConvert := &secrete.K8s2Res{}
	secretRes := secretConvert.SecretK8s2ResDetailConvert(*secretK8s)
	return secretRes, err
}

func (s *SecreteController) DeleteSecret(ctx context.Context, namespace string, name string) error {
	return s.KubeConfigSet.CoreV1().Secrets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
