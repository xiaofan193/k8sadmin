package types

import corev1 "k8s.io/api/core/v1"

type CreateOrUpdatePodReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
	} `json:"data"` // return data
}

type DeletePodByNameReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

type Namespace struct {
	Name              string `json:"name"`
	CreationTimestamp int64  `json:"creationTimestamp"`
	Status            string `json:"status"`
}

type PodListItem struct {
	Name     string `json:"name"`
	Ready    string `json:"ready"`
	Status   string `json:"status"`
	Restarts int32  `json:"restarts"`
	Age      int64  `json:"age"`
	IP       string `json:"IP"`
	Node     string `json:"node"`
}

type ListNamespacesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Namespaces []*Namespace `json:"namespaces"`
	} `json:"data"` // return data
}

type ListPodsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		PodList []*PodListItem `json:"podList"`
	} `json:"data"` // return data
}
type GetPodDetailReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Pod *Pod `json:"pod"`
	} `json:"data"` // return data
}

type GetPodListRequest struct {
	Namespace string `json:"namespace"`
	Keyword   string `json:"keyword"`
	NodeName  string `json:"nodeName"`
}

type GetPodDetailRequest struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type DeletedPodRequest struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

// node 请求列表请求
type NodeListRequest struct {
	KeyWord string `json:"keyWord"`
}

type NodeDetailRequest struct {
	NodeName string `json:"nodeName"`
}

type GetNodeDetailReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Node *Node `json:"node"`
	} `json:"data"` // return data
}

type ListNodeReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Nodes []*Node `json:"nodes"`
	} `json:"data"` // return data
}
type UpdatedLabelRequest struct {
	Name   string        `json:"name"`
	Labels []ListMapItem `json:"labels"`
}

type UpdateNodeLabelReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
	} `json:"data"` // return data
}

type UpdatedTaintRequest struct {
	Name   string         `json:"name"`
	Taints []corev1.Taint `json:"taints"`
}

type UpdatedTaintReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
	} `json:"data"` // return data
}
