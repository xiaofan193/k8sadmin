package types

type DeploymentBase struct {
	Name      string        `json:"name"`
	Namespace string        `json:"namespace"`
	Replicas  int32         `json:"replicas"`
	Labels    []ListMapItem `json:"labels"`
	Selector  []ListMapItem `json:"selector"`
}

type DeploymentRequest struct {
	Base     *DeploymentBase `json:"base"`
	Template *Pod            `json:"template"`
}
type DeploymentResponse struct {
	Base     *DeploymentBase `json:"base"`
	Template *Pod            `json:"template"`
}
type DeploymentRes struct {
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Replicas   int32  `json:"replicas"`
	Ready      int32  `json:"ready"`
	UpdateDate int32  `json:"updateDate"`
	Available  int32  `json:"available"`
	Age        int64  `json:"age"`
}

type DeploymentDetailReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Deployment *DeploymentResponse
	} `json:"data"` // return data
}

type DeploymentListReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		List []*DeploymentRes
	} `json:"data"` // return data
}
