package ingress

import (
	"github.com/xiaofan193/k8sadmin/internal/types"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

type IngressRouteRequest struct {
	Name      string              `json:"name"`
	Namespace string              `json:"namespace"`
	Labels    []types.ListMapItem `json:"labels"`
	IngressRouteSpec
}

// 得到k8s结构 -> response/request
type IngressRouteSpec struct {
	EntryPoints []string `json:"entryPoints"`
	Routes      []struct {
		Kind     string `json:"kind"`
		Match    string `json:"match"`
		Services []struct {
			Name string `json:"name"`
			Port int32  `json:"port"`
		} `json:"services"`
		Middlewares []struct {
			Name string `json:"name"`
		} `json:"middlewares"`
	} `json:"routes"`
	//配置tls证书 指针类型如果为空 tls->nil ->k8s忽略
	Tls *struct {
		SecretName string `json:"secretName"`
	} `json:"tls"`
}
type IngressRoute struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        metav1.ObjectMeta `json:"metadata"`
	Spec            IngressRouteSpec  `json:"spec"`
}
type IngressRouteList struct {
	Items           []IngressRoute `json:"items"`
	metav1.TypeMeta `json:",inline"`
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
}

type Middleware struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        metav1.ObjectMeta `json:"metadata"`
}

type MiddlewareList struct {
	Items           []Middleware `json:"items"`
	metav1.TypeMeta `json:",inline"`
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
}

type IngressRouteRes struct {
	Name      string              `json:"name"`
	Namespace string              `json:"namespace"`
	Labels    []types.ListMapItem `json:"labels"`
	IngressRouteSpec
	Age int64 `json:"age"`
}

type IngressRouteResReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		IngressRouteRes *IngressRouteRes
	} `json:"data"` // return data
}

type IngressRouteResListReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Data []*IngressRouteRes
	} `json:"data"` // return data
}

type IngressRouteMiddlewareListReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Data []string
	} `json:"data"` // return data
}
