package controller

import (
	"context"
	"github.com/xiaofan193/k8sadmin/internal/pkg/maputils"
	"github.com/xiaofan193/k8sadmin/internal/types/svc"
	"github.com/xiaofan193/k8sadmin/pkg/global"
	corve1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"strings"
	"sync"
)

var (
	svcInstance       *SvcController
	svcControllerOnce sync.Once
)

type SvcController struct {
	KubeConfigSet *kubernetes.Clientset
	CONF          *global.Server
}

func NewSvcController() *SvcController {
	svcControllerOnce.Do(func() {
		svcInstance = &SvcController{
			KubeConfigSet: global.GlobalKubeConfigSet,
			CONF:          global.CONF,
		}
	})
	return svcInstance
}

func (s *SvcController) CreateOrUpdateSvc(ctx context.Context, reqParam *svc.CreateorUpdateServiceRequest) error {
	serverPorts := make([]corve1.ServicePort, 0)

	for _, port := range reqParam.Ports {
		serverPorts = append(serverPorts, corve1.ServicePort{
			Name: port.Name,
			Port: port.Port,
			TargetPort: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: port.TargetPort,
			},
			NodePort: port.NodePort,
		})
	}

	service := corve1.Service{

		ObjectMeta: metav1.ObjectMeta{
			Name:      reqParam.Name,
			Namespace: reqParam.Namespace,
			Labels:    maputils.ToMap(reqParam.Labels),
		},
		Spec: corve1.ServiceSpec{
			Type:     reqParam.Type,
			Selector: maputils.ToMap(reqParam.Labels),
			Ports:    serverPorts,
		},
	}

	serverApi := s.KubeConfigSet.CoreV1().Services(service.Namespace)
	serviceK8s, err := serverApi.Get(ctx, service.Name, metav1.GetOptions{})
	if err == nil {
		serviceK8s.Spec = service.Spec
		_, err = serverApi.Update(ctx, serviceK8s, metav1.UpdateOptions{})
	} else {
		_, err = serverApi.Create(ctx, &service, metav1.CreateOptions{})
	}
	return err
}

func (s *SvcController) GetSvcDetail(ctx context.Context, namespace string, name string) (*svc.ServiceRes, error) {
	serverK8s, err := s.KubeConfigSet.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	servicePorts := make([]svc.ServicePort, 0)
	for _, port := range serverK8s.Spec.Ports {
		servicePorts = append(servicePorts, svc.ServicePort{
			Name:       port.Name,
			Port:       port.Port,
			TargetPort: port.TargetPort.IntVal,
			NodePort:   port.NodePort,
		})
	}

	svcRes := &svc.ServiceRes{
		Name:      serverK8s.Name,
		Namespace: serverK8s.Namespace,
		Labels:    maputils.ToList(serverK8s.Spec.Selector),
		Ports:     servicePorts,
	}
	return svcRes, err
}

func (s *SvcController) GetSvcList(ctx context.Context, namespace string, keyword string) ([]*svc.ServiceRes, error) {
	list, err := s.KubeConfigSet.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	serverList := make([]*svc.ServiceRes, 0)
	for _, item := range list.Items {
		if !strings.Contains(item.Name, keyword) {
			continue
		}
		serverList = append(serverList, &svc.ServiceRes{
			Name:       item.Name,
			Namespace:  item.Namespace,
			Type:       item.Spec.Type,
			ClusterIp:  item.Spec.ClusterIP,
			ExternalIp: item.Spec.ExternalIPs,
			Age:        item.CreationTimestamp.Unix(),
		})
	}
	return serverList, nil
}

func (s *SvcController) DeleteSvc(ctx context.Context, namespace string, name string) error {
	return s.KubeConfigSet.CoreV1().Services(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
