package initial

import (
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/xiaofan193/k8sadmin/internal/controller"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func InitKubeConfigSet() {
	kubeconfig := ".kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	controller.GlobalKubeConfigSet = clientset
	logger.Info("[k8sconfig statistics] was initialized")
}
