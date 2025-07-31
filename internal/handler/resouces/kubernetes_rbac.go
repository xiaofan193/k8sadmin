package resouces

import (
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/xiaofan193/k8sadmin/internal/controller"
	"github.com/xiaofan193/k8sadmin/internal/ecode"
	"github.com/xiaofan193/k8sadmin/internal/types/rbac"
)

type RBACHandler interface {
	GetServiceAccountList(c *gin.Context)
	CreateServiceAccount(c *gin.Context)
	DeleteServiceAccount(c *gin.Context)
	GetRoleDetail(c *gin.Context)
	GetRoleList(c *gin.Context)
}
type rbacHander struct {
}

func NewRBACHandler() RBACHandler {
	return &rbacHander{}
}

// GetServiceAccountList 获取serverAccount的列表
// @Summary GetServiceAccountList 获取serverAccount的列表
// @Description 获取serverAccount的列表
// @Tags rbac
// @Accept json
// @Produce json
// @Param namespace query string struct true "namespace"
// @Param name query string struct true "name"
// @Success 200 {object} types.ServiceAccountReply
// @Router /api/v1/k8s/sa/{namespace}
// @Security BearerAuth
func (r *rbacHander) GetServiceAccountList(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Query("name")

	if namespace == "" || name == "" {
		response.Error(c, ecode.InvalidParams, "namespace and name 不能为空")
		return
	}
	list, err := controller.NewRbacController().ServiceAccounts(c.Request.Context(), namespace, name)
	if err != nil {
		logger.Error("GetConfigMapDetail error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}
	resList := rbac.ServiceAccountReply{
		Data: struct{ List []*rbac.ServiceAccount }{List: list},
	}
	response.Success(c, resList)
}

// CreateServiceAccount 创建
// @Summary CreateServiceAccount 创建
// @Description  CreateServiceAccount 创建
// @Tags rbac
// @Accept json
// @Produce json
// @Param data body rbac.ServiceAccountRequest   true "请求参数"
// @Success 200 {object} types.DeleteStorageClassReply{}
// @Router /api/v1/k8s/sa/{name}/{namespace} [delete]
// @Security BearerAuth
func (r *rbacHander) CreateServiceAccount(c *gin.Context) {
	reqParam := &rbac.ServiceAccountRequest{}

	if err := c.ShouldBindJSON(reqParam); err != nil {
		response.Error(c, ecode.InvalidParams, err.Error())
		return
	}

	if err := controller.NewRbacController().CreateServiceAccount(c.Request.Context(), reqParam); err != nil {
		logger.Error("CreateServiceAccount error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}

	response.Success(c)
}

// DeleteServiceAccount 删除
// @Summary DeleteServiceAccount 删除sc
// @Description  DeleteServiceAccount 删除sc
// @Tags rbac
// @Accept json
// @Produce json
// @Param name query string   true "name"
// @Param namespace query string   true "命名空间"
// @Success 200 {object} types.DeleteStorageClassReply{}
// @Router /api/v1/k8s/sa/{name}/{namespace} [delete]
// @Security BearerAuth
func (h *rbacHander) DeleteServiceAccount(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	if namespace == "" || name == "" {
		response.Error(c, ecode.InvalidParams, "namespace and name 不能为空")
		return
	}

	err := controller.NewRbacController().DeleteServiceAccount(c.Request.Context(), namespace, name)
	if err != nil {
		logger.Error("DeleteServiceAccount error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}

	response.Success(c)
}

// GetRoleDetail role详情
// @Summary GetRoleDetail role详情
// @Description  GetRoleDetail role详情
// @Tags rbac
// @Accept json
// @Produce json
// @Param name query string   true "name"
// @Param namespace query string   true "命名空间"
// @Success 200 {object} rbac.RoleDetailReply struct
// @Router /api/v1/k8s/role/{name} [get]
// @Security BearerAuth
func (h *rbacHander) GetRoleDetail(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	if name == "" {
		response.Error(c, ecode.InvalidParams, "name 不能为空")
		return
	}
	rbacDetail, err := controller.NewRbacController().GetRoleDetail(c.Request.Context(), namespace, name)

	if err != nil {
		logger.Error("GetRoleDetail error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}
	resRabc := rbac.RoleDetailReply{
		Data: struct{ Role *rbac.RoleRes }{Role: rbacDetail},
	}
	response.Success(c, resRabc)
}

// GetRoleList 获取serverAccount的列表
// @Summary GetRoleList 获取serverAccount的列表
// @Description 获取serverAccount的列表
// @Tags rbac
// @Accept json
// @Produce json
// @Param namespace query string struct true "namespace"
// @Param name query string struct true "name"
// @Success 200 {object} rbac.RoleResListReply
// @Router /api/v1/k8s/roles
// @Security BearerAuth
func (h *rbacHander) GetRoleList(c *gin.Context) {
	namespace := c.Param("namespace")
	roleList, err := controller.NewRbacController().GetRoleList(c.Request.Context(), namespace)

	if err != nil {
		logger.Error("GetRoleDetail error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode(), err)
		return
	}
	roleResList := rbac.RoleResListReply{
		Data: struct{ List []*rbac.Role }{List: roleList},
	}
	response.Success(c, roleResList)
}
