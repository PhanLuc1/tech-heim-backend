package routes

import (
	"github.com/PhanLuc1/tech-heim-backend/controller"
	"github.com/gin-gonic/gin"
)

func ProductRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/product/", controller.GetProduct())
	incomingRoutes.GET("/product/:id", controller.GetProductDetail())
	incomingRoutes.GET("/productGroup", controller.GetProductGroup())
}
