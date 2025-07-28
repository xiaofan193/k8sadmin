package types

import corev1 "k8s.io/api/core/v1"

type Secret struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	DataNum   int               `json:"dataNum"`
	Age       int64             `json:"age"`
	Type      corev1.SecretType `json:"type"`
	Labels    []ListMapItem     `json:"labels"`
	Data      []ListMapItem     `json:"data"`
}
type CreateOrUpadteSecreteRequest struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Type      corev1.SecretType `json:"type"`
	Labels    []ListMapItem     `json:"labels"`
	Data      []ListMapItem     `json:"data"`
}

type CreateOrUpadteReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}
type ListSecretItemReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		List []*Secret `json:"list"`
	} `json:"data"` // return data
}

type GetSecretReply struct {
	Code int     `json:"code"` // return code
	Msg  string  `json:"msg"`  // return information description
	Data *Secret `json:"data"` // return data
}

type DeleteSecretReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
	} `json:"data"` // return data
}
