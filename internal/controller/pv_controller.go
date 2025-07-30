package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/xiaofan193/k8sadmin/internal/pkg/maputils"
	"github.com/xiaofan193/k8sadmin/internal/types"
	"github.com/xiaofan193/k8sadmin/pkg/global"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strconv"
	"strings"
	"sync"
)

var (
	pvInstance       *PvController
	pvControllerOnce sync.Once
)

type PvController struct {
	KubeConfigSet *kubernetes.Clientset
	CONF          *global.Server
}

func NewPvController() *PvController {
	pvControllerOnce.Do(func() {
		pvInstance = &PvController{
			KubeConfigSet: global.GlobalKubeConfigSet,
			CONF:          global.CONF,
		}
	})
	return pvInstance
}

func (c *PvController) Createpv(ctx context.Context, reqParm *types.PersistentVolumeRequest) error {
	var volumeSource corev1.PersistentVolumeSource
	switch reqParm.VolumeSource.Type {
	case "nfs":
		volumeSource.NFS = &corev1.NFSVolumeSource{
			Server:   reqParm.VolumeSource.NfsVolumeSource.NfsServer,
			Path:     reqParm.VolumeSource.NfsVolumeSource.NfsPath,
			ReadOnly: reqParm.VolumeSource.NfsVolumeSource.NfsReadOnly,
		}
	default:
		return errors.New("不支持的存储类型")
	}

	pv := corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:   reqParm.Name,
			Labels: maputils.ToMap(reqParm.Labels),
		},
		Spec: corev1.PersistentVolumeSpec{
			Capacity: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceStorage: resource.MustParse(strconv.Itoa(int(reqParm.Capacity)) + "Mi"),
			},
			AccessModes:                   reqParm.AccessModes,
			PersistentVolumeReclaimPolicy: reqParm.ReClaimPolicy,
			PersistentVolumeSource:        volumeSource,
		},
	}
	_, err := c.KubeConfigSet.CoreV1().PersistentVolumes().Create(ctx, &pv, metav1.CreateOptions{})
	return err
}

func (c *PvController) GetPvList(ctx context.Context, keyword string) ([]*types.PersistentVolumeRes, error) {
	pvList, err := c.KubeConfigSet.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}
	pvResList := make([]*types.PersistentVolumeRes, 0)
	for _, item := range pvList.Items {
		if !strings.Contains(item.Name, keyword) {
			continue
		}

		claim := ""

		if item.Spec.ClaimRef != nil {
			claim = item.Spec.ClaimRef.Name
		}

		pvRes := &types.PersistentVolumeRes{
			Name:          item.Name,
			Labels:        maputils.ToList(item.Labels),
			Capacity:      int32(item.Spec.Capacity.Storage().Value() / (1024 * 1024)),
			AccessModes:   item.Spec.AccessModes,
			ReClaimPolicy: item.Spec.PersistentVolumeReclaimPolicy,
			Status:        item.Status.Phase,
			Claim:         claim,

			// 当pv 是有sc创建时，就会有改字段
			StorageClassName: item.Spec.StorageClassName,
			Reason:           item.Status.Reason,
			Age:              item.CreationTimestamp.UnixMilli(),
		}
		pvResList = append(pvResList, pvRes)
	}
	return pvResList, nil
}

func (c *PvController) DeletePV(ctx context.Context, name string) error {
	err := c.KubeConfigSet.CoreV1().PersistentVolumes().Delete(ctx, name, metav1.DeleteOptions{})
	return err
}

func (c *PvController) CreatePVC(ctx context.Context, reqParam *types.PersistentVolumeClaimRequest) error {
	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      reqParam.Name,
			Namespace: reqParam.Namespace,
			Labels:    maputils.ToMap(reqParam.Labels),
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: maputils.ToMap(reqParam.Selector),
			},
			AccessModes: reqParam.AccessModes,
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(strconv.Itoa(int(reqParam.Capacity)) + "Mi"),
				},
			},
			StorageClassName: &reqParam.StorageClassName,
		},
	}

	if pvc.Spec.StorageClassName != nil {
		pvc.Spec.Selector = nil
	}

	_, err := c.KubeConfigSet.CoreV1().PersistentVolumeClaims(pvc.Namespace).Create(ctx, &pvc, metav1.CreateOptions{})
	return err
}

func (c *PvController) GetPVCList(ctx context.Context, namespace string, keyword string) ([]*types.PersistentVolumeClaimRes, error) {
	pvcResList := make([]*types.PersistentVolumeClaimRes, 0)
	list, err := c.KubeConfigSet.CoreV1().PersistentVolumeClaims(namespace).List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	for _, item := range list.Items {
		//if strings.Contains(item.Name,keyword) {
		//	continue
		//}
		matchLabels := make([]types.ListMapItem, 0)

		if item.Spec.Selector != nil {
			matchLabels = maputils.ToList(item.Spec.Selector.MatchLabels)
		}

		pvcResItem := &types.PersistentVolumeClaimRes{
			Name:        item.Name,
			Namespace:   item.Namespace,
			Status:      item.Status.Phase,
			Capacity:    int32(item.Spec.Resources.Requests.Storage().Value() / (1024 * 1024)),
			AccessModes: item.Spec.AccessModes,
			Age:         item.CreationTimestamp.UnixMilli(),
			Volume:      item.Spec.VolumeName,
			Labels:      maputils.ToList(item.Labels),
			Selector:    matchLabels,
		}

		pvcResList = append(pvcResList, pvcResItem)
	}
	return pvcResList, nil
}

func (c *PvController) DeletePVC(ctx context.Context, namespace string, name string) error {
	err := c.KubeConfigSet.CoreV1().PersistentVolumeClaims(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	return err
}

func (c *PvController) CreateSC(ctx context.Context, reqParam *types.StorageClassRequest) error {
	// 判断Provisioner是否在系统支持
	provisionerList := strings.Split(c.CONF.System.Provisioner, ",")
	var flag bool
	for _, val := range provisionerList {
		if reqParam.Provisioner == val {
			flag = true
		}
	}
	if !flag {
		return fmt.Errorf("当前K8S未支持！%s", reqParam.Provisioner)
	}

	sc := storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:   reqParam.Name,
			Labels: maputils.ToMap(reqParam.Labels),
		},
		Provisioner:          reqParam.Provisioner,
		MountOptions:         reqParam.MountOptions,
		VolumeBindingMode:    &reqParam.VolumeBindingMode,
		ReclaimPolicy:        &reqParam.ReclaimPolicy,
		AllowVolumeExpansion: &reqParam.AllowVolumeExpansion,
		Parameters:           maputils.ToMap(reqParam.Parameters),
	}

	_, err := c.KubeConfigSet.StorageV1().StorageClasses().Create(ctx, &sc, metav1.CreateOptions{})
	return err
}

func (c *PvController) GetSCList(ctx context.Context, keyword string) ([]*types.StorageClassRes, error) {
	list, err := c.KubeConfigSet.StorageV1().StorageClasses().List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}
	scResList := make([]*types.StorageClassRes, 0)

	for _, item := range list.Items {
		//if strings.Contains(item.Name, keyword) {
		//	continue
		//}

		// item ->response
		var allowVolumeExpansion bool
		if item.AllowVolumeExpansion != nil {
			allowVolumeExpansion = *item.AllowVolumeExpansion
		}

		mountOptions := make([]string, 0)
		if item.MountOptions != nil {
			mountOptions = item.MountOptions
		}
		scResItem := &types.StorageClassRes{
			Name:                 item.Name,
			Labels:               maputils.ToList(item.Labels),
			Provisioner:          item.Provisioner,
			MountOptions:         mountOptions,
			Parameters:           maputils.ToList(item.Parameters),
			ReclaimPolicy:        *item.ReclaimPolicy,
			AllowVolumeExpansion: allowVolumeExpansion,
			Age:                  item.CreationTimestamp.UnixMilli(),
			VolumeBindingMode:    *item.VolumeBindingMode,
		}
		scResList = append(scResList, scResItem)
	}
	return scResList, err
}

func (c *PvController) DeleteSC(ctx context.Context, name string) error {
	return c.KubeConfigSet.StorageV1().StorageClasses().Delete(ctx, name, metav1.DeleteOptions{})
}
