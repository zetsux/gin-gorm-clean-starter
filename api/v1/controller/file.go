package controller

import (
	"net/http"
	"os"
	"strings"

	"github.com/zetsux/gin-gorm-clean-starter/common/base"
	"github.com/zetsux/gin-gorm-clean-starter/common/constant"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/messages"

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
	fileID := ctx.Param("file_id")

	filePath := strings.Join([]string{constant.FileBasePath, dir, fileID}, "/")

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, base.CreateFailResponse(
			messages.MsgFileFetchFailed,
			err.Error(), http.StatusBadRequest,
		))
		return
	}

	ctx.File(filePath)
}
