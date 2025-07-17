package controller

import (
	"context"
	"fmt"
	"github.com/xiaofan193/k8sadmin/internal/pkg/pod"
	"github.com/xiaofan193/k8sadmin/internal/types"
	corev1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"strings"
	"sync"
	"time"
)

var (
	podInstance       *PodController
	podControllerOnce sync.Once
)

type PodController struct {
	KubeConfigSet *kubernetes.Clientset
	CONF          Server
}

func NewPodController() *PodController {
	podControllerOnce.Do(func() {
		podInstance = &PodController{
			KubeConfigSet: &kubernetes.Clientset{},
		}
	})
	return podInstance
}

func (p *PodController) CreateOrUpdatePod(ctx context.Context, podReq *types.Pod) (string, error) {
	// 把请求的pod结构数据转变为 k8s核心资源结构
	req2K8s := pod.Req2K8sConvert{}
	k8sPod := req2K8s.PodReq2K8s(podReq)
	podApi := p.KubeConfigSet.CoreV1().Pods(k8sPod.Namespace)
	k8sGetPod, err := podApi.Get(ctx, k8sPod.Name, metav1.GetOptions{})
	if err != nil {
		createdPod, err := podApi.Create(ctx, k8sPod, metav1.CreateOptions{})
		if err != nil {
			errMsg := fmt.Sprintf("Pod[namespace=%s,name=%s]创建失败，detail：%s", k8sPod.Namespace, k8sPod.Name, err.Error())
			return errMsg, err
		}

		successMsg := fmt.Sprintf("Pod[namespace=%s,name=%s]创建成功", createdPod.Namespace, createdPod.Name)
		return successMsg, err

	}

	// 更新
	// 先干运行尝试,确定参数合不合法
	k8sPodCopy := *k8sPod
	k8sPodCopy.Name = fmt.Sprintf("%s-validate", k8sPod.Name)
	_, err = podApi.Create(ctx, &k8sPodCopy, metav1.CreateOptions{DryRun: []string{metav1.DryRunAll}})

	if err != nil {
		errMsg := fmt.Sprintf("Pod[namespace=%s,name=%s]更新失败，detail：%s", k8sPod.Namespace, k8sPod.Name, err.Error())
		return errMsg, err
	}

	//比如pod处于terminating状态 监听pod删除完毕之后 才开始创建pod
	var labelSelector []string

	for k, v := range k8sGetPod.Labels {
		labelSelector = append(labelSelector, fmt.Sprintf("%=%s", k, v))
	}

	// label 的格式 app=test,app2=test2
	watcher, err := podApi.Watch(ctx, metav1.ListOptions{
		LabelSelector: strings.Join(labelSelector, ","),
	})

	if err != nil {
		errMsg := fmt.Sprintf("Pod[namespace = %s,name=%s]更新失败,desc:%s", k8sPod.Namespace, k8sPod.Name, err.Error())
		return errMsg, err
	}
	//删除 -- 强制删除
	background := metav1.DeletePropagationBackground
	var gracePeriodSeconds int64 = 0
	err = podApi.Delete(ctx, k8sPod.Name, metav1.DeleteOptions{
		GracePeriodSeconds: &gracePeriodSeconds,
		PropagationPolicy:  &background,
	})
	if err != nil {
		errMsg := fmt.Sprintf("Pod[namespace=%s,name=%s]更新失败，detail：%s", k8sPod.Namespace, k8sPod.Name, err.Error())
		return errMsg, err
	}

	for {
		select {
		case event := <-watcher.ResultChan():
			k8sPodChan := event.Object.(*corev1.Pod)
			//查询k8s 判断是否已经删除 那么就不用判断删除事件了
			if _, err := podApi.Get(ctx, k8sPod.Name, metav1.GetOptions{}); k8serror.IsNotFound(err) {
				//重新创建
				if createdPod, err := podApi.Create(ctx, k8sPod, metav1.CreateOptions{}); err != nil {
					errMsg := fmt.Sprintf("Pod[namespace=%s,name=%s]更新失败，detail：%s", k8sPod.Namespace, k8sPod.Name, err.Error())
					return errMsg, err
				} else {
					successMsg := fmt.Sprintf("Pod[namespace=%s,name=%s]更新成功", createdPod.Namespace, createdPod.Name)
					return successMsg, err
				}
			}
			switch event.Type {
			case watch.Deleted:
				if k8sPodChan.Name != k8sPod.Name {
					continue
				}
				//重新创建
				if createdPod, err := podApi.Create(ctx, k8sPod, metav1.CreateOptions{}); err != nil {
					errMsg := fmt.Sprintf("Pod[namespace=%s,name=%s]更新失败，detail：%s", k8sPod.Namespace, k8sPod.Name, err.Error())
					return errMsg, err
				} else {
					successMsg := fmt.Sprintf("Pod[namespace=%s,name=%s]更新成功", createdPod.Namespace, createdPod.Name)
					return successMsg, err
				}
			}
		case <-time.After(5 * time.Second):
			//重新创建
			if createdPod, err := podApi.Create(ctx, k8sPod, metav1.CreateOptions{}); err != nil {
				errMsg := fmt.Sprintf("Pod[namespace=%s,name=%s]更新失败，detail：%s", k8sPod.Namespace, k8sPod.Name, err.Error())
				return errMsg, err
			} else {
				successMsg := fmt.Sprintf("Pod[namespace=%s,name=%s]更新成功", createdPod.Namespace, createdPod.Name)
				return successMsg, err
			}
		}
	}
}

func (p *PodController) GetPodList(ctx context.Context, reqParam *types.GetPodListRequest) ([]*types.PodListItem, error) {
	list, err := p.KubeConfigSet.CoreV1().Pods(reqParam.Namespace).List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}
	podList := make([]*types.PodListItem, 0)
	k8s2req := &pod.K8s2ReqConvert{}
	for _, item := range list.Items {
		if reqParam.NodeName != "" && item.Spec.NodeName != reqParam.NodeName {
			continue
		}
		if strings.Contains(item.Name, reqParam.Keyword) {
			podItem := k8s2req.PodK8s2ItemRes(item)
			podList = append(podList, &podItem)
		}
	}
	return podList, err
}

func (p *PodController) GetPodDetail(ctx context.Context, reqParam *types.GetPodDetailRequest) (*types.Pod, error) {
	podApi := p.KubeConfigSet.CoreV1().Pods(reqParam.Namespace)
	k8sGetPod, err := podApi.Get(ctx, reqParam.Name, metav1.GetOptions{})

	if err != nil {
		return nil, fmt.Errorf("Pod[namespace=%s,name=%s]查询失败，detail：%s", reqParam.Namespace, reqParam.Name, err.Error())
	}

	// 将k8s pod 转为  pod response
	k8s2req := &pod.K8s2ReqConvert{}
	podRes := k8s2req.PodK8s2Req(*k8sGetPod)
	return &podRes, nil
}
