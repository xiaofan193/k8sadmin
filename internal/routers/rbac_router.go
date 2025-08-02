package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaofan193/k8sadmin/internal/handler/resouces"
)

// RBAC
func initRBACRouter(g *gin.RouterGroup) {
	rbac := resouces.NewRBACHandler()
	g.GET("/sa/:namespace", rbac.GetServiceAccountList)         //   [get] /api/v1/k8s/sa/:namespace
	g.POST("/sa", rbac.CreateServiceAccount)                    // [post]     /api/v1/k8s/sa
	g.DELETE("/sa/:name/:namespace", rbac.DeleteServiceAccount) // [delete] /api/v1/k8s/sa/:name/:namespace
	// 角色管理
	g.GET("/role/:namespace/:name", rbac.GetRoleDetail) //   [get] /api/v1/k8s/role/:namespace/:name
	g.GET("/role/:namespace", rbac.GetRoleList)         //   [get] /api/v1/k8s/roles/:namespace
	g.POST("/role", rbac.CreateOrUpdateRole)            // [post]     /api/v1/k8s/role

	// 角色绑定
	g.POST("/role/binding", rbac.CreateOrUpdateRoleBingding) // [post]     /api/v1/k8s/role/binding
	g.DELETE("/role/binding", rbac.DeleteRoleBingding)       // [delete] /api/v1/k8s/role/binding
	g.GET("/role/binding/:name", rbac.GetRolbingDetail)      // [get]     /api/v1/k8s/role/binding/:name
	g.GET("/role/binding", rbac.GetRolbingList)              // [get]     /api/v1/k8s/role/binding
}
