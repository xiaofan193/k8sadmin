package resouces

import (
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/xiaofan193/k8sadmin/internal/controller"
	"github.com/xiaofan193/k8sadmin/internal/ecode"
	"github.com/xiaofan193/k8sadmin/internal/types/ingress"
)

var _ IngressHandler = (*ingressHandler)(nil)

// ResoucesHandler defining the handler interface
type IngressHandler interface {
	CreateOrUpdateIngress(c *gin.Context)
	GetIngressDetail(c *gin.Context)
	GetIngressList(c *gin.Context)
	DeleteIngress(c *gin.Context)

	CreateOrUpdateIngRoute(c *gin.Context)
	GetIngRouteDetail(c *gin.Context)
	GetIngRouteList(c *gin.Context)
	GetIngRouteMiddlewareList(c *gin.Context)
	DeleteIngRoute(c *gin.Context)
}

type ingressHandler struct {
}

func NewIngressHandler() IngressHandler {
	return &ingressHandler{}
}

func (h *ingressHandler) CreateOrUpdateIngress(c *gin.Context) {
	reqParam := &ingress.CreateOrUpdateIngressRequest{}

	err := c.ShouldBindJSON(reqParam)

	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	err = controller.NewIngressController().CreateOrUpdateIngress(c.Request.Context(), reqParam)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("reqParam", reqParam), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	response.Success(c)
}

func (h *ingressHandler) GetIngressDetail(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	if namespace == "" || name == "" {
		response.Error(c, ecode.InvalidParams, "namespace  and name 不能为空")
		return
	}
	ingressDetail, err := controller.NewIngressController().GetIngressDetail(c.Request.Context(), namespace, name)
	if err != nil {
		logger.Error("GetIngressDetail error", logger.Err(err), logger.Any("reqParam", namespace), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	ingressRes := &ingress.IngressDetailReply{
		Data: struct{ IngressRes *ingress.IngressRes }{IngressRes: ingressDetail},
	}
	response.Success(c, ingressRes)
}

func (h *ingressHandler) GetIngressList(c *gin.Context) {
	namespace := c.Param("namespace")

	if namespace == "" {
		response.Error(c, ecode.InvalidParams, "namespace 不能为空")
		return
	}
	ingressList, err := controller.NewIngressController().GetIngressList(c.Request.Context(), namespace)
	if err != nil {
		logger.Error("GetIngressList error", logger.Err(err), logger.Any("reqParam", namespace), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	ingressRes := &ingress.IngressListReply{
		Data: struct{ List []*ingress.IngressRes }{List: ingressList},
	}
	response.Success(c, ingressRes)
}

func (h *ingressHandler) DeleteIngress(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	if namespace == "" || name == "" {
		response.Error(c, ecode.InvalidParams, "namespace  and name 不能为空")
		return
	}
	err := controller.NewIngressController().DetIngress(c.Request.Context(), namespace, name)

	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("reqParam", namespace), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

func (h *ingressHandler) CreateOrUpdateIngRoute(c *gin.Context) {
	reqParam := &ingress.IngressRouteRequest{}
	err := c.ShouldBind(reqParam)
	if err != nil {
		response.Error(c, ecode.InvalidParams, "namespace  and name 不能为空")
		return
	}
	err = controller.NewIngressController().CreateOrUpdateRoute(c.Request.Context(), reqParam)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("reqParam", reqParam), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)

}

func (h *ingressHandler) GetIngRouteDetail(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	if namespace == "" || name == "" {
		response.Error(c, ecode.InvalidParams, "namespace  and name 不能为空")
		return
	}
	ingressRoute, err := controller.NewIngressController().GetIngRouteDetail(c.Request.Context(), namespace, name)
	if err != nil {
		logger.Error("GetIngRouteDetail error", logger.Err(err), logger.Any("reqParam", name), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	res := &ingress.IngressRouteResReply{
		Data: struct{ IngressRouteRes *ingress.IngressRouteRes }{IngressRouteRes: ingressRoute},
	}
	response.Success(c, res)
}

func (h *ingressHandler) GetIngRouteList(c *gin.Context) {
	namespace := c.Param("namespace")
	keyword := c.Query("keyword")
	if namespace == "" {
		response.Error(c, ecode.InvalidParams, "namespace 不能为空")
		return
	}

	ingressRouteList, err := controller.NewIngressController().GetIngRouteList(c.Request.Context(), namespace, keyword)
	if err != nil {
		logger.Error("GetIngRouteDetail error", logger.Err(err), logger.Any("reqParam", namespace), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	res := ingress.IngressRouteResListReply{
		Data: struct{ Data []*ingress.IngressRouteRes }{Data: ingressRouteList},
	}
	response.Success(c, res)
}

func (h *ingressHandler) GetIngRouteMiddlewareList(c *gin.Context) {
	namespace := c.Param("namespace")
	list, err := controller.NewIngressController().GetIngRouteMiddlewareList(c.Request.Context(), namespace)
	if err != nil {
		logger.Error("GetIngRouteDetail error", logger.Err(err), logger.Any("reqParam", namespace), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	res := ingress.IngressRouteMiddlewareListReply{
		Data: struct{ Data []string }{Data: list},
	}
	response.Success(c, res)
}

func (h *ingressHandler) DeleteIngRoute(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	err := controller.NewIngressController().DeleteIngRoute(c.Request.Context(), namespace, name)
	if err != nil {
		logger.Error("DeleteIngRoute error", logger.Err(err), logger.Any("reqParam", namespace), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	response.Success(c, true)
}
