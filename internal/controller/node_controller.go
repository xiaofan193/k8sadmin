package controller

import (
	"context"
	"github.com/xiaofan193/k8sadmin/internal/pkg/node"
	"github.com/xiaofan193/k8sadmin/internal/types"
	"github.com/xiaofan193/k8sadmin/pkg/global"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/kubernetes"
	"strings"
	"sync"
)

var (
	nodeInstance       *NodeController
	nodeControllerOnce sync.Once
)

type NodeController struct {
	KubeConfigSet *kubernetes.Clientset
	CONF          Server
}

func NewNodeController() *NodeController {
	nodeControllerOnce.Do(func() {
		nodeInstance = &NodeController{
			KubeConfigSet: global.GlobalKubeConfigSet,
		}
	})
	return nodeInstance
}

func (n *NodeController) GetNodeDetail(ctx context.Context, reqparm *types.NodeDetailRequest) (*types.Node, error) {
	nodeK8s, err := n.KubeConfigSet.CoreV1().Nodes().Get(ctx, reqparm.NodeName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// node 类型转换
	nodeConvret := &node.NodeK8s2Res{}
	detail := nodeConvret.GetNodeDetail(nodeK8s)
	return detail, err
}

func (n *NodeController) GetNodeList(ctx context.Context, reqParam *types.NodeListRequest) ([]*types.Node, error) {
	list, err := n.KubeConfigSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	nodeConvret := &node.NodeK8s2Res{}
	nodeResList := make([]*types.Node, 0)
	for _, item := range list.Items {
		nodeRes := nodeConvret.GetNodeResItem(&item)
		if reqParam.KeyWord != "" {
			if strings.Contains(item.Name, reqParam.KeyWord) {
				nodeResList = append(nodeResList, nodeRes)
			}
		} else {
			nodeResList = append(nodeResList, nodeRes)
		}

	}
	return nodeResList, err
}

func (n *NodeController) UpdateNodeLabel(ctx context.Context, reqParam *types.UpdatedLabelRequest) error {
	labelsMap := make(map[string]string, 0)

	for _, label := range reqParam.Labels {
		labelsMap[label.Key] = label.Value
	}
	labelsMap["$patch"] = "replace"
	patchData := map[string]any{
		"metadata": map[string]any{
			"labels": labelsMap,
		},
	}

	patchDataBytes, _ := json.Marshal(&patchData)
	_, err := n.KubeConfigSet.CoreV1().Nodes().Patch(
		ctx,
		reqParam.Name,
		k8stypes.StrategicMergePatchType,
		patchDataBytes,
		metav1.PatchOptions{},
	)
	return err
}

func (n *NodeController) UpdateNodeTaint(ctx context.Context, reqParm *types.UpdatedTaintRequest) error {
	patchData := map[string]any{
		"spec": map[string]any{
			"taints": reqParm.Taints,
		},
	}

	patchDataBytes, _ := json.Marshal(&patchData)
	_, err := n.KubeConfigSet.CoreV1().Nodes().Patch(
		ctx,
		reqParm.Name,
		k8stypes.StrategicMergePatchType,
		patchDataBytes,
		metav1.PatchOptions{},
	)
	return err
}
