package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zetsux/gin-gorm-clean-starter/api/v1/controller"
)

func FileRouter(route *gin.Engine, fileController controller.FileController) {
	routes := route.Group("/api/v1/files")
	{
		routes.GET("/:dir/:file_id", fileController.GetFile)
	}
}
