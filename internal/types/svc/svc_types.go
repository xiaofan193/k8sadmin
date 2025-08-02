package svc

import (
	"github.com/xiaofan193/k8sadmin/internal/types"
	corev1 "k8s.io/api/core/v1"
)

type CreateorUpdateServiceRequest struct {
	Name      string              `json:"name"`
	Namespace string              `json:"namespace"`
	Labels    []types.ListMapItem `json:"labels"`
	Type      corev1.ServiceType  `json:"type"`
	Selector  []types.ListMapItem `json:"selector"`
	Ports     []ServicePort       `json:"ports"`
}

type ServicePort struct {
	Name       string `json:"name"`
	Port       int32  `json:"port"`
	TargetPort int32  `json:"targetPort"`
	NodePort   int32  `json:"nodePort"`
}

type CreateorUpdateServiceReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

type ServiceRes struct {
	Name       string              `json:"name"`
	Namespace  string              `json:"namespace"`
	Labels     []types.ListMapItem `json:"labels"`
	Type       corev1.ServiceType  `json:"type"`
	Selector   []types.ListMapItem `json:"selector"`
	Ports      []ServicePort       `json:"ports"`
	Age        int64               `json:"age"`
	ClusterIp  string              `json:"clusterIp"`
	ExternalIp []string            `json:"external"`
}

type ServiceResReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ServiceRes *ServiceRes
	} `json:"data"` // return data
}

type ServerResListReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		List []*ServiceRes
	} `json:"data"` // return data
}
