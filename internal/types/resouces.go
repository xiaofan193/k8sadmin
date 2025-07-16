package types

import (
	corev1 "k8s.io/api/core/v1"
)

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
type GetPodDetailReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Pod Pod `json:"pod"`
	} `json:"data"` // return data
}

// 以下是pod 的数据结构
type ListMapItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Base struct {
	//名字
	Name string `json:"name"`
	//标签
	Labels []ListMapItem `json:"labels"`
	//命名空间
	Namespace string `json:"namespace"`
	//重启策略 Always | Never | On-Failure
	RestartPolicy string `json:"restartPolicy"`
}
type NodeSelectorTermExpressions struct {
	Key      string                      `json:"key"`
	Operator corev1.NodeSelectorOperator `json:"operator"`
	Value    string                      `json:"value"`
}
type ConfigMapRefVolume struct {
	Name     string `json:"name"`
	Optional bool   `json:"optional"`
}
type SecretRefVolume struct {
	Name     string `json:"name"`
	Optional bool   `json:"optional"`
}
type NodeScheduling struct {
	//nodeName nodeSelector nodeAffinity
	Type         string                        `json:"type"`
	NodeName     string                        `json:"nodeName"`
	NodeSelector []ListMapItem                 `json:"nodeSelector"`
	NodeAffinity []NodeSelectorTermExpressions `json:"nodeAffinity"`
}
type Volume struct {
	Name string `json:"name"`
	//emptydir | configMap | secret | hostPath | downward | pvc
	Type               string             `json:"type"`
	ConfigMapRefVolume ConfigMapRefVolume `json:"configMapRefVolume"`
	SecretRefVolume    SecretRefVolume    `json:"secretRefVolume"`
	HostPathVolume     HostPathVolume     `json:"hostPathVolume"`
	DownwardAPIVolume  DownwardAPIVolume  `json:"downwardAPIVolume"`
	PVCVolume          PVCVolume          `json:"PVCVolume"`
}
type HostPathVolume struct {
	Type corev1.HostPathType `json:"type"`
	//宿主机路径
	Path string `json:"path"`
}

type DownwardAPIVolumeItem struct {
	Path         string `json:"path"`
	FieldRefPath string `json:"fieldRefPath"`
}

type DownwardAPIVolume struct {
	Items []DownwardAPIVolumeItem `json:"items"`
}
type PVCVolume struct {
	//pvc name
	Name string `json:"name"`
}
type DnsConfig struct {
	Nameservers []string `json:"nameservers"`
}
type NetWorking struct {
	HostNetwork bool          `json:"hostNetwork"`
	HostName    string        `json:"hostName"`
	DnsPolicy   string        `json:"dnsPolicy"`
	DnsConfig   DnsConfig     `json:"dnsConfig"`
	HostAliases []ListMapItem `json:"hostAliases"`
}
type ContainerPort struct {
	Name          string `json:"name"`
	ContainerPort int32  `json:"containerPort"`
	HostPort      int32  `json:"hostPort"`
}
type EnvVar struct {
	Name    string `json:"name"`
	RefName string `json:"refName"`
	Value   string `json:"value"`
	//configMap | secret | default(k/v形式)
	Type string `json:"type"`
}
type EnvVarFromResource struct {
	//资源名称
	Name string `json:"name"`
	//configMap | secret
	RefType string `json:"refType"`
	//用于表示环境变量前缀
	Prefix string `json:"prefix"`
}

type Resources struct {
	//是否配置容器的配额
	Enable bool `json:"enable"`
	//内存 Mi
	MemRequest int32 `json:"memRequest"`
	MemLimit   int32 `json:"memLimit"`
	//cpu m
	CpuRequest int32 `json:"cpuRequest"`
	CpuLimit   int32 `json:"cpuLimit"`
}
type VolumeMount struct {
	//挂载卷名称
	MountName string `json:"mountName"`
	//挂载卷->对应的容器内的路径
	MountPath string `json:"mountPath"`
	//是否只读
	ReadOnly bool `json:"readOnly"`
}
type ProbeTime struct {
	//初始化时间 初始化若干秒之后才开始探针
	InitialDelaySeconds int32 `json:"initialDelaySeconds"`
	//每隔若干秒之后 去探针
	PeriodSeconds int32 `json:"periodSeconds"`
	//探针等待时间 等待若干秒之后还没有返回 那么就是探测失败
	TimeoutSeconds int32 `json:"timeoutSeconds"`
	//探针若干次成功了 才认为这次探针成功
	SuccessThreshold int32 `json:"successThreshold"`
	//探测若干次 失败了 才认为这次探针失败
	FailureThreshold int32 `json:"failureThreshold"`
}
type ProbeHttpGet struct {
	//请求协议http / https
	Scheme string `json:"scheme"`
	//请求host 如果为空 那么就是Pod内请求
	Host string `json:"host"`
	//请求路径
	Path string `json:"path"`
	//请求端口
	Port int32 `json:"port"`
	//请求的header
	HttpHeaders []ListMapItem `json:"httpHeaders"`
}
type ProbeCommand struct {
	// cat /test/test.txt
	Command []string `json:"command"`
}
type ProbeTcpSocket struct {
	//请求host 如果为空 那么就是Pod内请求
	Host string `json:"host"`
	//探测端口
	Port int32 `json:"port"`
}
type ContainerProbe struct {
	//是否打开探针
	Enable bool `json:"enable"`
	//探针类型 tcp / http / exec
	Type      string         `json:"type"`
	HttpGet   ProbeHttpGet   `json:"httpGet"`
	Exec      ProbeCommand   `json:"exec"`
	TcpSocket ProbeTcpSocket `json:"tcpSocket"`
	ProbeTime
}
type Container struct {
	//容器名称
	Name string `json:"name"`
	//容器镜像
	Image string `json:"image"`
	//镜像拉取策略
	ImagePullPolicy string `json:"imagePullPolicy"`
	//是否开启伪终端
	Tty   bool            `json:"tty"`
	Ports []ContainerPort `json:"ports"`
	//工作目录
	WorkingDir string `json:"workingDir"`
	//执行命令
	Command []string `json:"command"`
	//参数
	Args []string `json:"args"`
	//环境变量 [{key:value}]
	Envs     []EnvVar             `json:"envs"`
	EnvsFrom []EnvVarFromResource `json:"envsFrom"`
	//是否开启模式
	Privileged bool `json:"privileged"`
	//容器申请配额
	Resources Resources `json:"resources"`
	//容器卷挂载
	VolumeMounts []VolumeMount `json:"volumeMounts"`
	//启动探针
	StartupProbe ContainerProbe `json:"startupProbe"`
	//存活探针
	LivenessProbe ContainerProbe `json:"livenessProbe"`
	//就绪探针
	ReadinessProbe ContainerProbe `json:"readinessProbe"`
}
type Pod struct {
	//基础定义信息
	Base Base `json:"base"`
	//pod 容忍
	Tolerations    []corev1.Toleration `json:"tolerations"`
	NodeScheduling NodeScheduling      `json:"nodeScheduling"`
	// 卷
	Volumes []Volume `json:"volumes"`
	//网络相关
	NetWorking NetWorking `json:"netWorking"`
	///init containers
	InitContainers []Container `json:"initContainers"`
	//containers
	Containers []Container `json:"containers"`
}
