package files

import (
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/Slava12/Computer_Market/logger"
)

// Save сохраняет файл на сервер
func Save(filePath string, fileHeader *multipart.FileHeader, fileName string) string {
	file, err := fileHeader.Open()

	bytesOfFile, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Warn(err, "Файл ", fileHeader.Filename, " не был прочитан!")
		return ""
	}
	extension := filepath.Ext(fileHeader.Filename)
	fullFilePath := filePath + fileName + extension

	fileInServer, err := os.Create(fullFilePath)
	if err != nil {
		logger.Warn(err, "Файл ", fullFilePath, " не был создан!")
		return ""
	}
	logger.Info("Файл ", fullFilePath, " был создан.")

	_, err = fileInServer.Write(bytesOfFile)
	if err != nil {
		logger.Warn(err, "Запись файла ", fullFilePath, " не удалась!")
		return ""
	}
	logger.Info("Файл ", fullFilePath, " был записан.")

	fileInServer.Close()
	logger.Info("Файл ", fullFilePath, " был закрыт.")
	return fullFilePath
}

func exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		//logger.Info("Указанный путь", path, " не существует.")
		return false
	}
	//logger.Info("Указанный путь", path, " существует.")
	return true
}

// CreateDirectory создаёт каталог, если не существует
func CreateDirectory(path string) {
	exist := exists(path)
	if exist == false {
		err := os.Mkdir(path, 0777)
		if err != nil {
			//logger.Warn("Директория", path, " не была создана!")
		} else {
			//logger.Info("Директория", path, " была создана.")
		}
	}
}

// RemoveAllFiles удаляет все файлы из указанной директории
func RemoveAllFiles(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		logger.Warn(err, "В директории ", path, " не удалось произвести удаление всех файлов!")
	} else {
		logger.Info("В директории ", path, " были удалены все файлы.")
	}
}
