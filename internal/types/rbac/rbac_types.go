package rbac

import "github.com/xiaofan193/k8sadmin/internal/types"
import rbacv1 "k8s.io/api/rbac/v1"

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

type RoleDetailReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Role *RoleRes
	} `json:"data"` // return data
}

type RoleRes struct {
	Name      string              `json:"name"`
	Namespace string              `json:"namespace"`
	Labels    []types.ListMapItem `json:"labels"`
	Rules     []rbacv1.PolicyRule `json:"rules"`
}

type RoleResListReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		List []*Role
	} `json:"data"` // return data
}

type RoleRequest struct {
	Name      string              `json:"name"`
	Namespace string              `json:"namespace"`
	Labels    []types.ListMapItem `json:"labels"`
	Rules     []rbacv1.PolicyRule `json:"rules"`
}

type RoleBindingRequest struct {
	Name      string              `json:"name"`
	Namespace string              `json:"namespace"`
	Labels    []types.ListMapItem `json:"labels"`
	//账号
	Subjects []ServiceAccount `json:"subjects"`
	//角色
	RoleRef string `json:"roleRef"`
}

type RoleBindingRes struct {
	Name      string              `json:"name"`
	Namespace string              `json:"namespace"`
	Labels    []types.ListMapItem `json:"labels"`
	//账号
	Subjects []ServiceAccount `json:"subjects"`
	//角色
	RoleRef string `json:"roleRef"`
	Age     int64  `json:"age"`
}
