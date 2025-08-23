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

var _ DaemonsetHandler = (*daemonsetHandler)(nil)

// ResoucesHandler defining the handler interface
type DaemonsetHandler interface {
	CreateOrUpdateDaemonset(c *gin.Context)
	GetDaemonsetDetail(c *gin.Context)
	GetDaemonsetList(c *gin.Context)
	DeleteDaemonset(c *gin.Context)
}

type daemonsetHandler struct {
}

func NewDaemonsethandler() DaemonsetHandler {
	return &daemonsetHandler{}
}

func (h *daemonsetHandler) CreateOrUpdateDaemonset(c *gin.Context) {
	reqParam := &types.DaemonsetReaqust{}
	err := c.ShouldBind(&reqParam)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return

	}

	err = controller.NewDaemonsetController().CreateOrDaemonset(c.Request.Context(), reqParam)
	if err != nil {
		logger.Error("CreateOrUpdateDaemonset error", logger.Err(err), logger.Any("form", reqParam), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, true)
}
func (h *daemonsetHandler) GetDaemonsetDetail(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	if namespace == "" || name == "" {
		if namespace == "" || name == "" {
			response.Error(c, ecode.InvalidParams, "namespace  and name 不能为空")
			return
		}
	}
	res, err := controller.NewDaemonsetController().GetDaemonsetDetail(c.Request.Context(), namespace, name)

	if err != nil {
		logger.Error("GetDaemonsetDetail error", logger.Err(err), logger.Any("reqParam", namespace), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	daemonRes := &types.DaemonsetDetailReply{
		Data: struct{ Daemonset *types.DaemonSetResonse }{Daemonset: res},
	}
	response.Success(c, daemonRes)
}
func (h *daemonsetHandler) GetDaemonsetList(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		response.Error(c, ecode.InvalidParams, "namespace  and name 不能为空")
		return
	}

	resList, err := controller.NewDaemonsetController().GetDaemonsetList(c.Request.Context(), namespace)

	if err != nil {
		logger.Error("GetDaemonsetList error", logger.Err(err), logger.Any("reqParam", namespace), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	resDeployList := types.DaemonsetListReply{
		Data: struct{ List []*types.DaemonSetRes }{List: resList},
	}

	response.Success(c, resDeployList)
}
func (h *daemonsetHandler) DeleteDaemonset(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	if namespace == "" || name == "" {
		if namespace == "" || name == "" {
			response.Error(c, ecode.InvalidParams, "namespace  and name 不能为空")
			return
		}
	}

	err := controller.NewDaemonsetController().DeleteDaemonset(c.Request.Context(), namespace, name)
	if err != nil {
		logger.Error("CreateOrUpdateDaemonset error", logger.Err(err), logger.Any("name", name), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, true)
}
