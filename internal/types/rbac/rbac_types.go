package rbac

import "github.com/xiaofan193/k8sadmin/internal/types"

type ServiceAccount struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Age       int64  `json:"age"`
}

type Role struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Age       int64  `json:"age"`
}

type RoleBinding struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Age       int64  `json:"age"`
}

type ServiceAccountReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		List []*ServiceAccount
	} `json:"data"` // return data
}

type ServiceAccountRequest struct {
	Name      string              `json:"name"`
	Namespace string              `json:"namespace"`
	Labels    []types.ListMapItem `json:"labels"`
}
