package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/xiaofan193/k8sadmin/internal/handler/resouces"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		kubernetsResoucesRouter(group, resouces.NewResourceHandler())
	})
}

func kubernetsResoucesRouter(group *gin.RouterGroup, h resouces.ResoucesHandler) {
	g := group.Group("/k8s")

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithExtraVerify(fn))
	//g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/pod", h.CreateOrUpdatePod)                 // [post] /api/v1/k8s/:namespace
	g.GET("/:namespace", h.GetPodList)                  // [get] /api/v1/k8s/:namespace/:id
	g.DELETE("/pod/:namepace/:name", h.DeletePod)       // [delete] /api/v1/k8s/pod/:namepace/:name
	g.POST("/namespace", h.GetNamespaceList)            // [post] /api/v1/k8s/namespae
	g.POST("/pod/:namespace/:name", h.GetNamespaceList) // [get] /api/v1/k8s/pod/:namespace/:name
}
