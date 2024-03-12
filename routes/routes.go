package routes

import (
	"github.com/PhanLuc1/tech-heim-backend/controller"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/user/signup", controller.Signup())
	incomingRoutes.POST("/user/login", controller.Login())
}
