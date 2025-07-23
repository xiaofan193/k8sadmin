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

var _ ConfigMapHandler = (*configMapHandler)(nil)

// ResoucesHandler defining the handler interface
type ConfigMapHandler interface {
	CreateOrUpdateConfigMap(c *gin.Context)
	GetConfigMapDetail(c *gin.Context)
	GetConfigMapList(c *gin.Context)
	DeleteConfigMap(c *gin.Context)
}

type configMapHandler struct {
}

func NewConfigMapHandler() ConfigMapHandler {
	return &configMapHandler{}
}

// CreateOrUpdateConfigMap 创建
// @Summary CreateOrUpdateConfigMap 打标签
// @Description 打标签
// @Tags configMap
// @Accept json
// @Produce json
// @Param data body types.CreateOrUpdateConfigMapRequest struct { true "请求参数"
// @Success 200 {object} types.CreateOrUpdateConfigMapReply{}
// @Router /api/v1/k8s/configmap [post]
// @Security BearerAuth
func (h *configMapHandler) CreateOrUpdateConfigMap(c *gin.Context) {
	reqParam := &types.CreateOrUpdateConfigMapRequest{}
	err := c.ShouldBind(reqParam)
	if err != nil {
		logger.Warn("CreateOrUpdateConfigMap error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	err = controller.NewConfigMapController().CreateOrUpdateConfigMap(c.Request.Context(), reqParam)
	if err != nil {
		logger.Error("CreateOrUpdateConfigMap error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}
	response.Success(c)
}

// GetConfigMapDetail 获取configMap的详情
// @Summary GetConfigMapDetail 获取configMap的详情
// @Description 获取configMap的详情
// @Tags configMap
// @Accept json
// @Produce json
// @Param namespace query string struct true "namespace"
// @Param name query string struct true "configname"
// @Success 200 {object} types.CreateOrUpdateConfigMapReply{}
// @Router /api/v1/k8s/configmap/{namespace}/{name} [get]
// @Security BearerAuth
func (h *configMapHandler) GetConfigMapDetail(c *gin.Context) {

	reqparam := types.GetConfigMapDetailORListRequest{}
	reqparam.Namespace = c.Param("namespace")
	reqparam.Name = c.Param("name")

	if reqparam.Namespace == "" || reqparam.Name == "" {
		response.Error(c, ecode.InvalidParams, "namespace and name 不能为空")
		return
	}
	cm, err := controller.NewConfigMapController().GetConfigMapDetail(c.Request.Context(), &reqparam)
	if err != nil {
		logger.Error("GetConfigMapDetail error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}

	response.Success(c, cm)

}

// GetConfigMapList 获取configMap的列表
// @Summary GetConfigMapList 获取configMap的详情
// @Description 获取configMap的详情
// @Tags configMap
// @Accept json
// @Produce json
// @Param namespace query string struct true "namespace"
// @Param keyword query string struct true "keyword"
// @Success 200 {object} []types.ConfigMapRe
// @Router /api/v1/k8s/configmap/:namespace [get]
// @Security BearerAuth
func (h *configMapHandler) GetConfigMapList(c *gin.Context) {
	reqparam := &types.GetConfigMapDetailORListRequest{}
	reqparam.Namespace = c.Param("namespace")
	reqparam.Namespace = c.Query("keyword")

	if reqparam.Namespace == "" {
		response.Error(c, ecode.InvalidParams, "namespace  不能为空")
		return
	}
	cmList, err := controller.NewConfigMapController().GetConfigMapList(c.Request.Context(), reqparam)
	if err != nil {
		logger.Error("GetConfigMapList error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}
	response.Success(c, cmList)
}

// GetConfigMapList 获取configMap的列表
// @Summary GetConfigMapList 获取configMap的详情
// @Description 获取configMap的详情
// @Tags configMap
// @Accept json
// @Produce json
// @Param namespace query string struct true "namespace"
// @Param name query string struct true "name"
// @Success 200 {object} types.CreateOrUpdateConfigMapReply{}
// @Router /api/v1/k8s/configmap/:namespace/:name
// @Security BearerAuth
func (h *configMapHandler) DeleteConfigMap(c *gin.Context) {
	reqParam := &types.DeleteConfigMapRequest{}
	reqParam.Namespace = c.Param("namespace")
	reqParam.Name = c.Param("name")

	if reqParam.Namespace == "" || reqParam.Name == "" {
		response.Error(c, ecode.InvalidParams, "namespace and name 不能为空")
		return
	}
	err := controller.NewConfigMapController().DeleteConfigMap(c.Request.Context(), reqParam)
	if err != nil {
		logger.Error("DeleteConfigMap error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}
	response.Success(c)
}
