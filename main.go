package main

import (
	"github.com/PhanLuc1/tech-heim-backend/middleware"
	"github.com/PhanLuc1/tech-heim-backend/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	routes.UserRoutes(router)
	routes.ProductRoutes(router)
	router.Run("0.0.0.0:8080")
}
