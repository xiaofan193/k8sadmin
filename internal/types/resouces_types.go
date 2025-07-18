package types

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
		Namespaces []Namespace `json:"namespaces"`
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
