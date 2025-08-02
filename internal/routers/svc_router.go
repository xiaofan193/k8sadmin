package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaofan193/k8sadmin/internal/handler/resouces"
)

func initSvcRouter(g *gin.RouterGroup) {
	svcApiGroup := resouces.NewSvcHandler()
	g.POST("/svc", svcApiGroup.CreateOrUpdateSvc)
	g.GET("/svc/:namespace/:name", svcApiGroup.GetSvcDetail)
	g.GET("/svc/:namespace", svcApiGroup.GetSvcList)
	g.DELETE("/svc/:namespace/:name", svcApiGroup.DeleteSvc)
}
