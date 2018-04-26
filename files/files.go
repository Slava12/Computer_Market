package files

import (
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

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
	//fileName := generateName(10)

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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

var src = rand.NewSource(time.Now().UnixNano())

func generateName(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
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
