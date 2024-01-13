package util

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/zetsux/gin-gorm-clean-starter/common/constant"
	errs "github.com/zetsux/gin-gorm-clean-starter/core/helper/errors"
)

func UploadFile(file *multipart.FileHeader, path string) error {
	subPath := strings.Split(path, "/")
	dirPath := fmt.Sprintf("%s/%s", constant.FileBasePath, subPath[0])

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
	filePath := fmt.Sprintf("%s/%s", constant.FileBasePath, path)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return errs.ErrFileNotFound
	}

	if err := os.Remove(filePath); err != nil {
		return errs.ErrFileDeleteFailed
	}

	return nil
}
