package resouces

import (
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/xiaofan193/k8sadmin/internal/controller"
	"github.com/xiaofan193/k8sadmin/internal/ecode"
	"github.com/xiaofan193/k8sadmin/internal/types"
)

var _ ResoucesHandler = (*resoucesHandler)(nil)

// ResoucesHandler defining the handler interface
type ResoucesHandler interface {
	CreateOrUpdatePod(c *gin.Context)
	GetPodList(c *gin.Context)
	GetPodDetail(c *gin.Context)
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
// @Param data body types.Pod true "pod information"
// @Success 200 {object} types.CreateOrUpdatePodReply{}
// @Router /api/v1/k8s/pod [post]
// @Security BearerAuth
func (h *resoucesHandler) CreateOrUpdatePod(c *gin.Context) {
	podReq := &types.Pod{}
	err := c.ShouldBind(&podReq)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	// 校验必填项
	podValidate := PodValidate{}
	if err = podValidate.Validate(podReq); err != nil {
		logger.Warn("Validate error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams, err)
		return
	}
	ctxg := c.Request.Context()
	msg, err := controller.NewPodController().CreateOrUpdatePod(ctxg, podReq)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("podReq", podReq), logger.String("msg", msg), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
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

// GetPodDetail get a detail by namespace and name
// @Summary Get a by namespace and name
// @Description Gets detailed information of a user specified by the given id in the path.
// @Tags user
// @Param namespace path string true "namespace"
// @Param name path string true "name"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetPodDetailReply{}
// @Router /api/v1/k8s/pod/{namespace}/{name} [get]
// @Security BearerAuth
func (h *resoucesHandler) GetPodDetail(c *gin.Context) {
	response.Success(c, gin.H{})
}
