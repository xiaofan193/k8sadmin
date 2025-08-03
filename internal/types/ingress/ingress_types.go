package ingress

import (
	"github.com/xiaofan193/k8sadmin/internal/types"
	networkingv1 "k8s.io/api/networking/v1"
)

type IngressRule struct {
	Host  string                        `json:"host"`
	Value networkingv1.IngressRuleValue `json:"value"`
}
type CreateOrUpdateIngressRequest struct {
	Name      string              `json:"name"`
	Namespace string              `json:"namespace"`
	Labels    []types.ListMapItem `json:"labels"`
	Rules     []IngressRule       `json:"rules"`
}

type IngressRes struct {
	Name      string              `json:"name"`
	Namespace string              `json:"namespace"`
	Labels    []types.ListMapItem `json:"labels"`
	Rules     []IngressRule       `json:"rules"`
	Class     string              `json:"class"`
	Hosts     string              `json:"hosts"`
	Age       int64               `json:"age"`
}

type IngressDetailReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		IngressRes *IngressRes
	} `json:"data"` // return data
}

type IngressListReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		List []*IngressRes
	} `json:"data"` // return data
}
