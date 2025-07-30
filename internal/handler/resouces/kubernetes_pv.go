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

var _ PvHandler = (*pvHandler)(nil)

// ResoucesHandler defining the handler interface
type PvHandler interface {
	CreatePv(c *gin.Context)
	GetPvList(c *gin.Context)
	DeletePV(c *gin.Context)
	CreatePVC(c *gin.Context)
	GetPVCList(c *gin.Context)
	DeletePVC(c *gin.Context)
	CreateSC(c *gin.Context)
	GetSCList(c *gin.Context)
	DeleteSC(c *gin.Context)
}

type pvHandler struct {
}

func NewPvHandler() PvHandler {

	return &pvHandler{}
}

// CreatePv 创建
// @Summary CreatePv 创建pv
// @Description  CreatePv 创建pv
// @Tags pv
// @Accept json
// @Produce json
// @Param data body types.PersistentVolumeRequest struct  true "请求参数"
// @Success 200 {object} types.CreatePersistentVolumeRequestReply{}
// @Router /api/v1/k8s/pv [post]
// @Security BearerAuth
func (h *pvHandler) CreatePv(c *gin.Context) {
	reqParam := &types.PersistentVolumeRequest{}
	err := c.ShouldBind(reqParam)
	if err != nil {
		logger.Warn("CreatePv error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams, err.Error())
		return
	}
	err = controller.NewPvController().Createpv(c.Request.Context(), reqParam)
	if err != nil {
		logger.Error("CreateOrUpdateConfigMap error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}
	response.Success(c)
}

// GetPvList 获取pv的列表
// @Summary GetPvList 获取pv的列表
// @Description 获取pv的列表
// @Tags pv
// @Accept json
// @Produce json
// @Param namespace query string  true "namespace"
// @Param keyword query string  true "keyword"
// @Success 200 {object} types.PersistentVolumeResListReply
// @Router /api/v1/k8s/pv/list [get]
// @Security BearerAuth
func (h *pvHandler) GetPvList(c *gin.Context) {
	keyword := c.Query("keyword")
	list, err := controller.NewPvController().GetPvList(c.Request.Context(), keyword)

	if err != nil {
		logger.Error("GetPvList error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}
	resList := types.PersistentVolumeResListReply{
		Data: struct{ List []*types.PersistentVolumeRes }{List: list},
	}
	response.Success(c, resList)
}

// DeletePV 删除
// @Summary DeletePV 删除pv
// @Description  DeletePV 删除pv
// @Tags pv
// @Accept json
// @Produce json
// @Param name query string   true "请求参数"
// @Success 200 {object} types.CreatePersistentVolumeClaimReply{}
// @Router /api/v1/k8s/pv [delete]
// @Security BearerAuth
func (h *pvHandler) DeletePV(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.Error(c, ecode.InvalidParams, " name 不能为空")
		return
	}
	err := controller.NewPvController().DeletePV(c.Request.Context(), name)
	if err != nil {
		logger.Error("DeletePV error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}
	response.Success(c)
}

// CreatePVC 创建
// @Summary CreatePVC 创建pv
// @Description  CreatePVC 创建pv
// @Tags pv
// @Accept json
// @Produce json
// @Param data body types.PersistentVolumeClaimRequest struct  true "请求参数"
// @Success 200 {object} types.CreatePersistentVolumeClaimReply{}
// @Router /api/v1/k8s/pvc [post]
// @Security BearerAuth
func (h *pvHandler) CreatePVC(c *gin.Context) {
	reqParam := &types.PersistentVolumeClaimRequest{}

	if err := c.ShouldBind(reqParam); err != nil {
		logger.Warn("CreatePVC error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams, err.Error())
		return
	}
	err := controller.NewPvController().CreatePVC(c.Request.Context(), reqParam)
	if err != nil {
		logger.Error("CreatePVC error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}
	response.Success(c)
}

// GetPVCList 获取pv的列表
// @Summary GetPVCList 获取pvc的列表
// @Description 获取pvc的列表
// @Tags pv
// @Accept json
// @Produce json
// @Param namespace query string  true "namespace"
// @Param keyword query string  true "keyword"
// @Success 200 {object} types.PersistentVolumeClaimResListReply
// @Router /api/v1/k8s/pvc/{namespace} [get]
// @Security BearerAuth
func (h *pvHandler) GetPVCList(c *gin.Context) {
	namespace := c.Param("namespace")

	if namespace == "" {
		response.Error(c, ecode.InvalidParams, " namespace 不能为空")
		return
	}

	keyworkd := c.Query("keyword")
	list, err := controller.NewPvController().GetPVCList(c.Request.Context(), namespace, keyworkd)
	if err != nil {
		logger.Error("GetPvList error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}
	resList := types.PersistentVolumeClaimResListReply{
		Data: struct {
			List []*types.PersistentVolumeClaimRes
		}{List: list},
	}
	response.Success(c, resList)
}

// DeletePVC 删除
// @Summary DeletePVC 删除pvc
// @Description  DeletePVC 删除pvc
// @Tags pv
// @Accept json
// @Produce json
// @Param name query string   true "name"
// @Param namespace query string   true "命名空间"
// @Success 200 {object} types.CreatePersistentVolumeClaimReply{}
// @Router /api/v1/k8s/pvc [delete]
// @Security BearerAuth
func (h *pvHandler) DeletePVC(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	if namespace == "" || name == "" {
		response.Error(c, ecode.InvalidParams, "name and  namespace 不能为空")
		return
	}

	err := controller.NewPvController().DeletePVC(c.Request.Context(), namespace, name)
	if err != nil {
		logger.Error("DeletePVC error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}

	response.Success(c)
}

// CreateSC 创建
// @Summary CreateSC 创建pv
// @Description  CreateSC 创建pv
// @Tags pv
// @Accept json
// @Produce json
// @Param data body types.StorageClassRequest struct  true "请求参数"
// @Success 200 {object} types.StorageClassReply{}
// @Router /api/v1/k8s/sc [post]
// @Security BearerAuth
func (h *pvHandler) CreateSC(c *gin.Context) {
	reqParam := &types.StorageClassRequest{}
	if err := c.ShouldBind(reqParam); err != nil {
		logger.Warn("CreateSC error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams, err.Error())
		return
	}
	err := controller.NewPvController().CreateSC(c.Request.Context(), reqParam)

	if err != nil {
		logger.Error("DeletePVC error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}

	response.Success(c)
}

// GetSCList 获取pv的列表
// @Summary GetSCList 获取pvc的列表
// @Description 获取pvc的列表
// @Tags pv
// @Accept json
// @Produce json
// @Param keyword query string  true "keyword"
// @Success 200 {object} types.StorageClassListReply
// @Router /api/v1/k8s/sc/list [get]
// @Security BearerAuth
func (h *pvHandler) GetSCList(c *gin.Context) {
	keyword := c.Query("keyword")
	list, err := controller.NewPvController().GetSCList(c.Request.Context(), keyword)
	if err != nil {
		logger.Error("GetSCList error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}

	resList := types.StorageClassListReply{
		Data: struct {
			List []*types.StorageClassRes `json:"list"`
		}{List: list},
	}
	response.Success(c, resList)
}

// DeleteSC 删除
// @Summary DeleteSC 删除sc
// @Description  DeleteSC 删除sc
// @Tags pv
// @Accept json
// @Produce json
// @Param name query string   true "name"
// @Param namespace query string   true "命名空间"
// @Success 200 {object} types.DeleteStorageClassReply{}
// @Router /api/v1/k8s/sc/{name} [delete]
// @Security BearerAuth
func (h *pvHandler) DeleteSC(c *gin.Context) {

	name := c.Param("name")

	if name == "" {
		response.Error(c, ecode.InvalidParams, " namespace 不能为空")
		return
	}
	err := controller.NewPvController().DeleteSC(c.Request.Context(), name)
	if err != nil {
		logger.Error("DeletePVC error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}

	response.Success(c)
}
