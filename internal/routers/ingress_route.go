package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaofan193/k8sadmin/internal/handler/resouces"
)

func initIngressRouter(g *gin.RouterGroup) {
	svcApiGroup := resouces.NewIngressHandler()

	g.POST("/ingress", svcApiGroup.CreateOrUpdateIngress)
	g.GET("/ingress/:namespace/:name", svcApiGroup.GetIngressDetail)
	g.GET("/ingress/:namespace", svcApiGroup.GetIngressList)
	g.DELETE("/ingress/:namespace/:name", svcApiGroup.DeleteIngress)

}
