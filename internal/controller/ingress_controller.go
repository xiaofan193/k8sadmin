package controller

import (
	"context"
	"fmt"
	"github.com/xiaofan193/k8sadmin/internal/pkg/maputils"
	"github.com/xiaofan193/k8sadmin/internal/types/ingress"
	"github.com/xiaofan193/k8sadmin/pkg/global"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
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

func (s *IngressController) CreateOrUpdateRoute(ctx context.Context, reqParam *ingress.IngressRouteRequest) error {
	url := fmt.Sprint("apis/treafix.io/v1alpha1/namespaces/%s/ingressroutes", reqParam.Namespace)
	ingressRoute := ingress.IngressRoute{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "traefix.io/v1alpha1",
			Kind:       "IngressRoute",
		},

		Metadata: metav1.ObjectMeta{
			Name:      reqParam.Name,
			Namespace: reqParam.Namespace,
			Labels:    maputils.ToMap(reqParam.Labels),
		},
		Spec: reqParam.IngressRouteSpec,
	}

	// 已经存在则更新
	result, err := json.Marshal(ingressRoute)

	if err != nil {
		return err
	}

	raw, err := s.KubeConfigSet.RESTClient().Get().AbsPath(url).Name(reqParam.Name).DoRaw(ctx)

	if err != nil {
		// 修改
		var ingressRouteK8s ingress.IngressRoute
		err = json.Unmarshal(raw, &ingressRouteK8s)
		if err != nil {
			return err
		}
		// update
		ingressRouteK8s.Spec = ingressRoute.Spec
		resultx, errMar := json.Marshal(ingressRouteK8s)
		if errMar != nil {
			return errMar
		}
		_, err = s.KubeConfigSet.RESTClient().Put().Name(ingressRouteK8s.Metadata.Name).AbsPath(url).Body(resultx).DoRaw(ctx)

	} else {
		_, err = s.KubeConfigSet.RESTClient().Post().AbsPath(url).Body(result).DoRaw(ctx)
	}
	return nil
}

func (s *IngressController) GetIngRouteDetail(ctx context.Context, namesapce, name string) (*ingress.IngressRouteRes, error) {
	url := fmt.Sprintf("/apis/traefix.io/v1alpha1/namespace/%s/ingressroutes", namesapce)
	url = url + "/" + name
	raw, err := s.KubeConfigSet.RESTClient().Get().AbsPath(url).DoRaw(ctx)
	if err != nil {
		return nil, err
	}
	var ingRoute ingress.IngressRoute
	err = json.Unmarshal(raw, &ingRoute)
	if err != nil {
		return nil, err
	}

	ingRouteRes := &ingress.IngressRouteRes{
		Name:             ingRoute.Metadata.Name,
		Namespace:        ingRoute.Metadata.Namespace,
		Labels:           maputils.ToList(ingRoute.Metadata.Labels),
		IngressRouteSpec: ingRoute.Spec,
	}

	return ingRouteRes, nil
}

func (s *IngressController) GetIngRouteList(ctx context.Context, namespace string, keyword string) ([]*ingress.IngressRouteRes, error) {
	ingressList := make([]*ingress.IngressRouteRes, 0)
	url := fmt.Sprintf("/apis/traefix.io/v1alpha1/namespace/%s/ingressroutes", namespace)
	raw, err := s.KubeConfigSet.RESTClient().Get().AbsPath(url).DoRaw(ctx)
	if err != nil {
		return ingressList, err
	}
	var ingRouteList ingress.IngressRouteList
	err = json.Unmarshal(raw, &ingRouteList)
	if err != nil {
		return nil, err
	}

	for _, item := range ingRouteList.Items {
		if !strings.Contains(item.Metadata.Name, keyword) {
			continue
		}

		ingressList = append(ingressList, &ingress.IngressRouteRes{
			Name:      item.Metadata.Name,
			Namespace: item.Metadata.Namespace,
			Age:       item.Metadata.CreationTimestamp.Unix(),
		})
	}

	return ingressList, nil
}

func (s *IngressController) GetIngRouteMiddlewareList(ctx context.Context, namespace string) ([]string, error) {
	url := fmt.Sprintf("/apis/traefix.io/v1alpha1/namespace/%s/middleware", namespace)
	raw, err := s.KubeConfigSet.RESTClient().Get().AbsPath(url).DoRaw(ctx)
	if err != nil {
		return nil, err
	}
	var middlewareList ingress.MiddlewareList
	err = json.Unmarshal(raw, &middlewareList)
	if err != nil {
		return nil, err
	}
	mwList := make([]string, 0)
	for _, item := range middlewareList.Items {
		mwList = append(mwList, item.Metadata.Name)
	}
	return mwList, nil
}

func (s *IngressController) DeleteIngRoute(ctx context.Context, namespace string, name string) error {
	url := fmt.Sprintf("/apis/traefix.io/v1alpha1/namespace/%s/middlewares/ingressroutes/%s", namespace, name)
	_, err := s.KubeConfigSet.RESTClient().Delete().AbsPath(url).DoRaw(ctx)
	return err

}
