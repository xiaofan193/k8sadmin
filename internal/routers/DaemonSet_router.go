package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaofan193/k8sadmin/internal/handler/resouces"
)

func initDaemonSetRouter(g *gin.RouterGroup) {
	deployApiGroup := resouces.NewDaemonsethandler()
	g.POST("/daemonset", deployApiGroup.CreateOrUpdateDaemonset)
	g.POST("/daemonset/:namespace:/:name", deployApiGroup.GetDaemonsetDetail)
	g.POST("/daemonset/:namespace:", deployApiGroup.GetDaemonsetList)
	g.POST("/daemonset/:namespace:", deployApiGroup.DeleteDaemonset)

}
