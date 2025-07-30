package controller

import (
	"context"
	"github.com/go-dev-frame/sponge/pkg/copier"
	"github.com/xiaofan193/k8sadmin/internal/pkg/configmap"
	"github.com/xiaofan193/k8sadmin/internal/types"
	"github.com/xiaofan193/k8sadmin/pkg/global"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sync"
)

var (
	configMapInstance       *ConfigMapController
	configMapControllerOnce sync.Once
)

type ConfigMapController struct {
	KubeConfigSet *kubernetes.Clientset
	CONF          global.Server
}

func NewConfigMapController() *ConfigMapController {
	configMapControllerOnce.Do(func() {
		configMapInstance = &ConfigMapController{
			KubeConfigSet: global.GlobalKubeConfigSet,
		}
	})
	return configMapInstance
}

func (c *ConfigMapController) CreateOrUpdateConfigMap(ctx context.Context, reqparam *types.CreateOrUpdateConfigMapRequest) error {
	configMap := &types.ConfigMap{}
	err := copier.Copy(configMap, reqparam)
	if err != nil {
		return err
	}

	// 将request 转为k8s结构
	req2K8s := &configmap.Req2K8s{}
	configMapObj := req2K8s.CmReq2K8sConvert(configMap)
	// 判断是否存在
	_, err = c.KubeConfigSet.CoreV1().ConfigMaps(configMap.Namespace).Get(ctx, configMap.Name, metav1.GetOptions{})
	if err != nil {
		_, err = c.KubeConfigSet.CoreV1().ConfigMaps(configMap.Namespace).Create(ctx, configMapObj, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		return nil
	}

	_, err = c.KubeConfigSet.CoreV1().ConfigMaps(configMap.Namespace).Update(ctx, configMapObj, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (c *ConfigMapController) GetConfigMapDetail(ctx context.Context, reqparam *types.GetConfigMapDetailORListRequest) (*types.ConfigMapRes, error) {
	confgiMapk8s, err := c.KubeConfigSet.CoreV1().ConfigMaps(reqparam.Namespace).Get(ctx, reqparam.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	k82Res := &configmap.K82Res{}
	cm := k82Res.GeCmReqDetail(confgiMapk8s)
	return cm, nil
}

func (c *ConfigMapController) GetConfigMapList(ctx context.Context, reqparam *types.GetConfigMapDetailORListRequest) ([]*types.ConfigMapRes, error) {
	list, err := c.KubeConfigSet.CoreV1().ConfigMaps(reqparam.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	k82Res := &configmap.K82Res{}
	configMapList := make([]*types.ConfigMapRes, 0)
	for _, cm := range list.Items {
		//if !strings.Contains(cm.Name,reqparam.Keyword) {
		//	continue
		//}
		cmRes := k82Res.GeCmReqDetail(&cm)
		configMapList = append(configMapList, cmRes)
	}
	return configMapList, nil
}

func (c *ConfigMapController) DeleteConfigMap(ctx context.Context, reqparam *types.DeleteConfigMapRequest) error {
	return c.KubeConfigSet.CoreV1().ConfigMaps(reqparam.Namespace).Delete(ctx, reqparam.Name, metav1.DeleteOptions{})
}
