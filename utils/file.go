package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/zetsux/gin-gorm-template-clean/common"
	"github.com/zetsux/gin-gorm-template-clean/internal/dto"
)

func UploadFile(file *multipart.FileHeader, path string) error {
	subPath := strings.Split(path, "/")
	dirPath := fmt.Sprintf("%s/%s", common.FILE_BASE_PATH, subPath[0])

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0777); err != nil {
			return err
		}
	}

	filePath := fmt.Sprintf("%s/%s", dirPath, subPath[1])
	uploadedFile, err := file.Open()
	if err != nil {
		return err
	}
	defer uploadedFile.Close()

	fileData, err := io.ReadAll(uploadedFile)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, fileData, 0666)
	if err != nil {
		return err
	}
	return nil
}

func DeleteFile(path string) error {
	filePath := fmt.Sprintf("%s/%s", common.FILE_BASE_PATH, path)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return dto.ErrFileNotFound
	}

	if err := os.Remove(filePath); err != nil {
		return dto.ErrFileDeleteFailed
	}

	return nil
}
