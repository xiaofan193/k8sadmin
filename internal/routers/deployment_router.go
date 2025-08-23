package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaofan193/k8sadmin/internal/handler/resouces"
)

func initDeploymentRouter(g *gin.RouterGroup) {
	deployApiGroup := resouces.NewDeploymentandler()
	g.POST("/deployment", deployApiGroup.CreateOrUpdateDeployment)
	g.POST("/deployment/:namespace:/:name", deployApiGroup.GetDeploymentDetail)
	g.POST("/deployment/:namespace:", deployApiGroup.GetDeploymentList)
	g.POST("/deployment/:namespace:", deployApiGroup.DeleteDeployment)

}
