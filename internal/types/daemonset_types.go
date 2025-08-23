package types

type DaemonsetBase struct {
	Name      string        `json:"name"`
	Namespace string        `json:"namespace"`
	Labels    []ListMapItem `json:"labels"`
	Selector  []ListMapItem `json:"selector"`
}

type DaemonsetReaqust struct {
	Base     *DaemonsetBase `json:"base"`
	Template *Pod           `json:"template"`
}

type DaemonSetResonse struct {
	Base     *DaemonsetBase `json:"base"`
	Template *Pod           `json:"template"`
}

type DaemonSetRes struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Desired   int32  `json:"desired"`   // 需要守护的进程副本数
	Current   int32  `json:"current"`   // 当前正在运行的进程副本
	Ready     int32  `json:"ready"`     // 表示已经就绪的进程副本数
	UpToData  int32  `json:"upToData"`  // 表示已经更新到最新的守护副本进程数
	Available int32  `json:"available"` // 标识可用的守护进程副本数
	Age       int64  `json:"age"`       //
}

type DaemonsetDetailReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Daemonset *DaemonSetResonse
	} `json:"data"` // return data
}

type DaemonsetListReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		List []*DaemonSetRes
	} `json:"data"` // return data
}
