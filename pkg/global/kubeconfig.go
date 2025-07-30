package global

import (
	"k8s.io/client-go/kubernetes"
)

var (
	CONF                *Server
	GlobalKubeConfigSet *kubernetes.Clientset
	//HarborClient        *harbor.Harbor
)
