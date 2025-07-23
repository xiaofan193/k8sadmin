package types

type ConfigMap struct {
	Name      string        `json:"name"`
	Namespace string        `json:"namespace"`
	Labels    []ListMapItem `json:"labels"`
	Data      []ListMapItem `json:"data"`
}

type ConfigMapRes struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	DataNum   int    `json:"dataNum"`
	Age       int64  `json:"age"`
	//查询configmap详情信息
	Data   []ListMapItem `json:"data"`
	Labels []ListMapItem `json:"labels"`
}

type CreateOrUpdateConfigMapRequest struct {
	Name      string        `json:"name"`
	Namespace string        `json:"namespace"`
	Labels    []ListMapItem `json:"labels"`
	Data      []ListMapItem `json:"data"`
}

type CreateOrUpdateConfigMapReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
	} `json:"data"` // return data
}

type GetConfigMapDetailORListRequest struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Keyword   string `json:"keyword"`
}

type DeleteConfigMapRequest struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}
