package files

import (
	"io/ioutil"
	"mime/multipart"
	"os"

	"github.com/Slava12/Computer_Market/logger"
)

// Save сохраняет файл на сервер
func Save(filePath string, fileHeader *multipart.FileHeader, fileName string) {
	file, err := fileHeader.Open()

	bytesOfFile, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Warn(err, "Файл", fileHeader.Filename, "не был прочитан!")
		return
	}
	fullFilePath := ""
	if fileName != "" {
		fullFilePath = filePath + fileName
	} else {
		fullFilePath = filePath + fileHeader.Filename
	}

	fileInServer, err := os.Create(fullFilePath)
	if err != nil {
		logger.Warn(err, "Файл", fullFilePath, "не был создан!")
		return
	}
	logger.Info("Файл", fullFilePath, "был создан.")

	_, err = fileInServer.Write(bytesOfFile)
	if err != nil {
		logger.Warn(err, "Запись файла", fullFilePath, "не удалась!")
		return
	}
	logger.Info("Файл", fullFilePath, "был записан.")

	fileInServer.Close()
	logger.Info("Файл", fullFilePath, "был закрыт.")
}
