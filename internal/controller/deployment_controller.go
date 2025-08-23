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
	deploymentInstance       *DeploymentController
	deploymentControllerOnce sync.Once
)

type DeploymentController struct {
	KubeConfigSet *kubernetes.Clientset
	CONF          global.Server
}

func NewDeploymentController() *DeploymentController {
	deploymentControllerOnce.Do(func() {
		deploymentInstance = &DeploymentController{
			KubeConfigSet: global.GlobalKubeConfigSet,
		}
	})
	return deploymentInstance
}

func (s *DeploymentController) CreateOrDeployment(ctx context.Context, reqParam *types.DeploymentRequest) error {
	// 转换为k8s结构
	podK8sConvert := pod.Req2K8sConvert{}
	podK8s := podK8sConvert.PodReq2K8s(reqParam.Template)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      reqParam.Base.Name,
			Namespace: reqParam.Base.Namespace,
			Labels:    maputils.ToMap(reqParam.Base.Labels),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &reqParam.Base.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: maputils.ToMap(reqParam.Base.Selector),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: podK8s.ObjectMeta,
				Spec:       podK8s.Spec,
			},
		},
	}

	deploymentApi := s.KubeConfigSet.AppsV1().Deployments(deployment.Namespace)
	deploymentK8s, err := deploymentApi.Get(ctx, deployment.Name, metav1.GetOptions{})
	if err == nil {
		deploymentK8s.Spec = deployment.Spec
		_, err = deploymentApi.Update(ctx, deploymentK8s, metav1.UpdateOptions{})
	} else {
		_, err = deploymentApi.Create(ctx, deployment, metav1.CreateOptions{})
	}
	return err
}

func (s *DeploymentController) GetDeploymentDetail(ctx context.Context, namespace string, name string) (*types.DeploymentResponse, error) {
	deploymentK8s, err := s.KubeConfigSet.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})

	if err != nil {
		return nil, err
	}

	podResConvert := pod.K8s2ReqConvert{}

	podRes := podResConvert.PodK8s2Req(corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels: deploymentK8s.Spec.Template.Labels,
		},
		Spec: deploymentK8s.Spec.Template.Spec,
	})

	deploymentRes := &types.DeploymentResponse{
		Base: &types.DeploymentBase{
			Name:      deploymentK8s.Name,
			Namespace: deploymentK8s.Namespace,
			Replicas:  *deploymentK8s.Spec.Replicas,
			Labels:    maputils.ToList(deploymentK8s.Labels),
			Selector:  maputils.ToList(deploymentK8s.Spec.Selector.MatchLabels),
		},
		Template: &podRes,
	}
	return deploymentRes, err
}

func (s *DeploymentController) GetDeploymentList(ctx context.Context, namespace string) ([]*types.DeploymentRes, error) {
	deploymentList := make([]*types.DeploymentRes, 0)
	list, err := s.KubeConfigSet.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return deploymentList, nil
	}

	for _, item := range list.Items {
		deploymentList = append(deploymentList, &types.DeploymentRes{
			Name:       item.Name,
			Namespace:  namespace,
			Age:        item.CreationTimestamp.Unix(),
			Replicas:   *item.Spec.Replicas,
			Available:  item.Status.AvailableReplicas,
			UpdateDate: item.Status.UpdatedReplicas,
		})
	}

	return deploymentList, err
}

func (s *DeploymentController) DeleteDeployment(ctx context.Context, namespace, name string) error {
	return s.KubeConfigSet.AppsV1().Deployments(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
