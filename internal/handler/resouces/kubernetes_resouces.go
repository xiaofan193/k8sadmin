package resouces

import (
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
)

var _ ResoucesHandler = (*resoucesHandler)(nil)

// ResoucesHandler defining the handler interface
type ResoucesHandler interface {
	CreateOrUpdatePod(c *gin.Context)
	GetPodList(c *gin.Context)
	DeletePod(c *gin.Context)
	GetNamespaceList(c *gin.Context)
}

type resoucesHandler struct {
}

// NewUserHandler creating the handler interface
func NewResourceHandler() ResoucesHandler {
	return &resoucesHandler{}
}

// CreateOrUpdatePod a new pod
// @Summary CreateOruUpdate a new pod
// @Description Creates a new pod entity using the provided data in the request body.
// @Tags pod
// @Accept json
// @Produce json
// @Param data body types.CreateUserRequest true "pod information"
// @Success 200 {object} types.CreateOrUpdatePodReply{}
// @Router /api/v1/k8s/pod [post]
// @Security BearerAuth
func (h *resoucesHandler) CreateOrUpdatePod(c *gin.Context) {
	response.Success(c, gin.H{})
}

// GetPodList 获取pod列表
// @Summary GetPodList 获取pod列表
// @Description get pod列表
// @Tags pod
// @Accept json
// @Produce json
// @Param namespace query types.CreateUserRequest true "pod information"
// @Success 200 {object} types.ListPodsReply{}
// @Router /api/v1/k8s/{namespace} [get]
// @Security BearerAuth
func (h *resoucesHandler) GetPodList(c *gin.Context) {
	response.Success(c, gin.H{})
}

// DeletePod delete a pod by name
// @Summary Delete a pod by name
// @Description Deletes a existing pod identified by the given name in the path.
// @Tags pod
// @Accept json
// @Produce json
// @Param namespace path string true "namespace"
// @Param name path string true "name"
// @Success 200 {object} types.DeletePodByNameReply{}
// @Router /api/v1/k8s/{namespace}/{name} [delete]
// @Security BearerAuth
func (h *resoucesHandler) DeletePod(c *gin.Context) {
	response.Success(c, gin.H{})
}

// GetNamespaceList 获取namespace列表
// @Summary GetNamespaceList 获取pod列表
// @Description get namespace列表
// @Tags pod
// @Accept json
// @Produce json
// @Success 200 {object} types.ListNamespacesReply{}
// @Router /api/v1/k8s/{namespace} [get]
// @Security BearerAuth
func (h *resoucesHandler) GetNamespaceList(c *gin.Context) {
	response.Success(c, gin.H{})
}
