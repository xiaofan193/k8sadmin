package resouces

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/xiaofan193/k8sadmin/internal/controller"
	"github.com/xiaofan193/k8sadmin/internal/ecode"
	"github.com/xiaofan193/k8sadmin/internal/types"
)

var _ SecreteMapHandler = (*scretetMapHandler)(nil)

// ResoucesHandler defining the handler interface
type SecreteMapHandler interface {
	CreateOrUpdateSecret(c *gin.Context)
	GetSecretDetail(c *gin.Context)
	GetSecretList(c *gin.Context)
	DeleteSecret(c *gin.Context)
}

type scretetMapHandler struct {
}

func NewSecretHandler() SecreteMapHandler {
	return &scretetMapHandler{}
}

// CreateOrUpdateSecret 创建
// @Summary CreateOrUpdateSecret 创建secret
// @Description  CreateOrUpdateSecret 创建secret
// @Tags secrete
// @Accept json
// @Produce json
// @Param data body types.CreateOrUpadteRequest struct  true "请求参数"
// @Success 200 {object} types.CreateOrUpdateConfigMapReply{}
// @Router /api/v1/k8s/secret [post]
// @Security BearerAuth
func (h *scretetMapHandler) CreateOrUpdateSecret(c *gin.Context) {
	reqParam := &types.CreateOrUpadteSecreteRequest{}
	if err := c.ShouldBindJSON(reqParam); err != nil {
		logger.Warn("CreateOrUpdateConfigMap error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	err := controller.NewSecreteController().CreateOrUpdateSecret(c.Request.Context(), reqParam)
	if err != nil {
		logger.Warn("CreateOrUpdateConfigMap error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}

	response.Success(c)
}

// GetSecretDetail 获取Secrete的详情
// @Summary GetSecretDetail 获取Secrete的详情
// @Description 获取Secrete的详情
// @Tags configMap
// @Accept json
// @Produce json
// @Param namespace query string struct true "namespace"
// @Param name query string struct true "name"
// @Success 200 {object} types.GetSecretReply{}
// @Router api/v1/k8s/secret/{namespace}/{name} [get]
// @Security BearerAuth
func (h *scretetMapHandler) GetSecretDetail(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	if name == "" || namespace == "" {
		err := fmt.Errorf("namespace or name 为空 %s", "")
		logger.Warn("CreateOrUpdateConfigMap error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return

	}
	res, err := controller.NewSecreteController().GetSecretDetail(c.Request.Context(), namespace, name)
	if err != nil {
		logger.Warn("CreateOrUpdateConfigMap error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}
	resDetail := types.GetSecretReply{
		Data: res,
	}
	response.Success(c, resDetail)
}

// GetSecretList 获取secretet的列表
// @Summary GetSecretList 获取secretet的列表
// @Description 获取secretet的列表
// @Tags secrete
// @Accept json
// @Produce json
// @Param namespace query string struct true "namespace"
// @Param keyword query string struct true "keyword"
// @Success 200 {object} []types.ListKeyItemReply
// @Router /api/v1/k8s/secret/:namespace [get]
// @Security BearerAuth
func (h *scretetMapHandler) GetSecretList(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		response.Error(c, ecode.InvalidParams, "namespace  不能为空")
		return
	}
	keyword := c.Query("keyword")
	list, err := controller.NewSecreteController().GetSecretList(c.Request.Context(), namespace, keyword)

	if err != nil {
		logger.Error("GetConfigMapList error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}
	resList := types.ListSecretItemReply{
		Data: struct {
			List []*types.Secret `json:"list"`
		}{List: list},
	}
	response.Success(c, resList)
}

// DeleteSecret 删除Secrete
// @Summary DeleteSecret 删除Secrete
// @Description 获取configMap的详情
// @Tags configMap
// @Accept json
// @Produce json
// @Param namespace query string struct true "namespace"
// @Param name query string struct true "name"
// @Success 200 {object} types.DeleteSecretReply{}
// @Router /api/v1/k8s/configmap/:namespace/:name [delete]
// @Security BearerAuth

func (h *scretetMapHandler) DeleteSecret(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if name == "" || namespace == "" {
		err := fmt.Errorf("namespace or name 为空 %s", "")
		logger.Warn("CreateOrUpdateConfigMap error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return

	}
	err := controller.NewSecreteController().DeleteSecret(c.Request.Context(), namespace, name)
	if err != nil {
		logger.Error("CreateOrUpdateConfigMap error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}
	response.Success(c)

}
