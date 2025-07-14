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
		Namespaces []Namespace `json:"namespaces"`
	} `json:"data"` // return data
}
