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
}
