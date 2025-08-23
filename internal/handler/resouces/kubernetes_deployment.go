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

var _ DeploymentHandler = (*deploymentHandler)(nil)

// ResoucesHandler defining the handler interface
type DeploymentHandler interface {
	CreateOrUpdateDeployment(c *gin.Context)
	GetDeploymentDetail(c *gin.Context)
	GetDeploymentList(c *gin.Context)
	DeleteDeployment(c *gin.Context)
}

type deploymentHandler struct {
}

func NewDeploymentandler() DeploymentHandler {
	return &deploymentHandler{}
}

func (h *deploymentHandler) CreateOrUpdateDeployment(c *gin.Context) {
	reqParam := &types.DeploymentRequest{}
	err := c.ShouldBind(&reqParam)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return

	}

	err = controller.NewDeploymentController().CreateOrDeployment(c.Request.Context(), reqParam)
	if err != nil {
		logger.Error("CreateOrUpdateDeployment error", logger.Err(err), logger.Any("form", reqParam), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, true)

}
func (h *deploymentHandler) GetDeploymentDetail(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	if namespace == "" || name == "" {
		if namespace == "" || name == "" {
			response.Error(c, ecode.InvalidParams, "namespace  and name 不能为空")
			return
		}
	}
	res, err := controller.NewDeploymentController().GetDeploymentDetail(c.Request.Context(), namespace, name)

	if err != nil {
		logger.Error("GetDeploymentDetail error", logger.Err(err), logger.Any("reqParam", namespace), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	deployRes := &types.DeploymentDetailReply{
		Data: struct{ Deployment *types.DeploymentResponse }{Deployment: res},
	}
	response.Success(c, deployRes)

}
func (h *deploymentHandler) GetDeploymentList(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		response.Error(c, ecode.InvalidParams, "namespace  and name 不能为空")
		return
	}

	resList, err := controller.NewDeploymentController().GetDeploymentList(c.Request.Context(), namespace)

	if err != nil {
		logger.Error("GetDeploymentDetail error", logger.Err(err), logger.Any("reqParam", namespace), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	resDeployList := types.DeploymentListReply{
		Data: struct{ List []*types.DeploymentRes }{List: resList},
	}

	response.Success(c, resDeployList)

}
func (h *deploymentHandler) DeleteDeployment(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	if namespace == "" || name == "" {
		if namespace == "" || name == "" {
			response.Error(c, ecode.InvalidParams, "namespace  and name 不能为空")
			return
		}
	}
	err := controller.NewDeploymentController().DeleteDeployment(c.Request.Context(), namespace, name)

	if err != nil {
		logger.Error("DeleteDeployment error", logger.Err(err), logger.Any("reqParam", namespace), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	response.Success(c, true)
}
