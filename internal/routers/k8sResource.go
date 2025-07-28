package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaofan193/k8sadmin/internal/handler/resouces"
)

func init() {
	ApiV1RouterFns = append(ApiV1RouterFns, func(group *gin.RouterGroup) {
		kubernetsResoucesRouter(group)
	})
}

func kubernetsResoucesRouter(group *gin.RouterGroup) {
	g := group.Group("/k8s")

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithExtraVerify(fn))
	//g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.
	h := resouces.NewResourceHandler()
	g.POST("/pod", h.CreateOrUpdatePod)            // [post] /api/v1/k8s/pod
	g.GET("/pod/:namespace", h.GetPodList)         // [get] /api/v1/k8s/pod/:namespace
	g.GET("/pod/:namespace/:name", h.GetPodDetail) // [get] /api/v1/k8s/pod/:namespace
	g.DELETE("/pod/:namespace/:name", h.DeletePod) // [delete] /api/v1/k8s/pod/:namepace/:name
	g.GET("/namespace", h.GetNamespaceList)        // [post] /api/v1/k8s/namespace

	// node调度

	nh := resouces.NewNodeHandler()
	g.GET("/node", nh.GetNodeList)           // [get] /api/v1/k8s/node
	g.GET("/node/:name", nh.GetNodeDetail)   // [get] /api/v1/k8s/node/:name
	g.PUT("/node/label", nh.UpdateNodeLabel) // [put] /api/v1/k8s/node/label
	g.PUT("/node/taint", nh.UpdateNodeTaint) // [put] /api/v1/k8s/node/taint

	// ConfigMap
	cm := resouces.NewConfigMapHandler()
	g.POST("/configmap", cm.CreateOrUpdateConfigMap)            // [post] /api/v1/k8s/configmap
	g.GET("/configmap/:namespace", cm.GetConfigMapList)         // [get] /api/v1/k8s/configmap/:namespace
	g.GET("/configmap/:namespace/:name", cm.GetConfigMapDetail) // [get] /api/v1/k8s/configmap/:namespace/:name
	g.DELETE("/configmap/:namespace/:name", cm.DeleteConfigMap) // [delete] /api/v1/k8s/configmap/:namespace/:name

	//  Secret
	sh := resouces.NewSecretHandler()
	g.POST("/secret", sh.CreateOrUpdateSecret)            // [post] /api/v1/k8s/secret
	g.GET("/secret/:namespace", sh.GetSecretList)         // [get] /api/v1/k8s/secret/:namespace
	g.GET("/secret/:namespace/:name", sh.GetSecretDetail) // [get] /api/v1/k8s/secret/:namespace/:name
	g.DELETE("/secret/:namespace/:name", sh.DeleteSecret) // [delete] /api/v1/k8s/secret/:namespace/:name
	// Service

}
