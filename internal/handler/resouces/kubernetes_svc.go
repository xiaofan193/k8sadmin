package resouces

import (
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/xiaofan193/k8sadmin/internal/controller"
	"github.com/xiaofan193/k8sadmin/internal/ecode"
	"github.com/xiaofan193/k8sadmin/internal/types/svc"
)

var _ SvcHandler = (*svcHandler)(nil)

// ResoucesHandler defining the handler interface
type SvcHandler interface {
	CreateOrUpdateSvc(c *gin.Context)
	GetSvcDetail(c *gin.Context)
	GetSvcList(c *gin.Context)
	DeleteSvc(c *gin.Context)
}

type svcHandler struct {
}

func NewSvcHandler() SvcHandler {

	return &svcHandler{}
}

// CreateOrUpdateSvc 创建svc
// @Summary CreateOrUpdateSvc  创建svc
// @Description 打标签
// @Tags svc
// @Accept json
// @Produce json
// @Param data body svc.CreateorUpdateServiceReply struct  true "请求参数"
// @Success 200 {object}
// @Router /api/v1/k8s/svc [post]
// @Security BearerAuth
func (h *svcHandler) CreateOrUpdateSvc(c *gin.Context) {
	reqparam := &svc.CreateorUpdateServiceRequest{}
	err := c.ShouldBind(reqparam)

	if err != nil {
		logger.Warn("CreateOrUpdateSvc error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	err = controller.NewSvcController().CreateOrUpdateSvc(c.Request.Context(), reqparam)

	if err != nil {
		logger.Error("CreateOrUpdateSvc error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}

	response.Success(c)
}

func (h *svcHandler) GetSvcDetail(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	if namespace == "" || name == "" {
		response.Error(c, ecode.InvalidParams, "namespace  and name 不能为空")
		return
	}
	server, err := controller.NewSvcController().GetSvcDetail(c.Request.Context(), namespace, name)

	if err != nil {
		logger.Error("GetSvcDetail error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}

	resServer := svc.ServiceResReply{
		Data: struct{ ServiceRes *svc.ServiceRes }{ServiceRes: server},
	}

	response.Success(c, resServer)
}

func (h *svcHandler) GetSvcList(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		response.Error(c, ecode.InvalidParams, "namespace  不能为空")
		return
	}
	keyword := c.Query("keyword")
	serverList, err := controller.NewSvcController().GetSvcList(c.Request.Context(), namespace, keyword)

	if err != nil {
		logger.Error("GetSvcDetail error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}

	resList := svc.ServerResListReply{
		Data: struct{ List []*svc.ServiceRes }{List: serverList},
	}
	response.Success(c, resList)
}

func (h *svcHandler) DeleteSvc(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	if namespace == "" {
		response.Error(c, ecode.InvalidParams, "namespace  不能为空")
		return
	}

	err := controller.NewSvcController().DeleteSvc(c.Request.Context(), namespace, name)
	if err != nil {
		logger.Error("DeleteSvc error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}

	response.Success(c)

}
