package controller

import (
	"net/http"
	"os"
	"strings"

	"github.com/zetsux/gin-gorm-template-clean/common/standard"
	"github.com/zetsux/gin-gorm-template-clean/internal/dto"

	"github.com/gin-gonic/gin"
)

type fileController struct{}

type FileController interface {
	GetFile(ctx *gin.Context)
}

func NewFileController() FileController {
	return &fileController{}
}

func (fc *fileController) GetFile(ctx *gin.Context) {
	dir := ctx.Param("dir")
	fileId := ctx.Param("file_id")

	filePath := strings.Join([]string{standard.FILE_BASE_PATH, dir, fileId}, "/")

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, standard.CreateFailResponse(
			dto.MESSAGE_FILE_FETCH_FAILED,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.File(filePath)
}
