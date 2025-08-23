package controller

import (
	"context"
	"github.com/xiaofan193/k8sadmin/internal/pkg/maputils"
	"github.com/xiaofan193/k8sadmin/internal/pkg/pod"
	"github.com/xiaofan193/k8sadmin/internal/types"
	"github.com/xiaofan193/k8sadmin/pkg/global"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sync"
)

var (
	daemonsetInstance       *Daemonsetontroller
	daemonsetControllerOnce sync.Once
)

type Daemonsetontroller struct {
	KubeConfigSet *kubernetes.Clientset
	CONF          global.Server
}

func NewDaemonsetController() *Daemonsetontroller {
	daemonsetControllerOnce.Do(func() {
		daemonsetInstance = &Daemonsetontroller{
			KubeConfigSet: global.GlobalKubeConfigSet,
		}
	})
	return daemonsetInstance
}

func (s *Daemonsetontroller) CreateOrDaemonset(ctx context.Context, reqParam *types.DaemonsetReaqust) error {
	// 转换为k8s结构
	podK8sConvert := pod.Req2K8sConvert{}
	podK8s := podK8sConvert.PodReq2K8s(reqParam.Template)
	daemonset := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      reqParam.Base.Name,
			Namespace: reqParam.Base.Namespace,
			Labels:    maputils.ToMap(reqParam.Base.Labels),
		},
		Spec: appsv1.DaemonSetSpec{

			Selector: &metav1.LabelSelector{
				MatchLabels: maputils.ToMap(reqParam.Base.Selector),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: podK8s.ObjectMeta,
				Spec:       podK8s.Spec,
			},
		},
	}

	daemonsetApi := s.KubeConfigSet.AppsV1().DaemonSets(daemonset.Namespace)
	daemonsetK8s, err := daemonsetApi.Get(ctx, daemonset.Name, metav1.GetOptions{})
	if err == nil {
		daemonsetK8s.Spec = daemonset.Spec
		_, err = daemonsetApi.Update(ctx, daemonsetK8s, metav1.UpdateOptions{})
	} else {
		_, err = daemonsetApi.Create(ctx, daemonset, metav1.CreateOptions{})
	}
	return err
}

func (s *Daemonsetontroller) GetDaemonsetDetail(ctx context.Context, namespace string, name string) (*types.DaemonSetResonse, error) {
	daemonsetK8s, err := s.KubeConfigSet.AppsV1().DaemonSets(namespace).Get(ctx, name, metav1.GetOptions{})

	if err != nil {
		return nil, err
	}

	podResConvert := pod.K8s2ReqConvert{}

	podRes := podResConvert.PodK8s2Req(corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels: daemonsetK8s.Spec.Template.Labels,
		},
		Spec: daemonsetK8s.Spec.Template.Spec,
	})

	daemonsetRes := &types.DaemonSetResonse{
		Base: &types.DaemonsetBase{
			Name:      daemonsetK8s.Name,
			Namespace: daemonsetK8s.Namespace,

			Labels:   maputils.ToList(daemonsetK8s.Labels),
			Selector: maputils.ToList(daemonsetK8s.Spec.Selector.MatchLabels),
		},
		Template: &podRes,
	}
	return daemonsetRes, err
}

func (s *Daemonsetontroller) GetDaemonsetList(ctx context.Context, namespace string) ([]*types.DaemonSetRes, error) {
	daemonsetList := make([]*types.DaemonSetRes, 0)
	list, err := s.KubeConfigSet.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return daemonsetList, nil
	}

	for _, item := range list.Items {
		daemonsetList = append(daemonsetList, &types.DaemonSetRes{
			Name:      item.Name,
			Namespace: namespace,
			Age:       item.CreationTimestamp.Unix(),
			Ready:     item.Status.NumberReady,
			Available: item.Status.NumberAvailable,
			UpToData:  item.Status.UpdatedNumberScheduled,
			Current:   item.Status.CurrentNumberScheduled,
		})
	}

	return daemonsetList, err
}

func (s *Daemonsetontroller) DeleteDaemonset(ctx context.Context, namespace, name string) error {
	return s.KubeConfigSet.AppsV1().DaemonSets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
