package controller

import (
	"context"
	"github.com/xiaofan193/k8sadmin/internal/pkg/maputils"
	"github.com/xiaofan193/k8sadmin/internal/types/ingress"
	"github.com/xiaofan193/k8sadmin/pkg/global"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"

	"sync"

	networkingv1 "k8s.io/api/networking/v1"
)

var (
	ingressInstance       *IngressController
	ingressControllerOnce sync.Once
)

type IngressController struct {
	KubeConfigSet *kubernetes.Clientset
	CONF          *global.Server
}

func NewIngressController() *IngressController {
	ingressControllerOnce.Do(func() {
		ingressInstance = &IngressController{
			KubeConfigSet: global.GlobalKubeConfigSet,
			CONF:          global.CONF,
		}
	})
	return ingressInstance
}

func (s *IngressController) CreateOrUpdateIngress(ctx context.Context, reqParam *ingress.CreateOrUpdateIngressRequest) error {
	ingressRules := make([]networkingv1.IngressRule, 0)
	for _, rule := range reqParam.Rules {
		ingressRules = append(ingressRules, networkingv1.IngressRule{
			Host:             rule.Host,
			IngressRuleValue: rule.Value,
		})
	}

	ingress := networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      reqParam.Name,
			Namespace: reqParam.Namespace,
		},
		Spec: networkingv1.IngressSpec{
			Rules: ingressRules,
		},
	}

	ingressApi := s.KubeConfigSet.NetworkingV1().Ingresses(ingress.Namespace)
	ingressK8s, err := ingressApi.Get(ctx, ingress.Name, metav1.GetOptions{})
	if err != nil {
		ingressK8s.Spec = ingress.Spec
		_, err = ingressApi.Update(ctx, ingressK8s, metav1.UpdateOptions{})
	} else {
		_, err = ingressApi.Create(ctx, &ingress, metav1.CreateOptions{})
	}
	return err
}

func (s *IngressController) GetIngressDetail(ctx context.Context, namespace string, name string) (*ingress.IngressRes, error) {
	ingressK8s, err := s.KubeConfigSet.NetworkingV1().Ingresses(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	rules := make([]ingress.IngressRule, 0)
	for _, rule := range ingressK8s.Spec.Rules {
		rules = append(rules, ingress.IngressRule{
			Host:  rule.Host,
			Value: rule.IngressRuleValue,
		})
	}

	ingressRes := &ingress.IngressRes{
		Name:      ingressK8s.Name,
		Namespace: ingressK8s.Namespace,
		Labels:    maputils.ToList(ingressK8s.Labels),
		Rules:     rules,
	}
	return ingressRes, nil
}

func (s *IngressController) GetIngressList(ctx context.Context, namespace string) ([]*ingress.IngressRes, error) {
	list, err := s.KubeConfigSet.NetworkingV1().Ingresses(namespace).List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	ingressList := make([]*ingress.IngressRes, 0)
	for _, item := range list.Items {
		hosts := make([]string, 0)
		for _, rule := range item.Spec.Rules {
			hosts = append(hosts, rule.Host)
		}
		ingressList = append(ingressList, &ingress.IngressRes{
			Name:      item.Name,
			Namespace: item.Namespace,
			Hosts:     strings.Join(hosts, ","),
			Age:       item.CreationTimestamp.Unix(),
		})
	}
	return ingressList, nil
}

func (s *IngressController) DetIngress(ctx context.Context, namespace string, name string) error {
	return s.KubeConfigSet.NetworkingV1().Ingresses(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
